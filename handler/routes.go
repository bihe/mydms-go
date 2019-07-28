package handler

import (
	"github.com/bihe/mydms/handler/appinfo"
	"github.com/bihe/mydms/handler/tags"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes defines the routes of the available handlers
func RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	ai := api.Group("/appinfo")
	ai.GET("", appinfo.GetAppInfo)

	t := api.Group("/tags")
	t.GET("", tags.GetAllTags)
	t.GET("/search", tags.SearchForTags)
}
