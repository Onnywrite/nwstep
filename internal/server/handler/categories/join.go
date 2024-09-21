package handlercateg

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CourseProvider interface {
	Course(context.Context, int) (*models.Course, error)
}

type GameSaver interface {
	SaveGame(context.Context, models.Game) (*models.Game, error)
}

type LobbyGameProvider interface {
	LobbyGame(ctx context.Context, courseId, needUsers int) (*models.Game, error)
}

type GameUserLinker interface {
	LinkUserGame(context.Context, models.UserInGame) error
}

type UsersInLobbyProvider interface {
	CountUsersInLobby(context.Context, int) (int, error)
}

type IsUserInLobbyProvider interface {
	IsUserInLobby(context.Context, uuid.UUID) bool
}

func PutJoin(courseProvider CourseProvider,
	ratingProvider RatingProvider,
	gameSaver GameSaver,
	lobbyProvider LobbyGameProvider,
	gameUserLinker GameUserLinker,
	usersLobbyProvider UsersInLobbyProvider,
	userInLobby IsUserInLobbyProvider,
) echo.HandlerFunc {
	type JoinedGame struct {
		PlayersCount   int `json:"playersCount"`
		PlayerRequired int `json:"playerRequired"`
		CourseId       int `json:"courseId"`
		GameId         int `json:"gameId"`
	}

	return func(c echo.Context) error {
		id := c.Get("id")

		uid, err := uuid.Parse(id.(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		courseIdStr := c.Param("course_id")
		courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid course id").SetInternal(err)
		}

		categoryIdStr := c.Param("category_id")
		categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid course id").SetInternal(err)
		}

		alreadyInLobby := userInLobby.IsUserInLobby(c.Request().Context(), uid)
		if alreadyInLobby {
			return echo.NewHTTPError(http.StatusConflict, "user already in a lobby")
		}

		rating, err := ratingProvider.Rating(c.Request().Context(), uid, int(categoryId))
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			rating = 0
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		course, err := courseProvider.Course(c.Request().Context(), int(courseId))
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "course not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		if course.MinRating > rating {
			return echo.NewHTTPError(http.StatusForbidden, "not enough rating")
		}

		game, err := lobbyProvider.LobbyGame(c.Request().Context(), course.Id, 5)
		if err != nil && !errors.Is(err, cuteql.ErrEmptyResult) {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		if game == nil {
			game = &models.Game{
				CourseId: course.Id,
			}

			game, err = gameSaver.SaveGame(c.Request().Context(), *game)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
			}
		}

		err = gameUserLinker.LinkUserGame(c.Request().Context(), models.UserInGame{
			GameId: game.Id,
			UserId: uid,
			Health: 20,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		playersCount, err := usersLobbyProvider.CountUsersInLobby(c.Request().Context(), game.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, JoinedGame{
			PlayersCount:   playersCount,
			PlayerRequired: 5,
			CourseId:       course.Id,
			GameId:         game.Id,
		})

		return nil
	}
}
