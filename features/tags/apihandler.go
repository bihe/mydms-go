package tags

import (
	"net/http"

	"github.com/bihe/mydms/internal/errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Tag is the json representation of the persistence entity
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Handler provides handler methods for tags
type Handler struct {
	R Repository
}

// GetAllTags godoc
// @Summary retrieve all tags
// @Description returns all available tags in alphabetical order
// @Tags tags
// @Produce  json
// @Success 200 {array} tags.Tag
// @Failure 401 {object} errors.ProblemDetail
// @Failure 403 {object} errors.ProblemDetail
// @Failure 404 {object} errors.ProblemDetail
// @Router /api/v1/tags [get]
func (h *Handler) GetAllTags(c echo.Context) error {
	var (
		tags    []TagEntity
		allTags []Tag
		err     error
	)
	log.Debug("return all available tags.")

	if tags, err = h.R.GetAllTags(); err != nil {
		return errors.NotFoundError{Err: err, Request: c.Request()}
	}

	for _, t := range tags {
		allTags = append(allTags, Tag{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allTags)
}

// SearchForTags godoc
// @Summary search for tags
// @Description returns all tags which match a given search-term
// @Tags tags
// @Produce  json
// @Param name query string true "SearchString"
// @Success 200 {array} tags.Tag
// @Failure 401 {object} errors.ProblemDetail
// @Failure 403 {object} errors.ProblemDetail
// @Failure 404 {object} errors.ProblemDetail
// @Router /api/v1/tags/search [get]
func (h *Handler) SearchForTags(c echo.Context) error {
	var (
		tags    []TagEntity
		allTags []Tag
		s       string
		err     error
	)
	s = c.QueryParam("name")

	log.Printf("search for tags which match '%s'.", s)

	if tags, err = h.R.SearchTags(s); err != nil {
		return errors.NotFoundError{Err: err, Request: c.Request()}
	}

	for _, t := range tags {
		allTags = append(allTags, Tag{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allTags)
}
