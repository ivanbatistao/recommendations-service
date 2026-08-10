# MiniStack Local Development Setup

## Overview

This project uses **MiniStack** to simulate AWS services locally, enabling development and testing without a real AWS account. MiniStack is a free, MIT-licensed alternative to LocalStack that provides 60+ AWS services without requiring a license key.

## Simulated Services

- **DynamoDB**: NoSQL database for storing recommendations
- **Kinesis**: Event stream for real-time processing

## Prerequisites

- Docker and Docker Compose installed
- AWS CLI installed (optional, for initialization scripts)

## Quick Start

### 1. Start MiniStack

```bash
docker-compose up -d ministack
```

This will start MiniStack at `http://localhost:4566` with DynamoDB and Kinesis enabled.

### 2. Initialize AWS Resources

Run the initialization script to create the DynamoDB table and Kinesis stream:

```bash
./scripts/init-ministack.sh
```

This will create:
- DynamoDB table: `Recommendations` (partition key: `UserID`)
- Kinesis stream: `recommendations-events` (1 shard)

### 3. Verify Functionality

```bash
./scripts/test-ministack.sh
```

This will execute connectivity tests with both services.

### 4. Start the Application

```bash
docker-compose up recommendation-service
```

The application will automatically connect to MiniStack using the configured environment variables.

## Configuration

### Environment Variables

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `AWS_REGION` | `us-east-1` | Simulated AWS region |
| `AWS_ACCESS_KEY_ID` | `test` | Fake credential for MiniStack |
| `AWS_SECRET_ACCESS_KEY` | `test` | Fake credential for MiniStack |
| `DYNAMODB_ENDPOINT` | `http://ministack:4566` | Local DynamoDB endpoint |
| `KINESIS_ENDPOINT` | `http://ministack:4566` | Local Kinesis endpoint |
| `USE_LOCAL_AWS` | `true` | Use local AWS services |

### Docker Compose

```yaml
ministack:
  image: ministackorg/ministack
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

The application automatically detects when to use MiniStack based on environment variables:

```go
if config.DynamoDBEndpoint != "" {
    // Use local DynamoDB (MiniStack)
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

### Event Generator with MiniStack

To use the Event Generator with local Kinesis:

```bash
./event-generator --kinesis --stream-name recommendations-events --endpoint http://localhost:4566
```

## Available Scripts

### `scripts/init-ministack.sh`
Initializes AWS resources in MiniStack:
- Creates DynamoDB table
- Creates Kinesis stream
- Waits for resources to be active

### `scripts/test-ministack.sh`
Tests connectivity with MiniStack:
- Lists DynamoDB tables
- Describes Recommendations table
- Lists Kinesis streams
- Writes/reads test data

## Troubleshooting

### MiniStack won't start
```bash
# Check logs
docker-compose logs ministack

# Restart
docker-compose restart ministack
```

### Table or stream don't exist
```bash
# Re-initialize resources
./scripts/init-ministack.sh
```

### Connection error
```bash
# Verify MiniStack is running
curl http://localhost:4566

# Check if port 4566 is in use
lsof -i :4566
```

### AWS permissions
MiniStack ignores real credentials but requires valid values:
```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

## MiniStack vs LocalStack

MiniStack is a free alternative to LocalStack with the following advantages:

- **Free forever**: No license key required
- **No telemetry**: No data collection
- **Drop-in replacement**: Compatible with existing AWS tools
- **60+ services**: More AWS services than LocalStack free tier
- **Real databases**: RDS runs actual Postgres/MySQL containers
- **Multi-account & multi-region**: Supports multiple AWS accounts and regions

## Limitations

MiniStack aims for high compatibility with AWS but has some limitations:

- Some advanced AWS features may not be fully implemented
- Performance may differ from real AWS
- No costs, but also no availability guarantees

## Next Steps

Once MiniStack is running:
1. ✅ Event Generator can send events to local Kinesis
2. ✅ Application can store recommendations in local DynamoDB
3. ✅ Ready for load testing with k6

## References

- [MiniStack Documentation](https://ministack.org)
- [MiniStack GitHub](https://github.com/ministackorg/ministack)
- [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)
- [Kinesis Developer Guide](https://docs.aws.amazon.com/streams/latest/dev/introduction.html)
