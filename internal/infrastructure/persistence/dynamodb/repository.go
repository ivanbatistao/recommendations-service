package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBRepository) GetByUserID(
	ctx context.Context,
	userID string,
) ([]recommendation.Recommendation, error) {
	// Query items by partition key (UserID)
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditionExpression: aws.String("UserID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query recommendations: %w", err)
	}

	// Convert DynamoDB items to domain entities
	recommendations := make([]recommendation.Recommendation, 0, len(result.Items))
	for _, item := range result.Items {
		var rec recommendation.Recommendation
		if err := attributevalue.UnmarshalMap(item, &rec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recommendation: %w", err)
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

func (r *DynamoDBRepository) Save(
	ctx context.Context,
	rec recommendation.Recommendation,
) error {
	// Convert domain entity to DynamoDB item
	item, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendation: %w", err)
	}

	// Put item (Create or Update)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save recommendation: %w", err)
	}

	return nil
}
