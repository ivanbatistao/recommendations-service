package workerpool

import (
	"context"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type RecommendationProcessorAdapter struct {
	service *recommendation.Service
}

func NewRecommendationProcessorAdapter(
	service *recommendation.Service,
) *RecommendationProcessorAdapter {
	return &RecommendationProcessorAdapter{
		service: service,
	}
}

func (a *RecommendationProcessorAdapter) ProcessEvent(
	ctx context.Context,
	ev event.Event,
) error {
	return a.service.ProcessEvent(ctx, ev)
}
