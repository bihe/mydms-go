package documents

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/features/filestore"
	"github.com/bihe/mydms/features/senders"
	"github.com/bihe/mydms/features/tags"
	"github.com/bihe/mydms/features/upload"
	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

const jsonTimeLayout = "2006-01-02T15:04:05+07:00"

// --------------------------------------------------------------------------
// JSON models
// --------------------------------------------------------------------------

// Document is the json representation of the persistence entity
type Document struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	AltID       string   `json:"alternativeId"`
	Amount      float32  `json:"amount"`
	Created     string   `json:"created"`
	Modified    string   `json:"modified,omitempty"`
	FileName    string   `json:"fileName"`
	PreviewLink string   `json:"previewLink,omitempty"`
	UploadToken string   `json:"uploadFileToken,omitempty"`
	Tags        []string `json:"tags"`
	Senders     []string `json:"senders"`
}

// PagedDcoument represents a paged result
type PagedDcoument struct {
	Documents    []Document `json:"documents"`
	TotalEntries int        `json:"totalEntries"`
}

// ActionResult is a code specifying a specific outcome/result
type ActionResult uint

const (
	// None is the default result
	None ActionResult = iota
	// Created indicates that an item was created
	Created
	// Updated indicates that an item was updated
	Updated
	// Deleted indicates that an item was deleted
	Deleted
	// Error indicates any error
	Error = 99
)

func (a ActionResult) mapToString() string {
	switch a {
	case Created:
		return "Created"
	case Updated:
		return "Updated"
	case Deleted:
		return "Deleted"
	case Error:
		return "Error"
	}
	return "None"
}

func (a ActionResult) mapToAction(s string) ActionResult {
	switch s {
	case "Created":
		return Created
	case "Updated":
		return Updated
	case "Deleted":
		return Deleted
	case "Error":
		return Error
	}
	return None
}

// MarshalJSON marshals the enum as a quoted json string
func (a ActionResult) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(a.mapToString())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (a *ActionResult) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*a = a.mapToAction(j)
	return nil
}

// Result is a generic result object
type Result struct {
	Message string       `json:"message"`
	Result  ActionResult `json:"result"`
}

// --------------------------------------------------------------------------
// Handler definition
// --------------------------------------------------------------------------

// Handler provides handler methods for documents
type Handler struct {
	docRepo    Repository
	tagRepo    tags.Repository
	senderRepo senders.Repository
	uploadRepo upload.Repository
	r          Repositories
	fs         filestore.FileService
	uc         upload.Config
	policy     *bluemonday.Policy
}

// Repositories combines necessary repositories for the document handler
type Repositories struct {
	DocRepo    Repository
	TagRepo    tags.Repository
	SenderRepo senders.Repository
	UploadRepo upload.Repository
}

// NewHandler returns a pointer to a new handler instance
func NewHandler(repos Repositories, fs filestore.FileService, config upload.Config) *Handler {
	return &Handler{
		docRepo:    repos.DocRepo,
		tagRepo:    repos.TagRepo,
		senderRepo: repos.SenderRepo,
		uploadRepo: repos.UploadRepo,
		fs:         fs,
		uc:         config,
		policy:     bluemonday.UGCPolicy(),
	}
}

// GetDocumentByID godoc
// @Summary get a document by id
// @Description use the supplied id to lookup the document from the store
// @Tags documents
// @Param id path string true "document ID"
// @Success 200 {object} documents.Document
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 404 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/documents/{id} [get]
func (h *Handler) GetDocumentByID(c echo.Context) error {
	var (
		d   DocumentEntity
		err error
	)
	id := c.Param("id")
	if d, err = h.docRepo.Get(id); err != nil {
		return core.NotFoundError{Err: err, Request: c.Request()}
	}

	return c.JSON(http.StatusOK, convert(h.policy, d))
}

