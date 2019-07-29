package handler

import (
	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/handler/appinfo"
	"github.com/bihe/mydms/handler/tags"
	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes defines the routes of the available handlers
func RegisterRoutes(e *echo.Echo, repo persistence.Repository, version core.VersionInfo) {
	api := e.Group("/api/v1")

	ai := api.Group("/appinfo")
	aih := &appinfo.Handler{VersionInfo: version}
	ai.GET("", aih.GetAppInfo)

	t := api.Group("/tags")
	th := &tags.Handler{Repo: repo}
	t.GET("", th.GetAllTags)
	t.GET("/search", th.SearchForTags)
}
