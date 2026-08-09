package commands

import (
	"context"
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestGenerateRecommendationsHandler(t *testing.T) {
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
	}

	service := recommendation.NewService(nil)
	handler := NewGenerateRecommendationsHandler(service)

	result, err := handler.Execute(
		context.Background(),
		GenerateRecommendationsCommand{
			UserID: "user-1",
			Events: events,
			Limit:  2,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(result))
	}

	if result[0].ProductID != "P20" {
		t.Fatalf("expected first recommendation P20, got %s", result[0].ProductID)
	}

	if result[0].Score != 8 {
		t.Fatalf("expected P20 score 8, got %.2f", result[0].Score)
	}
}

func TestGenerateRecommendationsHandlerErrors(t *testing.T) {
	service := recommendation.NewService(nil)
	handler := NewGenerateRecommendationsHandler(service)

	// Test empty userID
	_, err := handler.Execute(
		context.Background(),
		GenerateRecommendationsCommand{
			UserID: "",
			Events: []event.Event{},
			Limit:  10,
		},
	)
	if err != recommendation.ErrInvalidUserID {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}

	// Test invalid limit
	_, err = handler.Execute(
		context.Background(),
		GenerateRecommendationsCommand{
			UserID: "user-1",
			Events: []event.Event{},
			Limit:  0,
		},
	)
	if err != recommendation.ErrInvalidLimit {
		t.Fatalf("expected ErrInvalidLimit, got %v", err)
	}
}
