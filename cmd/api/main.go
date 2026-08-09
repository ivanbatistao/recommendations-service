package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivanbatistao/recommendations-service/configs"
	httpgin "github.com/ivanbatistao/recommendations-service/internal/infrastructure/http/gin"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/composition"
)

func main() {
	config := configs.Load()

	app, err := composition.NewApplication(config)
	if err != nil {
		os.Exit(1)
	}

	handler := httpgin.NewHandler(
		app.GetRecommendationsHandler,
		app.ProcessEventHandler,
		app.GenerateRecommendationsHandler,
	)

	router := httpgin.NewRouter(handler)

	server := httpgin.NewServer(config.Port, router)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Start()
	}()

	app.Logger.Info(
		"server running",
		slog.String("port", config.Port),
	)

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		app.Logger.Error(
			"server error",
			slog.String("error", err.Error()),
		)

	case signal := <-shutdownSignal:
		app.Logger.Info(
			"shutdown signal received",
			slog.String("signal", signal.String()),
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Error(
			"server shutdown error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	app.Shutdown()
	app.Logger.Info("server stopped")
}
