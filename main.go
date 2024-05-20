package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kuberise/todolist/internal/controller"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &controller.HTTPConfig{
		Port:            80,
		ShutdownTimeout: 10,
	}

	hc := controller.NewHTTPController(logger, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := hc.Run(ctx)
	if err != nil {
		logger.Error("http server error", "error", err)
	}
}
