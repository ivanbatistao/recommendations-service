package commands

import (
	"context"
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type mockRepository struct {
	recommendations []recommendation.Recommendation
	saveCalled      bool
	savedRecommendation recommendation.Recommendation
}

func (m *mockRepository) GetByUserID(
	_ context.Context,
	userID string,
) ([]recommendation.Recommendation, error) {
	if userID == "user-1" {
		return m.recommendations, nil
	}
	return []recommendation.Recommendation{}, nil
}

func (m *mockRepository) Save(
	_ context.Context,
	rec recommendation.Recommendation,
) error {
	m.saveCalled = true
	m.savedRecommendation = rec
	return nil
}

func TestProcessEventHandler(t *testing.T) {
	tests := []struct {
		name           string
		event          event.Event
		existingRecs   []recommendation.Recommendation
		expectSave     bool
		expectedScore  float64
	}{
		{
			name: "update existing recommendation",
			event: event.Event{
				EventID:   "event-1",
				EventType: event.ProductViewed,
				UserID:    "user-1",
				ProductID: "P10",
			},
			existingRecs: []recommendation.Recommendation{
				{
					UserID:    "user-1",
					ProductID: "P10",
					Score:     2.0,
				},
			},
			expectSave:    true,
			expectedScore: 3.0, // 2.0 + 1.0 (viewed)
		},
		{
			name: "create new recommendation",
			event: event.Event{
				EventID:   "event-1",
				EventType: event.ProductPurchased,
				UserID:    "user-1",
				ProductID: "P20",
			},
			existingRecs: []recommendation.Recommendation{
				{
					UserID:    "user-1",
					ProductID: "P10",
					Score:     2.0,
				},
			},
			expectSave:    true,
			expectedScore: 5.0, // new purchase
		},
		{
			name: "ignore unknown event type",
			event: event.Event{
				EventID:   "event-1",
				EventType: "unknown_type",
				UserID:    "user-1",
				ProductID: "P10",
			},
			existingRecs: []recommendation.Recommendation{
				{
					UserID:    "user-1",
					ProductID: "P10",
					Score:     2.0,
				},
			},
			expectSave: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				recommendations: tt.existingRecs,
			}

			service := recommendation.NewService(repo)
			handler := NewProcessEventHandler(service)

			err := handler.Execute(
				context.Background(),
				ProcessEventCommand{
					Event: tt.event,
				},
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.saveCalled != tt.expectSave {
				t.Fatalf("expected saveCalled=%v, got %v", tt.expectSave, repo.saveCalled)
			}

			if tt.expectSave && repo.savedRecommendation.Score != tt.expectedScore {
				t.Fatalf("expected score %.2f, got %.2f", tt.expectedScore, repo.savedRecommendation.Score)
			}
		})
	}
}
