package queries

import (
	"context"

	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type GetRecommendationsQuery struct {
	UserID string
}

type GetRecommendationsHandler struct {
	service *recommendation.Service
}

func NewGetRecommendationsHandler(
	service *recommendation.Service,
) *GetRecommendationsHandler {
	return &GetRecommendationsHandler{
		service: service,
	}
}

func (h *GetRecommendationsHandler) Execute(
	ctx context.Context,
	query GetRecommendationsQuery,
) ([]recommendation.Recommendation, error) {
	return h.service.GetByUserID(ctx, query.UserID)
}
