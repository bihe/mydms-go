package documents

import (
	"net/http"
	"strings"

	"github.com/bihe/mydms/core"
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

// Handler provides handler methods for documents
type Handler struct {
	dr Repository
	tr tags.Repository
	sr senders.Repository
	ur upload.Repository
}

// NewHandler returns a pointer to a new handler instance
func NewHandler(dr Repository, tr tags.Repository, sr senders.Repository, ur upload.Repository) *Handler {
	return &Handler{
		dr: dr,
		tr: tr,
		sr: sr,
		ur: ur,
	}
}

// GetDocumentByID godoc
// @Summary get a document by id
// @Description use the supplied id to lookup the document from the store
// @Tags documents
// @Param id path string true "Account ID"
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
	if d, err = h.dr.Get(id); err != nil {
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
