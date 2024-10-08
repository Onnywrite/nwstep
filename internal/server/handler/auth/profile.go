package handlerauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserByIdProvider interface {
	UserById(context.Context, uuid.UUID) (*models.User, error)
}

func GetProfile(provider UserByIdProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Get("id").(uuid.UUID)

		user, err := provider.UserById(c.Request().Context(), id)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, userProfile{
			Id:        user.Id,
			Login:     user.Login,
			Nickname:  user.Nickname,
			IsTeacher: user.IsTeacher,
			Pts:       user.Pts,
		})

		return nil
	}
}
