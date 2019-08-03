package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestErrorHandler(t *testing.T) {
	// Setup
	var (
		pd   ProblemDetail
		jerr error
		s    string
		req  *http.Request
		rec  *httptest.ResponseRecorder
		c    echo.Context
	)

	e := echo.New()

	// NotFoundError
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	nf := NotFoundError{Err: fmt.Errorf("error occured"), Request: c.Request()}
	CustomErrorHandler(nf, c)
	s = string(rec.Body.Bytes())
	if s == "" {
		t.Errorf("could not stringify result")
	}
	jerr = json.Unmarshal(rec.Body.Bytes(), &pd)
	if jerr != nil {
		t.Errorf("no result received from error handler")
	}
	if pd.Status != http.StatusNotFound {
		t.Errorf("Wrong status returned")
	}

	// BadRequestError
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	br := BadRequestError{Err: fmt.Errorf("error occured"), Request: c.Request()}
	CustomErrorHandler(br, c)
	s = string(rec.Body.Bytes())
	if s == "" {
		t.Errorf("could not stringify result")
	}
	jerr = json.Unmarshal(rec.Body.Bytes(), &pd)
	if jerr != nil {
		t.Errorf("no result received from error handler")
	}
	if pd.Status != http.StatusBadRequest {
		t.Errorf("Wrong status returned")
	}

	// RedirectError
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	rd := RedirectError{Err: fmt.Errorf("error occured"), Request: c.Request(), URL: "http://redirect"}
	CustomErrorHandler(rd, c)
	s = string(rec.Body.Bytes())
	if s == "" {
		t.Errorf("could not stringify result")
	}
	jerr = json.Unmarshal(rec.Body.Bytes(), &pd)
	if jerr != nil {
		t.Errorf("no result received from error handler")
	}
	if pd.Status != http.StatusTemporaryRedirect {
		t.Errorf("Wrong status returned")
	}
	if pd.Instance != rd.URL {
		t.Errorf("Wrong redirect URL")
	}

	// RedirectError - Browser Client
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Accept", "text/html")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	CustomErrorHandler(rd, c)
	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong status returned")
	}
	if rec.Header().Get("Location") != rd.URL {
		t.Errorf("Wrong redirect URL")
	}

	// error
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	CustomErrorHandler(fmt.Errorf("error occured"), c)
	s = string(rec.Body.Bytes())
	if s == "" {
		t.Errorf("could not stringify result")
	}
	jerr = json.Unmarshal(rec.Body.Bytes(), &pd)
	if jerr != nil {
		t.Errorf("no result received from error handler")
	}
	if pd.Status != http.StatusInternalServerError {
		t.Errorf("Wrong status returned")
	}

	// *echo.HTTPError
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	CustomErrorHandler(echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error occured")), c)
	s = string(rec.Body.Bytes())
	if s == "" {
		t.Errorf("could not stringify result")
	}
	jerr = json.Unmarshal(rec.Body.Bytes(), &pd)
	if jerr != nil {
		t.Errorf("no result received from error handler")
	}
	if pd.Status != http.StatusInternalServerError {
		t.Errorf("Wrong status returned")
	}
}

func TestContentNegotiation(t *testing.T) {

	tests := []struct {
		name   string
		header string
		want   content
	}{{
		name:   "empty",
		header: "",
		want:   JSON,
	}, {
		name:   "html",
		header: "text/html",
		want:   HTML,
	}, {
		name:   "json",
		header: "application/json",
		want:   JSON,
	}, {
		name:   "text",
		header: "text/plain",
		want:   TEXT,
	}, {
		name:   "complext",
		header: "text/plain; q=0.5, application/json, text/x-dvi; q=0.8, text/x-c",
		want:   JSON,
	}}

	e := echo.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderAccept, test.header)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			content := negotiateContent(c)

			if content != test.want {
				t.Errorf("Unexpected value\ngot:  %+v\nwant: %+v", content, test.want)
			}
		})
	}
}
