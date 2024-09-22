package handlercateg

import (
	"context"
	"errors"
	"net/http"

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

type UsersInGameProvider interface {
	CountUsersInGame(context.Context, int) (int, error)
}

type IsUserInLobbyProvider interface {
	IsUserInLobby(context.Context, uuid.UUID) bool
}

func PutJoin(playerRequired int,
	courseProvider CourseProvider,
	ratingProvider RatingProvider,
	gameSaver GameSaver,
	lobbyProvider LobbyGameProvider,
	gameUserLinker GameUserLinker,
	usersLobbyProvider UsersInGameProvider,
	userInLobby IsUserInLobbyProvider,
	randomQuestionsPicker RandomQuestionsPicker,
) echo.HandlerFunc {
	type JoinedGame struct {
		PlayersCount   int `json:"playersCount"`
		PlayerRequired int `json:"playerRequired"`
		CourseId       int `json:"courseId"`
		GameId         int `json:"gameId"`
	}

	return func(c echo.Context) error {
		id := c.Get("id").(uuid.UUID)

		courseId := c.Get("course_id").(int)
		categoryId := c.Get("category_id").(int)

		alreadyInLobby := userInLobby.IsUserInLobby(c.Request().Context(), id)
		if alreadyInLobby {
			return echo.NewHTTPError(http.StatusConflict, "user already in a lobby")
		}

		rating, err := ratingProvider.Rating(c.Request().Context(), id, int(categoryId))
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
			UserId: id,
			Health: 20,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		playersCount, err := usersLobbyProvider.CountUsersInGame(c.Request().Context(), game.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		if playersCount == playerRequired {
			err = pickRandomQuestions(c.Request().Context(), game.CourseId, game.Id, randomQuestionsPicker)
			if err != nil {
				return err
			}
		}

		c.JSON(http.StatusOK, JoinedGame{
			PlayersCount:   playersCount,
			PlayerRequired: playerRequired,
			CourseId:       course.Id,
			GameId:         game.Id,
		})

		return nil
	}
}

type RandomQuestionsPicker interface {
	PickRandomQuestions(ctx context.Context, gameId, courseId, count int) error
}

func pickRandomQuestions(ctx context.Context,
	gameId, courseId int,
	randomQuestionsPicker RandomQuestionsPicker,
) error {
	err := randomQuestionsPicker.PickRandomQuestions(ctx, gameId, courseId, 10)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
	}

	return nil
}
