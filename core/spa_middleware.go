package core

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// SpaConfig defines settings for the middleware
type SpaConfig struct {
	// Paths defines which URL paths should be handeled
	Paths []string
	// FilePath defines the filesystem path for index.html
	FilePath string
	// RedirectEmptyPath defines if the index.html file is returned for empty paths
	RedirectEmptyPath bool
}

// SpaWithConfig has the main purpose to return the contents of index.html for all
// paths which return a 404, are requested to be HTML and defined by the configuration
func SpaWithConfig(config SpaConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			err = next(c)
			return processSpaPath(config, c, err)
		}
	}
}

func processSpaPath(config SpaConfig, c echo.Context, err error) error {
	if err != nil && isHTMLContent(c) && isSpaPath(config, c.Path()) {
		log.Debugf("got path '%s'", c.Path())
		if eErr, ok := err.(*echo.HTTPError); ok {
			if eErr.Code == http.StatusNotFound {
				log.Debugf("got %d echo error %v", eErr.Code, eErr)
				return c.File(config.FilePath)
			}
		}
	}
	return err
}

func isHTMLContent(c echo.Context) bool {
	return negotiateContent(c) == HTML
}

func isSpaPath(config SpaConfig, path string) bool {
	if path == "" && config.RedirectEmptyPath {
		return true
	}
	for _, p := range config.Paths {
		if strings.Index(path, p) > -1 {
			return true
		}
	}
	return false
}
