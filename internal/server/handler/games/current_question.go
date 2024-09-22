package handlergames

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/labstack/echo/v4"
)

type UsersInGameProvider interface {
	CountUsersInGame(context.Context, int) (int, error)
}

type GameProvider interface {
	Game(context.Context, int) (*models.Game, error)
}

type GameUpdater interface {
	UpdateGame(context.Context, int, map[string]any) error
}

type GameQuestionProvider interface {
	QuestionByNumber(ctx context.Context, gameId, questionNumber int) (*models.Question, error)
}

type AnswersProvider interface {
	Answers(context.Context, int) ([]models.Answer, error)
}

func GetCurrentQuestion(playersRequired, intervalSeconds int,
	lobby UsersInGameProvider,
	gameProvider GameProvider,
	gameUpdater GameUpdater,
	questionProvider GameQuestionProvider,
	answersProvider AnswersProvider,
) echo.HandlerFunc {
	type Question struct {
		Id       int    `json:"id"`
		Question string `json:"question"`
	}

	type Answer struct {
		Id     int    `json:"id"`
		Answer string `json:"answer"`
	}

	return func(c echo.Context) error {
		gameId := c.Get("game_id").(int)

		// get count of users in the game
		usersCount, err := lobby.CountUsersInGame(c.Request().Context(), gameId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		// if users count is less than needed, we can't give a question
		if usersCount != playersRequired {
			return echo.NewHTTPError(http.StatusForbidden, "not enough players")
		}

		// otherwise, we need to get that game
		game, err := gameProvider.Game(c.Request().Context(), gameId)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "game not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		// and check if it's ended
		if game.EndAt != nil {
			return echo.NewHTTPError(http.StatusGone, "game is over")
		}

		// if it is the first time or not but intervalSeconds have ended, increment question number and reset time
		if game.LastQuestionTime == nil || time.Since(*game.LastQuestionTime) < time.Duration(intervalSeconds)*time.Second {
			game.LastQuestionNumber++

			err = gameUpdater.UpdateGame(c.Request().Context(), gameId, map[string]any{
				"last_question_number": game.LastQuestionNumber,
				"last_question_time":   time.Now(),
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
			}
		}

		// if it's the last question, end the game
		if game.LastQuestionNumber == 10 {
			err = gameUpdater.UpdateGame(c.Request().Context(), gameId, map[string]any{
				"end_at": time.Now(),
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
			}
		}

		// in any case we return current question
		question, err := questionProvider.QuestionByNumber(c.Request().Context(), gameId, game.LastQuestionNumber)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		answers, err := answersProvider.Answers(c.Request().Context(), question.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		outAnswers := make([]Answer, len(answers))

		for i, answer := range answers {
			outAnswers[i] = Answer{
				Id:     answer.Id,
				Answer: answer.Answer,
			}
		}

		c.JSON(http.StatusOK, echo.Map{
			"secondsLeft": max(0, intervalSeconds-int(time.Since(*game.LastQuestionTime).Seconds())),
			"question": Question{
				Id:       question.Id,
				Question: question.Question,
			},
			"answers": outAnswers,
		})

		return nil
	}
}
