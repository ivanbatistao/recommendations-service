package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ivanbatistao/recommendations-service/configs"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	httpgin "github.com/ivanbatistao/recommendations-service/internal/infrastructure/http/gin"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/logger"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/dynamodb"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/memory"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func main() {
	log := logger.New()

	config := configs.Load()

	// Choose repository based on environment
	var repository recommendation.Repository

	if config.UseDynamoDB {
		// Use DynamoDB (production or local development with MiniStack)
		var client *awsdynamodb.Client
		var err error

		if config.DynamoDBEndpoint != "" {
			// Local DynamoDB (MiniStack or DynamoDB Local)
			client, err = dynamodb.NewLocalDynamoDBClient(
				context.Background(),
				config.DynamoDBEndpoint,
			)
		} else {
			// AWS DynamoDB
			client, err = dynamodb.NewDynamoDBClient(
				context.Background(),
				config.AWSRegion,
			)
		}

		if err != nil {
			log.Error(
				"failed to initialize DynamoDB client",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}

		repository = dynamodb.NewDynamoDBRepository(client, config.DynamoDBTable)
		log.Info("using DynamoDB repository")
	} else {
		// Use memory repository (default for simple testing)
		repository = memory.NewMemoryRepository()
		log.Info("using memory repository")
	}

	service := recommendation.NewService(repository)

	getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
	processEventHandler := commands.NewProcessEventHandler(service)
	generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

	handler := httpgin.NewHandler(
		getRecommendationsHandler,
		processEventHandler,
		generateRecommendationsHandler,
	)

	log.Info("handlers created")

	router := httpgin.NewRouter(handler)

	log.Info("router created")

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
