# HTTP API - Decisiones de Implementación

## Módulo 4 — HTTP API con Gin

### Objetivo

Exponer los casos de uso mediante una API REST usando Gin.

### Endpoints Implementados

#### 1. GET /health
- **Propósito**: Health check del servicio
- **Response**: `{"status": "ok"}`
- **Status Code**: 200 OK

#### 2. GET /recommendations/:userId
- **Propósito**: Obtener recomendaciones para un usuario
- **Parámetros**: `userId` (path parameter)
- **Response**:
```json
{
  "recommendations": [
    {
      "user_id": "123",
      "product_id": "P10",
      "score": 0.95
    }
  ]
}
```
- **Status Codes**:
  - 200 OK: Recomendaciones encontradas
  - 400 Bad Request: userId vacío
  - 500 Internal Server Error: Error del servicio

#### 3. POST /events
- **Propósito**: Procesar un evento de interacción
- **Request Body**:
```json
{
  "event_id": "event-1",
  "event_type": "product_viewed",
  "user_id": "123",
  "product_id": "P10",
  "product_category": "electronics",
  "product_brand": "Apple",
  "metadata": {
    "device": "mobile",
    "country": "US"
  },
  "occurred_at": "2026-08-08T22:47:17Z"
}
```
- **Response**: `{"status": "event processed"}`
- **Status Codes**:
  - 202 Accepted: Evento procesado
  - 400 Bad Request: JSON inválido o datos inválidos
  - 500 Internal Server Error: Error del servicio

#### 4. POST /recommendations/generate
- **Propósito**: Generar recomendaciones desde un batch de eventos
- **Request Body**:
```json
{
  "user_id": "123",
  "events": [
    {
      "event_id": "event-1",
      "event_type": "product_viewed",
      "user_id": "123",
      "product_id": "P10",
      "occurred_at": "2026-08-08T22:47:17Z"
    }
  ],
  "limit": 10
}
```
- **Response**:
```json
{
  "recommendations": [
    {
      "user_id": "123",
      "product_id": "P10",
      "score": 1.0
    }
  ]
}
```
- **Status Codes**:
  - 200 OK: Recomendaciones generadas
  - 400 Bad Request: JSON inválido o datos inválidos
  - 500 Internal Server Error: Error del servicio

### Handler Structure

#### Handler Struct
```go
type Handler struct {
    getRecommendationsHandler    *queries.GetRecommendationsHandler
    processEventHandler         *commands.ProcessEventHandler
    generateRecommendationsHandler *commands.GenerateRecommendationsHandler
}
```

**Decisiones**:
- Handler contiene referencias a los handlers de aplicación
- Constructor injection para testabilidad
- No hay lógica de negocio en los handlers HTTP

### Validaciones HTTP

#### GetRecommendations
- Valida que `userId` no esté vacío
- Retorna 400 si el parámetro falta

#### ProcessEvent
- Valida el body JSON con `ShouldBindJSON`
- Valida el timestamp en `ToEventDomain`
- Retorna 400 si el JSON es inválido

#### GenerateRecommendations
- Valida el body JSON con `ShouldBindJSON`
- Valida que `user_id`, `events` y `limit` estén presentes
- Valida cada evento en `ToEventDomain`
- Retorna 400 si el JSON es inválido

### Conversiones

#### Request → DTO → Domain
- Los requests HTTP se convierten a DTOs
- Los DTOs se convierten a entidades del dominio
- La conversión maneja timestamps (string RFC3339 → time.Time)

#### Domain → DTO → Response
- Las entidades del dominio se convierten a DTOs
- Los DTOs se serializan a JSON
- La conversión usa funciones helper en el paquete `dto`

### Error Handling

#### Errores del Dominio
- Los errores del dominio se propagan sin wrapping
- Se retornan como 500 Internal Server Error con el mensaje
- Errores de validación se convierten a 400 Bad Request

