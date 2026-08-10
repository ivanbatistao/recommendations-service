# Pending Refactorings

This document contains identified but not yet implemented refactorings, prioritized by impact and complexity.

## High Priority 🔴

### 1. Move Composition Root to `internal/app/composition/`

**Status:** ✅ Completed

**Previous location:** `internal/infrastructure/composition/root.go`

**New location:** `internal/app/composition/root.go`

**Reason:**
- Composition root is *application assembly*, not AWS infrastructure
- More semantically correct than "infrastructure"
- Better follows Clean Architecture pattern

**Changes Made:**
- ✅ Moved file to `internal/app/composition/root.go`
- ✅ Updated imports in `cmd/api/main.go`
- ✅ Updated imports in `cmd/lambda/main.go`
- ✅ Removed old `internal/infrastructure/composition/` directory

**Impact:** 
- Location change, no functional change
- Better architectural clarity
- All tests passing after refactoring

**Complexity:** Low (move file and update imports)

---

## Medium Priority 🟡

### 2. Improve Lambda Handler Routing

**Status:** ✅ Kept as-is (Working Well)

**Location:** `cmd/lambda/main.go` - HandleRequest method

**Decision:** 
- Current switch/case routing is simple and functional
- Router pattern would require significant complexity for handler context
- Current approach is maintainable for the current number of routes
- Can be revisited if number of routes grows significantly

**Impact:** No change, maintains simplicity

**Complexity:** Medium (refactor routing structure)

---

### 3. Add HTTP Response Helpers

**Status:** ✅ Completed

**Location:** `cmd/lambda/main.go` - HTTP Response Helpers added

**Changes Made:**
- ✅ Added `jsonResponse()` helper for JSON responses
- ✅ Added `jsonError()` helper for error responses
- ✅ Eliminated code duplication in all handlers
- ✅ Cleaner, more maintainable handler code

**Impact:** Code duplication reduction, better maintainability

**Complexity:** Low (add helpers)

---

## Low Priority 🟢

### 4. Extract Request Validation

**Status:** ✅ Completed

**Location:** `cmd/lambda/main.go` - Validation helpers added

**Changes Made:**
- ✅ Added `validateUserID()` helper
- ✅ Added `validateRequestBody()` helper
- ✅ Updated handlers to use validation helpers
- ✅ Better testability and reusability

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

## Completed Refactorings Summary

**✅ High Priority:**
1. ✅ Move Composition Root to `internal/app/composition/` - Better architectural clarity

**✅ Medium Priority:**
2. ✅ Keep Lambda Handler Routing as-is - Simple and functional for current needs
3. ✅ Add HTTP Response Helpers - Eliminated code duplication
4. ✅ Extract Request Validation - Better testability and reusability
5. ✅ Extract Lambda Helpers and Handlers to Separate Files - Better code organization

**Total Refactorings Completed:** 5 out of 9 prioritized items

## Remaining Refactorings (Optional)

**Medium Priority:**
5. Implement Centralized Error Handling - Could improve error response consistency

**Low Priority:**
6. Add Context Propagation - Better observability and debugging
7. Remove Dead Code - Code cleanup
8. Improve Variable Naming - Better readability

**Pending Decisions:**
9. Use gin-gonic/lambdaadapter - Not investigated

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
