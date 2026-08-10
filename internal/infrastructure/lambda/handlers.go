package lambda

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ivanbatistao/recommendations-service/configs"
	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/dto"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/app/composition"
)

type LambdaHandler struct {
	*composition.Application
}

func NewLambdaHandler() *LambdaHandler {
	config := configs.Load()

	app, err := composition.NewApplication(config)
	if err != nil {
		panic(err)
	}

	return &LambdaHandler{
		Application: app,
	}
}

func (h *LambdaHandler) HandleRequest(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	h.Logger.Info(
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
		return jsonError(404, "not found")
	}

	return jsonError(404, "not found")
}

func (h *LambdaHandler) handleHealth(
	ctx context.Context,
) (events.APIGatewayProxyResponse, error) {
	return jsonResponse(200, map[string]string{"status": "ok"})
}

func (h *LambdaHandler) handleGetRecommendations(
	ctx context.Context,
	userID string,
) (events.APIGatewayProxyResponse, error) {
	if err := validateUserID(userID); err != nil {
		return jsonError(400, err.Error())
	}

	result, err := h.GetRecommendationsHandler.Execute(
		ctx,
		queries.GetRecommendationsQuery{
			UserID: userID,
		},
	)

	if err != nil {
		h.Logger.Error(
			"failed to get recommendations",
			slog.String("error", err.Error()),
		)
		return jsonError(500, "internal server error")
	}

	dtos := dto.FromDomainSlice(result)
	return jsonResponse(200, map[string]interface{}{
		"recommendations": dtos,
	})
}

func (h *LambdaHandler) handleProcessEvent(
	ctx context.Context,
	body string,
) (events.APIGatewayProxyResponse, error) {
	if err := validateRequestBody(body); err != nil {
		return jsonError(400, err.Error())
	}

	var eventDTO dto.EventDTO
	if err := json.Unmarshal([]byte(body), &eventDTO); err != nil {
		return jsonError(400, "invalid request body")
	}

	event, err := dto.ToEventDomain(eventDTO)
	if err != nil {
		return jsonError(400, "invalid event data")
	}

	err = h.ProcessEventHandler.Execute(
		ctx,
		commands.ProcessEventCommand{
			Event: event,
		},
	)

	if err != nil {
		h.Logger.Error(
			"failed to process event",
			slog.String("error", err.Error()),
		)
		return jsonError(500, "internal server error")
	}

	return jsonResponse(202, map[string]string{"status": "event processed"})
}

func (h *LambdaHandler) handleGenerateRecommendations(
	ctx context.Context,
	body string,
) (events.APIGatewayProxyResponse, error) {
	if err := validateRequestBody(body); err != nil {
		return jsonError(400, err.Error())
	}

	var request struct {
		UserID string          `json:"user_id"`
		Events []dto.EventDTO  `json:"events"`
		Limit  int             `json:"limit"`
	}

	if err := json.Unmarshal([]byte(body), &request); err != nil {
		return jsonError(400, "invalid request body")
	}

	eventList := make([]event.Event, len(request.Events))
	for i, eventDTO := range request.Events {
		e, err := dto.ToEventDomain(eventDTO)
		if err != nil {
			return jsonError(400, "invalid event data")
		}
		eventList[i] = e
	}

	result, err := h.GenerateRecommendationsHandler.Execute(
		ctx,
		commands.GenerateRecommendationsCommand{
			UserID: request.UserID,
			Events: eventList,
			Limit:  request.Limit,
		},
	)

	if err != nil {
		h.Logger.Error(
			"failed to generate recommendations",
			slog.String("error", err.Error()),
		)
		return jsonError(500, "internal server error")
	}

	dtos := dto.FromDomainSlice(result)
	return jsonResponse(200, map[string]interface{}{
		"recommendations": dtos,
	})
}
