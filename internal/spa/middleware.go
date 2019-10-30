package spa

import (
	"net/http"
	"strings"

	"github.com/bihe/mydms/internal"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// Config defines settings for the middleware
type Config struct {
	// Paths defines which URL paths should be handeled
	Paths []string
	// FilePath defines the filesystem path for index.html
	FilePath string
	// RedirectEmptyPath defines if the index.html file is returned for empty paths
	RedirectEmptyPath bool
}

// WithConfig has the main purpose to return the contents of index.html for all
// paths which return a 404, are requested to be HTML and defined by the configuration
func WithConfig(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			err = next(c)
			return processSpaPath(config, c, err)
		}
	}
}

func processSpaPath(config Config, c echo.Context, err error) error {
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
	return internal.NegotiateContent(c) == internal.HTML
}

func isSpaPath(config Config, path string) bool {
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
