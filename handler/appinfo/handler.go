package appinfo

import (
	"log"
	"net/http"

	"github.com/bihe/mydms/config"
	"github.com/bihe/mydms/security"
	"github.com/labstack/echo/v4"
)

// AppInfo provides information of the authenticated user and application meta-data
type AppInfo struct {
	UserInfo    UserInfo    `json:"userInfo"`
	VersionInfo VersionInfo `json:"versionInfo"`
}

// UserInfo provides information about authenticated user
type UserInfo struct {
	// DisplayName of authenticated user
	DisplayName string `json:"displayName"`
	// UserID of authenticated user
	UserID string `json:"userId"`
	// UserName of authenticated user
	UserName string `json:"userName"`
	// Email of authenticated user
	Email string `json:"email"`
	// Roles the authenticated user possesses
	Roles []string `json:"roles"`
}

// Claim defines a permission information for a given URL containing a specific role
type Claim struct {
	// Name of the application
	Name string `json:"name"`
	// URL of the application
	URL string `json:"url"`
	// Role as a form of permission
	Role string `json:"rol"`
}

// VersionInfo contains application meta-data
type VersionInfo struct {
	// Version of the application
	Version string `json:"version"`
	// BuildNumber defines the specific build
	BuildNumber string `json:"buildNumber"`
	// BuildDate specifies the date of the build
	BuildDate string `json:"buildDate"`
}

// GetAppInfo provides information about the application
func GetAppInfo(c echo.Context) error {
	sc := c.(*security.ServerContext)
	app := c.Get(config.APP).(*config.App)
	id := sc.Identity
	log.Printf("Got user: %s, email: %s", id.Username, id.Email)

	a := AppInfo{
		UserInfo: UserInfo{
			DisplayName: id.DisplayName,
			UserID:      id.UserID,
			UserName:    id.Username,
			Email:       id.Email,
			Roles:       id.Roles,
		},
		VersionInfo: VersionInfo{
			Version:     app.V.Version,
			BuildNumber: app.V.Build,
			BuildDate:   app.V.BuildDate,
		},
	}

	return c.JSON(http.StatusOK, a)
}
