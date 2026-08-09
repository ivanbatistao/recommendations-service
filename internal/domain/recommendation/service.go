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
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	return s.repository.GetByUserID(ctx, userID)
}

func (s *Service) GenerateRecommendations(
	userID string,
	events []event.Event,
	limit int,
) ([]Recommendation, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	interest := CalculateInterest(userID, events)

	return Rank(userID, interest, limit), nil
}

func (s *Service) ProcessEvent(
	ctx context.Context,
	currentEvent event.Event,
) error {
	score := ScoreEvent(currentEvent.EventType)

	if score == 0 {
		return nil
	}

	recommendations, err := s.repository.GetByUserID(
		ctx,
		currentEvent.UserID,
	)
	if err != nil {
		return err
	}

	for _, recommendation := range recommendations {
		if recommendation.ProductID != currentEvent.ProductID {
			continue
		}

		recommendation.Score += score

		return s.repository.Save(ctx, recommendation)
	}

	return s.repository.Save(
		ctx,
		Recommendation{
			UserID:    currentEvent.UserID,
			ProductID: currentEvent.ProductID,
			Score:     score,
		},
	)
}