package recommendation

import "github.com/ivanbatistao/recommendations-service/internal/domain/event"

func CalculateInterest(
	userID string,
	events []event.Event,
) map[string]float64 {
	interest := make(map[string]float64)

	for _, currentEvent := range events {
		if currentEvent.UserID != userID {
			continue
		}

		score := ScoreEvent(currentEvent.EventType)

		if score == 0 {
			continue
		}

		interest[currentEvent.ProductID] += score
	}

	return interest
}