// DeleteDocumentByID godoc
// @Summary delete a document by id
// @Description use the supplied id to delete the document from the store
// @Tags documents
// @Param id path string true "document ID"
// @Success 200 {object} documents.Result
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/documents/{id} [delete]
func (h *Handler) DeleteDocumentByID(c echo.Context) (err error) {
	id := c.Param("id")

	atomic, err := h.startAtomic(c)
	if err != nil {
		return
	}

	// complete the atomic method
	defer func() {
		err = persistence.HandleTX(true, &atomic, err)
	}()

	fileName, err := h.docRepo.Exists(id, atomic)
	if err != nil {
		log.Warnf("the document '%s' is not available, %v", id, err)
		err = fmt.Errorf("document '%s' not available", id)
		return core.NotFoundError{Err: err, Request: c.Request()}
	}

	err = h.docRepo.Delete(id, atomic)
	if err != nil {
		log.Warnf("error during delete operation of '%s', %v", id, err)
		err = fmt.Errorf("could not delete '%s', %v", id, err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	// also remove the file payload stored in the backend store
	err = h.fs.DeleteFile(fileName)
	if err != nil {
		log.Errorf("could not delete file in backend store '%s', %v", fileName, err)
		err = fmt.Errorf("could not delete '%s', %v", id, err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	return c.JSON(http.StatusOK, Result{
		Message: fmt.Sprintf("Document with id '%s' was deleted.", id),
		Result:  Deleted,
	})
}

// SearchDocuments godoc
// @Summary search for documents
// @Description use filters to search for docments. the result is a paged set
// @Tags documents
// @Param title query string false "title search"
// @Param tag query string false "tag search"
// @Param sender query string false "sender search"
// @Param from query string false "start date"
// @Param to query string false "end date"
// @Param limit query int false "limit max results"
// @Param skip query int false "skip N results"
// @Success 200 {object} documents.PagedDcoument
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/documents/search [get]
func (h *Handler) SearchDocuments(c echo.Context) (err error) {
	var (
		title     string
		tag       string
		sender    string
		fromDate  string
		untilDate string
		limit     int
		skip      int
		order     []OrderBy
	)

	title = c.QueryParam("title")
	tag = c.QueryParam("tag")
	sender = c.QueryParam("sender")
	fromDate = c.QueryParam("from")
	untilDate = c.QueryParam("to")

	// defaults
	limit = parseIntVal(c.QueryParam("limit"), 20)
	skip = parseIntVal(c.QueryParam("skip"), 0)
	orderByTitle := OrderBy{Field: "title", Order: ASC}
	orderByCreated := OrderBy{Field: "created", Order: DESC}

	docs, err := h.docRepo.Search(DocSearch{
		Title:  title,
		Tag:    tag,
		Sender: sender,
		From:   parseDateTime(fromDate),
		Until:  parseDateTime(untilDate),
		Limit:  limit,
		Skip:   skip,
	}, append(order, orderByCreated, orderByTitle))

	if err != nil {
		log.Warnf("could not search for documents, %v", err)
		err = fmt.Errorf("error searching documents, %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	pDoc := PagedDcoument{
		TotalEntries: docs.Count,
		Documents:    convertList(h.policy, docs.Documents),
	}

	return c.JSON(http.StatusOK, pDoc)
}

// SaveDocument godoc
// @Summary save a document
// @Description use the supplied document payload and store the data
// @Tags documents
// @Accept  json
// @Produce  json
// @Param document body documents.Document true "document payload"
// @Success 200 {object} documents.Result
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/documents [post]
func (h *Handler) SaveDocument(c echo.Context) (err error) {
	atomic, err := h.startAtomic(c)
	if err != nil {
		return
	}

	// complete the atomic method
	defer func() {
		err = persistence.HandleTX(true, &atomic, err)
	}()

	d := new(Document)
	if err = c.Bind(d); err != nil {
		log.Warnf("could not bind supplied payload, %v", err)
		err = fmt.Errorf("could not bind supplied data: %v", err)
		return core.BadRequestError{Err: err, Request: c.Request()}
	}

	d = sanitize(h.policy, d)

	d.FileName, err = h.procssUploadFile(d.UploadToken, d.FileName, atomic)
	if err != nil {
		log.Warnf("could not process the uploaded file, %v", err)
		err = fmt.Errorf("upload-file error: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	tagIds, tagList, err := h.processTags(d.Tags, atomic)
	if err != nil {
		log.Warnf("could not process the supplied tags, %v", err)
		err = fmt.Errorf("tag error: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	senderIds, senderList, err := h.processSenders(d.Senders, atomic)
	if err != nil {
		log.Warnf("could not process the supplied senders, %v", err)
		err = fmt.Errorf("senders error: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	var doc DocumentEntity
	newDoc := true
	if d.ID == "" {
		doc = initDocument(d, senderList, tagList)
	} else {
		// supplied ID needs to be checked if exists
		doc, err = h.docRepo.Get(d.ID)
		if err != nil {
			log.Warnf("cannot find document by ID '%s' - create a new entry, %v", d.ID, err)
			doc = initDocument(d, senderList, tagList)
		} else {
			newDoc = false
			log.Infof("will update existing document ID '%s'", d.ID)
			doc.Title = d.Title
			doc.FileName = d.FileName
			doc.PreviewLink = sql.NullString{String: base64.StdEncoding.EncodeToString([]byte(d.FileName)), Valid: true}
			doc.Amount = d.Amount
			doc.SenderList = senderList
			doc.TagList = tagList
		}
	}

	doc, err = h.docRepo.Save(doc, atomic)
	if err != nil {
		log.Errorf("could not save document: %v", err)
		err = fmt.Errorf("error while saving document: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	err = h.docRepo.SaveReferences(doc.ID, tagIds, senderIds, atomic)
	if err != nil {
		log.Errorf("failed to save tags/senders references: %v", err)
		err = fmt.Errorf("error while saving references: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	var r Result
	var code int
	if newDoc {
		r = Result{
			Message: fmt.Sprintf("Created new document '%s' (%s)", doc.Title, doc.ID),
			Result:  Created,
		}
		code = http.StatusCreated
	} else {
		r = Result{
			Message: fmt.Sprintf("Updated existing document '%s' (%s)", doc.Title, doc.ID),
			Result:  Updated,
		}
		code = http.StatusOK
	}
	c.JSON(code, r)
	return
}

// --------------------------------------------------------------------------
// helpers and internal functions
// --------------------------------------------------------------------------

func (h *Handler) procssUploadFile(token, fileName string, atomic persistence.Atomic) (string, error) {
	if token == "" || token == "-" {
		return fileName, nil
	}

	u, err := h.uploadRepo.Read(token)
	if err != nil {
		log.Errorf("could not read upload-file for token '%s', %v", token, err)
		return "", fmt.Errorf("upload token error: %v", err)
	}

	log.Infof("use uploaded file identified by token '%s'", token)

	now := time.Now().UTC()
	folder := now.Format("2006_01_02")

	ext := filepath.Ext(fileName)
	uploadFile := filepath.Join(h.uc.UploadPath, token+ext)
	payload, err := ioutil.ReadFile(uploadFile)
	if err != nil {
		log.Errorf("could not read upload file '%s', %v", uploadFile, err)
		return "", fmt.Errorf("error reading upload-file: %v", err)
	}

	log.Debugf("got upload file '%s' with payload size '%d'!", uploadFile, len(payload))

	item := filestore.FileItem{
		FileName:   fileName,
		FolderName: folder,
		MimeType:   u.MimeType,
		Payload:    payload,
	}
	err = h.fs.SaveFile(item)
	if err != nil {
		log.Errorf("could not save file '%s', %v", uploadFile, err)
		return "", fmt.Errorf("error while saving file: %v", err)
	}

	err = os.Remove(uploadFile)
	if err != nil {
		// this error is ignored, does not invalidate the overall operation
		log.Warnf("could not delete upload-file '%s', %v", uploadFile, err)
	}
	err = h.uploadRepo.Delete(token, atomic)
	if err != nil {
		// this error is ignored, does not invalidate the overall operation
		log.Errorf("could not delete the upload-item by id '%s', %v", token, err)
	}

	return fmt.Sprintf("/%s/%s", folder, fileName), nil
}

func (h *Handler) processTags(tList []string, atomic persistence.Atomic) ([]int, string, error) {
	var tagList []int
	var stringList string
	for _, t := range tList {
		if stringList != "" {
			stringList += ";"
		}
		tag, err := h.tagRepo.GetTagByName(t)
		if err != nil {
			log.Warnf("could not find tag '%s', %v", t, err)
			if tag, err = h.tagRepo.CreateTag(t, atomic); err != nil {
				return nil, "", err
			}
			tagList = append(tagList, tag.ID)
			stringList += t
		} else {
			tagList = append(tagList, tag.ID)
			stringList += tag.Name
		}
	}
	return tagList, stringList, nil
}

func (h *Handler) processSenders(sList []string, atomic persistence.Atomic) ([]int, string, error) {
	var senderList []int
	var stringList string
	for _, s := range sList {
		if stringList != "" {
			stringList += ";"
		}
		sender, err := h.senderRepo.GetSenderByName(s)
		if err != nil {
			log.Warnf("could not find sender '%s', %v", s, err)
			if sender, err = h.senderRepo.CreateSender(s, atomic); err != nil {
				return nil, "", err
			}
			senderList = append(senderList, sender.ID)
			stringList += s
		} else {
			senderList = append(senderList, sender.ID)
			stringList += sender.Name
		}
	}
	return senderList, stringList, nil
}

func (h *Handler) startAtomic(c echo.Context) (persistence.Atomic, error) {
	atomic, err := h.docRepo.CreateAtomic()
	if err != nil {
		log.Errorf("failed to start transaction: %v", err)
		err = fmt.Errorf("could not start atomic operation: %v", err)
		return persistence.Atomic{}, core.ServerError{Err: err, Request: c.Request()}
	}
	return atomic, nil
}

func initDocument(d *Document, sList, tList string) DocumentEntity {
	return DocumentEntity{
		Title:       d.Title,
		FileName:    d.FileName,
		PreviewLink: sql.NullString{String: base64.StdEncoding.EncodeToString([]byte(d.FileName)), Valid: true},
		Amount:      d.Amount,
		SenderList:  sList,
		TagList:     tList,
	}
}

func sanitize(policy *bluemonday.Policy, d *Document) *Document {
	doc := Document{
		ID:     d.ID,
		Amount: d.Amount,
	}
	doc.Title = policy.Sanitize(d.Title)
	doc.AltID = policy.Sanitize(d.AltID)
	doc.Created = policy.Sanitize(d.Created)
	doc.Modified = policy.Sanitize(d.Modified)
	doc.FileName = policy.Sanitize(d.FileName)
	doc.PreviewLink = policy.Sanitize(d.PreviewLink)
	doc.UploadToken = policy.Sanitize(d.UploadToken)

	for _, t := range d.Tags {
		doc.Tags = append(doc.Tags, policy.Sanitize(t))
	}
	for _, s := range d.Senders {
		doc.Senders = append(doc.Senders, policy.Sanitize(s))
	}
	return &doc
}

func convert(policy *bluemonday.Policy, d DocumentEntity) Document {
	var (
		tags    []string
		senders []string
		cre     string
		mod     string
	)

	p := d.PreviewLink
	preview := ""
	if p.Valid {
		preview = p.String
	}
	tags = strings.Split(d.TagList, ";")
	senders = strings.Split(d.SenderList, ";")
	cre = d.Created.Format(jsonTimeLayout)
	if d.Modified.Valid {
		mod = d.Modified.Time.Format(jsonTimeLayout)
	}
	doc := sanitize(policy, &Document{
		ID:          d.ID,
		Title:       d.Title,
		AltID:       d.AltID,
		Amount:      d.Amount,
		Created:     cre,
		Modified:    mod,
		FileName:    d.FileName,
		PreviewLink: preview,
		Tags:        tags,
		Senders:     senders,
	})
	return *doc
}

func convertList(policy *bluemonday.Policy, ds []DocumentEntity) []Document {
	var (
		doc  Document
		docs []Document
	)
	for _, d := range ds {
		doc = convert(policy, d)
		docs = append(docs, doc)
	}
	return docs
}

func parseIntVal(input string, def int) int {
	v, err := strconv.Atoi(input)
	if err != nil {
		return def
	}
	return v
}

func parseDateTime(input string) time.Time {
	t, err := time.Parse(jsonTimeLayout, input)
	if err != nil {
		return time.Time{}
	}
	return t
}
