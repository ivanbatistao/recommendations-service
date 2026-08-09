package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ivanbatistao/recommendations-service/internal/application/commands"
	"github.com/ivanbatistao/recommendations-service/internal/application/queries"
	"github.com/ivanbatistao/recommendations-service/internal/application/dto"
	httpgin "github.com/ivanbatistao/recommendations-service/internal/infrastructure/http/gin"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/persistence/memory"
	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestGetRecommendations(t *testing.T) {
	repository := memory.NewMemoryRepository()
	service := recommendation.NewService(repository)

	repository.Save(nil, recommendation.Recommendation{
		UserID:    "user-123",
		ProductID: "P10",
		Score:     0.95,
	})

	getHandler := queries.NewGetRecommendationsHandler(service)
	handler := httpgin.NewHandler(getHandler, nil, nil)
	router := httpgin.NewRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/recommendations/user-123", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Recommendations []dto.RecommendationDTO `json:"recommendations"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(response.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(response.Recommendations))
	}

	if response.Recommendations[0].ProductID != "P10" {
		t.Fatalf("expected product P10, got %s", response.Recommendations[0].ProductID)
	}
}

func TestProcessEvent(t *testing.T) {
	repository := memory.NewMemoryRepository()
	service := recommendation.NewService(repository)

	processHandler := commands.NewProcessEventHandler(service)
	handler := httpgin.NewHandler(nil, processHandler, nil)
	router := httpgin.NewRouter(handler)

	eventDTO := dto.EventDTO{
		EventID:   "event-1",
		EventType: string(event.ProductViewed),
		UserID:    "user-123",
		ProductID: "P10",
		Metadata: map[string]string{
			"device":  "mobile",
			"country": "US",
		},
		OccurredAt: time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(eventDTO)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, recorder.Code)
	}

	recs, _ := repository.GetByUserID(nil, "user-123")
	if len(recs) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recs))
	}

	if recs[0].ProductID != "P10" {
		t.Fatalf("expected product P10, got %s", recs[0].ProductID)
	}
}

func TestGenerateRecommendations(t *testing.T) {
	service := recommendation.NewService(nil)
	generateHandler := commands.NewGenerateRecommendationsHandler(service)
	handler := httpgin.NewHandler(nil, nil, generateHandler)
	router := httpgin.NewRouter(handler)

	request := struct {
		UserID string          `json:"user_id"`
		Events []dto.EventDTO  `json:"events"`
		Limit  int             `json:"limit"`
	}{
		UserID: "user-123",
		Events: []dto.EventDTO{
			{
				EventID:   "event-1",
				EventType: string(event.ProductViewed),
				UserID:    "user-123",
				ProductID: "P10",
				OccurredAt: time.Now().Format(time.RFC3339),
			},
			{
				EventID:   "event-2",
				EventType: string(event.ProductPurchased),
				UserID:    "user-123",
				ProductID: "P20",
				OccurredAt: time.Now().Format(time.RFC3339),
			},
		},
		Limit: 2,
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/recommendations/generate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Recommendations []dto.RecommendationDTO `json:"recommendations"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(response.Recommendations) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(response.Recommendations))
	}

	if response.Recommendations[0].ProductID != "P20" {
		t.Fatalf("expected first product P20, got %s", response.Recommendations[0].ProductID)
	}
}
