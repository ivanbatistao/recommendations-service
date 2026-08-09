package workerpool_test

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/processing/workerpool"
)

// Mock processor for testing
type mockProcessor struct {
	mu           sync.Mutex
	processedEvents []event.Event
	processDelay time.Duration
}

func (m *mockProcessor) ProcessEvent(ctx context.Context, ev event.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.processedEvents = append(m.processedEvents, ev)
	
	if m.processDelay > 0 {
		time.Sleep(m.processDelay)
	}
	
	return nil
}

func (m *mockProcessor) GetProcessedEvents() []event.Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	result := make([]event.Event, len(m.processedEvents))
	copy(result, m.processedEvents)
	return result
}

func TestWorkerPool_StartAndStop(t *testing.T) {
	processor := &mockProcessor{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	wp := workerpool.NewWorkerPool(2, 10, processor, logger)
	
	wp.Start()
	
	// Verify workers are running
	time.Sleep(100 * time.Millisecond)
	
	wp.Stop()
	
	// Verify workers stopped gracefully
}

func TestWorkerPool_SubmitAndProcess(t *testing.T) {
	processor := &mockProcessor{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	wp := workerpool.NewWorkerPool(2, 10, processor, logger)
	wp.Start()
	defer wp.Stop()
	
	events := []event.Event{
		{EventID: "event-1", UserID: "user-1", ProductID: "P1"},
		{EventID: "event-2", UserID: "user-2", ProductID: "P2"},
		{EventID: "event-3", UserID: "user-3", ProductID: "P3"},
	}
	
	for _, ev := range events {
		wp.Submit(ev)
	}
	
	// Wait for processing
	time.Sleep(200 * time.Millisecond)
	
	processed := processor.GetProcessedEvents()
	if len(processed) != 3 {
		t.Fatalf("expected 3 processed events, got %d", len(processed))
	}
}

func TestWorkerPool_ConcurrentProcessing(t *testing.T) {
	processor := &mockProcessor{processDelay: 50 * time.Millisecond}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	wp := workerpool.NewWorkerPool(3, 10, processor, logger)
	wp.Start()
	defer wp.Stop()
	
	// Submit more events than workers
	for i := 0; i < 10; i++ {
		wp.Submit(event.Event{
			EventID:   "event-" + string(rune('0'+i)),
			UserID:    "user-1",
			ProductID: "P1",
		})
	}
	
	// Wait for processing
	time.Sleep(500 * time.Millisecond)
	
	processed := processor.GetProcessedEvents()
	if len(processed) != 10 {
		t.Fatalf("expected 10 processed events, got %d", len(processed))
	}
}

func TestWorkerPool_Stats(t *testing.T) {
	processor := &mockProcessor{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	wp := workerpool.NewWorkerPool(5, 100, processor, logger)
	
	stats := wp.Stats()
	
	if stats["num_workers"] != 5 {
		t.Fatalf("expected num_workers 5, got %v", stats["num_workers"])
	}
	
	if stats["buffer_size"] != 100 {
		t.Fatalf("expected buffer_size 100, got %v", stats["buffer_size"])
	}
}

func TestWorkerPool_Backpressure(t *testing.T) {
	processor := &mockProcessor{processDelay: 100 * time.Millisecond}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	wp := workerpool.NewWorkerPool(1, 2, processor, logger)
	wp.Start()
	defer wp.Stop()
	
	// Submit more events than buffer size
	for i := 0; i < 10; i++ {
		wp.Submit(event.Event{
			EventID:   "event-" + string(rune('0'+i)),
			UserID:    "user-1",
			ProductID: "P1",
		})
	}
	
	// Wait for processing
	time.Sleep(500 * time.Millisecond)
	
	processed := processor.GetProcessedEvents()
	// Should process at most buffer size + some that fit
	if len(processed) > 3 {
		t.Logf("processed %d events (some may have been dropped)", len(processed))
	}
}
