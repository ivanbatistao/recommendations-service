package recommendation

import "github.com/ivanbatistao/recommendations-service/internal/domain/event"

const (
	viewedScore    = 1.0
	searchScore    = 2.0
	addedCartScore = 3.0
	purchasedScore = 5.0
)

func ScoreEvent(eventType event.Type) float64 {
	switch eventType {
	case event.ProductViewed:
		return viewedScore
	case event.SearchPerformed:
		return searchScore
	case event.ProductAddedCart:
		return addedCartScore
	case event.ProductPurchased:
		return purchasedScore
	default:
		return 0
	}
}
