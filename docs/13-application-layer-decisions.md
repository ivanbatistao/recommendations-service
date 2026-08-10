# Application Layer - Implementation Decisions

## Module 3 — Application Layer

### Objective

Implement application use cases by orchestrating the domain without depending on concrete infrastructure.

### Architecture

The application layer follows the CQRS (Command Query Responsibility Segregation) pattern:

```
Commands (Write)
├── ProcessEventCommand
└── GenerateRecommendationsCommand

Queries (Read)
└── GetRecommendationsQuery

DTOs (Data Transfer Objects)
├── RecommendationDTO
└── EventDTO
```

### Implemented Use Cases

#### 1. Get Recommendations (Query)
- **Location**: `internal/application/queries/get_recommendations.go`
- **Handler**: `GetRecommendationsHandler`
- **Purpose**: Get existing recommendations for a user
- **Flow**:
  1. Receives `GetRecommendationsQuery` with `UserID`
  2. Delegates to `recommendation.Service.GetByUserID`
  3. Returns slice of `Recommendation` (domain entities)

**Decisions**:
- Query pattern for read operations
- Doesn't use DTOs in handler (returns domain entities)
- DTO conversion will be done in HTTP layer

#### 2. Process Event (Command)
- **Location**: `internal/application/commands/process_event.go`
- **Handler**: `ProcessEventHandler`
- **Purpose**: Process an interaction event and update recommendations
- **Flow**:
  1. Receives `ProcessEventCommand` with an `Event`
  2. Delegates to `recommendation.Service.ProcessEvent`
  3. Service updates or creates recommendations via Repository

**Decisions**:
- Command pattern for write operations
- Uses domain entities directly
- Validation is done in the domain

#### 3. Generate Recommendations (Command)
- **Location**: `internal/application/commands/generate_recommendations.go`
- **Handler**: `GenerateRecommendationsHandler`
- **Purpose**: Generate recommendations from a batch of events
- **Flow**:
  1. Receives `GenerateRecommendationsCommand` with `UserID`, `Events`, `Limit`
  2. Delegates to `recommendation.Service.GenerateRecommendations`
  3. Returns generated slice of `Recommendation`

**Decisions**:
- Command pattern even though it's a read operation (for consistency)
- Useful for offline batch processing
- Doesn't persist results (in-memory calculation only)

### Implemented DTOs

#### 1. RecommendationDTO
- **Location**: `internal/application/dto/recommendation_dto.go`
- **Purpose**: Transfer recommendation data between layers
- **Fields**:
  - `UserID`: string
  - `ProductID`: string
  - `Score`: float64

**Functions**:
- `FromDomain`: Converts domain entity to DTO
- `ToDomain`: Converts DTO to domain entity
- `FromDomainSlice`: Converts entity slice to DTO slice

**Decisions**:
- Fields identical to domain entity (for now)
- Conversion functions to facilitate testing
- Slice conversion for HTTP response usage

#### 2. EventDTO
- **Location**: `internal/application/dto/event_dto.go`
- **Purpose**: Transfer event data between layers
- **Fields**:
  - `EventID`: string
  - `EventType`: string (not enum for JSON)
  - `UserID`: string
  - `ProductID`: string
  - `ProductCategory`: string (optional)
  - `ProductBrand`: string (optional)
  - `Metadata`: map[string]string (for JSON)
  - `OccurredAt`: string (ISO 8601)

**Functions**:
- `FromEventDomain`: Converts domain entity to DTO
- `ToEventDomain`: Converts DTO to domain entity with timestamp parsing

**Decisions**:
- `EventType` as string for JSON serialization
- `Metadata` as map for flexibility
- `OccurredAt` as RFC3339 string for HTTP transport
- Timestamp parsing in `ToEventDomain` with error handling

### Dependency Injection

All handlers receive the `recommendation.Service` in the constructor:

```go
func NewGetRecommendationsHandler(service *recommendation.Service) *GetRecommendationsHandler
func NewProcessEventHandler(service *recommendation.Service) *ProcessEventHandler
func NewGenerateRecommendationsHandler(service *recommendation.Service) *GenerateRecommendationsHandler
```

**Decisions**:
- Constructor injection for testability
- No DI frameworks (manual injection in composition root)
- Service is shared between handlers (singleton)

### Error Handling

Handlers propagate domain errors without transformation:

```go
func (h *Handler) Execute(ctx context.Context, cmd Command) (Result, error) {
    // Delegates to domain service
    return h.service.SomeMethod(ctx, cmd.Params)
}
```

**Decisions**:
- Domain errors flow without wrapping
- HTTP layer will handle conversion to status codes
- `context.Context` is passed for cancellation and timeouts

### Implemented Tests

#### 1. GetRecommendations Tests
- **Location**: `internal/application/queries/get_recommendations_test.go`
- **Coverage**:
  - Success case with fake repository
  - Verification of sorted results

#### 2. ProcessEvent Tests
- **Location**: `internal/application/commands/process_event_test.go`
- **Coverage**:
  - Update existing recommendation
  - Create new recommendation
  - Ignore events with unknown type
  - Mock repository to verify calls

#### 3. GenerateRecommendations Tests
- **Location**: `internal/application/commands/generate_recommendations_test.go`
- **Coverage**:
  - Success case with test events
  - Validation of empty userID
  - Validation of invalid limit

#### 4. DTO Tests
- **Location**: `internal/application/dto/recommendation_dto_test.go`
- **Coverage**:
  - Domain → DTO conversion
  - DTO → Domain conversion
  - Slice conversion

### Validations Performed

- [x] `internal/application` doesn't depend on Gin
- [x] `internal/application` doesn't depend on AWS SDK
- [x] `internal/application` doesn't depend on concrete infrastructure
- [x] All unit tests pass
- [x] Tests with race detector pass
- [x] Go vet reports no errors
- [x] Formatted with gofmt

### Trade-offs Considered

1. **Commands vs Functions**: Used Command/Handler pattern instead of simple functions for CQRS consistency and future extensibility (middleware, logging, etc.).

2. **DTOs vs Domain Entities**: Command/query handlers return domain entities, not DTOs. DTO conversion is delegated to HTTP layer to keep the application protocol-agnostic.

3. **GenerateRecommendations as Command**: Although it's a read operation, implemented as Command for consistency with offline batch processing use case.

### Module 3 Status

- [x] Closed
