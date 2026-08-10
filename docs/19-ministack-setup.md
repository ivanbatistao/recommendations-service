# MiniStack Local Development Setup

## Overview

This project uses **LocalStack** to simulate AWS services locally, enabling development and testing without a real AWS account.

## Simulated Services

- **DynamoDB**: NoSQL database for storing recommendations
- **Kinesis**: Event stream for real-time processing

## Prerequisites

- Docker and Docker Compose installed
- AWS CLI installed (optional, for initialization scripts)

## Quick Start

### 1. Start LocalStack

```bash
docker-compose up -d localstack
```

This will start LocalStack at `http://localhost:4566` with DynamoDB and Kinesis enabled.

### 2. Initialize AWS Resources

Run the initialization script to create the DynamoDB table and Kinesis stream:

```bash
./scripts/init-localstack.sh
```

This will create:
- DynamoDB table: `Recommendations` (partition key: `UserID`)
- Kinesis stream: `recommendations-events` (1 shard)

### 3. Verify Functionality

```bash
./scripts/test-localstack.sh
```

This will execute connectivity tests with both services.

### 4. Start the Application

```bash
docker-compose up recommendation-service
```

The application will automatically connect to LocalStack using the configured environment variables.

## Configuration

### Environment Variables

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `AWS_REGION` | `us-east-1` | Simulated AWS region |
| `AWS_ACCESS_KEY_ID` | `test` | Fake credential for LocalStack |
| `AWS_SECRET_ACCESS_KEY` | `test` | Fake credential for LocalStack |
| `DYNAMODB_ENDPOINT` | `http://localstack:4566` | Local DynamoDB endpoint |
| `KINESIS_ENDPOINT` | `http://localstack:4566` | Local Kinesis endpoint |
| `USE_LOCAL_AWS` | `true` | Use local AWS services |

### Docker Compose

```yaml
localstack:
  image: localstack/localstack:latest
  ports:
    - "4566:4566"
  environment:
    - SERVICES=dynamodb,kinesis
    - AWS_DEFAULT_REGION=us-east-1
    - AWS_ACCESS_KEY_ID=test
    - AWS_SECRET_ACCESS_KEY=test
```

## Usage

### Local Development

The application automatically detects when to use LocalStack based on environment variables:

```go
if config.DynamoDBEndpoint != "" {
    // Use local DynamoDB (LocalStack)
    client, err = dynamodb.NewLocalDynamoDBClient(
        context.Background(),
        config.DynamoDBEndpoint,
    )
} else {
    // Use real AWS DynamoDB
    client, err = dynamodb.NewDynamoDBClient(
        context.Background(),
        config.AWSRegion,
    )
}
```

### Event Generator with LocalStack

To use the Event Generator with local Kinesis:

```bash
./event-generator --kinesis --stream-name recommendations-events --endpoint http://localhost:4566
```

## Available Scripts

### `scripts/init-localstack.sh`
Initializes AWS resources in LocalStack:
- Creates DynamoDB table
- Creates Kinesis stream
- Waits for resources to be active

### `scripts/test-localstack.sh`
Tests connectivity with LocalStack:
- Lists DynamoDB tables
- Describes Recommendations table
- Lists Kinesis streams
- Writes/reads test data

## Troubleshooting

### LocalStack won't start
```bash
# Check logs
docker-compose logs localstack

# Restart
docker-compose restart localstack
```

### Table or stream don't exist
```bash
# Re-initialize resources
./scripts/init-localstack.sh
```

### Connection error
```bash
# Verify LocalStack is running
curl http://localhost:4566

# Check if port 4566 is in use
lsof -i :4566
```

### AWS permissions
LocalStack ignores real credentials but requires valid values:
```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

## Local Data

LocalStack data is saved in `./localstack_data/`:
```yaml
volumes:
  - "./localstack_data:/tmp/localstack/data"
```

To clean data:
```bash
rm -rf localstack_data/*
docker-compose restart localstack
```

## Limitations

LocalStack is not 100% compatible with real AWS. Known limitations:

- Some advanced DynamoDB features are not implemented
- Kinesis has some behavioral differences
- No costs, but also no availability guarantees

## Next Steps

Once LocalStack is running:
1. ✅ Event Generator can send events to local Kinesis
2. ✅ Application can store recommendations in local DynamoDB
3. ✅ Ready for load testing with k6

## References

- [LocalStack Documentation](https://docs.localstack.cloud/)
- [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)
- [Kinesis Developer Guide](https://docs.aws.amazon.com/streams/latest/dev/introduction.html)
