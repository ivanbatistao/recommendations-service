package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpgin "github.com/ivanbatistao/recommendations-service/internal/infrastructure/http/gin"
)

func TestHealth(t *testing.T) {
	router := httpgin.NewRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	expectedBody := `{"status":"ok"}`

	if recorder.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %q, got %q",
			expectedBody,
			recorder.Body.String(),
		)
	}

	requestID := recorder.Header().Get("X-Request-ID")

	if requestID == "" {
		t.Fatal("expected X-Request-ID header")
	}
}

func TestHealthWithRequestID(t *testing.T) {
	router := httpgin.NewRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	req.Header.Set("X-Request-ID", "test-request-123")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	requestID := recorder.Header().Get("X-Request-ID")

	if requestID != "test-request-123" {
		t.Fatalf(
			"expected request ID %q, got %q",
			"test-request-123",
			requestID,
		)
	}
}
