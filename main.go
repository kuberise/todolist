package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuberise/todolist/internal/config"
	"github.com/kuberise/todolist/internal/controller"
	"github.com/kuberise/todolist/internal/gateway/relational"
	"github.com/kuberise/todolist/pkg/postgres"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.New("/configs.yaml")
	if err != nil {
		logger.Error("error parsing config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		logger.Error("error connecting to postgres", "error", err)
		os.Exit(1)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (title VARCHAR(255))")
	if err != nil {
		logger.Error("error creating table", "error", err)
		os.Exit(1)
	}

	repository := relational.NewRepository(db)

	hc := controller.NewHTTPController(logger, &cfg.HTTP, repository)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = hc.Run(ctx)
	if err != nil {
		logger.Error("http server error", "error", err)
	}
}
