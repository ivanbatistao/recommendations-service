package composition

import (
	"context"
	"log/slog"

	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ivanbatistao/recommendations-service/configs"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/logger"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/dynamodb"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/memory"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type Application struct {
	Service                         *recommendation.Service
	GetRecommendationsHandler       *queries.GetRecommendationsHandler
	ProcessEventHandler            *commands.ProcessEventHandler
	GenerateRecommendationsHandler *commands.GenerateRecommendationsHandler
	Logger                          *slog.Logger
}

func NewApplication(config configs.Config) (*Application, error) {
	log := logger.New()

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
			return nil, err
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

	return &Application{
		Service:                         service,
		GetRecommendationsHandler:       getRecommendationsHandler,
		ProcessEventHandler:            processEventHandler,
		GenerateRecommendationsHandler: generateRecommendationsHandler,
		Logger:                          log,
	}, nil
}

func NewApplicationWithLogger(config configs.Config, logger *slog.Logger) (*Application, error) {
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
			logger.Error(
				"failed to initialize DynamoDB client",
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		repository = dynamodb.NewDynamoDBRepository(client, config.DynamoDBTable)
		logger.Info("using DynamoDB repository")
	} else {
		// Use memory repository (default for simple testing)
		repository = memory.NewMemoryRepository()
		logger.Info("using memory repository")
	}

	service := recommendation.NewService(repository)

	getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
	processEventHandler := commands.NewProcessEventHandler(service)
	generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

	return &Application{
		Service:                         service,
		GetRecommendationsHandler:       getRecommendationsHandler,
		ProcessEventHandler:            processEventHandler,
		GenerateRecommendationsHandler: generateRecommendationsHandler,
		Logger:                          logger,
	}, nil
}

func (app *Application) Shutdown() {
	// Graceful shutdown logic if needed
	app.Logger.Info("application shutdown")
}
