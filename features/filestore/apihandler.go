package filestore

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/bihe/mydms/core"
	"github.com/labstack/echo/v4"
)

// Handler defines the filestore API
type Handler struct {
	fs FileService
}

// NewHandler returns a pointer to a new handler instance
func NewHandler(config S3Config) *Handler {
	fs := NewService(config)
	return &Handler{fs}
}

// GetFile returns
func (h *Handler) GetFile(c echo.Context) error {
	path := c.QueryParam("path")
	decodedPath, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		return core.BadRequestError{
			Err:     fmt.Errorf("the supplied path param cannot be decoded. %v", err),
			Request: c.Request()}
	}

	file, err := h.fs.GetFile(string(decodedPath))
	if err != nil {
		return core.NotFoundError{
			Err:     fmt.Errorf("file not found '%s'. %v", decodedPath, err),
			Request: c.Request(),
		}
	}

	return c.Blob(http.StatusOK, file.MimeType, file.Payload)
}
