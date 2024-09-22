package handlercateg

import (
	"context"
	"errors"
	"net/http"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/domain/single"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CoursesProvider interface {
	Courses(context.Context, int) ([]models.Course, error)
}

type RatingProvider interface {
	Rating(context.Context, uuid.UUID, int) (int, error)
}

type QuestionsCoursesProvider interface {
	Questions(context.Context, int) ([]models.Question, error)
}

type AnswersProvider interface {
	Answers(context.Context, int) ([]models.Answer, error)
}

type CourseDeleter interface {
	DeleteCourse(context.Context, int) error
}

func DeleteCourse(deleter CourseDeleter) echo.HandlerFunc {
	return func(c echo.Context) error {
		courseId := c.Get("course_id").(int)

		err := deleter.DeleteCourse(c.Request().Context(), courseId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.NoContent(http.StatusNoContent)

		return nil
	}
}

func GetQuestions(provider QuestionsCoursesProvider, answProvider AnswersProvider) echo.HandlerFunc {
	type Answer struct {
		Id        int    `json:"id"`
		Answer    string `json:"answer"`
		IsCorrect bool   `json:"isCorrect"`
	}

	type Question struct {
		Id       int      `json:"id"`
		Question string   `json:"question"`
		CourseId int      `json:"courseId"`
		Answers  []Answer `json:"answers"`
	}

	return func(c echo.Context) error {
		courseId := c.Get("course_id").(int)

		questions, err := provider.Questions(c.Request().Context(), courseId)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "course not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		outQuestions := make([]Question, len(questions))

		for i, question := range questions {
			answers, err := answProvider.Answers(c.Request().Context(), question.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
			}

			outAnswers := make([]Answer, len(answers))

			for j, answer := range answers {
				outAnswers[j] = Answer{
					Id:        answer.Id,
					Answer:    answer.Answer,
					IsCorrect: answer.IsCorrect,
				}
			}

			outQuestions[i] = Question{
				Id:       question.Id,
				Question: question.Question,
				CourseId: question.CourseId,
				Answers:  outAnswers,
			}
		}

		c.JSON(http.StatusOK, outQuestions)

		return nil
	}
}

type QuestionsSaver interface {
	SaveQuestion(context.Context, models.Question) (*models.Question, error)
}

type AnswersSaver interface {
	SaveAnswers(context.Context, ...models.Answer) ([]models.Answer, error)
}

func PostQuestion(saver QuestionsSaver, answSaver AnswersSaver) echo.HandlerFunc {
	type Answer struct {
		Answer    string `json:"answer" validate:"required,max=200"`
		IsCorrect bool   `json:"isCorrect"`
	}

	type OutAnswer struct {
		Id        int    `json:"id"`
		Answer    string `json:"answer"`
		IsCorrect bool   `json:"isCorrect"`
	}

	type Question struct {
		Question string   `json:"question" validate:"required,max=200"`
		Answers  []Answer `json:"answers" validate:"required,min=2,dive"`
	}

	return func(c echo.Context) error {
		var q Question
		if err := c.Bind(&q); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request").SetInternal(err)
		}

		if err := single.V.Struct(q); err != nil {
			return single.ValidationError(err)
		}

		courseId := c.Get("course_id").(int)

		savedQ, err := saver.SaveQuestion(c.Request().Context(), models.Question{
			Question: q.Question,
			CourseId: courseId,
		})
		switch {
		case errors.Is(err, cuteql.ErrUnique):
			return echo.NewHTTPError(http.StatusConflict, "question already exist")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		modelsAnswers := make([]models.Answer, len(q.Answers))
		for i, answer := range q.Answers {
			modelsAnswers[i] = models.Answer{
				QuestionId: savedQ.Id,
				Answer:     answer.Answer,
				IsCorrect:  answer.IsCorrect,
			}
		}

		savedAnswers, err := answSaver.SaveAnswers(c.Request().Context(), modelsAnswers...)
		switch {
		case errors.Is(err, cuteql.ErrUnique):
			return echo.NewHTTPError(http.StatusConflict, "answer already exist")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		outAnswers := make([]OutAnswer, len(savedAnswers))
		for i, answer := range savedAnswers {
			outAnswers[i] = OutAnswer{
				Id:        answer.Id,
				Answer:    answer.Answer,
				IsCorrect: answer.IsCorrect,
			}
		}

		c.JSON(http.StatusOK, echo.Map{
			"id":       savedQ.Id,
			"question": savedQ.Question,
			"courseId": savedQ.CourseId,
			"answers":  outAnswers,
		})

		return nil
	}
}

type CourseSaver interface {
	SaveCourse(context.Context, models.Course) (*models.Course, error)
}

func PostCourse(saver CourseSaver) echo.HandlerFunc {
	type Course struct {
		Name          string `json:"name" validate:"required,max=100"`
		Description   string `json:"description" validate:"required,max=200"`
		MinRating     int    `json:"minRating" validate:"gte=0"`
		OptimalRating int    `json:"optimalRating" validate:"gte=0"`
		PhotoUrl      string `json:"photoUrl" validate:"omitempty,max=200,url"`
	}

	return func(c echo.Context) error {
		var course Course
		if err := c.Bind(&course); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request").SetInternal(err)
		}

		if err := single.V.Struct(course); err != nil {
			return single.ValidationError(err)
		}

		categoryId := c.Get("category_id").(int)

		saved, err := saver.SaveCourse(c.Request().Context(), models.Course{
			Name:          course.Name,
			Description:   course.Description,
			MinRating:     course.MinRating,
			OptimalRating: course.OptimalRating,
			CategoryId:    categoryId,
			PhotoUrl:      course.PhotoUrl,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, echo.Map{
			"id":            saved.Id,
			"name":          saved.Name,
			"description":   saved.Description,
			"minRating":     saved.MinRating,
			"optimalRating": saved.OptimalRating,
			"categoryId":    saved.CategoryId,
			"photoUrl":      saved.PhotoUrl,
		})

		return nil
	}
}

func GetCourses(provider CoursesProvider, ratingProvider RatingProvider) echo.HandlerFunc {
	type Course struct {
		Id            int    `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		MinRating     int    `json:"minRating"`
		OptimalRating int    `json:"optimalRating"`
		IsAwailable   bool   `json:"isAwailable"`
		PhotoUrl      string `json:"photoUrl"`
	}

	type Courses struct {
		Courses []Course `json:"courses"`
	}

	return func(c echo.Context) error {
		id := c.Get("id").(uuid.UUID)
		categoryId := c.Get("category_id").(int)

		courses, err := provider.Courses(c.Request().Context(), categoryId)
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			return echo.NewHTTPError(http.StatusNotFound, "category not found")
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		rating, err := ratingProvider.Rating(c.Request().Context(), id, int(categoryId))
		switch {
		case errors.Is(err, cuteql.ErrEmptyResult):
			rating = 0
		case err != nil:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		outCourses := make([]Course, len(courses))

		for i, course := range courses {
			outCourses[i] = Course{
				Id:            course.Id,
				Name:          course.Name,
				Description:   course.Description,
				MinRating:     course.MinRating,
				OptimalRating: course.OptimalRating,
				IsAwailable:   course.MinRating <= rating,
				PhotoUrl:      course.PhotoUrl,
			}
		}

		c.JSON(http.StatusOK, Courses{outCourses})

		return nil
	}
}
