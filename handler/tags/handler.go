package tags

import (
	"log"
	"net/http"

	"github.com/bihe/mydms/config"
	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
)

// Tag is the json representation of the persistence entity
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetAllTags godoc
// @Summary retrieve all tags
// @Description returns all available tags in alphabetical order
// @Tags tags
// @Produce  json
// @Success 200 {array} tags.Tag
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/tags [get]
func GetAllTags(c echo.Context) error {
	var (
		tags    []persistence.Tag
		allTags []Tag
		err     error
	)
	app := c.Get(config.APP).(*config.App)
	log.Printf("return all available tags.")

	if tags, err = app.DB.GetAllTags(); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
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
// @Success 200 {array} tags.Tag
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/tags/search [get]
func SearchForTags(c echo.Context) error {
	var (
		tags    []persistence.Tag
		allTags []Tag
		s       string
		err     error
	)
	app := c.Get(config.APP).(*config.App)
	s = c.QueryParam("name")

	log.Printf("search for tags which match '%s'.", s)

	if tags, err = app.DB.SearchTags(s); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	for _, t := range tags {
		allTags = append(allTags, Tag{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allTags)
}