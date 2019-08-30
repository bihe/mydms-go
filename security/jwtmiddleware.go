package security

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/core"
	"github.com/labstack/echo/v4"
)

// JwtWithConfig returns the configured JWT Auth middleware
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
					return core.RedirectError{
						Status:  http.StatusUnauthorized,
						Err:     fmt.Errorf("invalid authentication, no JWT token present"),
						Request: c.Request(),
						URL:     options.RedirectURL,
					}
				}
				token = cookie.Value
			}

			// to speed up processing use the cache for token lookups
			var user User
			u := cache.get(token)
			if u != nil {
				// cache hit, put the user in the context
				log.Debug("Cache HIT!")
				sc := &ServerContext{Context: c, Identity: *u}
				return next(sc)
			}

			log.Debug("Cache MISS!")

			var payload JwtTokenPayload
			if payload, err = ParseJwtToken(token, options.JwtSecret, options.JwtIssuer); err != nil {
				log.Warnf("Could not decode the JWT token payload: %s", err)
				return core.RedirectError{
					Status:  http.StatusUnauthorized,
					Err:     fmt.Errorf("invalid authentication, could not parse the JWT token: %v", err),
					Request: c.Request(),
					URL:     options.RedirectURL,
				}
			}
			var roles []string
			if roles, err = Authorize(options.RequiredClaim, payload.Claims); err != nil {
				log.Warnf("Insufficient permissions to access the resource: %s", err)
				return core.RedirectError{
					Status:  http.StatusForbidden,
					Err:     fmt.Errorf("Invalid authorization: %v", err),
					Request: c.Request(),
					URL:     options.RedirectURL,
				}
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
