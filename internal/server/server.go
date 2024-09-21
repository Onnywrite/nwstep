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
	users      UserRepo
	categories CategoriesRepo
}

type UserRepo interface {
	handlerauth.UserSaver
	handlerauth.UserProvider
	handlerauth.UserByIdProvider
}

type CategoriesRepo interface {
	handlercateg.CategoriesProvider
	handlercateg.CoursesProvider
	handlercateg.RatingProvider
}

func New(port uint32, users UserRepo, categories CategoriesRepo) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS(), middleware.Recover(), middleware.Logger())

	server := &Server{
		e:          e,
		address:    fmt.Sprintf(":%d", port),
		users:      users,
		categories: categories,
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

		auth.POST("/register", handlerauth.PostRegister(s.users, "secret"))
		auth.POST("/sign-in", handlerauth.PostSignIn(s.users, "secret"))
		auth.GET("/profile", handlerauth.GetProfile(s.users), mw.Auth("secret"))
	}

	{
		categories := api.Group(("/categories"), mw.Auth("secret"))

		categories.GET("", handlercateg.GetCategories(s.categories))
		categories.GET("/:category_id/courses", handlercateg.GetCourses(s.categories, s.categories))
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
