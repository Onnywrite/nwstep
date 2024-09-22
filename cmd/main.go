package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Onnywrite/nwstep/internal/app"
	"github.com/Onnywrite/nwstep/internal/config"
)

func main() {
	conf := config.MustLoad("/etc/nwstep/conf.yaml")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app.New(conf).MustRun(ctx)
}