#### Errores HTTP
- JSON parsing errors → 400 Bad Request
- Binding errors → 400 Bad Request
- Internal errors → 500 Internal Server Error

### Memory Repository

#### Implementación
- **Ubicación**: `internal/infrastructure/persistence/memory/repository.go`
- **Propósito**: Repository en memoria para desarrollo local
- **Características**:
  - Thread-safe con `sync.RWMutex`
  - Mapa de `userID → []Recommendation`
  - Actualización in-place si la recomendación existe
  - Creación nueva si no existe

**Decisiones**:
- Usado para desarrollo local y testing
- Reemplazable por DynamoDB en producción
- No persiste entre reinicios del servidor

### Composition Root

#### main.go
El `main.go` actúa como composition root:

```go
repository := memory.NewMemoryRepository()
service := recommendation.NewService(repository)

getRecommendationsHandler := queries.NewGetRecommendationsHandler(service)
processEventHandler := commands.NewProcessEventHandler(service)
generateRecommendationsHandler := commands.NewGenerateRecommendationsHandler(service)

handler := httpgin.NewHandler(
    getRecommendationsHandler,
    processEventHandler,
    generateRecommendationsHandler,
)

router := httpgin.NewRouter(handler)
server := httpgin.NewServer(config.Port, router)
```

**Decisiones**:
- Inyección manual de dependencias
- No usa framework de DI
- Fácil de testear con fakes/mocks
- Fácil de cambiar a DynamoDB en el futuro

### Tests HTTP

#### TestHealth
- Verifica que `/health` retorna 200
- Verifica el body `{"status":"ok"}`
- Verifica el header `X-Request-ID`

#### TestHealthWithRequestID
- Verifica que el `X-Request-ID` se propaga
- Envía un header y verifica que se retorne el mismo

#### TestGetRecommendations
- Crea una recomendación en el repository
- Llama al endpoint GET /recommendations/:userId
- Verifica que retorne la recomendación correcta

#### TestProcessEvent
- Llama al endpoint POST /events
- Verifica que la recomendación se cree en el repository
- Verifica el status code 202 Accepted

#### TestGenerateRecommendations
- Llama al endpoint POST /recommendations/generate
- Verifica que se generen las recomendaciones correctas
- Verifica el orden por score descendente

### Middleware

#### Request ID
- Genera un UUID si no está presente
- Propaga el header a la respuesta
- Útil para tracing y debugging

#### Logger
- Gin's built-in logger
- Registra todas las requests
- Formato: `[GIN] method path status latency`

#### Recovery
- Gin's built-in recovery
- Captura panics y retorna 500
- Evita que el servidor se caiga

### Validaciones Realizadas

- [x] `internal/infrastructure/http/gin` depende de Gin (intencional)
- [x] `internal/infrastructure/http/gin` no depende de AWS SDK
- [x] Los handlers no contienen lógica de negocio
- [x] Todos los tests HTTP pasan
- [x] Tests con race detector pasan
- [x] Go vet no reporta errores
- [x] Formateo con gofmt aplicado
- [x] El servidor arranca correctamente
- [x] Todos los endpoints responden correctamente

### Trade-offs Considerados

1. **Memory Repository vs DynamoDB**: Se usó memory repository para desarrollo local. Facilita testing y desarrollo rápido. Será reemplazado por DynamoDB en el Módulo 5.

2. **Handler Composition vs Single Handler**: Se usó un struct Handler con múltiples métodos en lugar de handlers separados. Facilita la inyección de dependencias y mantiene consistencia.

3. **Error Responses Simples vs Structured**: Se usaron responses simples con mensajes de error. Futuras mejoras podrían incluir códigos de error estructurados para mejor debugging.

4. **Context Propagation**: Se usa `c.Request.Context()` para propagar context a los handlers de aplicación. Permite cancelación y timeouts en el futuro.

### Estado del Módulo 4

- [x] Cerrado
