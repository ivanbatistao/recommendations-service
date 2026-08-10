#!/bin/bash

# Test script for MiniStack integration
# Tests DynamoDB and Kinesis connectivity

set -e

MINISTACK_ENDPOINT="http://localhost:4566"
AWS_REGION="us-east-1"

echo "🧪 Testing MiniStack integration..."

# Test 1: List DynamoDB tables
echo "📦 Test 1: Listing DynamoDB tables..."
aws dynamodb list-tables --endpoint-url $MINISTACK_ENDPOINT --region $AWS_REGION

# Test 2: Describe Recommendations table
echo ""
echo "📦 Test 2: Describing Recommendations table..."
aws dynamodb describe-table --table-name Recommendations --endpoint-url $MINISTACK_ENDPOINT --region $AWS_REGION

# Test 3: List Kinesis streams
echo ""
echo "🌊 Test 3: Listing Kinesis streams..."
aws kinesis list-streams --endpoint-url $MINISTACK_ENDPOINT --region $AWS_REGION

# Test 4: Describe Kinesis stream
echo ""
echo "🌊 Test 4: Describing recommendations-events stream..."
aws kinesis describe-stream --stream-name recommendations-events --endpoint-url $MINISTACK_ENDPOINT --region $AWS_REGION

# Test 5: Put item to DynamoDB
echo ""
echo "📦 Test 5: Writing test item to DynamoDB..."
aws dynamodb put-item \
  --table-name Recommendations \
  --item '{"UserID": {"S": "test-user-123"}, "Recommendations": {"S": "[{\"product_id\":\"P1\",\"score\":100},{\"product_id\":\"P2\",\"score\":90}]"}}' \
  --endpoint-url $MINISTACK_ENDPOINT \
  --region $AWS_REGION

# Test 6: Get item from DynamoDB
echo ""
echo "📦 Test 6: Reading test item from DynamoDB..."
aws dynamodb get-item \
  --table-name Recommendations \
  --key '{"UserID": {"S": "test-user-123"}}' \
  --endpoint-url $MINISTACK_ENDPOINT \
  --region $AWS_REGION

# Test 7: Put record to Kinesis
echo ""
echo "🌊 Test 7: Writing test record to Kinesis..."
aws kinesis put-record \
  --stream-name recommendations-events \
  --data '{"event_id":"test-event-123","event_type":"product_viewed","user_id":"test-user-123","product_id":"P1"}' \
  --partition-key "test-user-123" \
  --endpoint-url $MINISTACK_ENDPOINT \
  --region $AWS_REGION

# Test 8: Get records from Kinesis
echo ""
echo "🌊 Test 8: Reading records from Kinesis stream..."
SHARD_ITERATOR=$(aws kinesis get-shard-iterator \
  --stream-name recommendations-events \
  --shard-id shardId-000000000000 \
  --shard-iterator-type TRIM_HORIZON \
  --endpoint-url $MINISTACK_ENDPOINT \
  --region $AWS_REGION \
  --query 'ShardIterator' \
  --output text)

aws kinesis get-records \
  --shard-iterator $SHARD_ITERATOR \
  --endpoint-url $MINISTACK_ENDPOINT \
  --region $AWS_REGION

echo ""
echo "✅ All MiniStack tests passed!"
