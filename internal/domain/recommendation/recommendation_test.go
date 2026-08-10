package recommendation_test

import (
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestRank(t *testing.T) {
	interest := map[string]float64{
		"P10": 5,
		"P20": 10,
		"P30": 3,
		"P40": 7,
	}

	result := recommendation.Rank("user-1", interest, 3)

	if len(result) != 3 {
		t.Fatalf(
			"expected 3 recommendations, got %d",
			len(result),
		)
	}

	expected := []struct {
		productID string
		score     float64
	}{
		{"P20", 10},
		{"P40", 7},
		{"P10", 5},
	}

	for i, item := range expected {
		if result[i].ProductID != item.productID {
			t.Fatalf(
				"position %d: expected product %s, got %s",
				i,
				item.productID,
				result[i].ProductID,
			)
		}

		if result[i].Score != item.score {
			t.Fatalf(
				"position %d: expected score %.2f, got %.2f",
				i,
				item.score,
				result[i].Score,
			)
		}
	}
}

func TestRankLimitGreaterThanResults(t *testing.T) {
	interest := map[string]float64{
		"P10": 5,
		"P20": 10,
	}

	result := recommendation.Rank("user-1", interest, 10)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 recommendations, got %d",
			len(result),
		)
	}
}

func TestRankWithZeroLimit(t *testing.T) {
	interest := map[string]float64{
		"P10": 5,
	}

	result := recommendation.Rank("user-1", interest, 0)

	if len(result) != 0 {
		t.Fatalf(
			"expected 0 recommendations, got %d",
			len(result),
		)
	}
}

func TestRankWithEmptyInterest(t *testing.T) {
	result := recommendation.Rank(
		"user-1",
		map[string]float64{},
		3,
	)

	if len(result) != 0 {
		t.Fatalf(
			"expected 0 recommendations, got %d",
			len(result),
		)
	}
}
