package event

import "time"

type Event struct {
	EventID         string
	EventType       Type
	UserID          string
	ProductID       string
	ProductCategory string
	ProductBrand    string
	Metadata        Metadata
	OccurredAt      time.Time
}
