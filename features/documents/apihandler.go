package documents

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/features/filestore"
	"github.com/bihe/mydms/features/senders"
	"github.com/bihe/mydms/features/tags"
	"github.com/bihe/mydms/features/upload"
	"github.com/labstack/echo/v4"
)

const jsonTimeLayout = "2006-01-02T15:04:05+07:00"

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

func (a ActionResult) String() string {
	switch a {
	case None:
		return "None"
	case Created:
		return "Created"
	case Updated:
		return "Updated"
	case Deleted:
		return "Deleted"
	case Error:
		return "Error"
	}
	return ""
}

// Result is a generic result object
type Result struct {
	Message string       `json:"message"`
	Result  ActionResult `json:"result"`
}

// Handler provides handler methods for documents
type Handler struct {
	r  Repositories
	fs filestore.FileService
}

// Repositories combines necessary repositories for the document handler
type Repositories struct {
	DocRepo    Repository
	TagRepo    tags.Repository
	SenderRepo senders.Repository
	UploadRepo upload.Repository
}

// NewHandler returns a pointer to a new handler instance
func NewHandler(repos Repositories, fs filestore.FileService) *Handler {
	return &Handler{r: repos, fs: fs}
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
	if d, err = h.r.DocRepo.Get(id); err != nil {
		return core.NotFoundError{Err: err, Request: c.Request()}
	}

	return c.JSON(http.StatusOK, convert(d))
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

	atomic, err := h.r.DocRepo.CreateAtomic()
	if err != nil {
		log.Errorf("failed to start transaction: %v", err)
		err = fmt.Errorf("could not start atomic operation: %v", err)
		return core.ServerError{Err: err, Request: c.Request()}
	}

	// complete the atomic method
	defer func() {
		switch err {
		case nil:
			err = atomic.Commit()
		default:
			log.Errorf("could not complete the transaction: %v", err)
			if e := atomic.Rollback(); e != nil {
				err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}()

	fileName, err := h.r.DocRepo.Exists(id, atomic)
	if err != nil {
		log.Warnf("the document '%s' is not available, %v", id, err)
		err = fmt.Errorf("document '%s' not available", id)
		return core.NotFoundError{Err: err, Request: c.Request()}
	}

	err = h.r.DocRepo.Delete(id, atomic)
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

	docs, err := h.r.DocRepo.Search(DocSearch{
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
		Documents:    convertList(docs.Documents),
	}

	return c.JSON(http.StatusOK, pDoc)
}

func convert(d DocumentEntity) Document {
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
	return Document{
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
	}
}

func convertList(ds []DocumentEntity) []Document {
	var (
		doc  Document
		docs []Document
	)
	for _, d := range ds {
		doc = convert(d)
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
	const jsDateFormat = "2006-01-02T15:04:05+01:00"
	t, err := time.Parse(jsDateFormat, input)
	if err != nil {
		return time.Time{}
	}
	return t
}
