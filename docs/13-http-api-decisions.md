# HTTP API - Implementation Decisions

## Module 4 — HTTP API with Gin

### Objective

Expose use cases through a REST API using Gin.

### Implemented Endpoints

#### 1. GET /health
- **Purpose**: Service health check
- **Response**: `{"status": "ok"}`
- **Status Code**: 200 OK

#### 2. GET /recommendations/:userId
- **Purpose**: Get recommendations for a user
- **Parameters**: `userId` (path parameter)
- **Response**:
```json
{
  "recommendations": [
    {
      "user_id": "123",
      "product_id": "P10",
      "score": 0.95
    }
  ]
}
```
- **Status Codes**:
  - 200 OK: Recommendations found
  - 400 Bad Request: Empty userId
  - 500 Internal Server Error: Service error

#### 3. POST /events
- **Purpose**: Process an interaction event
- **Request Body**:
```json
{
  "event_id": "event-1",
  "event_type": "product_viewed",
  "user_id": "123",
  "product_id": "P10",
  "product_category": "electronics",
  "product_brand": "Apple",
  "metadata": {
    "device": "mobile",
    "country": "US"
  },
  "occurred_at": "2026-08-08T22:47:17Z"
}
```
- **Response**: `{"status": "event processed"}`
- **Status Codes**:
  - 202 Accepted: Event processed
  - 400 Bad Request: Invalid JSON or invalid data
  - 500 Internal Server Error: Service error

#### 4. POST /recommendations/generate
- **Purpose**: Generate recommendations from a batch of events
- **Request Body**:
```json
{
  "user_id": "123",
  "events": [
    {
      "event_id": "event-1",
      "event_type": "product_viewed",
      "user_id": "123",
      "product_id": "P10",
      "occurred_at": "2026-08-08T22:47:17Z"
    }
  ],
  "limit": 10
}
```
- **Response**:
```json
{
  "recommendations": [
    {
      "user_id": "123",
      "product_id": "P10",
      "score": 1.0
    }
  ]
}
```
- **Status Codes**:
  - 200 OK: Recommendations generated
  - 400 Bad Request: Invalid JSON or invalid data
  - 500 Internal Server Error: Service error

### Handler Structure

#### Handler Struct
```go
type Handler struct {
    getRecommendationsHandler    *queries.GetRecommendationsHandler
    processEventHandler         *commands.ProcessEventHandler
    generateRecommendationsHandler *commands.GenerateRecommendationsHandler
}
```

**Decisions**:
- Handler contains references to application handlers
- Constructor injection for testability
- No business logic in HTTP handlers

### HTTP Validations

#### GetRecommendations
- Validates that `userId` is not empty
- Returns 400 if parameter is missing

#### ProcessEvent
- Validates JSON body with `ShouldBindJSON`
- Validates timestamp in `ToEventDomain`
- Returns 400 if JSON is invalid

#### GenerateRecommendations
- Validates JSON body with `ShouldBindJSON`
- Validates that `user_id`, `events`, and `limit` are present
- Validates each event in `ToEventDomain`
- Returns 400 if JSON is invalid

### Conversions

#### Request → DTO → Domain
- HTTP requests are converted to DTOs
- DTOs are converted to domain entities
- Conversion handles timestamps (RFC3339 string → time.Time)

#### Domain → DTO → Response
- Domain entities are converted to DTOs
- DTOs are serialized to JSON
- Conversion uses helper functions in `dto` package

### Error Handling

#### Domain Errors
- Domain errors propagate without wrapping
- Returned as 500 Internal Server Error with message
- Validation errors converted to 400 Bad Request

#### HTTP Errors
- JSON parsing errors → 400 Bad Request
- Binding errors → 400 Bad Request
- Internal errors → 500 Internal Server Error

### Memory Repository

#### Implementation
- **Location**: `internal/infrastructure/persistence/memory/repository.go`
- **Purpose**: In-memory repository for local development
- **Features**:
  - Thread-safe with `sync.RWMutex`
  - Map of `userID → []Recommendation`
  - In-place update if recommendation exists
  - New creation if doesn't exist

**Decisions**:
- Used for local development and testing
- Replaceable by DynamoDB in production
- Doesn't persist between server restarts

### Composition Root

#### main.go
The `main.go` acts as composition root:

```go
repository := memory.NewMemoryRepository()
service := recommendation.NewService(repository)

getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
processEventHandler := commands.NewProcessEventHandler(service)
generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

handler := httpgin.NewHandler(
    getRecommendationsHandler,
    processEventHandler,
    generateRecommendationsHandler,
)

router := httpgin.NewRouter(handler)
server := httpgin.NewServer(config.Port, router)
```

**Decisions**:
- Manual dependency injection
- No DI framework
- Easy to test with fakes/mocks
- Easy to switch to DynamoDB in the future

### HTTP Tests

#### TestHealth
- Verifies that `/health` returns 200
- Verifies body `{"status":"ok"}`
- Verifies header `X-Request-ID`

#### TestHealthWithRequestID
- Verifies that `X-Request-ID` propagates
- Sends a header and verifies the same is returned

#### TestGetRecommendations
- Creates a recommendation in repository
- Calls GET /recommendations/:userId endpoint
- Verifies correct recommendation is returned

#### TestProcessEvent
- Calls POST /events endpoint
- Verifies recommendation is created in repository
- Verifies status code 202 Accepted

#### TestGenerateRecommendations
- Calls POST /recommendations/generate endpoint
- Verifies correct recommendations are generated
- Verifies descending score order

### Middleware

#### Request ID
- Generates UUID if not present
- Propagates header to response
- Useful for tracing and debugging

#### Logger
- Gin's built-in logger
- Logs all requests
- Format: `[GIN] method path status latency`

#### Recovery
- Gin's built-in recovery
- Catches panics and returns 500
- Prevents server from crashing

### Validations Performed

- [x] `internal/infrastructure/http/gin` depends on Gin (intentional)
- [x] `internal/infrastructure/http/gin` doesn't depend on AWS SDK
- [x] Handlers don't contain business logic
- [x] All HTTP tests pass
- [x] Tests with race detector pass
- [x] Go vet reports no errors
- [x] Formatted with gofmt
- [x] Server starts correctly
- [x] All endpoints respond correctly

### Trade-offs Considered

1. **Memory Repository vs DynamoDB**: Used memory repository for local development. Facilitates testing and rapid development. Will be replaced by DynamoDB in Module 5.

2. **Handler Composition vs Single Handler**: Used a Handler struct with multiple methods instead of separate handlers. Facilitates dependency injection and maintains consistency.

3. **Simple vs Structured Error Responses**: Used simple responses with error messages. Future improvements could include structured error codes for better debugging.

4. **Context Propagation**: Uses `c.Request.Context()` to propagate context to application handlers. Enables cancellation and timeouts in the future.

### Module 4 Status

- [x] Closed
