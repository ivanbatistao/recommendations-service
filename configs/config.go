package configs

import (
	"os"
	"strconv"
)

type Config struct {
	Port              string
	UseDynamoDB       bool
	DynamoDBTable     string
	DynamoDBEndpoint  string
	AWSRegion         string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	useDynamoDB, _ := strconv.ParseBool(os.Getenv("USE_DYNAMODB"))
	dynamoDBTable := os.Getenv("DYNAMODB_TABLE")
	if dynamoDBTable == "" {
		dynamoDBTable = "Recommendations"
	}

	dynamoDBEndpoint := os.Getenv("DYNAMODB_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}

	return Config{
		Port:              port,
		UseDynamoDB:       useDynamoDB,
		DynamoDBTable:     dynamoDBTable,
		DynamoDBEndpoint:  dynamoDBEndpoint,
		AWSRegion:         awsRegion,
	}
}
