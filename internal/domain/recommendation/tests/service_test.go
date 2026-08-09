package recommendation_test

import (
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

	result := service.GenerateRecommendations(
		"user-1",
		events,
		2,
	)

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
