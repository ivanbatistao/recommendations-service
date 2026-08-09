package queries

import (
	"context"
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type fakeRecommendationRepository struct {
	recommendations []recommendation.Recommendation
}

func (f *fakeRecommendationRepository) GetByUserID(
	_ context.Context,
	_ string,
) ([]recommendation.Recommendation, error) {
	return f.recommendations, nil
}

func (f *fakeRecommendationRepository) Save(
	_ context.Context,
	_ recommendation.Recommendation,
) error {
	return nil
}

func TestGetRecommendations(t *testing.T) {
	repository := &fakeRecommendationRepository{
		recommendations: []recommendation.Recommendation{
			{
				UserID:    "123",
				ProductID: "P20",
				Score:     0.95,
			},
			{
				UserID:    "123",
				ProductID: "P40",
				Score:     0.82,
			},
		},
	}

	service := recommendation.NewService(repository)

	handler := NewGetRecommendationsHandler(service)

	result, err := handler.Execute(
		context.Background(),
		GetRecommendationsQuery{
			UserID: "123",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(result))
	}

	if result[0].ProductID != "P20" {
		t.Fatalf(
			"expected first recommendation P20, got %s",
			result[0].ProductID,
		)
	}

	if result[1].ProductID != "P40" {
		t.Fatalf(
			"expected second recommendation P40, got %s",
			result[1].ProductID,
		)
	}
}
