// +build !prod

package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/core"
)

// InitLogger performs a setup for the logging mechanism
func InitLogger(conf core.LogConfig, e *echo.Echo) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339_nano}] (${id}) ${method} '${uri}' [${status}] Host: ${host}, IP: ${remote_ip}, error: '${error}', (latency: ${latency_human}) \n",
	}))
}
