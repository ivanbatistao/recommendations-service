package event

import "time"

type Event struct {
	EventID         string    `json:"event_id"`
	EventType       Type      `json:"event_type"`
	UserID          string    `json:"user_id"`
	ProductID       string    `json:"product_id"`
	ProductCategory string    `json:"product_category,omitempty"`
	ProductBrand    string    `json:"product_brand,omitempty"`
	Metadata        Metadata  `json:"metadata"`
	OccurredAt      time.Time `json:"occurred_at"`
}
