package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	// Setup
	var (
		pd  ProblemDetail
		err error
		s   string
		req *http.Request
		rec *httptest.ResponseRecorder
		c   echo.Context
	)

	e := echo.New()

	testcases := []struct {
		Name   string
		Status int
		URL    string
	}{
		{
			Name:   "NotFoundError",
			Status: http.StatusNotFound,
		},
		{
			Name:   "BadRequestError",
			Status: http.StatusBadRequest,
		},
		{
			Name:   "RedirectError",
			Status: http.StatusTemporaryRedirect,
			URL:    "http://redirect",
		},
		{
			Name:   "RedirectErrorBrowser",
			Status: http.StatusTemporaryRedirect,
			URL:    "http://redirect",
		},
		{
			Name:   "error",
			Status: http.StatusInternalServerError,
		},
		{
			Name:   "*echo.HTTPError",
			Status: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, "/", nil)
			rec = httptest.NewRecorder()
			c = e.NewContext(req, rec)

			switch tc.Name {
			case "NotFoundError":
				err = NotFoundError{Err: fmt.Errorf("error occured"), Request: c.Request()}
				break
			case "BadRequestError":
				err = BadRequestError{Err: fmt.Errorf("error occured"), Request: c.Request()}
				break
			case "RedirectError":
				err = RedirectError{Err: fmt.Errorf("error occured"), Request: c.Request(), URL: tc.URL, Status: http.StatusTemporaryRedirect}
				break
			case "RedirectErrorBrowser":
				req.Header.Add("Accept", "text/html")
				err = RedirectError{Err: fmt.Errorf("error occured"), Request: c.Request(), URL: tc.URL, Status: http.StatusTemporaryRedirect}
				break
			case "error":
				err = fmt.Errorf("error occured")
				break
			case "*echo.HTTPError":
				err = echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error occured"))
				break
			}

			CustomErrorHandler(err, c)
			assert.Equal(t, tc.Status, rec.Code)

			if tc.Name == "RedirectErrorBrowser" {
				assert.Equal(t, tc.URL, rec.Header().Get("Location"))
				return
			}

			s = string(rec.Body.Bytes())
			assert.NotEqual(t, "", s)
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &pd))

			assert.Equal(t, tc.Status, pd.Status)

			if tc.Name == "RedirectError" {
				assert.Equal(t, tc.URL, pd.Instance)
			}
		})
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
