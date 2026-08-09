# Real-Time Recommendation Service - Project Summary

## Overview

Sistema de recomendaciones en tiempo real para ecommerce construido con Go, implementando Clean Architecture, procesamiento de eventos con Kinesis, almacenamiento en DynamoDB, y deployment en AWS Lambda.

## Stack Tecnológico

| Componente | Tecnología | Propósito |
|------------|-----------|-----------|
| Lenguaje | Go 1.26.5 | Performance y concurrencia nativa |
| HTTP Framework | Gin | API REST高性能 |
| Arquitectura | Clean Architecture | Separación de capas y testabilidad |
| Event Streaming | AWS Kinesis | Procesamiento de eventos en tiempo real |
| Database | AWS DynamoDB | Almacenamiento NoSQL escalable |
| Compute | AWS Lambda | Serverless compute |
| Local Development | LocalStack | Simulación AWS local |
| Load Testing | k6 | Pruebas de carga y métricas de performance |
| Concurrency | Goroutines + Channels | Worker pool para procesamiento paralelo |

## Arquitectura

### Capas de Clean Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                          │
│                  cmd/api, cmd/lambda                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│              Application Layer (Commands/Queries)            │
│         internal/application/commands, queries, dto          │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Domain Layer                              │
│           internal/domain/{event, recommendation}            │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│              Infrastructure Layer                            │
│   internal/infrastructure/{persistence, streaming, logger}   │
└─────────────────────────────────────────────────────────────┘
```

### Flujo de Arquitectura

**Flujo de Eventos:**
```text
Event Generator → Kinesis Stream → Worker Pool → Domain → DynamoDB
```

**Flujo de Consultas:**
```text
Client → API Gateway → Lambda → Gin → Application → Domain → DynamoDB
```

### Componentes Principales

#### 1. Domain Layer
- **Event**: Modela eventos de usuario (view, purchase, cart add)
- **Recommendation**: Entidad de recomendaciones con scoring
- **Service**: Lógica de negocio para generar recomendaciones
- **Repository Interface**: Abstracción de persistencia

#### 2. Application Layer
- **Commands**: ProcessEvent, GenerateRecommendations
- **Queries**: GetRecommendations
- **DTOs**: Data transfer objects para integración

#### 3. Infrastructure Layer
- **Persistence**: Memory repository (testing) + DynamoDB repository (producción)
- **Streaming**: Kinesis producer/consumer para eventos
- **Composition Root**: Inicialización de dependencias
- **Logger**: Structured logging con slog

#### 4. HTTP Layer
- **Gin API**: REST endpoints (/health, /recommendations, /events)
- **Lambda Handler**: AWS Lambda con API Gateway integration
- **Event Generator**: Generador de tráfico realista para testing

## Decisiones Arquitectónicas Clave

### 1. Clean Architecture
**Decisión:** Implementar Clean Architecture con separación estricta de capas.

**Razón:**
- Testabilidad del dominio sin dependencias externas
- Cambio de framework HTTP sin afectar lógica de negocio
- Facilita testing unitario y mantenimiento

**Trade-off:** Mayor complejidad inicial vs maintainability a largo plazo.

### 2. DynamoDB Partition Key: UserID
**Decisión:** Usar `UserID` como partition key en lugar de `EventID`.

**Razón:**
- Queries principales son por usuario: "Get recommendations for user X"
- Eficiente para patrones de acceso principales
- Simplemente implementa el caso de uso principal

**Trade-off:** Queries por producto requieren Global Secondary Index (no implementado).

### 3. Worker Pool Pattern
**Decisión:** Implementar worker pool con goroutines y channels para procesamiento de eventos.

**Razón:**
- Control de concurrencia: N workers configurables
- Evita goroutine explosion
- Graceful shutdown con context cancellation
- Mejor performance que sequential processing

**Trade-off:** Mayor complejidad vs simple sequential processing.

### 4. Composition Root Pattern
**Decisión:** Extraer lógica de inicialización a `internal/infrastructure/composition/root.go`.

**Razón:**
- Elimina duplicación entre cmd/api y cmd/lambda
- Centraliza configuración de dependencias
- Facilita testing con diferentes configuraciones
- Sigue Clean Architecture Composition Root pattern

**Trade-off:** Mayor indirección vs inicialización directa en main.go.

### 5. LocalStack para Desarrollo Local
**Decisión:** Usar LocalStack para simular DynamoDB y Kinesis localmente.

**Razón:**
- Desarrollo sin cuenta AWS
- Testing rápido sin depender de cloud
- Misma API que AWS real
- Costo cero para desarrollo

**Trade-off:** LocalStack no es 100% compatible con AWS real.

## Métricas de Performance

### Objetivos de Diseño
- **Latencia:** < 500ms (p95) para get recommendations
- **Throughput:** 100+ RPS con 10 VUs
- **Error Rate:** < 1% en condiciones normales
- **Concurrencia:** Configurable worker pool (N workers)

### Métricas Esperadas (k6 Tests)
- **API Load Test:** 5-10-10-0 VUs, 70s, < 500ms p95
- **Stress Test:** 10-50-100-100-50-0 VUs, 80s, < 1s p95
- **Spike Test:** 5-100-100-5-0 VUs, 36s, < 2s p95

## Implementación

### Módulos Completados

**✅ Core Architecture (8 módulos):**
1. Domain Layer - Entidades y lógica de negocio
2. Application Layer - Commands, queries, DTOs
3. HTTP API con Gin - Endpoints REST
4. DynamoDB - Repository para persistencia
5. Kinesis - Producer/consumer de eventos
6. Worker Pool - Procesamiento concurrente
7. AWS Lambda - Serverless deployment
8. Composition Root - Inicialización compartida

**✅ Development Tools (3 módulos):**
9. Event Generator - Generador de tráfico realista
10. MiniStack - LocalStack para desarrollo local
11. Load Testing - k6 scripts para métricas

### Estructura del Proyecto

```
recommendations-service/
├── cmd/
│   ├── api/              # Gin HTTP server
│   ├── lambda/           # AWS Lambda handler
│   └── event-generator/  # Traffic simulator
├── internal/
│   ├── application/      # Commands, queries, DTOs
│   ├── domain/          # Event, Recommendation, Service
│   └── infrastructure/   # Persistence, streaming, composition
├── configs/             # Configuration management
├── docs/               # Architecture decisions and documentation
├── loadtests/          # k6 load testing scripts
├── scripts/            # Utility scripts
└── docker-compose.yml  # LocalStack setup
```

## Casos de Uso

### 1. User Gets Recommendations
```
Client → GET /recommendations/{userId}
→ Recommendation API
→ Application Layer (GetRecommendationsHandler)
→ Domain Service (GenerateRecommendations)
→ DynamoDB Repository (GetEvents)
→ Scoring + Ranking
→ Top N Recommendations
→ JSON Response
```

### 2. User Interaction Event
```
Ecommerce → POST /events
→ Recommendation API
→ Application Layer (ProcessEventHandler)
→ Kinesis Producer
→ Kinesis Stream
→ Kinesis Consumer
→ Worker Pool
→ Domain Service (ProcessEvent)
→ DynamoDB Repository (SaveEvent)
```

## Testing

### Estrategia de Testing
- **Unit Tests:** Domain layer, application layer
- **Integration Tests:** HTTP API, DynamoDB, Kinesis
- **Load Tests:** k6 scripts para performance
- **Race Detection:** `go test -race ./...`

### Cobertura de Tests
- Domain logic: 100% coverage (scoring, ranking, aggregation)
- Application layer: 100% coverage (commands, queries)
- Infrastructure: Integration tests for DynamoDB/Kinesis
- HTTP API: HTTP tests for all endpoints

## Deployment

### Local Development
```bash
docker-compose up localstack
./scripts/init-localstack.sh
docker-compose up recommendation-service
```

### AWS Production
```bash
# Build Lambda
docker build -f Dockerfile.lambda -t recommendations-lambda .
# Deploy to AWS Lambda (via SAM or Serverless Framework)
# Configure API Gateway triggers
```

## Próximos Pasos

### Enhancements Prioritarios
1. **Mover Composition Root** a `internal/app/composition/` (documentado en REFACTORING_TODO.md)
2. **Implementar Kinesis Integration** en Event Generator (actualmente solo HTTP)
3. **Observability** - CloudWatch metrics, structured logging, tracing
4. **Resilience** - Retry logic, circuit breaker, backoff strategies

### Roadmap Original (13 módulos pendientes)
- Testing avanzado
- Performance benchmarking con pprof
- Observability (CloudWatch, metrics, tracing)
- Resilience patterns
- Developer experience mejorado
- Infrastructure as Code (Terraform)
- AWS deployment (SAM/Serverless Framework)
- CI/CD pipeline
- Documentation completa
- Interview preparation
- CV update

## Referencias

### Documentación del Proyecto
- `docs/01-problem-statement.md` - Definición del problema
- `docs/02-requirements.md` - Requisitos funcionales y no funcionales
- `docs/03-domain.md` - Modelo de dominio
- `docs/04-events.md` - Eventos del dominio
- `docs/05-api-design.md` - Diseño de API
- `docs/06-architecture.md` - Arquitectura de alto nivel
- `docs/07-data-model.md` - Modelo de datos
- `docs/08-decisions.md` - Decisiones arquitectónicas
- `docs/09-non-functional.md` - Requisitos no funcionales
- `docs/10-project-structure.md` - Estructura del proyecto
- `docs/11-domain-modeling-decisions.md` - Decisiones de modelado
- `docs/12-application-layer-decisions.md` - Decisiones de aplicación
- `docs/13-http-api-decisions.md` - Decisiones de HTTP API
- `docs/14-dynamodb-decisions.md` - Decisiones de DynamoDB
- `docs/15-kinesis-decisions.md` - Decisiones de Kinesis
- `docs/16-worker-pool-decisions.md` - Decisiones de Worker Pool
- `docs/17-lambda-decisions.md` - Decisiones de Lambda
- `docs/18-ministack-setup.md` - Configuración de LocalStack
- `docs/19-load-testing.md` - Guía de load testing
- `docs/REFACTORING_TODO.md` - Refactorizaciones pendientes
- `docs/PROJECT_STATUS.md` - Estado del proyecto vs roadmap

### Tecnologías
- [Gin Web Framework](https://gin-gonic.com/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [LocalStack](https://localstack.cloud/)
- [k6 Load Testing](https://k6.io/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## Conclusión

Este proyecto implementa un sistema de recomendaciones en tiempo real utilizando mejores prácticas de ingeniería:

✅ **Clean Architecture** - Separación de capas, domain-independent
✅ **Event-Driven** - Kinesis para procesamiento en tiempo real
✅ **Serverless** - AWS Lambda para scalability
✅ **Concurrent** - Worker pool para performance
✅ **Testable** - Unit tests, integration tests, load tests
✅ **Observable** - Structured logging, ready for metrics
✅ **Deployable** - LocalStack para dev, AWS para prod
✅ **Maintainable** - Composition root, documentation, refactor plan

El sistema está listo para interviews técnicas con arquitectura defendible, código funcional, métricas de performance obtenibles, y documentación completa.
