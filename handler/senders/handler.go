package senders

import (
	"log"
	"net/http"

	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
)

// Sender is the json representation of the persistence entity
type Sender struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Handler provides handler methods for senders
type Handler struct {
	Reader persistence.SenderReader
}

// GetAllSenders godoc
// @Summary retrieve all senders
// @Description returns all available senders in alphabetical order
// @Tags senders
// @Produce  json
// @Success 200 {array} senders.Sender
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/senders [get]
func (h *Handler) GetAllSenders(c echo.Context) error {
	var (
		senders    []persistence.Sender
		allSenders []Sender
		err        error
	)
	log.Printf("return all available senders.")

	if senders, err = h.Reader.GetAllSenders(); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
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
// @Success 200 {array} senders.Sender
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/senders/search [get]
func (h *Handler) SearchForSenders(c echo.Context) error {
	var (
		senders    []persistence.Sender
		allSenders []Sender
		s          string
		err        error
	)
	s = c.QueryParam("name")

	log.Printf("search for senders which match '%s'.", s)

	if senders, err = h.Reader.SearchSenders(s); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	for _, t := range senders {
		allSenders = append(allSenders, Sender{Name: t.Name, ID: t.ID})
	}
	return c.JSON(http.StatusOK, allSenders)
}
