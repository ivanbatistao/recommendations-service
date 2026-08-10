package recommendation_test

import (
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestScoreEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType event.Type
		expected  float64
	}{
		{
			name:      "product viewed",
			eventType: event.ProductViewed,
			expected:  1.0,
		},
		{
			name:      "search performed",
			eventType: event.SearchPerformed,
			expected:  2.0,
		},
		{
			name:      "product added to cart",
			eventType: event.ProductAddedCart,
			expected:  3.0,
		},
		{
			name:      "product purchased",
			eventType: event.ProductPurchased,
			expected:  5.0,
		},
		{
			name:      "unknown event",
			eventType: event.Type("unknown"),
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := recommendation.ScoreEvent(tt.eventType)

			if result != tt.expected {
				t.Fatalf(
					"expected score %.2f, got %.2f",
					tt.expected,
					result,
				)
			}
		})
	}
}
