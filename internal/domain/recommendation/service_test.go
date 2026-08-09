package recommendation_test

import (
	"context"
	"testing"

	"github.com/ivanbatistao/recommendation-service/internal/domain/recommendation"
)

type fakeRepository struct {
	recommendations []recommendation.Recommendation
}

func (f *fakeRepository) GetByUserID(
	_ context.Context,
	_ string,
) ([]recommendation.Recommendation, error) {
	return f.recommendations, nil
}

func (f *fakeRepository) Save(
	_ context.Context,
	_ recommendation.Recommendation,
) error {
	return nil
}

func TestServiceGetByUserID(t *testing.T) {
	expected := []recommendation.Recommendation{
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
	}

	repository := &fakeRepository{
		recommendations: expected,
	}

	service := recommendation.NewService(repository)

	result, err := service.GetByUserID(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf(
			"expected %d recommendations, got %d",
			len(expected),
			len(result),
		)
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Fatalf(
				"expected recommendation %+v, got %+v",
				expected[i],
				result[i],
			)
		}
	}
}
