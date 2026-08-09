package kinesis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type Consumer struct {
	client     *kinesis.Client
	streamName string
}

func NewConsumer(client *kinesis.Client, streamName string) *Consumer {
	return &Consumer{
		client:     client,
		streamName: streamName,
	}
}

func (c *Consumer) GetRecords(
	ctx context.Context,
	shardIterator string,
	limit int,
) ([]event.Event, string, error) {
	input := &kinesis.GetRecordsInput{
		ShardIterator: aws.String(shardIterator),
		Limit:         aws.Int32(int32(limit)),
	}

	result, err := c.client.GetRecords(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get records: %w", err)
	}

	events := make([]event.Event, 0, len(result.Records))
	for _, record := range result.Records {
		var ev event.Event
		if err := json.Unmarshal(record.Data, &ev); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal event: %w", err)
		}
		events = append(events, ev)
	}

	nextIterator := ""
	if result.NextShardIterator != nil {
		nextIterator = *result.NextShardIterator
	}

	return events, nextIterator, nil
}

func (c *Consumer) GetShardIterator(
	ctx context.Context,
	shardID string,
	iteratorType types.ShardIteratorType,
) (string, error) {
	input := &kinesis.GetShardIteratorInput{
		StreamName:      aws.String(c.streamName),
		ShardId:         aws.String(shardID),
		ShardIteratorType: iteratorType,
	}

	result, err := c.client.GetShardIterator(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get shard iterator: %w", err)
	}

	return *result.ShardIterator, nil
}

func (c *Consumer) ListShards(
	ctx context.Context,
) ([]string, error) {
	input := &kinesis.ListShardsInput{
		StreamName: aws.String(c.streamName),
	}

	result, err := c.client.ListShards(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list shards: %w", err)
	}

	shards := make([]string, len(result.Shards))
	for i, shard := range result.Shards {
		shards[i] = *shard.ShardId
	}

	return shards, nil
}
