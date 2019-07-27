package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bihe/mydms/config"
	"github.com/bihe/mydms/handler"
	"github.com/bihe/mydms/security"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// Version exports the application version
	Version = "2.0.0"
	// Build provides information about the application build
	Build = "1-local"
	// BuildDate provides information when the application was built
	BuildDate = "2019.07.27 13:45:00"
)

func main() {
	api, addr := setupAPIServer()

	// Initialize the app context that's passed around.
	app := &config.App{
		V: config.VersionInfo{
			Build:     Build,
			Version:   Version,
			BuildDate: BuildDate,
		},
	}

	// Register app (*App) to be injected into all HTTP handlers.
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(config.APP, app)
			return next(c)
		}
	})

	handler.RegisterRoutes(api)

	// Start server
	go func() {
		if err := api.Start(addr); err != nil {
			api.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		api.Logger.Fatal(err)
	}
}

func parseFlags() *ServerArgs {
	c := NewServerArgs()
	flag.StringVar(&c.HostName, "hostname", "localhost", "the server hostname")
	flag.IntVar(&c.Port, "port", 3000, "network port to listen")
	flag.StringVar(&c.ConfigFile, "c", "application.json", "path to the application c file")
	flag.Parse()
	return c
}

func configFromFile(configFileName string) Configuration {
	f, err := os.Open(configFileName)
	if err != nil {
		panic(fmt.Sprintf("Could not open specific config file '%s': %v", configFileName, err))
	}
	defer f.Close()

	c, err := GetSettings(f)
	if err != nil {
		panic(fmt.Sprintf("Could not get server config values from file '%s': %v", configFileName, err))
	}
	return *c
}

func setupAPIServer() (*echo.Echo, string) {
	args := parseFlags()
	c := configFromFile(args.ConfigFile)
	InitLogger(c.Log)

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339_nano}] (${id}) ${method} '${uri}' [${status}] Host: ${host}, IP: ${remote_ip}, error: '${error}', (latency: ${latency_human}) \n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(security.JwtWithConfig(security.JwtOptions{
		JwtSecret:  c.Sec.JwtSecret,
		JwtIssuer:  c.Sec.JwtIssuer,
		CookieName: c.Sec.CookieName,
		RequiredClaim: security.Claim{
			Name:  c.Sec.Claim.Name,
			URL:   c.Sec.Claim.URL,
			Roles: c.Sec.Claim.Roles,
		},
		RedirectURL:   c.Sec.LoginRedirect,
		CacheDuration: c.Sec.CacheDuration,
	}))
	e.Static(c.FS.URLPath, c.FS.Path)

	return e, fmt.Sprintf("%s:%d", args.HostName, args.Port)
}