```mermaid
flowchart TD
    CLIENT["Client / Ecommerce"]
    
    CLIENT --> API["API Gateway / HTTP"]
    CLIENT --> EVENT_GEN["Event Generator"]
    
    API --> LAMBDA["AWS Lambda"]
    API --> GIN["Gin HTTP Server"]
    
    LAMBDA --> APP["Application Layer"]
    GIN --> APP
    
    APP --> COMMANDS["Commands"]
    APP --> QUERIES["Queries"]
    
    COMMANDS --> PROCESS_EVENT["Process Event"]
    COMMANDS --> GEN_REC["Generate Recommendations"]
    QUERIES --> GET_REC["Get Recommendations"]
    
    PROCESS_EVENT --> KINESIS_PROD["Kinesis Producer"]
    GEN_REC --> DOMAIN["Domain Service"]
    GET_REC --> DOMAIN
    
    KINESIS_PROD --> KINESIS["Kinesis Stream"]
    
    KINESIS --> KINESIS_CONS["Kinesis Consumer"]
    KINESIS_CONS --> WORKER_POOL["Worker Pool"]
    
    WORKER_POOL --> DOMAIN
    DOMAIN --> DYNAMODB["DynamoDB"]
    
    EVENT_GEN --> KINESIS_PROD
    
    classDef client fill:#fbbf24,stroke:#b45309,color:#78350f
    classDef api fill:#60a5fa,stroke:#1d4ed8,color:#1e3a8a
    classDef lambda fill:#818cf8,stroke:#4f46e5,color:#312e81
    classDef app fill:#34d399,stroke:#059669,color:#064e3b
    classDef domain fill:#f472b6,stroke:#db2777,color:#831843
    classDef infrastructure fill:#93c5fd,stroke:#0284c7,color:#0c4a6e
    classDef streaming fill:#a78bfa,stroke:#7c3aed,color:#4c1d95
    
    class CLIENT client
    class API,GIN api
    class LAMBDA lambda
    class APP,COMMANDS,QUERIES app
    class DOMAIN,PROCESS_EVENT,GEN_REC,GET_REC domain
    class KINESIS_PROD,KINESIS_CONS,KINESIS infrastructure
    class WORKER_POOL,DYNAMODB streaming
```

## Application Flow Overview

### 1. User Gets Recommendations (Read Path)
```
Client → API Gateway → Lambda/Gin → GetRecommendations Query → Domain Service → DynamoDB → Recommendations
```

### 2. User Interaction Event (Write Path)
```
Ecommerce → Event Generator → Kinesis Producer → Kinesis Stream → Kinesis Consumer → Worker Pool → Process Event Command → Domain Service → DynamoDB
```

### 3. Offline Recommendation Generation
```
Batch Events → Generate Recommendations Command → Domain Service → Recommendations (in-memory)
```

## Component Responsibilities

### HTTP Layer
- **Gin HTTP Server**: REST API endpoints for local development
- **AWS Lambda**: Serverless function for production deployment
- **API Gateway**: AWS API Gateway for Lambda integration

### Application Layer
- **Commands**: Write operations (Process Event, Generate Recommendations)
- **Queries**: Read operations (Get Recommendations)
- **DTOs**: Data transfer objects for layer boundaries

### Domain Layer
- **Recommendation Service**: Core business logic
- **Event Processing**: Scoring, interest calculation, ranking
- **Repository Interface**: Abstraction for data persistence

### Infrastructure Layer
- **Kinesis Producer**: Publishes events to stream
- **Kinesis Consumer**: Reads events from stream
- **Worker Pool**: Concurrent event processing
- **DynamoDB Repository**: Data persistence
- **Memory Repository**: In-memory persistence for testing

### External Services
- **DynamoDB**: NoSQL database for recommendations
- **Kinesis**: Event streaming platform
- **LocalStack**: Local AWS simulation for development
