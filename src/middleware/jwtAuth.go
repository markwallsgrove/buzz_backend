package middleware

import (
	"net/http"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/markwallsgrove/muzz_devops/src/database"
	"github.com/markwallsgrove/muzz_devops/src/models/security"
	"go.uber.org/zap"
)

type JWTAuth struct {
	DB     database.Database
	Secret string
	Logger *zap.Logger
}

type Error struct {
	Error string `json:"error"`
}

var UserIdKey = "userId"
var unauthenticated = Error{Error: "Unauthorised access"}

// Process middleware handler to authenticate the user based on the
// Authorization header.
//
// If the authentication fails the HTTP call  will be ended and the
// next handler will not be called.
//
// If the authentication was successful the next handler will be
// called and the context will contain a key `userId` which holds
// the user ID.
func (a *JWTAuth) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization, found := c.Request().Header["Authorization"]
		if !found || len(authorization) != 1 {
			c.JSON(http.StatusUnauthorized, unauthenticated)
			return nil
		}

		parts := strings.SplitN(authorization[0], " ", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, unauthenticated)
			return nil
		}

		if parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, unauthenticated)
			return nil
		}

		claims, err := security.ValidateToken(parts[1], a.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, unauthenticated)
			return nil
		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			a.Logger.Error("cannot convert subject to user id from jwt", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "internal server error")
		}

		c.Set(UserIdKey, userId)
		return next(c)
	}
}
