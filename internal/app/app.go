package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Onnywrite/nwstep/internal/config"
	"github.com/Onnywrite/nwstep/internal/lib/logger/slogpretty"
	"github.com/Onnywrite/nwstep/internal/server"
	"github.com/Onnywrite/nwstep/internal/storage"
)

type Application struct {
	cfg *config.Config
}

func New(conf *config.Config) *Application {
	slog.SetDefault(getPrettySlog())

	return &Application{
		cfg: conf,
	}
}

func (a *Application) MustRun(ctx context.Context) {
	slog.Debug("starting application")

	database, err := storage.Connect(a.cfg.Conn)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := database.Disconnect(); err != nil {
			slog.Error("failed to disconnect from database", "error", err)
		}

		slog.Info("disconnected from database")
	}()

	server := server.New(a.cfg.Port, database)

	err = server.Start()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := server.Stop(ctx); err != nil {
			slog.Error("failed to stop server", "error", err)
		}

		slog.Info("server stopped")
	}()

	slog.Info("application started")

	<-ctx.Done()
	slog.Debug("stopping application")
}

func getPrettySlog() *slog.Logger {
	//nolint: exhaustruct
	opts := slogpretty.PrettyHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
