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
		s   string
		req *http.Request
		rec *httptest.ResponseRecorder
		c   echo.Context
	)

	e := echo.New()
	errReq := httptest.NewRequest(http.MethodGet, "/", nil)
	redirect := "http://redirect"
	testcases := []struct {
		Name   string
		Status int
		Error  error
	}{
		{
			Name:   "NotFoundError",
			Status: http.StatusNotFound,
			Error:  NotFoundError{Err: fmt.Errorf("error occured"), Request: errReq},
		},
		{
			Name:   "BadRequestError",
			Status: http.StatusBadRequest,
			Error:  BadRequestError{Err: fmt.Errorf("error occured"), Request: errReq},
		},
		{
			Name:   "RedirectError",
			Status: http.StatusTemporaryRedirect,
			Error:  RedirectError{Err: fmt.Errorf("error occured"), Request: errReq, URL: redirect, Status: http.StatusTemporaryRedirect},
		},
		{
			Name:   "RedirectErrorBrowser",
			Status: http.StatusTemporaryRedirect,
			Error:  RedirectError{Err: fmt.Errorf("error occured"), Request: errReq, URL: redirect, Status: http.StatusTemporaryRedirect},
		},
		{
			Name:   "error",
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("error occured"),
		},
		{
			Name:   "*echo.HTTPError",
			Status: http.StatusInternalServerError,
			Error:  echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error occured")),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, "/", nil)
			rec = httptest.NewRecorder()
			c = e.NewContext(req, rec)

			if tc.Name == "RedirectErrorBrowser" {
				req.Header.Add("Accept", "text/html")
			}
			CustomErrorHandler(tc.Error, c)
			assert.Equal(t, tc.Status, rec.Code)

			if tc.Name == "RedirectErrorBrowser" {
				assert.Equal(t, redirect, rec.Header().Get("Location"))
				return
			}

			s = string(rec.Body.Bytes())
			assert.NotEqual(t, "", s)
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &pd))

			assert.Equal(t, tc.Status, pd.Status)

			if tc.Name == "RedirectError" {
				assert.Equal(t, redirect, pd.Instance)
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
