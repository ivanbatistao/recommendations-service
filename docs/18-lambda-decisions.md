# AWS Lambda - Implementation Decisions

## Module 8 — AWS Lambda Integration

### Objective

Adapt the code to run on AWS Lambda as a serverless function.

### What is AWS Lambda?

**AWS Lambda** is a serverless service that:

- **Executes code** without server provisioning
- **Scales automatically** based on traffic
- **Charges for usage** (execution milliseconds + memory)
- **Integrates** with other AWS services (API Gateway, Kinesis, etc.)

### Architecture with Lambda

```text
API Gateway
      |
      v
AWS Lambda
      |
  Application Layer
      |
  Domain Layer
      |
  DynamoDB
```

### Lambda Architecture Decision

#### Not Using Gin in Lambda

**Decision:** Do not use Gin HTTP framework directly in Lambda

**Reasons:**
1. **Simplicity**: Gin is designed for persistent HTTP servers
2. **Performance**: Lambda has short timeouts, Gin adds overhead
3. **Maintenance**: Two different codebases (HTTP server vs Lambda)
4. **Testing**: Easier to test direct handlers than Gin adapters

**Alternative:**
- Reuse application handlers directly
- Convert `APIGatewayProxyRequest` → Commands/Queries
- Convert Responses → `APIGatewayProxyResponse`

### Lambda Handler Implementation

#### LambdaHandler

**Location**: `cmd/lambda/main.go`

**Structure:**
```go
type LambdaHandler struct {
    getRecommendationsHandler    *queries.GetRecommendationsHandler
    processEventHandler         *commands.ProcessEventHandler
    generateRecommendationsHandler *commands.GenerateRecommendationsHandler
    logger                       *slog.Logger
}
```

**Components:**
- Reuses existing application handlers
- Doesn't use Gin directly
- Structured logging for CloudWatch

#### HandleRequest - API Gateway Router

```go
func (h *LambdaHandler) HandleRequest(
    ctx context.Context,
    req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error)
```

**API Gateway → Handlers conversion:**
- API Gateway sends `APIGatewayProxyRequest`
- Manual routing by path and method
- Calls application handlers
- Returns `APIGatewayProxyResponse`

**Supported routes:**
- `GET /health` → Health check
- `GET /recommendations/{userId}` → Get recommendations
- `POST /events` → Process event
- `POST /recommendations/generate` → Generate recommendations

#### Handler Methods

##### 1. handleHealth
```go
func (h *LambdaHandler) handleHealth(ctx context.Context) (events.APIGatewayProxyResponse, error)
```

**Implementation:**
- Returns simple `{"status":"ok"}`
- Useful for API Gateway health checks

##### 2. handleGetRecommendations
```go
func (h *LambdaHandler) handleGetRecommendations(ctx context.Context, userID string) (events.APIGatewayProxyResponse, error)
```

**Implementation:**
- Extracts userID from path parameters
- Calls `GetRecommendationsHandler`
- Converts to DTOs and JSON
- Returns 200 with recommendations or 400/500 on error

##### 3. handleProcessEvent
```go
func (h *LambdaHandler) handleProcessEvent(ctx context.Context, body string) (events.APIGatewayProxyResponse, error)
```

**Implementation:**
- Parses JSON body to EventDTO
- Converts to domain Event
- Calls `ProcessEventHandler`
- Returns 202 Accepted or 400/500 on error

##### 4. handleGenerateRecommendations
```go
func (h *LambdaHandler) handleGenerateRecommendations(ctx context.Context, body string) (events.APIGatewayProxyResponse, error)
```

**Implementation:**
- Parses JSON body with userID, events, limit
- Converts EventDTOs to domain Events
- Calls `GenerateRecommendationsHandler`
- Converts to DTOs and JSON
- Returns 200 with recommendations or 400/500 on error

### Initialization in Lambda

#### NewLambdaHandler

```go
func NewLambdaHandler() *LambdaHandler
```

**Similar to main.go:**
- Loads configuration from environment variables
- Selects repository (DynamoDB or memory)
- Creates application handlers
- Returns initialized LambdaHandler

**Differences from main.go:**
- Doesn't create HTTP server
- Doesn't handle graceful shutdown
- Doesn't listen on port

### Dockerfile for Lambda

**Location**: `Dockerfile.lambda`

```dockerfile
FROM golang:1.26.5 AS builder
# ... build steps ...

FROM gcr.io/distroless/static-debian12
COPY --from=builder /bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"]
```

**Features:**
- Multi-stage build (compilation + runtime)
- Standalone binary: `/bootstrap`
- Distroless base for security
- Compatible with AWS Lambda custom runtimes

### Lambda Environment Variables

#### AWS Lambda Configuration

| Variable            | Default           | Description                              |
| ------------------- | ----------------- | ---------------------------------------- |
| `USE_DYNAMODB`      | `true`            | Use DynamoDB in Lambda                  |
| `DYNAMODB_TABLE`    | `Recommendations`  | DynamoDB table name                      |
| `AWS_REGION`        | From Lambda       | AWS region (configured by Lambda)       |
| `GIN_MODE`          | N/A               | Not applicable in Lambda                 |

#### Local Configuration for Testing

```bash
USE_DYNAMODB=false  # Use memory repository for local tests
```

### Implemented Tests

#### 1. TestLambdaHandler_HandleHealth
- Verifies simple health check
- Verifies status 200 and correct body

#### 2. TestLambdaHandler_HandleGetRecommendations
- Verifies GET /recommendations/{userId}
- Verifies path parameter parsing
- Verifies response structure

#### 3. TestLambdaHandler_HandleProcessEvent
- Verifies POST /events
- Verifies JSON body parsing
- Verifies conversion to domain Event
- Verifies status 202 Accepted

#### 4. TestLambdaHandler_HandleGenerateRecommendations
- Verifies POST /recommendations/generate
- Verifies complex request parsing
- Verifies multiple event conversion
- Verifies response structure

#### 5. TestLambdaHandler_HandleNotFound
- Verifies non-existent routes
- Verifies status 404

### Trade-offs Considered

1. **Gin vs Direct Handlers**
   - **Choice**: Direct handlers in Lambda
   - **Trade-off**: Code sharing vs simplicity
   - **Decision**: Application handlers are shared, only adapter is different

2. **Custom Runtime vs SAM/Serverless Framework**
   - **Choice**: Custom runtime with Docker
   - **Trade-off**: Control vs convenience
   - **Decision**: Custom runtime gives more control and flexibility

3. **Single vs Multiple Lambdas**
   - **Choice**: Single Lambda with manual routing
   - **Trade-off**: Simplicity vs specialization
   - **Decision**: Single lambda is simpler for this use case

4. **Memory vs DynamoDB in Lambda**
   - **Choice**: DynamoDB by default in Lambda
   - **Trade-off**: Persistence vs simplicity
   - **Decision**: Lambda is for production, use real DynamoDB

### Validations Performed

- [x] Code compiles without errors
- [x] Lambda handler tests pass
- [x] Dockerfile for Lambda created
- [x] Reuse of application handlers
- [x] API Gateway ↔ Domain conversion
- [x] Appropriate error handling
- [x] Structured logging for CloudWatch
- [x] Documentation of decisions

### Module 8 Status

- [x] Closed (deployment would require additional AWS configuration)
