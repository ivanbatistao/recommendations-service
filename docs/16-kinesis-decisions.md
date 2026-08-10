# Kinesis - Implementation Decisions

## Module 6 — Event Processing with Kinesis

### Objective

Implement real-time event processing using AWS Kinesis Data Streams.

### What is Kinesis Data Streams?

**AWS Kinesis Data Streams** is a real-time data streaming service that:

- **Receives** events from multiple sources (producers)
- **Stores** events temporarily (up to 365 days)
- **Allows** multiple consumers to read events
- **Scales** automatically based on data volume

### System Architecture with Kinesis

```text
Ecommerce / Event Generator
            |
            v
       Kinesis Stream
            |
    Shard (Partition Key: UserID)
            |
 vvvvvvvvvvvvvvvvvvvvvv
 Recommendation Processor
            |
       Worker Pool
            |
            v
      DynamoDB
```

### Kinesis Stream Design

#### Stream Configuration

**Stream name**: `recommendation-events`

**Initial configuration:**
- **Shards**: 1 (for development)
- **Retention period**: 24 hours
- **Partition key**: `UserID`

#### Partition Key Decision

**Choice:** UserID as partition key

**Reasons:**
1. **User ordering**: All events from a user go to the same shard
2. **Avoids race conditions**: Sequential processing per user
3. **Natural distribution**: If there are many users, they distribute well across shards
4. **Primary use case**: Event processing by user

**Trade-offs:**
- **Advantage**: Ordered processing by user
- **Disadvantage**: If a user is very active, it can overload a shard
- **Mitigation**: Automatic shard scaling based on throughput

### Producer Implementation

#### Producer

**Location**: `internal/infrastructure/streaming/kinesis/producer.go`

**Structure:**
```go
type Producer struct {
    client    *kinesis.Client
    streamName string
}
```

**Methods:**

##### 1. PublishEvent
```go
func (p *Producer) PublishEvent(
    ctx context.Context,
    ev event.Event,
) error
```

**Implementation:**
- Converts event to JSON
- Uses `PutRecord` operation
- Uses `UserID` as partition key

**Analogy:**
- Like publishing a message to a queue
- But with partition key for distribution
- Events with same partition key go to the same shard

##### 2. PublishBatch
```go
func (p *Producer) PublishBatch(
    ctx context.Context,
    events []event.Event,
) error
```

**Implementation:**
- Converts multiple events to JSON
- Uses `PutRecords` operation
- Supports up to 500 records per call

**Advantages:**
- **More efficient**: One HTTP call instead of many
- **Limit**: Maximum 500 records per call
- **Cost**: Fewer API calls

### Consumer Implementation

#### Consumer

**Location**: `internal/infrastructure/streaming/kinesis/consumer.go`

**Structure:**
```go
type Consumer struct {
    client     *kinesis.Client
    streamName string
}
```

**Methods:**

##### 1. GetShardIterator
```go
func (c *Consumer) GetShardIterator(
    ctx context.Context,
    shardID string,
    iteratorType types.ShardIteratorType,
) (string, error)
```

**Implementation:**
- Gets a "cursor" to read events from a shard
- The iterator type determines where to start

**Iterator Types:**
- **TRIM_HORIZON**: Reads from the oldest available event
- **LATEST**: Reads only new events (after the iterator)
- **AT_SEQUENCE_NUMBER**: Reads from a specific sequence number

##### 2. GetRecords
```go
func (c *Consumer) GetRecords(
    ctx context.Context,
    shardIterator string,
    limit int,
) ([]event.Event, string, error)
```

**Implementation:**
- Reads events from the shard iterator
- Converts JSON to domain events
- Returns events and the next iterator

**Features:**
- **Pagination**: Allows sequential event reading
- **Limit**: Maximum 1000 records per call
- **Next iterator**: To continue reading

##### 3. ListShards
```go
func (c *Consumer) ListShards(
    ctx context.Context,
) ([]string, error)
```

**Implementation:**
- Lists all shards in the stream
- Returns shard IDs

