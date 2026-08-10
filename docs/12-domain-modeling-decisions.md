# Domain Modeling Decisions

## Module 2 — Domain Layer

### Domain Entities

#### 1. Event
- **Location**: `internal/domain/event/entity.go`
- **Purpose**: Represents a user interaction event with products
- **Fields**:
  - `EventID`: Unique event identifier
  - `EventType`: Event type (viewed, search, added_to_cart, purchased)
  - `UserID`: ID of the user who generated the event
  - `ProductID`: ID of the involved product
  - `ProductCategory`: Product category (optional for future filtering)
  - `ProductBrand`: Product brand (optional for future filtering)
  - `Metadata`: Additional information (device, country)
  - `OccurredAt`: Event timestamp

**Decisions**:
- String used for `EventID` for AWS Kinesis compatibility
- `EventType` is an enum type for type safety
- `Metadata` is an optional struct to extend information without changing the main entity

#### 2. Recommendation
- **Location**: `internal/domain/recommendation/entity.go`
- **Purpose**: Represents a calculated recommendation for a user
- **Fields**:
  - `UserID`: ID of the target user
  - `ProductID`: ID of the recommended product
  - `Score`: Calculated relevance score

**Decisions**:
- Score is `float64` to allow sum and average calculations
- No timestamps included because score reflects current state
- Minimalist structure to optimize DynamoDB storage

### Domain Algorithms

#### 1. ScoreEvent
- **Location**: `internal/domain/recommendation/scoring.go`
- **Purpose**: Assign weight to each event type
- **Scores**:
  - `ProductViewed`: 1.0
  - `SearchPerformed`: 2.0
  - `ProductAddedCart`: 3.0
  - `ProductPurchased`: 5.0

**Decisions**:
- Scores are simple constants but configurable in the future
- Weight increases with purchase intent
- Unknown events return 0 to ignore them

#### 2. CalculateInterest
- **Location**: `internal/domain/recommendation/interest.go`
- **Purpose**: Aggregate scores by product for a specific user
- **Logic**:
  - Filters events that don't belong to the user
  - Sums scores by product
  - Ignores events with score 0

**Decisions**:
- Explicit filtering by `userID` to avoid contamination between users
- Returns a map for O(1) lookup in ranking
- Pure function without side effects

#### 3. Rank
- **Location**: `internal/domain/recommendation/ranking.go`
- **Purpose**: Sort products by score and apply limit
- **Logic**:
  - Converts interest map to recommendation slice
  - Sorts by score descending
  - Tie-breaker by `ProductID` ascending (deterministic)
  - Applies N limit (Top N)

**Decisions**:
- Deterministic sorting for reproducible results
- Tie-breaker by alphabetical ID for consistency
- Returns empty slice if not enough products
- Adjusts limit if fewer products than requested

### Domain Errors

#### Recommendation Errors
- **Location**: `internal/domain/recommendation/errors.go`
- **Defined errors**:
  - `ErrInvalidUserID`: Empty or invalid UserID
  - `ErrInvalidProductID`: Empty or invalid ProductID
  - `ErrInvalidEventType`: Unrecognized event type
  - `ErrInvalidLimit`: Limit <= 0
  - `ErrUserNotFound`: User doesn't exist in repository
  - `ErrProductNotFound`: Product doesn't exist in repository
  - `ErrRepositoryError`: Generic repository error
  - `ErrRecommendationNotFound`: Recommendation not found

**Decisions**:
- Simple errors without wrapping for direct comparison
- Separation between validation errors and infrastructure errors
- Specific errors for each use case

#### Event Errors
- **Location**: `internal/domain/event/errors.go`
- **Defined errors**:
  - `ErrInvalidEventID`: Empty or invalid EventID
  - `ErrInvalidUserID`: Empty or invalid UserID
  - `ErrInvalidProductID`: Empty or invalid ProductID
  - `ErrInvalidEventType`: Unrecognized event type
  - `ErrInvalidTimestamp`: Invalid timestamp
  - `ErrInvalidMetadata`: Invalid metadata
  - `ErrEventNotFound`: Event not found

**Decisions**:
- Consistent with recommendation errors
- Prepared for future event validation

### Domain Interfaces

#### Repository
- **Location**: `internal/domain/recommendation/repository.go`
- **Methods**:
  - `GetByUserID(ctx, userID)`: Get recommendations by user
  - `Save(ctx, recommendation)`: Save or update recommendation

**Decisions**:
- Minimalist interface to facilitate multiple implementations
- Uses `context.Context` for cancellation and timeouts
- No deletion methods (recommendations are updated)

### Domain Service

#### Service
- **Location**: `internal/domain/recommendation/service.go`
- **Methods**:
  - `GetByUserID`: Gets recommendations with validation
  - `GenerateRecommendations`: Generates recommendations from events
  - `ProcessEvent`: Processes an event and updates recommendations

**Decisions**:
- `GenerateRecommendations` is functional (doesn't depend on repository)
- `ProcessEvent` is imperative (uses repository for persistence)
- Parameter validation at service level
- `ProcessEvent` ignores events with score 0

### Trade-offs Considered

1. **Simple score vs complex model**: Chose simple additive score for simplicity and speed. Future improvements could include temporal decay or contextual factors.

2. **Explicit vs implicit UserID**: Pass `userID` explicitly in `GenerateRecommendations` to avoid ambiguity and enable batch processing.

3. **Deterministic vs random ranking**: Uses `ProductID` tie-breaker for reproducible results in tests and debugging.

4. **Simple errors vs error wrapping**: Uses simple errors without wrapping for direct comparison (`errors.Is`). Future complex errors could use `fmt.Errorf` with wrapping.

### Validations Performed

- [x] `internal/domain` doesn't depend on Gin
- [x] `internal/domain` doesn't depend on AWS SDK
- [x] `internal/domain` doesn't depend on infrastructure
- [x] All unit tests pass
- [x] Tests with race detector pass
- [x] Go vet reports no errors
- [x] Formatted with gofmt

### Module 2 Status

- [x] Closed
