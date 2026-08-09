# AWS Lambda - Decisiones de Implementación

## Módulo 8 — AWS Lambda Integration

### Objetivo

Adaptar el código para ejecutar en AWS Lambda como función serverless.

### ¿Qué es AWS Lambda?

**AWS Lambda** es un servicio serverless que:

- **Ejecuta código** sin provisioning de servidores
- **Escala automáticamente** según el tráfico
- **Cobra por uso** (milisegundos de ejecución + memoria)
- **Se integra** con otros servicios AWS (API Gateway, Kinesis, etc.)

### Arquitectura con Lambda

```text
API Gateway
      |
      v
AWS Lambda
      |
  Application Layer
      |
  Domain Layer
      |
  DynamoDB
```

### Decisión de Arquitectura Lambda

#### No usar Gin en Lambda

**Decisión:** No usar Gin HTTP framework directamente en Lambda

**Razones:**
1. **Simplicidad**: Gin está diseñado para servidores HTTP persistentes
2. **Performance**: Lambda tiene timeouts cortos, Gin adds overhead
3. **Mantenimiento**: Dos codebases diferentes (HTTP server vs Lambda)
4. **Testing**: Más fácil testear handlers directos que adaptadores Gin

**Alternativa:**
- Reutilizar handlers de aplicación directamente
- Convertir `APIGatewayProxyRequest` → Commands/Queries
- Convertir Responses → `APIGatewayProxyResponse`

### Implementación del Handler Lambda

#### LambdaHandler

**Ubicación**: `cmd/lambda/main.go`

**Estructura:**
```go
type LambdaHandler struct {
    getRecommendationsHandler    *queries.GetRecommendationsHandler
    processEventHandler         *commands.ProcessEventHandler
    generateRecommendationsHandler *commands.GenerateRecommendationsHandler
    logger                       *slog.Logger
}
```

**Componentes:**
- Reutiliza handlers de aplicación existentes
- No usa Gin directamente
- Logger estructurado para CloudWatch

#### HandleRequest - Router API Gateway

```go
func (h *LambdaHandler) HandleRequest(
    ctx context.Context,
    req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error)
```

**Conversión API Gateway → Handlers:**
- API Gateway envía `APIGatewayProxyRequest`
- Rutea manual por path y method
- Llama a los handlers de aplicación
- Retorna `APIGatewayProxyResponse`

**Rutas soportadas:**
- `GET /health` → Health check
- `GET /recommendations/{userId}` → Get recommendations
- `POST /events` → Process event
- `POST /recommendations/generate` → Generate recommendations

#### Métodos de Handlers

##### 1. handleHealth
```go
func (h *LambdaHandler) handleHealth(ctx context.Context) (events.APIGatewayProxyResponse, error)
```

**Implementación:**
- Retorna `{"status":"ok"}` simple
- Útil para health checks de API Gateway

##### 2. handleGetRecommendations
```go
func (h *LambdaHandler) handleGetRecommendations(ctx context.Context, userID string) (events.APIGatewayProxyResponse, error)
```

**Implementación:**
- Extrae userID de path parameters
- Llama a `GetRecommendationsHandler`
- Convierte a DTOs y JSON
- Retorna 200 con recommendations o 400/500 en error

##### 3. handleProcessEvent
```go
func (h *LambdaHandler) handleProcessEvent(ctx context.Context, body string) (events.APIGatewayProxyResponse, error)
```

**Implementación:**
- Parsea JSON body a EventDTO
- Convierte a domain Event
- Llama a `ProcessEventHandler`
- Retorna 202 Accepted o 400/500 en error

##### 4. handleGenerateRecommendations
```go
func (h *LambdaHandler) handleGenerateRecommendations(ctx context.Context, body string) (events.APIGatewayProxyResponse, error)
```

**Implementación:**
- Parsea JSON body con userID, events, limit
- Convierte EventDTOs a domain Events
- Llama a `GenerateRecommendationsHandler`
- Convierte a DTOs y JSON
- Retorna 200 con recommendations o 400/500 en error

