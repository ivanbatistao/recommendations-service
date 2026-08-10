# Real-Time Recommendation Service

Real-time recommendation microservice built with Go, Gin and AWS services.

## Requirements

- Go
- Docker
- Docker Compose

## Run locally

```bash
go run ./cmd/api
```

## Documentation

The project documentation is organized in the `docs/` directory:

- [Problem Statement](docs/00-problem-statement.md) - Overview of the problem and objectives
- [Functional Requirements](docs/01-functional-requirements.md) - Detailed functional requirements
- [Domain](docs/02-domain.md) - Domain analysis and concepts
- [Events](docs/03-events.md) - Event definitions and flows
- [API Design](docs/04-api-design.md) - HTTP API specification
- [Events Storming](docs/05-events-storming.md) - Event storming results
- [Architecture](docs/06-architecture.md) - System architecture overview
- [Data Model](docs/07-data-model.md) - Data model definitions
- [Non-Functional Requirements](docs/08-non-functional-requirements.md) - Performance, security, and scalability requirements
- [Stack](docs/09-stack.md) - Technology stack decisions
- [Project Structure](docs/10-project-structure.md) - Code organization
- [Application Flow](docs/11-app-flow.md) - Application flow diagrams
- [Domain Modeling Decisions](docs/12-domain-modeling-decisions.md) - Domain modeling rationale
- [Application Layer Decisions](docs/13-application-layer-decisions.md) - Application layer design decisions
- [HTTP API Decisions](docs/14-http-api-decisions.md) - HTTP API implementation decisions
- [DynamoDB Decisions](docs/15-dynamodb-decisions.md) - DynamoDB implementation details
- [Kinesis Decisions](docs/16-kinesis-decisions.md) - Kinesis streaming decisions
- [Worker Pool Decisions](docs/17-worker-pool-decisions.md) - Worker pool implementation
- [Lambda Decisions](docs/18-lambda-decisions.md) - AWS Lambda integration
- [Ministack Setup](docs/19-ministack-setup.md) - Local development environment setup
- [Load Testing](docs/20-load-testing.md) - Load testing strategies and results
- [Project Summary](docs/21-project-summary.md) - Overall project summary
- [Stakeholders](docs/22-stakeholders.md) - Stakeholder analysis