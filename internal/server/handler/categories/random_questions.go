package handlercateg

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/labstack/echo/v4"
)

type RandomQuestionsProvider interface {
	GetRandomQuestions(ctx context.Context, courseId, count int) ([]models.Question, error)
}

func GetQuestionsRand(provider RandomQuestionsProvider, answProv AnswersProvider) echo.HandlerFunc {
	type Question struct {
		Id       int    `json:"id"`
		Question string `json:"question"`
		CourseId int    `json:"courseId"`
	}

	type Answer struct {
		Id        int    `json:"id"`
		Answer    string `json:"answer"`
		IsCorrect bool   `json:"isCorrect"`
	}

	type QuestionAnswers struct {
		Question Question `json:"question"`
		Answers  []Answer `json:"answers"`
	}

	return func(c echo.Context) error {
		courseId := c.Get("course_id").(int)

		countStr := c.QueryParam("count")
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil || count <= 0 {
			count = 10
		}

		randomQuestions, err := provider.GetRandomQuestions(c.Request().Context(), courseId, int(count))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		questions := make([]QuestionAnswers, len(randomQuestions))

		for i, question := range randomQuestions {
			questions[i] = QuestionAnswers{
				Question: Question{
					Id:       question.Id,
					Question: question.Question,
					CourseId: question.CourseId,
				},
			}

			answers, err := answProv.Answers(c.Request().Context(), question.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
			}

			questions[i].Answers = make([]Answer, len(answers))
			for j, answer := range answers {
				questions[i].Answers[j] = Answer{
					Id:        answer.Id,
					Answer:    answer.Answer,
					IsCorrect: answer.IsCorrect,
				}
			}
		}

		c.JSON(http.StatusOK, echo.Map{
			"questions": questions,
		})

		return nil
	}
}
