```
recommendation-service/
│
├── cmd/
│   ├── api/                    # Gin HTTP server
│   ├── lambda/                 # AWS Lambda handler
│   └── event-generator/        # Traffic simulator
│
├── internal/
│   ├── domain/                 # Domain layer (business logic)
│   │   ├── recommendation/     # Recommendation entity, service, errors
│   │   ├── event/              # Event entity, types, errors
│   │   └── product/            # Product entity
│   │
│   ├── application/            # Application layer (use cases)
│   │   ├── commands/           # Write operations (ProcessEventHandler, GenerateRecommendationsHandler)
│   │   ├── queries/            # Read operations (GetRecommendationsHandler)
│   │   └── dto/                # Data transfer objects (RecommendationDTO, EventDTO)
│   │
│   └── infrastructure/         # Infrastructure layer (external services)
│       ├── http/               # HTTP layer
│       │   └── gin/            # Gin handlers and routing
│       ├── persistence/        # Data persistence
│       │   ├── memory/         # In-memory repository
│       │   └── dynamodb/       # DynamoDB repository
│       ├── streaming/          # Event streaming
│       │   └── kinesis/        # Kinesis producer/consumer
│       ├── processing/         # Event processing
│       │   └── workerpool/     # Worker pool implementation
│       ├── composition/        # Dependency injection (Composition Root)
│       └── logger/             # Structured logging
│
├── configs/                    # Configuration management
├── deploy/                     # Deployment configurations
├── scripts/                    # Utility scripts (init-ministack.sh, test-ministack.sh, run-load-tests.sh)
├── loadtests/                  # k6 load testing scripts
├── benchmarks/                 # Performance benchmarks
├── docs/                       # Project documentation
└── go.mod                      # Go module definition
```

## Directory Structure Explanation

### cmd/
Entry points for different execution contexts:
- `api/`: HTTP server with Gin framework
- `lambda/`: AWS Lambda function handler
- `event-generator/`: Traffic generation tool for testing

### internal/domain/
Core business logic, independent of frameworks and infrastructure:
- `recommendation/`: Recommendation entity, service logic, errors
- `event/`: Event entity, event types, validation
- `product/`: Product entity and metadata

### internal/application/
Application use cases orchestrating the domain:
- `commands/`: Write operations (CQRS pattern)
- `queries/`: Read operations (CQRS pattern)
- `dto/`: Data transfer objects for layer boundaries

### internal/infrastructure/
External service implementations:
- `http/gin/`: HTTP layer with Gin framework
- `persistence/`: Data storage (memory, DynamoDB)
- `streaming/kinesis/`: Event streaming with Kinesis
- `processing/workerpool/`: Concurrent event processing
- `composition/`: Dependency injection (Composition Root)
- `logger/`: Structured logging configuration

### configs/
Configuration management and environment variables

### scripts/
Utility scripts for development and testing:
- `init-ministack.sh`: Initialize MiniStack resources
- `test-ministack.sh`: Test MiniStack connectivity
- `run-load-tests.sh`: Execute k6 load tests

### loadtests/
k6 load testing scripts:
- `api-load-test.js`: Standard API load test
- `stress-test.js`: High concurrency stress test
- `spike-test.js`: Sudden traffic spike test

### docs/
Project documentation and architectural decisions:
- `11-domain-modeling-decisions.md`: Domain layer decisions
- `12-application-layer-decisions.md`: Application layer decisions
- `13-http-api-decisions.md`: HTTP API implementation decisions
- `14-dynamodb-decisions.md`: DynamoDB integration decisions
- `15-kinesis-decisions.md`: Kinesis streaming decisions
- `16-worker-pool-decisions.md`: Worker pool concurrency decisions
- `17-lambda-decisions.md`: AWS Lambda integration decisions
- `18-ministack-setup.md`: MiniStack development setup
- `19-load-testing.md`: Load testing with k6 guide
- `20-project-summary.md`: Complete project summary
- `REFACTORING_TODO.md`: Pending refactoring improvements
- `PROJECT_STATUS.md`: Roadmap vs implementation status
