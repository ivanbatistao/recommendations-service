# DynamoDB - Implementation Decisions

## Module 5 — Persistence Layer with DynamoDB

### Objective

Implement real persistence using DynamoDB for production and local development with MiniStack.

### DynamoDB Table Design

#### Table Schema

```text
Table Name: Recommendations
Partition Key: UserID (String)
Sort Key: ProductID (String)
Attributes:
  - Score (Number)
```

#### Partition Key Decision

**Choice:** UserID as Partition Key, ProductID as Sort Key

**Reasons:**
1. **Primary use case**: `GET /recommendations/:userId` - queries by user
2. **Performance**: Query by partition key is O(1)
3. **Simplicity**: Direct implementation of Repository interface
4. **Cost**: No GSI, more economical

**Trade-offs:**
- **Advantage**: User queries are very fast
- **Disadvantage**: A user with thousands of products may have a large partition
- **Mitigation**: DynamoDB handles large partitions with automatic pagination

### Repository Implementation

#### DynamoDBRepository

**Location**: `internal/infrastructure/persistence/dynamodb/repository.go`

**Structure:**
```go
type DynamoDBRepository struct {
    client    *dynamodb.Client
    tableName string
}
```

**Methods:**

##### 1. GetByUserID
```go
func (r *DynamoDBRepository) GetByUserID(
    ctx context.Context,
    userID string,
) ([]recommendation.Recommendation, error)
```

**Implementation:**
- Uses `Query` operation with KeyConditionExpression
- Expression: `UserID = :uid`
- Converts AttributeValues to Go structs with `UnmarshalMap`

**SQL Analogy:**
```sql
SELECT * FROM Recommendations WHERE UserID = :uid
```

##### 2. Save
```go
func (r *DynamoDBRepository) Save(
    ctx context.Context,
    rec recommendation.Recommendation,
) error
```

**Implementation:**
- Uses `PutItem` operation
- Converts Go struct to AttributeValues with `MarshalMap`
- **Idempotent**: If exists (same UserID + ProductID), replaces it

**Behavior:**
- **Create**: Item doesn't exist → creates new
- **Update**: Item exists → complete replacement
- **Upsert**: Automatic upsert operation

### DynamoDB ↔ Go Mapping

#### dynamodbav Tags

```go
type Recommendation struct {
    UserID    string `dynamodbav:"UserID"`
    ProductID string `dynamodbav:"ProductID"`
    Score     float64 `dynamodbav:"Score"`
}
```

**Purpose:**
- Explicit control of mapping between Go structs and DynamoDB attributes
- Avoids naming inconsistencies
- Allows attribute names different from Go field names

**Without tags vs With tags:**
- **Without tags**: SDK attempts automatic mapping (error-prone)
- **With tags**: Explicit and robust mapping

### DynamoDB Client Configuration

#### NewDynamoDBClient (AWS Real)

```go
func NewDynamoDBClient(ctx context.Context, region string) (*dynamodb.Client, error)
```

**Features:**
- Uses `config.LoadDefaultConfig` to automatically load credentials
- Looks for credentials in:
  1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  2. ~/.aws/credentials file
  3. IAM roles (if in EC2/Lambda)
- **Usage**: Production in AWS

#### NewLocalDynamoDBClient (Local Development)

```go
func NewLocalDynamoDBClient(ctx context.Context, endpoint string) (*dynamodb.Client, error)
```

**Features:**
- Override endpoint with `BaseEndpoint`
- Points to local: `http://localhost:8000`
- **Usage**: Development with MiniStack or DynamoDB Local

### Environment Variables

#### Configuration

```go
type Config struct {
    Port              string
    UseDynamoDB       bool
    DynamoDBTable     string
    DynamoDBEndpoint  string
    AWSRegion         string
}
```

#### Environment Variables

| Variable            | Default           | Description                              |
| ------------------- | ----------------- | ---------------------------------------- |
| `PORT`              | `8080`            | HTTP server port                         |
| `USE_DYNAMODB`      | `false`           | Enable DynamoDB                          |
| `DYNAMODB_TABLE`    | `Recommendations`  | DynamoDB table name                      |
| `DYNAMODB_ENDPOINT` | `""`              | Endpoint for DynamoDB Local              |
| `AWS_REGION`        | `us-east-1`       | AWS region                               |

#### Usage Examples

**Development with Memory Repository:**
```bash
# No additional environment variables
go run cmd/api/main.go
```

**Development with DynamoDB Local:**
```bash
export USE_DYNAMODB=true
export DYNAMODB_ENDPOINT=http://localhost:8000
go run cmd/api/main.go
```

**Production in AWS:**
```bash
export USE_DYNAMODB=true
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
go run cmd/api/main.go
```

### Dynamic Repository Selection

#### Logic in main.go

```go
if config.UseDynamoDB {
    // Configure DynamoDB client
    if config.DynamoDBEndpoint != "" {
        client = dynamodb.NewLocalDynamoDBClient(...)
    } else {
        client = dynamodb.NewDynamoDBClient(...)
    }
    repository = dynamodb.NewDynamoDBRepository(client, config.DynamoDBTable)
} else {
    repository = memory.NewMemoryRepository()
}
```

**Benefits:**
- **Rapid development**: Memory repository by default
- **Local testing**: DynamoDB Local with endpoint override
- **Production**: DynamoDB with AWS credentials
- **Same code**: Transparent change without modifying business logic

### Installed AWS SDK Packages

1. **`github.com/aws/aws-sdk-go-v2`**: Main AWS SDK
2. **`github.com/aws/aws-sdk-go-v2/config`**: Credential and region configuration
3. **`github.com/aws/aws-sdk-go-v2/service/dynamodb`**: DynamoDB specific client
4. **`github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue`**: Conversion between Go and DynamoDB

### Tests

#### Current Status

DynamoDB tests are in `Skip` because:

1. **SDK doesn't use interfaces**: DynamoDB client is a concrete struct, not an interface
2. **Difficult mocking**: Cannot easily mock without wrappers
3. **Integration tests**: Require DynamoDB Local or MiniStack

#### Test Plan

Integration tests will be implemented in **Module 9 - MiniStack** when we configure:
- DynamoDB Local for testing
- End-to-end tests with real DynamoDB
- Schema and operations validation

### Trade-offs Considered

1. **Partition Key Strategy**
   - **Choice**: UserID as PK
   - **Trade-off**: Fast queries vs large partitions
   - **Decision**: User queries are the primary use case

2. **Memory vs DynamoDB Repository**
   - **Choice**: Both with dynamic selection
   - **Trade-off**: Complexity vs flexibility
   - **Decision**: Flexibility for different environments

3. **dynamodbav Tags**
   - **Choice**: Explicit tags on all fields
   - **Trade-off**: Verbosity vs robustness
   - **Decision**: Robustness and explicit control

4. **Unit vs Integration Tests**
   - **Choice**: Integration tests in Module 9
   - **Trade-off**: Early coverage vs real infrastructure
   - **Decision**: Integration tests with DynamoDB Local are more valuable

### Validations Performed

- [x] Code compiles without errors
- [x] All existing tests pass
- [x] Environment variable configuration
- [x] Dynamic repository selection
- [x] DynamoDB ↔ Go mapping with tags
- [x] DynamoDB client for AWS and local
- [x] Documentation of decisions

### Module 5 Status

- [x] Closed (integration tests pending in Module 9)
