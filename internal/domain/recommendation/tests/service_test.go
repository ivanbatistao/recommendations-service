package recommendation_test

import (
	"context"
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestServiceGenerateRecommendations(t *testing.T) {
	events := []event.Event{
		{
			EventID:   "event-1",
			EventType: event.ProductViewed,
			UserID:    "user-1",
			ProductID: "P10",
		},
		{
			EventID:   "event-2",
			EventType: event.ProductAddedCart,
			UserID:    "user-1",
			ProductID: "P20",
		},
		{
			EventID:   "event-3",
			EventType: event.ProductPurchased,
			UserID:    "user-1",
			ProductID: "P20",
		},
		{
			EventID:   "event-4",
			EventType: event.ProductViewed,
			UserID:    "user-1",
			ProductID: "P30",
		},
	}

	service := recommendation.NewService(nil)

	result, err := service.GenerateRecommendations(
		"user-1",
		events,
		2,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 recommendations, got %d",
			len(result),
		)
	}

	if result[0].ProductID != "P20" {
		t.Fatalf(
			"expected first recommendation P20, got %s",
			result[0].ProductID,
		)
	}

	if result[0].Score != 8 {
		t.Fatalf(
			"expected P20 score 8, got %.2f",
			result[0].Score,
		)
	}

	if result[1].ProductID != "P10" {
		t.Fatalf(
			"expected second recommendation P10, got %s",
			result[1].ProductID,
		)
	}
}

func TestServiceGenerateRecommendationsErrors(t *testing.T) {
	service := recommendation.NewService(nil)

	// Test empty userID
	_, err := service.GenerateRecommendations("", []event.Event{}, 10)
	if err != recommendation.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}

	// Test invalid limit
	_, err = service.GenerateRecommendations("user-1", []event.Event{}, 0)
	if err != recommendation.ErrInvalidLimit {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}

	_, err = service.GenerateRecommendations("user-1", []event.Event{}, -5)
	if err != recommendation.ErrInvalidLimit {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
}

func TestServiceGetByUserIDErrors(t *testing.T) {
	service := recommendation.NewService(nil)

	// Test empty userID
	_, err := service.GetByUserID(nil, "")
	if err != recommendation.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

type fakeRepository struct {
	recommendations []recommendation.Recommendation
}

func (f *fakeRepository) GetByUserID(
	_ context.Context,
	userID string,
) ([]recommendation.Recommendation, error) {
	var result []recommendation.Recommendation

	for _, recommendation := range f.recommendations {
		if recommendation.UserID == userID {
			result = append(result, recommendation)
		}
	}

	return result, nil
}

func (f *fakeRepository) Save(
	_ context.Context,
	recommendation recommendation.Recommendation,
) error {
	for i, current := range f.recommendations {
		if current.UserID == recommendation.UserID &&
			current.ProductID == recommendation.ProductID {
			f.recommendations[i] = recommendation
			return nil
		}
	}

	f.recommendations = append(
		f.recommendations,
		recommendation,
	)

	return nil
}

func TestServiceProcessEvent(t *testing.T) {
	repository := &fakeRepository{}

	service := recommendation.NewService(repository)

	err := service.ProcessEvent(
		context.Background(),
		event.Event{
			EventID:   "event-1",
			EventType: event.ProductViewed,
			UserID:    "123",
			ProductID: "P10",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	recommendations, err := repository.GetByUserID(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(recommendations) != 1 {
		t.Fatalf(
			"expected 1 recommendation, got %d",
			len(recommendations),
		)
	}

	if recommendations[0].ProductID != "P10" {
		t.Fatalf(
			"expected P10, got %s",
			recommendations[0].ProductID,
		)
	}

	if recommendations[0].Score != 1.0 {
		t.Fatalf(
			"expected score 1.0, got %.2f",
			recommendations[0].Score,
		)
	}
}
