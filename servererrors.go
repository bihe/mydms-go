package main

import (
	"fmt"
	"net/http"

	"github.com/bihe/mydms/security"
	"github.com/labstack/echo/v4"
	"github.com/markusthoemmes/goautoneg"
)

type content int

const (
	// TEXT content-type requested by client
	TEXT content = iota
	// JSON content-type requested by client
	JSON
	// HTML content-type requested by cleint
	HTML
)

func customHTTPErrorHandler(err error, c echo.Context) {
	// decide based on error-type what to do
	content := negotiateContent(c)

	if he, ok := err.(*echo.HTTPError); ok {
		switch content {
		case JSON:
			c.JSON(he.Code, he)
			break
		default:
			c.String(he.Code, fmt.Sprintf("%v", he))
			break
		}
		return
	}

	if re, ok := err.(security.RedirectError); ok {
		switch content {
		case JSON:
			c.JSON(re.Code, re)
			break
		case HTML:
			c.Redirect(http.StatusTemporaryRedirect, re.URL)
			break
		default:
			c.String(re.Code, fmt.Sprintf("%s; %v", re.Err, re.URL))
			break
		}
		return
	}

	// any other case
	switch content {
	case JSON:
		c.JSON(http.StatusInternalServerError, err.Error())
		break
	default:
		c.String(http.StatusInternalServerError, err.Error())
		break
	}
}

func negotiateContent(c echo.Context) content {
	header := c.Request().Header.Get("Accept")
	if header == "" {
		return JSON // default
	}

	accept := goautoneg.ParseAccept(header)
	if len(accept) == 0 {
		return JSON // default
	}

	// use the first element, because this has the highest priority
	switch accept[0].SubType {
	case "html":
		return HTML
	case "json":
		return JSON
	case "plain":
		return TEXT
	default:
		return JSON
	}
}
