package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type EventGenerator struct {
	productIDs []string
	userIDs    []string
	rate       int // events per second
	kinesis    bool // true = send to Kinesis, false = send to HTTP API
	httpURL    string // HTTP API endpoint (if kinesis=false)
	streamName string // Kinesis stream name (if kinesis=true)
	logger     *slog.Logger
}

func NewEventGenerator(productIDs, userIDs []string, rate int, kinesis bool, httpURL, streamName string) *EventGenerator {
	return &EventGenerator{
		productIDs: productIDs,
		userIDs:    userIDs,
		rate:       rate,
		kinesis:    kinesis,
		httpURL:    httpURL,
		streamName: streamName,
		logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (g *EventGenerator) GenerateEvent() event.Event {
	productID := g.productIDs[rand.Intn(len(g.productIDs))]
	userID := g.userIDs[rand.Intn(len(g.userIDs))]
	eventType := g.randomEventType()

	ev := event.Event{
		EventID:         fmt.Sprintf("event-%d", time.Now().UnixNano()),
		EventType:       eventType,
		UserID:          userID,
		ProductID:       productID,
		ProductCategory: g.randomCategory(),
		ProductBrand:    g.randomBrand(),
		Metadata: event.Metadata{
			Device:  g.randomDevice(),
			Country: g.randomCountry(),
		},
		OccurredAt:      time.Now(),
	}

	return ev
}

func (g *EventGenerator) randomEventType() event.Type {
	types := []event.Type{
		event.ProductViewed,
		event.SearchPerformed,
		event.ProductAddedCart,
		event.ProductPurchased,
	}
	return types[rand.Intn(len(types))]
}

func (g *EventGenerator) randomCategory() string {
	categories := []string{"electronics", "clothing", "home", "books", "sports"}
	return categories[rand.Intn(len(categories))]
}

func (g *EventGenerator) randomBrand() string {
	brands := []string{"Apple", "Samsung", "Nike", "Adidas", "Sony", "LG"}
	return brands[rand.Intn(len(brands))]
}

func (g *EventGenerator) randomDevice() string {
	devices := []string{"mobile", "desktop", "tablet"}
	return devices[rand.Intn(len(devices))]
}

func (g *EventGenerator) randomCountry() string {
	countries := []string{"US", "ES", "UK", "DE", "FR", "IT"}
	return countries[rand.Intn(len(countries))]
}

func (g *EventGenerator) Start(ctx context.Context, duration time.Duration) {
	g.logger.Info(
		"event generator started",
		slog.Int("rate", g.rate),
		slog.Bool("kinesis", g.kinesis),
		slog.String("http_url", g.httpURL),
		slog.Duration("duration", duration),
	)

	ticker := time.NewTicker(time.Second / time.Duration(g.rate))
	defer ticker.Stop()

	endTime := time.Now().Add(duration)
	eventCount := 0

	for {
		select {
		case <-ctx.Done():
			g.logger.Info("event generator stopped", slog.Int("events_generated", eventCount))
			return

		case <-ticker.C:
			if time.Now().Before(endTime) {
				ev := g.GenerateEvent()
				g.sendEvent(ctx, ev)
				eventCount++
			} else {
				g.logger.Info("event generator finished", slog.Int("events_generated", eventCount))
				return
			}
		}
	}
}

func (g *EventGenerator) sendEvent(ctx context.Context, ev event.Event) {
	if g.kinesis {
		// TODO: Send to Kinesis when we implement Kinesis integration
		g.logger.Debug("would send to kinesis", slog.String("event_id", ev.EventID))
	} else {
		// Send to HTTP API
		g.sendToHTTP(ctx, ev)
	}
}

func (g *EventGenerator) sendToHTTP(ctx context.Context, ev event.Event) {
	body, err := json.Marshal(ev)
	if err != nil {
		g.logger.Error("failed to marshal event", slog.String("error", err.Error()))
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.httpURL, bytes.NewBuffer(body))
	if err != nil {
		g.logger.Error("failed to create request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		g.logger.Error("failed to send event", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		g.logger.Error("event rejected", slog.Int("status", resp.StatusCode))
		return
	}

	g.logger.Debug("event sent successfully", slog.String("event_id", ev.EventID))
}

func main() {
	// Configuration for realistic ecommerce data
	var (
		rate        = 10      // events per second
		duration    = 1 * time.Minute
		useKinesis  = false
		httpURL     = "http://localhost:8080/events"
		streamName  = "recommendations-events"
		productIDs  = []string{
			"PROD-001", "PROD-002", "PROD-003", "PROD-004", "PROD-005",
			"PROD-006", "PROD-007", "PROD-008", "PROD-009", "PROD-010",
			"PROD-011", "PROD-012", "PROD-013", "PROD-014", "PROD-015",
			"PROD-016", "PROD-017", "PROD-018", "PROD-019", "PROD-020",
		}
		userIDs = []string{
			"user-001", "user-002", "user-003", "user-004", "user-005",
			"user-006", "user-007", "user-008", "user-009", "user-010",
			"user-011", "user-012", "user-013", "user-014", "user-015",
		}
	)

	// TODO: Parse command line flags for configuration

	generator := NewEventGenerator(
		productIDs,
		userIDs,
		rate,
		useKinesis,
		httpURL,
		streamName,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		generator.logger.Info("received shutdown signal")
		cancel()
	}()

	// Run generator
	generator.Start(ctx, duration)
}
