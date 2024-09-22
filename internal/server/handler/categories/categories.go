package handlercateg

import (
	"context"
	"errors"
	"net/http"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/domain/single"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/labstack/echo/v4"
)

type CategoriesProvider interface {
	Categories(context.Context) ([]models.Category, error)
}

type CategorySaver interface {
	SaveCategory(context.Context, models.Category) (*models.Category, error)
}

func PostCategory(saver CategorySaver) echo.HandlerFunc {
	type Category struct {
		Name          string `json:"name" validate:"required,max=100"`
		Description   string `json:"description" validate:"required,max=200"`
		PhotoUrl      string `json:"photoUrl" validate:"required,max=200"`
		BackgroundUrl string `json:"backgroundUrl" validate:"required,max=200"`
	}

	return func(c echo.Context) error {
		var cat Category
		if err := c.Bind(&cat); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request").SetInternal(err)
		}

		if err := single.V.Struct(cat); err != nil {
			return single.ValidationError(err)
		}

		saved, err := saver.SaveCategory(c.Request().Context(), models.Category{
			Name:          cat.Name,
			Description:   cat.Description,
			PhotoUrl:      cat.PhotoUrl,
			BackgroundUrl: cat.BackgroundUrl,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, saved)

		return nil
	}
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

type CategoryProvider interface {
	CategoryById(context.Context, int) (*models.Category, error)
}

func GetCategory(provider CategoryProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryId := c.Get("category_id").(int)

		category, err := provider.CategoryById(c.Request().Context(), categoryId)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "category not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, category)

		return nil
	}
}
