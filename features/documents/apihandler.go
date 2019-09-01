package documents

import (
	"fmt"
	"net/http"
	"strings"

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
// @Failure 400 {object} core.ProblemDetail
// @Failure 404 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/documents/{id} [get]
func (h *Handler) GetDocumentByID(c echo.Context) error {
	var (
		d       DocumentEntity
		err     error
		tags    []string
		senders []string
		cre     string
		mod     string
	)
	id := c.Param("id")
	if d, err = h.r.DocRepo.Get(id); err != nil {
		return core.NotFoundError{Err: err, Request: c.Request()}
	}

	// prepare some values for the API
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
	doc := Document{
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
	return c.JSON(http.StatusOK, doc)
}

// DeleteDocumentByID godoc
// @Summary delete a document by id
// @Description use the supplied id to delete the document from the store
// @Tags documents
// @Param id path string true "document ID"
// @Success 200 {object} documents.Result
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 400 {object} core.ProblemDetail
// @Failure 404 {object} core.ProblemDetail
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

	// before we delete the entry, check if it is really available!
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
