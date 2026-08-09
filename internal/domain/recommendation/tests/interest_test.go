package recommendation_test

import (
	"testing"
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestCalculateInterest(t *testing.T) {
	events := []event.Event{
		{
			EventID:    "event-1",
			EventType:  event.ProductViewed,
			UserID:     "user-1",
			ProductID:  "P10",
			OccurredAt: time.Now(),
		},
		{
			EventID:    "event-2",
			EventType:  event.ProductViewed,
			UserID:     "user-1",
			ProductID:  "P10",
			OccurredAt: time.Now(),
		},
		{
			EventID:    "event-3",
			EventType:  event.ProductAddedCart,
			UserID:     "user-1",
			ProductID:  "P10",
			OccurredAt: time.Now(),
		},
		{
			EventID:    "event-4",
			EventType:  event.ProductPurchased,
			UserID:     "user-1",
			ProductID:  "P20",
			OccurredAt: time.Now(),
		},
	}

	result := recommendation.CalculateInterest(
		"user-1",
		events,
	)

	if result["P10"] != 5 {
		t.Fatalf(
			"expected P10 score 5, got %.2f",
			result["P10"],
		)
	}

	if result["P20"] != 5 {
		t.Fatalf(
			"expected P20 score 5, got %.2f",
			result["P20"],
		)
	}
}

func TestCalculateInterestIgnoresOtherUsers(t *testing.T) {
	events := []event.Event{
		{
			EventType: event.ProductPurchased,
			UserID:    "user-1",
			ProductID: "P10",
		},
		{
			EventType: event.ProductPurchased,
			UserID:    "user-2",
			ProductID: "P20",
		},
	}

	result := recommendation.CalculateInterest(
		"user-1",
		events,
	)

	if result["P10"] != 5 {
		t.Fatalf(
			"expected P10 score 5, got %.2f",
			result["P10"],
		)
	}

	if _, exists := result["P20"]; exists {
		t.Fatal("P20 should not be included")
	}
}
