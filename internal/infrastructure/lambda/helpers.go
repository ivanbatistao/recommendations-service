package lambda

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

// HTTP Response Helpers
func jsonResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return jsonError(500, "failed to marshal response")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(data),
	}, nil
}

func jsonError(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error":"` + message + `"}`,
	}, nil
}

// Request Validation Helpers
func validateUserID(userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func validateRequestBody(body string) error {
	if body == "" {
		return errors.New("request body is required")
	}
	return nil
}
