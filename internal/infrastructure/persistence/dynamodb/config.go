package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBClient(ctx context.Context, region string) (*dynamodb.Client, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	log.Printf("DynamoDB client initialized for region: %s", region)

	return client, nil
}

func NewLocalDynamoDBClient(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	// Create custom configuration for local DynamoDB
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	// Override endpoint for local development
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	log.Printf("DynamoDB client initialized for local endpoint: %s", endpoint)

	return client, nil
}
