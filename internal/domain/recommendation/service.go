package recommendation

import (
	"context"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByUserID(
	ctx context.Context,
	userID string,
) ([]Recommendation, error) {
	return s.repository.GetByUserID(ctx, userID)
}

func (s *Service) GenerateRecommendations(
	userID string,
	events []event.Event,
	limit int,
) []Recommendation {
	interest := CalculateInterest(userID, events)

	return Rank(userID, interest, limit)
}
