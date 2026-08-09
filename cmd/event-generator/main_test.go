package main

import (
	"context"
	"testing"
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

func TestEventGenerator_GenerateEvent(t *testing.T) {
	productIDs := []string{"P1", "P2", "P3"}
	userIDs := []string{"user-1", "user-2"}

	generator := NewEventGenerator(productIDs, userIDs, 10, false, "http://localhost:8080/events", "test-stream")

	ev := generator.GenerateEvent()

	// Verify event has required fields
	if ev.EventID == "" {
		t.Error("EventID should not be empty")
	}

	if ev.UserID == "" {
		t.Error("UserID should not be empty")
	}

	if ev.ProductID == "" {
		t.Error("ProductID should not be empty")
	}

	// Verify UserID is from our list
	validUserID := false
	for _, uid := range userIDs {
		if ev.UserID == uid {
			validUserID = true
			break
		}
	}
	if !validUserID {
		t.Errorf("UserID %s should be from the list", ev.UserID)
	}

	// Verify ProductID is from our list
	validProductID := false
	for _, pid := range productIDs {
		if ev.ProductID == pid {
			validProductID = true
			break
		}
	}
	if !validProductID {
		t.Errorf("ProductID %s should be from the list", ev.ProductID)
	}

	// Verify event type is valid
	validEventType := false
	validTypes := []event.Type{
		event.ProductViewed,
		event.SearchPerformed,
		event.ProductAddedCart,
		event.ProductPurchased,
	}
	for _, et := range validTypes {
		if ev.EventType == et {
			validEventType = true
			break
		}
	}
	if !validEventType {
		t.Errorf("EventType %s should be valid", ev.EventType)
	}
}

func TestEventGenerator_RandomEventType(t *testing.T) {
	productIDs := []string{"P1"}
	userIDs := []string{"user-1"}
	generator := NewEventGenerator(productIDs, userIDs, 10, false, "http://localhost:8080/events", "test-stream")

	// Generate many events and verify they're all valid
	for i := 0; i < 100; i++ {
		ev := generator.GenerateEvent()
		validTypes := []event.Type{
			event.ProductViewed,
			event.SearchPerformed,
			event.ProductAddedCart,
			event.ProductPurchased,
		}

		valid := false
		for _, et := range validTypes {
			if ev.EventType == et {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("Event type %s is not valid", ev.EventType)
		}
	}
}

func TestEventGenerator_Start(t *testing.T) {
	productIDs := []string{"P1", "P2"}
	userIDs := []string{"user-1", "user-2"}
	generator := NewEventGenerator(productIDs, userIDs, 10, false, "http://localhost:8080/events", "test-stream")

	ctx := context.Background()
	duration := 100 * time.Millisecond

	// Should complete without error
	generator.Start(ctx, duration)
}

func TestEventGenerator_SendToHTTP(t *testing.T) {
	productIDs := []string{"P1"}
	userIDs := []string{"user-1"}
	
	// Use a non-existent URL - should log error but not panic
	generator := NewEventGenerator(productIDs, userIDs, 10, false, "http://localhost:9999/events", "test-stream")

	ev := generator.GenerateEvent()
	ctx := context.Background()

	// Should not panic
	generator.sendToHTTP(ctx, ev)
}
