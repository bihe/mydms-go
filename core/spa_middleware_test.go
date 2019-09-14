package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

const htmlPayload = `<html>
<body><p>Hello, World<p></body>
</html>`

const noErrExpected = "no error expected: %v"
const htmlStart = "<html>"
const expectedHTML = "expected a html result"

func TestJwtMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ui", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	tempPath := getTempPath()
	file := "index.html"
	filePath := filepath.Join(tempPath, file)
	ioutil.WriteFile(filePath, []byte(htmlPayload), 0644)
	defer func() {
		os.Remove(filePath)
	}()

	config := SpaConfig{
		Paths:             []string{"/ui", "/abc"},
		FilePath:          filePath,
		RedirectEmptyPath: true,
	}

	h := SpaWithConfig(config)(func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "404")
	})
	req.Header.Set("Accept", "text/html")

	c.SetPath("/ui")
	err := h(c)
	if err != nil {
		t.Errorf(noErrExpected, err)
	}
	result := string(rec.Body.Bytes())
	if result == "" && strings.Index(result, htmlStart) == -1 {
		t.Errorf(expectedHTML)
	}

	// is handled because of config.RedirectEmptyPath
	c.SetPath("")
	err = h(c)
	if err != nil {
		t.Errorf(noErrExpected, err)
	}
	result = string(rec.Body.Bytes())
	if result == "" && strings.Index(result, htmlStart) == -1 {
		t.Errorf(expectedHTML)
	}

	// non-matched path -- resturn a 404
	c.SetPath("/api/v1")
	err = h(c)
	if err == nil {
		t.Errorf("error expected")
	}
	if echoErr, ok := err.(*echo.HTTPError); ok {
		if echoErr.Code != 404 {
			t.Errorf("code 404 expected")
		}
	} else {
		t.Errorf("HTTPError expected")
	}
}

func getTempPath() string {
	dir, _ := ioutil.TempDir("", "mydms")
	return dir
}
