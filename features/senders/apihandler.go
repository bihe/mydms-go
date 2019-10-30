package senders

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/internal/errors"
	"github.com/labstack/echo/v4"
)

// Sender is the json representation of the persistence entity
type Sender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Handler provides handler methods for senders
type Handler struct {
	R Repository
}

// GetAllSenders godoc
// @Summary retrieve all senders
// @Description returns all available senders in alphabetical order
// @Tags senders
// @Produce  json
// @Success 200 {array} senders.Sender
// @Failure 401 {object} errors.ProblemDetail
// @Failure 403 {object} errors.ProblemDetail
// @Failure 404 {object} errors.ProblemDetail
// @Router /api/v1/senders [get]
func (h *Handler) GetAllSenders(c echo.Context) error {
	var (
		senders    []SenderEntity
		allSenders []Sender
		err        error
	)
	log.Debug("return all available senders.")

	if senders, err = h.R.GetAllSenders(); err != nil {
		return errors.NotFoundError{Err: err, Request: c.Request()}
	}

	for _, t := range senders {
		allSenders = append(allSenders, Sender{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allSenders)
}

// SearchForSenders godoc
// @Summary search for senders
// @Description returns all senders which match a given search-term
// @Tags senders
// @Produce  json
// @Param name query string true "SearchString"
// @Success 200 {array} senders.Sender
// @Failure 401 {object} errors.ProblemDetail
// @Failure 403 {object} errors.ProblemDetail
// @Failure 404 {object} errors.ProblemDetail
// @Router /api/v1/senders/search [get]
func (h *Handler) SearchForSenders(c echo.Context) error {
	var (
		senders    []SenderEntity
		allSenders []Sender
		s          string
		err        error
	)
	s = c.QueryParam("name")

	log.Debugf("search for senders which match '%s'.", s)

	if senders, err = h.R.SearchSenders(s); err != nil {
		return errors.NotFoundError{Err: err, Request: c.Request()}
	}

	for _, t := range senders {
		allSenders = append(allSenders, Sender{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allSenders)
}
