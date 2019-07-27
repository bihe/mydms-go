package security

import "github.com/labstack/echo/v4"

// JwtOptions defines presets for the Authentication handler
// by the default the JWT token is fetched from the Authentication header
// as a fallback it is possible to fetch the token from a specific cookie
type JwtOptions struct {
	// JwtSecret is the jwt signing key
	JwtSecret string
	// JwtIssuer specifies identifies the principal that issued the token
	JwtIssuer string
	// CookieName spedifies the HTTP cookie holding the token
	CookieName string
	// RequiredClaim to access the application
	RequiredClaim Claim
	// RedirectURL forwards the request to an external authentication service
	RedirectURL string
	// CacheDuration defines the duration to cache the JWT token result
	CacheDuration string
}

// User is the authenticated principal extracted from the JWT token
type User struct {
	Username      string
	Roles         []string
	Email         string
	UserID        string
	DisplayName   string
	Authenticated bool
}

// Claim defines the authorization requiremenets
type Claim struct {
	// Name of the applicatiion
	Name string
	// URL of the application
	URL string
	// Roles possible roles
	Roles []string
}

// ServerContext is a application specific context implementation
type ServerContext struct {
	echo.Context
	Identity User
}

// RedirectError is a specific error indicating a necessary redirect
type RedirectError struct {
	// Code is a HTTP status code
	Code int
	// Err is the error message
	Err string
	// URL ist the redirect URL
	URL string
}

func (e RedirectError) Error() string {
	return e.Err
}
