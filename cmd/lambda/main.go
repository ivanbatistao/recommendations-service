package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ivanbatistao/recommendations-service/configs"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/dto"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/logger"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/dynamodb"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/memory"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type LambdaHandler struct {
	getRecommendationsHandler    *queries.GetRecommendationsHandler
	processEventHandler         *commands.ProcessEventHandler
	generateRecommendationsHandler *commands.GenerateRecommendationsHandler
	logger                       *slog.Logger
}

func NewLambdaHandler() *LambdaHandler {
	log := logger.New()
	config := configs.Load()

	// Choose repository based on environment
	var repository recommendation.Repository

	if config.UseDynamoDB {
		// Use DynamoDB (AWS Lambda environment)
		var client interface{}
		var err error

		if config.DynamoDBEndpoint != "" {
			// Local DynamoDB (for testing)
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

		dynamoDBClient := client.(*awsdynamodb.Client)
		repository = dynamodb.NewDynamoDBRepository(dynamoDBClient, config.DynamoDBTable)
		log.Info("using DynamoDB repository")
	} else {
		// Use memory repository (for testing)
		repository = memory.NewMemoryRepository()
		log.Info("using memory repository")
	}

	service := recommendation.NewService(repository)

	getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
	processEventHandler := commands.NewProcessEventHandler(service)
	generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

	return &LambdaHandler{
		getRecommendationsHandler:    getRecommendationsHandler,
		processEventHandler:         processEventHandler,
		generateRecommendationsHandler: generateRecommendationsHandler,
		logger:                       log,
	}
}

func (h *LambdaHandler) HandleRequest(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	h.logger.Info(
		"lambda request",
		slog.String("path", req.Path),
		slog.String("method", req.HTTPMethod),
	)

	switch req.Path {
	case "/health":
		return h.handleHealth(ctx)

	case "/recommendations/" + req.PathParameters["userId"]:
		return h.handleGetRecommendations(ctx, req.PathParameters["userId"])

	case "/events":
		if req.HTTPMethod == "POST" {
			return h.handleProcessEvent(ctx, req.Body)
		}

	case "/recommendations/generate":
		if req.HTTPMethod == "POST" {
			return h.handleGenerateRecommendations(ctx, req.Body)
		}

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"not found"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"error":"not found"}`,
	}, nil
}

func (h *LambdaHandler) handleHealth(
	ctx context.Context,
) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"status":"ok"}`,
	}, nil
}

func (h *LambdaHandler) handleGetRecommendations(
	ctx context.Context,
	userID string,
) (events.APIGatewayProxyResponse, error) {
	if userID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"user_id is required"}`,
		}, nil
	}

	result, err := h.getRecommendationsHandler.Execute(
		ctx,
		queries.GetRecommendationsQuery{
			UserID: userID,
		},
	)

	if err != nil {
		h.logger.Error(
			"failed to get recommendations",
			slog.String("error", err.Error()),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"internal server error"}`,
		}, nil
	}

	dtos := dto.FromDomainSlice(result)
	body, _ := json.Marshal(map[string]interface{}{
		"recommendations": dtos,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func (h *LambdaHandler) handleProcessEvent(
	ctx context.Context,
	body string,
) (events.APIGatewayProxyResponse, error) {
	var eventDTO dto.EventDTO
	if err := json.Unmarshal([]byte(body), &eventDTO); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"invalid request body"}`,
		}, nil
	}

	event, err := dto.ToEventDomain(eventDTO)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"invalid event data"}`,
		}, nil
	}

	err = h.processEventHandler.Execute(
		ctx,
		commands.ProcessEventCommand{
			Event: event,
		},
	)

	if err != nil {
		h.logger.Error(
			"failed to process event",
			slog.String("error", err.Error()),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 202,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"status":"event processed"}`,
	}, nil
}

func (h *LambdaHandler) handleGenerateRecommendations(
	ctx context.Context,
	body string,
) (events.APIGatewayProxyResponse, error) {
	var request struct {
		UserID string          `json:"user_id"`
		Events []dto.EventDTO  `json:"events"`
		Limit  int             `json:"limit"`
	}

	if err := json.Unmarshal([]byte(body), &request); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"invalid request body"}`,
		}, nil
	}

	eventList := make([]event.Event, len(request.Events))
	for i, eventDTO := range request.Events {
		e, err := dto.ToEventDomain(eventDTO)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"error":"invalid event data"}`,
			}, nil
		}
		eventList[i] = e
	}

	result, err := h.generateRecommendationsHandler.Execute(
		ctx,
		commands.GenerateRecommendationsCommand{
			UserID: request.UserID,
			Events: eventList,
			Limit:  request.Limit,
		},
	)

	if err != nil {
		h.logger.Error(
			"failed to generate recommendations",
			slog.String("error", err.Error()),
		)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"internal server error"}`,
		}, nil
	}

	dtos := dto.FromDomainSlice(result)
	responseBody, _ := json.Marshal(map[string]interface{}{
		"recommendations": dtos,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
}
