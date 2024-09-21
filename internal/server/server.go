package server

import (
	"context"
	"fmt"
	"net/http"

	handlerauth "github.com/Onnywrite/nwstep/internal/server/handler/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e       *echo.Echo
	address string
	users   UserRepo
}

type UserRepo interface {
	handlerauth.UserSaver
	handlerauth.UserProvider
}

func New(port uint32, users UserRepo) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS(), middleware.Recover(), middleware.Logger())

	server := &Server{
		e:       e,
		address: fmt.Sprintf(":%d", port),
		users:   users,
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
