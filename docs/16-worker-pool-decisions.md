# Worker Pool - Implementation Decisions

## Module 7 — Worker Pool with Go Concurrency

### Objective

Implement concurrent event processing using a worker pool with goroutines and channels.

### What is a Worker Pool?

A **Worker Pool** is a concurrency pattern that:

- **Creates a fixed number of workers** (goroutines)
- **Uses a channel** as a task queue
- **Workers process tasks** in parallel
- **Controls resource usage** by limiting concurrency

### Why Do We Need a Worker Pool?

**Problem:**
- Kinesis can send many events quickly
- If we process each event in a new goroutine, we can:
  - Exhaust memory
  - Overload the database
  - Cause rate limiting

**Solution:**
- Worker pool with fixed number of workers
- Concurrency control
- Efficient processing without overload

### Worker Pool Architecture

```text
Kinesis Consumer
      |
      v
  Event Channel (Buffered)
      |
      v
  Worker Pool (N workers)
      |
  +---+---+---+---+
  |   |   |   |   |  (Goroutines)
  +---+---+---+---+
      |
      v
Recommendation Service
      |
      v
   DynamoDB
```

### Worker Pool Implementation

#### WorkerPool

**Location**: `internal/infrastructure/processing/workerpool/worker_pool.go`

**Structure:**
```go
type WorkerPool struct {
    numWorkers   int                    // Number of workers
    eventChan    chan event.Event       // Event channel (buffered)
    processor    EventProcessor         // Interface for event processing
    wg           sync.WaitGroup         // To wait for workers to finish
    logger       *slog.Logger          // Logger
    ctx          context.Context       // Context for cancellation
    cancel       context.CancelFunc     // Function to cancel
}
```

**Key components:**
- **Buffered channel**: Event queue with limited capacity
- **NumWorkers**: Fixed number of worker goroutines
- **EventProcessor**: Interface to decouple processing logic
- **Context**: For graceful shutdown

#### EventProcessor Interface

```go
type EventProcessor interface {
    ProcessEvent(ctx context.Context, ev event.Event) error
}
```

**Why an interface?**
- **Decoupling**: Worker pool doesn't know how to process events
- **Testability**: We can use mocks in tests
- **Flexibility**: We can easily change the implementation

#### Worker Pool Methods

##### 1. Start - Start Workers
```go
func (wp *WorkerPool) Start() {
    for i := 0; i < wp.numWorkers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}
```

**What does it do?**
- Creates `numWorkers` goroutines
- Each goroutine executes `worker(i)`
- Uses `WaitGroup` to track active workers

##### 2. Worker - Individual Goroutine
```go
func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case <-wp.ctx.Done():
            return  // Close worker
        case ev := <-wp.eventChan:
            wp.processor.ProcessEvent(wp.ctx, ev)  // Process event
        }
    }
}
```

**Classic worker pattern:**
- **Infinite loop**: Worker is always listening
- **Select**: Waits for either context canceled or new event
- **Context.Done()**: Graceful shutdown
- **Channel receive**: Processes event when it arrives

##### 3. Submit - Send Events
```go
func (wp *WorkerPool) Submit(ev event.Event) {
    select {
    case wp.eventChan <- ev:
        // Event sent successfully
    default:
        // Channel full, drop event
    }
}
```

**Why `select` with `default`?**
- **Non-blocking**: Doesn't wait if channel is full
- **Backpressure**: If channel is full, drops events
- **Trade-off**: Better to drop events than block the producer

**Alternative:** Could use blocking send, but that could saturate the producer.

##### 4. Stop - Graceful Shutdown
```go
func (wp *WorkerPool) Stop() {
    wp.cancel()           // Cancel context
    close(wp.eventChan)   // Close channel
    wp.wg.Wait()         // Wait for workers to finish
}
```

**Graceful shutdown:**
1. **Cancel context**: Workers receive stop signal
2. **Close channel**: Unblocks workers waiting on channel
3. **WaitGroup**: Waits for all workers to finish

##### 5. Stats - Pool Metrics
```go
func (wp *WorkerPool) Stats() map[string]interface{} {
    return map[string]interface{}{
        "num_workers":  wp.numWorkers,
        "buffer_size":  cap(wp.eventChan),
        "buffer_used":  len(wp.eventChan),
    }
}
```

**Available metrics:**
- **num_workers**: Number of active workers
- **buffer_size**: Channel capacity
- **buffer_used**: Events currently in channel

### RecommendationProcessorAdapter

**Location**: `internal/infrastructure/processing/workerpool/processor_adapter.go`

**Purpose:**
- Connect the `recommendation.Service` with the `WorkerPool`
- Implement the `EventProcessor` interface

**Implementation:**
```go
type RecommendationProcessorAdapter struct {
    service *recommendation.Service
}

func (a *RecommendationProcessorAdapter) ProcessEvent(
    ctx context.Context,
    ev event.Event,
) error {
    return a.service.ProcessEvent(ctx, ev)
}
```

**Why an adapter?**
- **Interface segregation**: Worker pool only knows `EventProcessor`
- **Single responsibility**: Adapter connects two layers
- **Testability**: We can easily mock `EventProcessor`

### Worker Pool Configuration

#### Configurable Parameters

```go
NewWorkerPool(
    numWorkers int,           // Number of workers (goroutines)
    bufferSize int,            // Channel buffer size
    processor EventProcessor,  // Event processor
    logger *slog.Logger,       // Logger
)
```

**Recommendations:**
- **numWorkers**: Based on number of CPUs or DB limits
- **bufferSize**: 2-3x numWorkers to avoid blocking
- **Example**: 4 workers, buffer of 10-20 events

### Implemented Tests

#### 1. TestWorkerPool_StartAndStop
- Verifies pool starts and stops workers correctly
- Verifies graceful shutdown

#### 2. TestWorkerPool_SubmitAndProcess
- Verifies events are processed correctly
- Verifies all events are processed

#### 3. TestWorkerPool_ConcurrentProcessing
- Verifies concurrent processing with multiple workers
- Verifies more events than workers are processed correctly

#### 4. TestWorkerPool_Stats
- Verifies statistics are correct
- Verifies num_workers and buffer_size

#### 5. TestWorkerPool_Backpressure
- Verifies behavior when channel is full
- Verifies events are dropped when no capacity

### Trade-offs Considered

1. **Non-blocking vs Blocking Submit**
   - **Choice**: Non-blocking with backpressure
   - **Trade-off**: Event loss vs producer saturation
   - **Decision**: Better to drop events than block the producer

2. **Fixed vs Dynamic Workers**
   - **Choice**: Fixed number of workers
   - **Trade-off**: Simplicity vs auto-scaling
   - **Decision**: Fixed number is simpler and more predictable

3. **Buffer Size**
   - **Choice**: 2-3x numWorkers
   - **Trade-off**: Memory vs throughput
   - **Decision**: Reasonable balance for most cases

4. **Graceful vs Immediate Shutdown**
   - **Choice**: Graceful shutdown with context
   - **Trade-off**: Latency vs completeness
   - **Decision**: Completeness is more important for events

### Validations Performed

- [x] Code compiles without errors
- [x] All worker pool tests pass
- [x] Concurrency tests pass
- [x] Mock processor for testing
- [x] Graceful shutdown implemented
- [x] Backpressure handling implemented
- [x] Stats function implemented
- [x] Documentation of decisions

### Module 7 Status

- [x] Closed
