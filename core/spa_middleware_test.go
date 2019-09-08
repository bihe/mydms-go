package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

const htmlPayload = `<html>
<body><p>Hello, World<p></body>
</html>`

func TestJwtMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ui", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/ui")
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

	err := h(c)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}

	c.SetPath("")
	err = h(c)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}

	c.SetPath("/api/v1")
	err = h(c)
	if err != nil {
		t.Errorf("no error expected: %v", err)
	}
}

func getTempPath() string {
	dir, _ := ioutil.TempDir("", "mydms")
	return dir
}
