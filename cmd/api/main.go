package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivanbatistao/recommendations-service/configs"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	httpgin "github.com/ivanbatistao/recommendations-service/internal/infrastructure/http/gin"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/logger"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/memory"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func main() {
	log := logger.New()

	config := configs.Load()

	repository := memory.NewMemoryRepository()
	service := recommendation.NewService(repository)

	getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
	processEventHandler := commands.NewProcessEventHandler(service)
	generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

	handler := httpgin.NewHandler(
		getRecommendationsHandler,
		processEventHandler,
		generateRecommendationsHandler,
	)

	router := httpgin.NewRouter(handler)

	server := httpgin.NewServer(config.Port, router)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Start()
	}()

	log.Info(
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
		log.Error(
			"server error",
			slog.String("error", err.Error()),
		)

	case signal := <-shutdownSignal:
		log.Info(
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
		log.Error(
			"server shutdown error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	log.Info("server stopped")
}
