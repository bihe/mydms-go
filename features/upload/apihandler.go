package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/persistence"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Result represents status of the upload opeation
type Result struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Handler defines the upload API
type Handler struct {
	rw     ReaderWriter
	config Config
}

// NewHandler returns a pointer to a new handler instance
func NewHandler(rw ReaderWriter, config Config) *Handler {
	return &Handler{rw: rw, config: config}
}

// UploadFile godoc
// @Summary upload a document
// @Description temporarily stores a file and creates a item in the repository
// @Tags upload
// @Consumes multipart/form-data
// @Produce  json
// @Param file formData file true "file to upload"
// @Success 200 {object} upload.Result
// @Failure 401 {object} core.ProblemDetail
// @Failure 403 {object} core.ProblemDetail
// @Failure 500 {object} core.ProblemDetail
// @Router /api/v1/uploads/file [post]
func (h *Handler) UploadFile(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return core.BadRequestError{Err: fmt.Errorf("no file provided: %v", err), Request: c.Request()}
	}

	if file.Size > h.config.MaxUploadSize {
		return core.BadRequestError{
			Err:     fmt.Errorf("the upload exceeds the maximum size of %d - filesize is: %d", h.config.MaxUploadSize, file.Size),
			Request: c.Request()}
	}

	ext := filepath.Ext(file.Filename)
	ext = strings.Replace(ext, ".", "", 1)
	var typeAllowed = false
	for _, t := range h.config.AllowedFileTypes {
		if t == ext {
			typeAllowed = true
			break
		}
	}
	if !typeAllowed {
		return core.BadRequestError{
			Err:     fmt.Errorf("the uploaded file-type '%s' is not allowed, only use: '%s'", ext, strings.Join(h.config.AllowedFileTypes, ",")),
			Request: c.Request()}
	}
	mimeType := file.Header.Get("Content-Type")

	src, err := file.Open()
	if err != nil {
		return core.BadRequestError{Err: fmt.Errorf("could not open upload file: %v", err), Request: c.Request()}
	}
	defer src.Close()

	// Destination
	id := uuid.New().String()
	var tempFileName = id + "." + ext
	uploadPath := path.Join(h.config.UploadPath, tempFileName)
	dst, err := os.Create(uploadPath)
	if err != nil {
		return core.ServerError{Err: fmt.Errorf("could not create file: %v", err), Request: c.Request()}
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return core.ServerError{Err: fmt.Errorf("could not copy file: %v", err), Request: c.Request()}
	}

	u := Upload{
		ID:       id,
		FileName: file.Filename,
		MimeType: mimeType,
		Created:  time.Now().UTC(),
	}
	err = h.rw.Write(u, persistence.Atomic{})
	if err != nil {
		ioerr := os.Remove(uploadPath)
		if ioerr != nil {
			log.Printf("Clean-Up file-upload. Could not delete temp file: '%s': %v", uploadPath, ioerr)
		}
		return core.ServerError{Err: fmt.Errorf("could not save upload item in store: %v", err), Request: c.Request()}
	}
	c.JSON(http.StatusCreated, Result{Token: id, Message: fmt.Sprintf("File '%s' was uploaded successfully!", file.Filename)})

	return nil
}
