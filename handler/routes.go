package handler

import (
	"fmt"

	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/handler/appinfo"
	"github.com/bihe/mydms/handler/senders"
	"github.com/bihe/mydms/handler/tags"
	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes defines the routes of the available handlers
func RegisterRoutes(e *echo.Echo, repoConn persistence.RepositoryConnection, version core.VersionInfo) {
	var (
		err error
		tr  persistence.TagReader
		sr  persistence.SenderReader
	)

	api := e.Group("/api/v1")

	ai := api.Group("/appinfo")
	aih := &appinfo.Handler{VersionInfo: version}
	ai.GET("", aih.GetAppInfo)

	t := api.Group("/tags")
	tr, err = persistence.NewTagReader(repoConn)
	if err != nil {
		panic(fmt.Sprintf("Could not create persistence tags repository: %v", err))
	}
	th := &tags.Handler{Reader: tr}
	t.GET("", th.GetAllTags)
	t.GET("/search", th.SearchForTags)

	s := api.Group("/senders")
	sr, err = persistence.NewSenderReader(repoConn)
	if err != nil {
		panic(fmt.Sprintf("Could not create persistence senders repository: %v", err))
	}
	sh := &senders.Handler{Reader: sr}
	s.GET("", sh.GetAllSenders)
	s.GET("/search", sh.SearchForSenders)
}
