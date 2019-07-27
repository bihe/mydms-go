package security

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// JwtWithConfig returns the configured JWT Autch middleware
func JwtWithConfig(options JwtOptions) echo.MiddlewareFunc {
	cache := newMemCache(parseDuration(options.CacheDuration))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			var token string
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				token = strings.Replace(authHeader, "Bearer ", "", 1)
			}
			if token == "" {
				// fallback to get the token via the cookie
				var cookie *http.Cookie
				if cookie, err = c.Request().Cookie(options.CookieName); err != nil {
					// neither the header nor the cookie supplied a jwt token
					return RedirectError{http.StatusUnauthorized, "Invalid authentication, no JWT token present!", options.RedirectURL}
				}
				token = cookie.Value
			}

			// to speed up processing use the cache for token lookups
			var user User
			u := cache.get(token)
			if u != nil {
				// cache hit, put the user in the context
				log.Printf("Cache HIT!")
				sc := &ServerContext{Context: c, Identity: *u}
				return next(sc)
			}

			log.Printf("Cache MISS!")

			var payload JwtTokenPayload
			if payload, err = ParseJwtToken(token, options.JwtSecret, options.JwtIssuer); err != nil {
				log.Printf("Could not decode the JWT token payload: %s", err)
				return RedirectError{http.StatusUnauthorized, fmt.Sprintf("Invalid authentication, could not parse the JWT token: %v", err), options.RedirectURL}
			}
			var roles []string
			if roles, err = Authorize(options.RequiredClaim, payload.Claims); err != nil {
				log.Printf("Insufficient permissions to access the resource: %s", err)
				return RedirectError{http.StatusForbidden, fmt.Sprintf("Invalid authorization: %v", err), options.RedirectURL}
			}

			user = User{
				DisplayName:   payload.DisplayName,
				Email:         payload.Email,
				Roles:         roles,
				UserID:        payload.UserID,
				Username:      payload.UserName,
				Authenticated: true,
			}
			cache.set(token, &user)
			sc := &ServerContext{Context: c, Identity: user}
			return next(sc)
		}
	}
}

func parseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic(fmt.Sprintf("wrong value, cannot parse duration: %v", err))
	}
	return d
}