**Usage:**
- Necessary to know which shards to consume
- Useful for monitoring and scaling

### Kinesis Client Configuration

#### NewKinesisClient (AWS Real)

```go
func NewKinesisClient(ctx context.Context, region string) (*kinesis.Client, error)
```

**Features:**
- Uses `config.LoadDefaultConfig` to automatically load credentials
- Looks for credentials in:
  1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  2. ~/.aws/credentials file
  3. IAM roles (if in EC2/Lambda)
- **Usage**: Production in AWS

#### NewLocalKinesisClient (Local Development)

```go
func NewLocalKinesisClient(ctx context.Context, endpoint string) (*kinesis.Client, error)
```

**Features:**
- Override endpoint with `BaseEndpoint`
- Points to local: `http://localhost:4566` (LocalStack)
- **Usage**: Development with LocalStack

### JSON Serialization

#### Tags in Event Entity

```go
type Event struct {
    EventID         string    `json:"event_id"`
    EventType       Type      `json:"event_type"`
    UserID          string    `json:"user_id"`
    ProductID       string    `json:"product_id"`
    ProductCategory string    `json:"product_category,omitempty"`
    ProductBrand    string    `json:"product_brand,omitempty"`
    Metadata        Metadata  `json:"metadata"`
    OccurredAt      time.Time `json:"occurred_at"`
}
```

**Purpose:**
- Explicit control of JSON mapping
- Snake_case names for JSON (HTTP convention)
- PascalCase names in Go (Go convention)

**Without tags vs With tags:**
- **Without tags**: JSON tags would be eventID, eventType (not HTTP convention)
- **With tags**: event_id, event_type (standard HTTP convention)

### Installed AWS SDK Packages

1. **`github.com/aws/aws-sdk-go-v2/service/kinesis`**: Kinesis specific client
2. **`github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream`**: Streaming protocol

### Tests

#### Current Status

Kinesis tests are in `Skip` because:

1. **SDK doesn't use interfaces**: Kinesis client is a concrete struct
2. **Difficult mocking**: Cannot easily mock without wrappers
3. **Integration tests**: Require Kinesis Local or LocalStack

#### Test Plan

Integration tests will be implemented in **Module 9 - MiniStack** when we configure:
- Kinesis Local for testing
- End-to-end tests with real Kinesis
- Validation of producer and consumer

### System Workflow

#### 1. Event Publishing
```text
Ecommerce API
    |
    v
POST /events
    |
    v
Process Event Command
    |
    v
Kinesis Producer
    |
    v
Kinesis Stream (recommendation-events)
```

#### 2. Event Consumption
```text
Kinesis Stream
    |
    v
Kinesis Consumer
    |
    v
Worker Pool (Module 7)
    |
    v
Recommendation Service
    |
    v
DynamoDB
```

### Trade-offs Considered

1. **Partition Key Strategy**
   - **Choice**: UserID as partition key
   - **Trade-off**: User ordering vs uniform distribution
   - **Decision**: User ordering is more important for the use case

2. **Single Record vs Batch**
   - **Choice**: Both methods available
   - **Trade-off**: Simplicity vs efficiency
   - **Decision**: Batch for high throughput, single for low latency

3. **JSON vs Binary Serialization**
   - **Choice**: JSON for Kinesis
   - **Trade-off**: Size vs readability
   - **Decision**: JSON is more readable and sufficient for this use case

4. **Unit vs Integration Tests**
   - **Choice**: Integration tests in Module 9
   - **Trade-off**: Early coverage vs real infrastructure
   - **Decision**: Integration tests with Kinesis Local are more valuable

### Validations Performed

- [x] Code compiles without errors
- [x] All existing tests pass
- [x] JSON tags in Event entity
- [x] Producer implemented (single and batch)
- [x] Consumer implemented (GetRecords, GetShardIterator, ListShards)
- [x] Kinesis client for AWS and local
- [x] Documentation of decisions

### Module 6 Status

- [x] Closed (integration tests pending in Module 9)
