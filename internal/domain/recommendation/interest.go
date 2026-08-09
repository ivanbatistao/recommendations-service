package recommendation

import "github.com/ivanbatistao/recommendations-service/internal/domain/event"

func CalculateInterest(events []event.Event) map[string]float64 {
	interest := make(map[string]float64)

	for _, currentEvent := range events {
		score := ScoreEvent(currentEvent.EventType)

		if score == 0 {
			continue
		}

		interest[currentEvent.ProductID] += score
	}

	return interest
}
