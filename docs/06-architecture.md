# Architecture Decision Records (ADR)


## ADR-001 — Database
**Decision:** Use DynamoDB.

**Reason:** Queries are performed by key (userId) and require low latency. DynamoDB provides a key-value model that is well suited to this access pattern and scales horizontally.

## ADR-002 — Asynchronous Communication
**Decision:** Use Kinesis.

**Reason:** It decouples event generation from processing, allows the system to absorb traffic spikes, and makes it possible to scale consumers independently.

## ADR-003 — Go
**Decision:** Implement the service in Go.

**Reason:** Excellent performance, straightforward concurrency using goroutines and channels, lightweight binaries, and startup times suitable for Lambda.

## Diagram

```mermaid
flowchart TD
    A["Event Generator"]:::source
    B["Kinesis Stream"]:::aws
    C["Recommendation Processor"]:::service
    D["Worker Pool<br/>(Concurrent)"]:::worker
    E["DynamoDB"]:::database
    F["Recommendation API"]:::api
    G["Frontend"]:::frontend

    A --> B
    B --> C
    C --> D
    D --> E
    F --> E
    G --> F

    classDef source fill:#e3f2fd,stroke:#1976d2,stroke-width:2px,color:#111
    classDef aws fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#111
    classDef service fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#111
    classDef worker fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#111
    classDef database fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#111
    classDef api fill:#e0f7fa,stroke:#00838f,stroke-width:2px,color:#111
    classDef frontend fill:#fffde7,stroke:#f9a825,stroke-width:2px,color:#111
```
