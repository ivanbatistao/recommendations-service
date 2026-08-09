package recommendation

import (
	"sort"
)

func Rank(
	userID string,
	interest map[string]float64,
	limit int,
) []Recommendation {
	if limit <= 0 || len(interest) == 0 {
		return []Recommendation{}
	}

	recommendations := make([]Recommendation, 0, len(interest))

	for productID, score := range interest {
		recommendations = append(
			recommendations,
			Recommendation{
				UserID:    userID,
				ProductID: productID,
				Score:     score,
			},
		)
	}

	sort.Slice(
		recommendations,
		func(i, j int) bool {
			return recommendations[i].Score >
				recommendations[j].Score
		},
	)

	if limit > len(recommendations) {
		limit = len(recommendations)
	}

	return recommendations[:limit]
}
