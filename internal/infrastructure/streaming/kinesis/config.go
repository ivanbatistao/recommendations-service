package kinesis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func NewKinesisClient(ctx context.Context, region string) (*kinesis.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := kinesis.NewFromConfig(cfg)

	log.Printf("Kinesis client initialized for region: %s", region)

	return client, nil
}

func NewLocalKinesisClient(ctx context.Context, endpoint string) (*kinesis.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	client := kinesis.NewFromConfig(cfg, func(o *kinesis.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	log.Printf("Kinesis client initialized for local endpoint: %s", endpoint)

	return client, nil
}