### Inicialización en Lambda

#### NewLambdaHandler

```go
func NewLambdaHandler() *LambdaHandler
```

**Similar a main.go:**
- Carga configuración de variables de entorno
- Selecciona repository (DynamoDB o memory)
- Crea handlers de aplicación
- Retorna LambdaHandler inicializado

**Diferencias con main.go:**
- No crea servidor HTTP
- No maneja graceful shutdown
- No escucha en puerto

### Dockerfile para Lambda

**Ubicación**: `Dockerfile.lambda`

```dockerfile
FROM golang:1.26.5 AS builder
# ... build steps ...

FROM gcr.io/distroless/static-debian12
COPY --from=builder /bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"]
```

**Características:**
- Multi-stage build (compilación + runtime)
- Binary standalone: `/bootstrap`
- Distroless base para seguridad
- Compatible con AWS Lambda custom runtimes

### Variables de Entorno Lambda

#### Configuración AWS Lambda

| Variable            | Default           | Descripción                              |
| ------------------- | ----------------- | ---------------------------------------- |
| `USE_DYNAMODB`      | `true`            | Usar DynamoDB en Lambda                  |
| `DYNAMODB_TABLE`    | `Recommendations`  | Nombre de la tabla DynamoDB               |
| `AWS_REGION`        | Desde Lambda      | Región AWS (configurada por Lambda)     |
| `GIN_MODE`          | N/A               | No aplica en Lambda                       |

#### Configuración Local para Testing

```bash
USE_DYNAMODB=false  # Usar memory repository para tests locales
```

### Tests Implementados

#### 1. TestLambdaHandler_HandleHealth
- Verifica health check simple
- Verifica status 200 y body correcto

#### 2. TestLambdaHandler_HandleGetRecommendations
- Verifica GET /recommendations/{userId}
- Verifica parsing de path parameters
- Verifica response structure

#### 3. TestLambdaHandler_HandleProcessEvent
- Verifica POST /events
- Verifica parsing de JSON body
- Verifica conversión a domain Event
- Verifica status 202 Accepted

#### 4. TestLambdaHandler_HandleGenerateRecommendations
- Verifica POST /recommendations/generate
- Verifica parsing de request complejo
- Verifica conversión de múltiples eventos
- Verifica response structure

#### 5. TestLambdaHandler_HandleNotFound
- Verifica rutas no existentes
- Verifica status 404

### Trade-offs Considerados

1. **Gin vs Direct Handlers**
   - **Elección**: Direct handlers en Lambda
   - **Trade-off**: Compartición de código vs simplicidad
   - **Decisión**: Handlers de aplicación son compartidos, solo adaptador es diferente

2. **Custom Runtime vs SAM/Serverless Framework**
   - **Elección**: Custom runtime con Docker
   - **Trade-off**: Control vs conveniencia
   - **Decisión**: Custom runtime da más control y flexibilidad

3. **Single Lambda vs Multiple Lambdas**
   - **Elección**: Single Lambda con routing manual
   - **Trade-off**: Simplicidad vs especialización
   - **Decisión**: Single lambda es más simple para este caso de uso

4. **Memory Repository vs DynamoDB en Lambda**
   - **Elección**: DynamoDB por defecto en Lambda
   - **Trade-off**: Persistencia vs simplicidad
   - **Decisión**: Lambda es para producción, usar DynamoDB real

### Validaciones Realizadas

- [x] Código compila sin errores
- [x] Tests del handler Lambda pasan
- [x] Dockerfile para Lambda creado
- [x] Reutilización de handlers de aplicación
- [x] Conversión API Gateway ↔ Domain
- [x] Manejo de errores apropiado
- [x] Logging estructurado para CloudWatch
- [x] Documentación de decisiones

### Estado del Módulo 8

- [x] Cerrado (deployment requeriría configuración AWS adicional)
