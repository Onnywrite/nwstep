package handlercateg

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CategoryTopProvider interface {
	CategoryTop(ctx context.Context, catId, limit int) ([]models.TopUser, error)
}

type UserTopProvider interface {
	UserTopPosition(ctx context.Context, categoryId int, userId uuid.UUID) (*models.TopUser, error)
}

func GetTop(provider CategoryTopProvider, userTop UserTopProvider) echo.HandlerFunc {
	type Top struct {
		You models.TopUser   `json:"you"`
		Top []models.TopUser `json:"top"`
	}

	return func(c echo.Context) error {
		id := c.Get("id")

		uid, err := uuid.Parse(id.(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		categoryIdStr := c.Param("category_id")

		categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid category id").SetInternal(err)
		}

		topLimitStr := c.QueryParam("top")

		topLimit, err := strconv.ParseInt(topLimitStr, 10, 64)
		if err != nil || topLimit <= 0 {
			topLimit = 10
		}

		top, err := provider.CategoryTop(c.Request().Context(), int(categoryId), int(topLimit))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		userTop, err := userTop.UserTopPosition(c.Request().Context(), int(categoryId), uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, Top{
			You: *userTop,
			Top: top,
		})

		return nil
	}
}
