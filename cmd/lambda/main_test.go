package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ivanbatistao/recommendations-service/internal/infrastructure/lambda"
)

func TestLambdaHandler_HandleHealth(t *testing.T) {
	handler := lambda.NewLambdaHandler()

	req := events.APIGatewayProxyRequest{
		Path:       "/health",
		HTTPMethod: "GET",
	}

	resp, err := handler.HandleRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}
}

func TestLambdaHandler_HandleGetRecommendations(t *testing.T) {
	handler := lambda.NewLambdaHandler()

	req := events.APIGatewayProxyRequest{
		Path:           "/recommendations/user-123",
		HTTPMethod:     "GET",
		PathParameters: map[string]string{"userId": "user-123"},
	}

	resp, err := handler.HandleRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if _, ok := body["recommendations"]; !ok {
		t.Fatal("expected recommendations field in response")
	}
}

func TestLambdaHandler_HandleProcessEvent(t *testing.T) {
	handler := lambda.NewLambdaHandler()

	eventData := map[string]interface{}{
		"event_id":   "event-1",
		"event_type": "product_viewed",
		"user_id":    "user-123",
		"product_id": "P10",
		"metadata": map[string]string{
			"device":  "mobile",
			"country": "US",
		},
		"occurred_at": "2026-08-09T00:00:00Z",
	}

	body, _ := json.Marshal(eventData)

	req := events.APIGatewayProxyRequest{
		Path:       "/events",
		HTTPMethod: "POST",
		Body:       string(body),
	}

	resp, err := handler.HandleRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 202 {
		t.Fatalf("expected status 202, got %d", resp.StatusCode)
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Body), &responseBody); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if responseBody["status"] != "event processed" {
		t.Fatalf("expected status 'event processed', got %v", responseBody["status"])
	}
}

func TestLambdaHandler_HandleGenerateRecommendations(t *testing.T) {
	handler := lambda.NewLambdaHandler()

	requestData := map[string]interface{}{
		"user_id": "user-456",
		"events": []map[string]interface{}{
			{
				"event_id":   "event-1",
				"event_type": "product_viewed",
				"user_id":    "user-456",
				"product_id": "P20",
				"occurred_at": "2026-08-09T00:00:00Z",
			},
		},
		"limit": 5,
	}

	body, _ := json.Marshal(requestData)

	req := events.APIGatewayProxyRequest{
		Path:       "/recommendations/generate",
		HTTPMethod: "POST",
		Body:       string(body),
	}

	resp, err := handler.HandleRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Body), &responseBody); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if _, ok := responseBody["recommendations"]; !ok {
		t.Fatal("expected recommendations field in response")
	}
}

func TestLambdaHandler_HandleNotFound(t *testing.T) {
	handler := lambda.NewLambdaHandler()

	req := events.APIGatewayProxyRequest{
		Path:       "/unknown",
		HTTPMethod: "GET",
	}

	resp, err := handler.HandleRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}
