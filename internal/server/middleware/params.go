package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func IntParams(params ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, param := range params {
				valueStr := c.Param(param)

				value, err := strconv.ParseInt(valueStr, 10, 64)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "invalid "+param+" param").SetInternal(err)
				}

				c.Set(param, int(value))
			}

			return next(c)
		}
	}
}
