# Pending Refactorings

This document contains identified but not yet implemented refactorings, prioritized by impact and complexity.

## High Priority 🔴

### 1. Move Composition Root to `internal/app/composition/`

**Status:** Identified but not implemented

**Current location:** `internal/infrastructure/composition/root.go`

**Suggested location:** `internal/app/composition/root.go`

**Reason:**
- Composition root is *application assembly*, not AWS infrastructure
- More semantically correct than "infrastructure"
- Better follows Clean Architecture pattern

**Impact:** 
- Location change, no functional change
- Better architectural clarity
- Requires updating imports in cmd/api and cmd/lambda

**Complexity:** Low (move file and update imports)

---

## Medium Priority 🟡

### 2. Improve Lambda Handler Routing

**Status:** Functional but can be improved

**Location:** `cmd/lambda/main.go` - HandleRequest method

**Current problem:**
- Verbose switch/case for routing
- Duplicated HTTP response code
- Not easily scalable

**Suggestion 1 - Router Pattern:**
```go
type Route struct {
    Path    string
    Method  string
    Handler func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

routes := []Route{
    {"/health", "GET", h.handleHealth},
    {"/recommendations/{userId}", "GET", h.handleGetRecommendations},
    // ...
}
```

**Suggestion 2 - Handler Map:**
```go
handlers := map[string]func(ctx context.Context, body string) (events.APIGatewayProxyResponse, error){
    "GET_/health": h.handleHealth,
    "GET_/recommendations/{userId}": h.handleGetRecommendations,
    // ...
}
```

**Impact:** Better maintainability if number of routes grows

**Complexity:** Medium (refactor routing structure)

---

### 3. Add HTTP Response Helpers

**Status:** Duplicated code in Lambda handlers

**Location:** `cmd/lambda/main.go` - multiple methods

**Current problem:**
```go
return events.APIGatewayProxyResponse{
    StatusCode: 200,
    Headers: map[string]string{"Content-Type": "application/json"},
    Body: body,
}, nil
```

**Suggestion:**
```go
func jsonResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {
    data, _ := json.Marshal(body)
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: string(data),
    }, nil
}

func jsonError(statusCode int, message string) events.APIGatewayProxyResponse {
    return jsonResponse(statusCode, map[string]string{"error": message})
}
```

**Impact:** Code duplication reduction

**Complexity:** Low (add helpers)

---

## Low Priority 🟢

### 4. Extract Request Validation

**Status:** Inline validation in each handler

**Location:** `cmd/lambda/main.go` - handleGetRecommendations, handleProcessEvent, etc.

**Current problem:**
```go
if userID == "" {
    return events.APIGatewayProxyResponse{
        StatusCode: 400,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: `{"error":"user_id is required"}`,
    }, nil
}
```

**Suggestion:**
```go
func validateUserID(userID string) error {
    if userID == "" {
        return errors.New("user_id is required")
    }
    return nil
}
```

**Impact:** Better testability and reusability

**Complexity:** Low (extract validation functions)

---

### 5. Implement Centralized Error Handling

**Status:** Dispersed error handling in each handler

**Location:** Multiple handlers in cmd/lambda/main.go

**Current problem:**
```go
if err != nil {
    h.Logger.Error("failed to get recommendations", slog.String("error", err.Error()))
    return events.APIGatewayProxyResponse{
        StatusCode: 500,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: `{"error":"internal server error"}`,
    }, nil
}
```

**Suggestion:**
```go
func (h *LambdaHandler) handleError(err error, context string) events.APIGatewayProxyResponse {
    h.Logger.Error(context, slog.String("error", err.Error()))
    
    var statusCode int
    var message string
    
    switch {
    case errors.Is(err, recommendation.ErrInvalidUserID):
        statusCode = 400
        message = "invalid user_id"
    default:
        statusCode = 500
        message = "internal server error"
    }
    
    return jsonError(statusCode, message)
}
```

**Impact:** Better error response consistency

**Complexity:** Medium (requires error type analysis)

---

### 6. Add Context Propagation

**Status:** Context not fully propagated

**Location:** Various handlers in cmd/lambda/main.go

**Current problem:**
- Request ID not propagated from API Gateway
- No distributed tracing

**Suggestion:**
```go
func (h *LambdaHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Extract request ID from headers
    requestID := req.Headers["X-Request-ID"]
    if requestID == "" {
        requestID = generateRequestID()
    }
    
    ctx = context.WithValue(ctx, "requestID", requestID)
    ctx = h.Logger.With("request_id", requestID)
    
    // Rest of code with enriched ctx
}
```

**Impact:** Better observability and debugging

**Complexity:** Medium (requires adding tracing)

---

## Minor Implementation Details ⚪

### 7. Remove Dead Code

**Status:** Possible unused code after refactorings

**Places to review:**
- Check for unused imports
- Check for unused functions after refactorings
- Clean up tests that no longer apply

**Impact:** Code cleanup

**Complexity:** Low (analysis and cleanup)

---

### 8. Improve Variable Naming

**Status:** Some names could be more descriptive

**Examples:**
- `eventList` could be `events` (already used in domain)
- `app` could be `application` in some contexts

**Impact:** Better readability

**Complexity:** Low (rename variables)

---

## Pending Decisions

### 9. Use gin-gonic/lambdaadapter

**Status:** Not investigated

**Description:** Evaluate using `github.com/awslabs/aws-lambda-go-api-proxy` to integrate Gin with Lambda

**Advantages:**
- Reuse exact Gin code
- No need for manual adapter
- Proven solution

**Disadvantages:**
- Additional dependency
- Possible overhead
- Less control over request/response

**Impact:** Simplify Lambda code

**Complexity:** Medium (evaluate library and refactor)

---

## Prioritization Criteria

**High Priority 🔴:**
- Significant architectural impact
- Improves long-term maintainability
- Low implementation complexity

**Medium Priority 🟡:**
- Significantly improves code quality
- Reduces duplication
- Medium complexity

**Low Priority 🟢:**
- Cosmetic improvements
- Minor optimizations
- Nice-to-have features

---

## Process for Implementing Refactorings

1. **Create branch:** `refactor/refactoring-name`
2. **Implement:** Make the proposed changes
3. **Tests:** Ensure all tests pass
4. **Document:** Update relevant documentation
5. **Commit:** With descriptive message
6. **Merge:** To main after review

---

## Notes

- This document is live - add new refactorings as they are identified
- Before implementing, evaluate if it really adds value
- Consider cost/benefit of each refactoring
- Prioritize functionality over perfection
