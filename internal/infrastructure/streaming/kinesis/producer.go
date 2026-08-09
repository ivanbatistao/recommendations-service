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

type Producer struct {
	client    *kinesis.Client
	streamName string
}

func NewProducer(client *kinesis.Client, streamName string) *Producer {
	return &Producer{
		client:    client,
		streamName: streamName,
	}
}

func (p *Producer) PublishEvent(
	ctx context.Context,
	ev event.Event,
) error {
	// Convert event to JSON
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// PutRecord operation
	input := &kinesis.PutRecordInput{
		StreamName:   aws.String(p.streamName),
		Data:         data,
		PartitionKey: aws.String(ev.UserID), // Use UserID as partition key
	}

	_, err = p.client.PutRecord(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put record: %w", err)
	}

	return nil
}

func (p *Producer) PublishBatch(
	ctx context.Context,
	events []event.Event,
) error {
	if len(events) == 0 {
		return nil
	}

	// Kinesis supports up to 500 records per PutRecords request
	records := make([]types.PutRecordsRequestEntry, len(events))

	for i, ev := range events {
		data, err := json.Marshal(ev)
		if err != nil {
			return fmt.Errorf("failed to marshal event %d: %w", i, err)
		}

		records[i] = types.PutRecordsRequestEntry{
			Data:         data,
			PartitionKey: aws.String(ev.UserID),
		}
	}

	input := &kinesis.PutRecordsInput{
		StreamName: aws.String(p.streamName),
		Records:    records,
	}

	_, err := p.client.PutRecords(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put records: %w", err)
	}

	return nil
}
