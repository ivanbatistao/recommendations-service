#!/bin/bash

# LocalStack initialization script
# Creates DynamoDB table and Kinesis stream for local development

set -e

LOCALSTACK_ENDPOINT="http://localhost:4566"
AWS_REGION="us-east-1"
TABLE_NAME="Recommendations"
STREAM_NAME="recommendations-events"

echo "🚀 Initializing LocalStack resources..."

# Create DynamoDB table
echo "📦 Creating DynamoDB table: $TABLE_NAME"
aws dynamodb create-table \
  --table-name $TABLE_NAME \
  --attribute-definitions AttributeName=UserID,AttributeType=S \
  --key-schema AttributeName=UserID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --endpoint-url $LOCALSTACK_ENDPOINT \
  --region $AWS_REGION \
  || echo "Table might already exist"

# Wait for table to be created
echo "⏳ Waiting for table to be active..."
aws dynamodb wait table-exists \
  --table-name $TABLE_NAME \
  --endpoint-url $LOCALSTACK_ENDPOINT \
  --region $AWS_REGION

echo "✅ DynamoDB table is ready"

# Create Kinesis stream
echo "🌊 Creating Kinesis stream: $STREAM_NAME"
aws kinesis create-stream \
  --stream-name $STREAM_NAME \
  --shard-count 1 \
  --endpoint-url $LOCALSTACK_ENDPOINT \
  --region $AWS_REGION \
  || echo "Stream might already exist"

# Wait for stream to be created
echo "⏳ Waiting for stream to be active..."
aws kinesis wait stream-exists \
  --stream-name $STREAM_NAME \
  --endpoint-url $LOCALSTACK_ENDPOINT \
  --region $AWS_REGION

echo "✅ Kinesis stream is ready"

# List resources
echo ""
echo "📋 Current resources:"
echo "DynamoDB Tables:"
aws dynamodb list-tables --endpoint-url $LOCALSTACK_ENDPOINT --region $AWS_REGION

echo ""
echo "Kinesis Streams:"
aws kinesis list-streams --endpoint-url $LOCALSTACK_ENDPOINT --region $AWS_REGION

echo ""
echo "✨ LocalStack initialization complete!"
