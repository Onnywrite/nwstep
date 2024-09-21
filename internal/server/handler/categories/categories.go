package handlercateg

import (
	"context"
	"errors"
	"net/http"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/labstack/echo/v4"
)

type CategoriesProvider interface {
	Categories(context.Context) ([]models.Category, error)
}

func GetCategories(provider CategoriesProvider) echo.HandlerFunc {
	type Categories struct {
		Categories []models.Category `json:"categories"`
	}

	return func(c echo.Context) error {
		categories, err := provider.Categories(c.Request().Context())
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNoContent, "no categories found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, Categories{categories})

		return nil
	}
}
