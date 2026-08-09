package commands

import (
	"context"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type GenerateRecommendationsCommand struct {
	UserID string
	Events []event.Event
	Limit  int
}

type GenerateRecommendationsHandler struct {
	service *recommendation.Service
}

func NewGenerateRecommendationsHandler(
	service *recommendation.Service,
) *GenerateRecommendationsHandler {
	return &GenerateRecommendationsHandler{
		service: service,
	}
}

func (h *GenerateRecommendationsHandler) Execute(
	ctx context.Context,
	command GenerateRecommendationsCommand,
) ([]recommendation.Recommendation, error) {
	return h.service.GenerateRecommendations(
		command.UserID,
		command.Events,
		command.Limit,
	)
}
