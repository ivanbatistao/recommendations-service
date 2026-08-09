package dto

import (
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type EventDTO struct {
	EventID         string            `json:"event_id"`
	EventType       string            `json:"event_type"`
	UserID          string            `json:"user_id"`
	ProductID       string            `json:"product_id"`
	ProductCategory string            `json:"product_category,omitempty"`
	ProductBrand    string            `json:"product_brand,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	OccurredAt      string            `json:"occurred_at"`
}

func FromEventDomain(e event.Event) EventDTO {
	return EventDTO{
		EventID:         e.EventID,
		EventType:       string(e.EventType),
		UserID:          e.UserID,
		ProductID:       e.ProductID,
		ProductCategory: e.ProductCategory,
		ProductBrand:    e.ProductBrand,
		Metadata: map[string]string{
			"device":  e.Metadata.Device,
			"country": e.Metadata.Country,
		},
		OccurredAt: e.OccurredAt.Format(time.RFC3339),
	}
}

func ToEventDomain(dto EventDTO) (event.Event, error) {
	occurredAt, err := time.Parse(time.RFC3339, dto.OccurredAt)
	if err != nil {
		return event.Event{}, err
	}

	return event.Event{
		EventID:         dto.EventID,
		EventType:       event.Type(dto.EventType),
		UserID:          dto.UserID,
		ProductID:       dto.ProductID,
		ProductCategory: dto.ProductCategory,
		ProductBrand:    dto.ProductBrand,
		Metadata: event.Metadata{
			Device:  dto.Metadata["device"],
			Country: dto.Metadata["country"],
		},
		OccurredAt: occurredAt,
	}, nil
}
