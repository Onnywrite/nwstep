package server

import (
	"context"
	"fmt"
	"net/http"

	handlerauth "github.com/Onnywrite/nwstep/internal/server/handler/auth"
	handlercateg "github.com/Onnywrite/nwstep/internal/server/handler/categories"
	handlergames "github.com/Onnywrite/nwstep/internal/server/handler/games"
	mw "github.com/Onnywrite/nwstep/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e          *echo.Echo
	address    string
	secret     string
	users      UserRepo
	categories CategoryRepo
	games      GameRepo
	questions  QuestionRepo
}

type UserRepo interface {
	handlerauth.UserSaver
	handlerauth.UserProvider
	handlerauth.UserByIdProvider
}

type CategoryRepo interface {
	handlercateg.CategoriesProvider
	handlercateg.CategoryProvider
	handlercateg.CoursesProvider
	handlercateg.CourseProvider
	handlercateg.RatingProvider
	handlercateg.CategoryTopProvider
	handlercateg.UserTopProvider
	// idc
	handlercateg.CategorySaver
	handlercateg.CourseSaver
}

type GameRepo interface {
	handlercateg.LobbyGameProvider
	handlercateg.GameSaver
	handlercateg.GameUserLinker
	handlercateg.UsersInGameProvider
	handlercateg.IsUserInLobbyProvider
}

type QuestionRepo interface {
	handlercateg.RandomQuestionsPicker
	handlergames.UsersInGameProvider
	handlergames.GameProvider
	handlergames.GameUpdater
	handlergames.GameQuestionProvider
	handlergames.AnswersProvider
	handlercateg.QuestionsSaver
	handlercateg.AnswersSaver
}

func New(port uint32, secret string, users UserRepo, categories CategoryRepo, games GameRepo, questions QuestionRepo) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS(), middleware.Recover(), middleware.Logger())

	server := &Server{
		e:          e,
		address:    fmt.Sprintf(":%d", port),
		secret:     secret,
		users:      users,
		categories: categories,
		games:      games,
		questions:  questions,
	}

	server.initApi()

	return server
}

func (s *Server) initApi() {
	api := s.e.Group("/api")

	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	{
		auth := api.Group("/auth")

		auth.POST("/register", handlerauth.PostRegister(s.users, s.secret))
		auth.POST("/sign-in", handlerauth.PostSignIn(s.users, s.secret))
		auth.GET("/profile", handlerauth.GetProfile(s.users), mw.Auth(s.secret))
	}

	{
		categories := api.Group(("/categories"), mw.Auth(s.secret))

		categories.POST("", handlercateg.PostCategory(s.categories), mw.Auth(s.secret, true))
		categories.GET("", handlercateg.GetCategories(s.categories))

		categories.GET("/:category_id",
			handlercateg.GetCategory(s.categories),
			mw.IntParams("category_id"))
		categories.GET("/:category_id/courses",
			handlercateg.GetCourses(s.categories, s.categories),
			mw.IntParams("category_id"))
		categories.POST("/:category_id/courses", handlercateg.PostCourse(s.categories),
			mw.Auth(s.secret, true),
			mw.IntParams("category_id"))
		categories.PUT("/:category_id/courses/:course_id/join",
			handlercateg.PutJoin(5, s.categories, s.categories, s.games, s.games,
				s.games, s.games, s.games, s.questions),
			mw.IntParams("category_id", "course_id"))
		categories.GET("/:category_id/top", handlercateg.GetTop(s.categories, s.categories),
			mw.IntParams("category_id"))
	}

	api.POST("/courses/:course_id/questions", handlercateg.PostQuestion(s.questions, s.questions),
		mw.IntParams("course_id"), mw.Auth(s.secret, true))

	{
		games := api.Group("/games", mw.Auth(s.secret))

		games.GET("/:game_id/currentQuestion", handlergames.GetCurrentQuestion(5, 15,
			s.games, s.questions, s.questions, s.questions, s.questions),
			mw.IntParams("game_id"))
	}
}

func (s *Server) Start() error {
	err := s.e.Start(s.address)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.e.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
