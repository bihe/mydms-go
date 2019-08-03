package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// rather small PDF payload
// https://stackoverflow.com/questions/17279712/what-is-the-smallest-possible-valid-pdf
const pdfPayload = `%PDF-1.0
1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj 2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj 3 0 obj<</Type/Page/MediaBox[0 0 3 3]>>endobj
xref
0 4
0000000000 65535 f
0000000010 00000 n
0000000053 00000 n
0000000102 00000 n
trailer<</Size 4/Root 1 0 R>>
startxref
149
%EOF
`

// Write(item Upload, a persistence.Atomic) (err error)
// Read(id string) (Upload, error)
// Delete(id string, a persistence.Atomic) (err error)
type mockReaderWriter struct{}

func (m mockReaderWriter) Write(item Upload, a persistence.Atomic) (err error) {
	if item.FileName == "test.pdf" {
		return nil
	}
	return fmt.Errorf("error")
}

func (m mockReaderWriter) Read(id string) (Upload, error) {
	return Upload{}, nil
}

func (m mockReaderWriter) Delete(id string, a persistence.Atomic) (err error) {
	return nil
}

func TestUpload(t *testing.T) {
	// Setup
	var err error

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Errorf("could not create multipart: %v", err)
		return
	}
	io.Copy(part, bytes.NewBuffer([]byte(pdfPayload)))
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	ctype := writer.FormDataContentType()
	req.Header.Add("Content-Type", ctype)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	rw := mockReaderWriter{}

	tmp := ""
	if os.PathSeparator == '\\' {
		tmp = os.Getenv("TEMP")
	} else {
		tmp = "/tmp"
	}

	h := NewHandler(rw, Config{
		AllowedFileTypes: []string{"png", "pdf"},
		MaxUploadSize:    10000,
		UploadPath:       tmp,
	})

	if assert.NoError(t, h.UploadFile(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var r Result
		err := json.Unmarshal(rec.Body.Bytes(), &r)
		if err != nil {
			t.Errorf("could not get valid json: %v", err)
		}
	}
}

func TestUploadFail(t *testing.T) {
	// Setup
	var err error

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Errorf("could not create multipart: %v", err)
		return
	}
	io.Copy(part, bytes.NewBuffer([]byte(pdfPayload)))
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	ctype := writer.FormDataContentType()
	req.Header.Add("Content-Type", ctype)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	rw := mockReaderWriter{}

	tmp := ""
	if os.PathSeparator == '\\' {
		tmp = os.Getenv("TEMP")
	} else {
		tmp = "/tmp"
	}

	// upload size
	h := NewHandler(rw, Config{
		AllowedFileTypes: []string{"png", "pdf"},
		MaxUploadSize:    1,
		UploadPath:       tmp,
	})
	err = h.UploadFile(c)
	if err == nil {
		t.Errorf("expected error upload size!")
	}

	// upload size
	h = NewHandler(rw, Config{
		AllowedFileTypes: []string{"png"},
		MaxUploadSize:    1000,
		UploadPath:       tmp,
	})
	err = h.UploadFile(c)
	if err == nil {
		t.Errorf("expected error file type!")
	}

	// upload destination
	h = NewHandler(rw, Config{
		AllowedFileTypes: []string{"pdf"},
		MaxUploadSize:    1000,
		UploadPath:       "/NOTAVAIL/",
	})
	err = h.UploadFile(c)
	if err == nil {
		t.Errorf("expected error file type!")
	}
}

func TestMissingUploadFile(t *testing.T) {
	// Setup
	var err error

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("FILE__", "test.pdf")
	if err != nil {
		t.Errorf("could not create multipart: %v", err)
		return
	}
	io.Copy(part, bytes.NewBuffer([]byte(pdfPayload)))
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	ctype := writer.FormDataContentType()
	req.Header.Add("Content-Type", ctype)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	rw := mockReaderWriter{}

	tmp := ""
	if os.PathSeparator == '\\' {
		tmp = os.Getenv("TEMP")
	} else {
		tmp = "/tmp"
	}

	// missing file
	h := NewHandler(rw, Config{
		AllowedFileTypes: []string{"png", "pdf"},
		MaxUploadSize:    1000,
		UploadPath:       tmp,
	})
	err = h.UploadFile(c)
	if err == nil {
		t.Errorf("expected error missing file!")
	}
}

func TestUploadPeristenceFail(t *testing.T) {
	// Setup
	var err error

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "error.pdf")
	if err != nil {
		t.Errorf("could not create multipart: %v", err)
		return
	}
	io.Copy(part, bytes.NewBuffer([]byte(pdfPayload)))
	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", body)
	ctype := writer.FormDataContentType()
	req.Header.Add("Content-Type", ctype)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	rw := mockReaderWriter{}

	tmp := ""
	if os.PathSeparator == '\\' {
		tmp = os.Getenv("TEMP")
	} else {
		tmp = "/tmp"
	}

	// missing file
	h := NewHandler(rw, Config{
		AllowedFileTypes: []string{"png", "pdf"},
		MaxUploadSize:    1000,
		UploadPath:       tmp,
	})
	err = h.UploadFile(c)
	if err == nil {
		t.Errorf("expected error persistence!")
	}
}
