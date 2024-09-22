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

type CoursesProvider interface {
	Courses(context.Context, int) ([]models.Course, error)
}

type RatingProvider interface {
	Rating(context.Context, uuid.UUID, int) (int, error)
}

type QuestionsSaver interface {
	SaveQuestion(context.Context, models.Question) (*models.Question, error)
}

func PostQuestions(saver QuestionsSaver) echo.HandlerFunc {
	type Question struct {
		Question string   `json:"question"`
		Answers  []string `json:"answers"`
		Correct  int      `json:"correct"`
	}

	return func(c echo.Context) error {
		return nil
	}
}

type CourseSaver interface {
	SaveCourse(context.Context, models.Course) (*models.Course, error)
}

func PostCourse(saver CourseSaver) echo.HandlerFunc {
	type Course struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		MinRating     int    `json:"minRating"`
		OptimalRating int    `json:"optimalRating"`
		CategoryId    int    `json:"categoryId"`
		PhotoUrl      string `json:"photoUrl"`
	}

	return func(c echo.Context) error {
		var course Course
		if err := c.Bind(&c); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request").SetInternal(err)
		}

		saved, err := saver.SaveCourse(c.Request().Context(), models.Course{
			Name:          course.Name,
			Description:   course.Description,
			MinRating:     course.MinRating,
			OptimalRating: course.OptimalRating,
			CategoryId:    course.CategoryId,
			PhotoUrl:      course.PhotoUrl,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error").SetInternal(err)
		}

		c.JSON(http.StatusOK, saved)

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
