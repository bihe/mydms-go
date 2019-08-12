package main

import (
	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/features/appinfo"
	"github.com/bihe/mydms/features/filestore"
	"github.com/bihe/mydms/features/senders"
	"github.com/bihe/mydms/features/tags"
	"github.com/bihe/mydms/features/upload"
	"github.com/bihe/mydms/persistence"
	"github.com/labstack/echo/v4"
)

// registerRoutes defines the routes of the available handlers
func registerRoutes(e *echo.Echo, con persistence.Connection, config core.Configuration, version core.VersionInfo) (err error) {
	var (
		tr  tags.Reader
		sr  senders.Reader
		urw upload.ReaderWriter
	)

	urw, err = upload.NewReaderWriter(con)
	if err != nil {
		return
	}
	tr, err = tags.NewReader(con)
	if err != nil {
		return
	}
	sr, err = senders.NewReader(con)
	if err != nil {
		return
	}

	// global API path
	api := e.Group("/api/v1")

	// appinfo
	ai := api.Group("/appinfo")
	aih := &appinfo.Handler{VersionInfo: version}
	ai.GET("", aih.GetAppInfo)

	// tags
	t := api.Group("/tags")
	th := &tags.Handler{Reader: tr}
	t.GET("", th.GetAllTags)
	t.GET("/search", th.SearchForTags)

	// senders
	s := api.Group("/senders")
	sh := &senders.Handler{Reader: sr}
	s.GET("", sh.GetAllSenders)
	s.GET("/search", sh.SearchForSenders)

	// upload
	u := api.Group("/upload")
	uh := upload.NewHandler(urw, upload.Config{
		AllowedFileTypes: config.UP.AllowedFileTypes,
		MaxUploadSize:    config.UP.MaxUploadSize,
		UploadPath:       config.UP.UploadPath,
	})
	u.POST("/file", uh.UploadFile)

	// file
	f := api.Group("/file")
	fh := filestore.NewHandler(filestore.S3Config{
		Region: config.Store.Region,
		Bucket: config.Store.Bucket,
		Key:    config.Store.Key,
		Secret: config.Store.Secret,
	})
	f.GET("", fh.GetFile)
	f.GET("/", fh.GetFile)

	return
}