package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func Auth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header["Authorization"]
			if len(auth) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			bearerToken := strings.Split(auth[0], " ")
			if bearerToken[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}
			tokenString := bearerToken[1]

			claims := make(jwt.MapClaims, 3)

			parser := jwt.Parser{
				SkipClaimsValidation: true,
			}

			_, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrUnexpectedSigningMethod
				}

				return []byte(secret), nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			exp := claims["exp"].(float64)
			if float64(time.Now().Unix()) > exp {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			c.Set("login", claims["login"].(string))
			c.Set("id", claims["id"].(string))
			c.Set("tchr", claims["tchr"].(bool))

			return next(c)
		}
	}
}
