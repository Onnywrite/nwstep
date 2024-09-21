package server

import (
	"context"
	"fmt"
	"net/http"

	handlerauth "github.com/Onnywrite/nwstep/internal/server/handler/auth"
	handlercateg "github.com/Onnywrite/nwstep/internal/server/handler/categories"
	mw "github.com/Onnywrite/nwstep/internal/server/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e          *echo.Echo
	address    string
	secret     string
	users      UserRepo
	categories CategoriesRepo
	games      GameRepo
}

type UserRepo interface {
	handlerauth.UserSaver
	handlerauth.UserProvider
	handlerauth.UserByIdProvider
}

type CategoriesRepo interface {
	handlercateg.CategoriesProvider
	handlercateg.CoursesProvider
	handlercateg.CourseProvider
	handlercateg.RatingProvider
}

type GameRepo interface {
	handlercateg.LobbyGameProvider
	handlercateg.GameSaver
	handlercateg.GameUserLinker
	handlercateg.UsersInLobbyProvider
	handlercateg.IsUserInLobbyProvider
}

func New(port uint32, secret string, users UserRepo, categories CategoriesRepo, games GameRepo) *Server {
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

		categories.GET("", handlercateg.GetCategories(s.categories))
		categories.GET("/:category_id/courses", handlercateg.GetCourses(s.categories, s.categories))
		categories.PUT("/:category_id/courses/:course_id/join",
			handlercateg.PutJoin(s.categories, s.categories, s.games, s.games, s.games, s.games, s.games))
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
