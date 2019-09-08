package main

import (
	"github.com/bihe/mydms/core"
	"github.com/bihe/mydms/features/appinfo"
	"github.com/bihe/mydms/features/documents"
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
		tr tags.Repository
		sr senders.Repository
		ur upload.Repository
		dr documents.Repository
	)

	ur, err = upload.NewRepository(con)
	if err != nil {
		return
	}
	tr, err = tags.NewRepository(con)
	if err != nil {
		return
	}
	sr, err = senders.NewRepository(con)
	if err != nil {
		return
	}
	dr, err = documents.NewRepository(con)
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
	th := &tags.Handler{R: tr}
	t.GET("", th.GetAllTags)
	t.GET("/search", th.SearchForTags)

	// senders
	s := api.Group("/senders")
	sh := &senders.Handler{R: sr}
	s.GET("", sh.GetAllSenders)
	s.GET("/search", sh.SearchForSenders)

	// upload
	u := api.Group("/upload")
	uploadConfig := upload.Config{
		AllowedFileTypes: config.UP.AllowedFileTypes,
		MaxUploadSize:    config.UP.MaxUploadSize,
		UploadPath:       config.UP.UploadPath,
	}
	uh := upload.NewHandler(ur, uploadConfig)
	u.POST("/file", uh.UploadFile)

	// file
	storeSvc := filestore.NewService(filestore.S3Config{
		Region: config.Store.Region,
		Bucket: config.Store.Bucket,
		Key:    config.Store.Key,
		Secret: config.Store.Secret,
	})
	f := api.Group("/file")
	fh := filestore.NewHandler(storeSvc)
	f.GET("", fh.GetFile)
	f.GET("/", fh.GetFile)

	// documents
	d := api.Group("/documents")
	dh := documents.NewHandler(documents.Repositories{
		DocRepo:    dr,
		TagRepo:    tr,
		SenderRepo: sr,
		UploadRepo: ur,
	}, storeSvc, uploadConfig)
	d.GET("/:id", dh.GetDocumentByID)
	d.DELETE("/:id", dh.DeleteDocumentByID)
	d.GET("/search", dh.SearchDocuments)
	d.POST("", dh.SaveDocument)
	d.POST("/", dh.SaveDocument)

	return
}
