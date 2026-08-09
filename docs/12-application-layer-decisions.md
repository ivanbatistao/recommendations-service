# Application Layer - Decisiones de Implementación

## Módulo 3 — Application Layer

### Objetivo

Implementar los casos de uso de la aplicación orquestando el dominio sin depender de infraestructura concreta.

### Arquitectura

La capa de aplicación sigue el patrón CQRS (Command Query Responsibility Segregation):

```
Commands (Write)
├── ProcessEventCommand
└── GenerateRecommendationsCommand

Queries (Read)
└── GetRecommendationsQuery

DTOs (Data Transfer Objects)
├── RecommendationDTO
└── EventDTO
```

### Casos de Uso Implementados

#### 1. Get Recommendations (Query)
- **Ubicación**: `internal/application/queries/get_recommendations.go`
- **Handler**: `GetRecommendationsHandler`
- **Propósito**: Obtener recomendaciones existentes para un usuario
- **Flujo**:
  1. Recibe `GetRecommendationsQuery` con `UserID`
  2. Delega a `recommendation.Service.GetByUserID`
  3. Retorna slice de `Recommendation` (entidades del dominio)

**Decisiones**:
- Query pattern para operaciones de lectura
- No usa DTOs en el handler (retorna entidades del dominio)
- La conversión a DTOs se hará en la capa HTTP

#### 2. Process Event (Command)
- **Ubicación**: `internal/application/commands/process_event.go`
- **Handler**: `ProcessEventHandler`
- **Propósito**: Procesar un evento de interacción y actualizar recomendaciones
- **Flujo**:
  1. Recibe `ProcessEventCommand` con un `Event`
  2. Delega a `recommendation.Service.ProcessEvent`
  3. El servicio actualiza o crea recomendaciones vía Repository

**Decisiones**:
- Command pattern para operaciones de escritura
- Usa entidades del dominio directamente
- Validación se hace en el dominio

#### 3. Generate Recommendations (Command)
- **Ubicación**: `internal/application/commands/generate_recommendations.go`
- **Handler**: `GenerateRecommendationsHandler`
- **Propósito**: Generar recomendaciones desde un batch de eventos
- **Flujo**:
  1. Recibe `GenerateRecommendationsCommand` con `UserID`, `Events`, `Limit`
  2. Delega a `recommendation.Service.GenerateRecommendations`
  3. Retorna slice de `Recommendation` generado

**Decisiones**:
- Command pattern aunque es de lectura (para consistencia)
- Útil para procesamiento por lotes offline
- No persiste resultados (solo cálculo en memoria)

### DTOs Implementados

#### 1. RecommendationDTO
- **Ubicación**: `internal/application/dto/recommendation_dto.go`
- **Propósito**: Transferir datos de recomendaciones entre capas
- **Campos**:
  - `UserID`: string
  - `ProductID`: string
  - `Score`: float64

**Funciones**:
- `FromDomain`: Convierte entidad del dominio a DTO
- `ToDomain`: Convierte DTO a entidad del dominio
- `FromDomainSlice`: Convierte slice de entidades a slice de DTOs

**Decisiones**:
- Campos idénticos a la entidad del dominio (por ahora)
- Funciones de conversión para facilitar testing
- Slice conversion para uso en HTTP responses

#### 2. EventDTO
- **Ubicación**: `internal/application/dto/event_dto.go`
- **Propósito**: Transferir datos de eventos entre capas
- **Campos**:
  - `EventID`: string
  - `EventType`: string (no enum para JSON)
  - `UserID`: string
  - `ProductID`: string
  - `ProductCategory`: string (opcional)
  - `ProductBrand`: string (opcional)
  - `Metadata`: map[string]string (para JSON)
  - `OccurredAt`: string (ISO 8601)

**Funciones**:
- `FromEventDomain`: Convierte entidad del dominio a DTO
- `ToEventDomain`: Convierte DTO a entidad del dominio con parsing de timestamp

**Decisiones**:
- `EventType` como string para serialización JSON
- `Metadata` como map para flexibilidad
- `OccurredAt` como string RFC3339 para transporte HTTP
- Parsing de timestamp en `ToEventDomain` con manejo de errores

### Inyección de Dependencias

Todos los handlers reciben el `recommendation.Service` en el constructor:

```go
func NewGetRecommendationsHandler(service *recommendation.Service) *GetRecommendationsHandler
func NewProcessEventHandler(service *recommendation.Service) *ProcessEventHandler
func NewGenerateRecommendationsHandler(service *recommendation.Service) *GenerateRecommendationsHandler
```

**Decisiones**:
- Constructor injection para testabilidad
- No usa frameworks de DI (inyección manual en composition root)
- Service es compartido entre handlers (singleton)

### Manejo de Errores

Los handlers propagan errores del dominio sin transformación:

```go
func (h *Handler) Execute(ctx context.Context, cmd Command) (Result, error) {
    // Delega al servicio del dominio
    return h.service.SomeMethod(ctx, cmd.Params)
}
```

**Decisiones**:
- Errores del dominio fluyen sin wrapping
- La capa HTTP manejará la conversión a status codes
- `context.Context` se pasa para cancelación y timeouts

### Tests Implementados

#### 1. GetRecommendations Tests
- **Ubicación**: `internal/application/queries/get_recommendations_test.go`
- **Coverage**:
  - Caso exitoso con fake repository
  - Verificación de resultados ordenados

#### 2. ProcessEvent Tests
- **Ubicación**: `internal/application/commands/process_event_test.go`
- **Coverage**:
  - Actualización de recomendación existente
  - Creación de nueva recomendación
  - Ignorar eventos con tipo desconocido
  - Mock repository para verificar llamadas

#### 3. GenerateRecommendations Tests
- **Ubicación**: `internal/application/commands/generate_recommendations_test.go`
- **Coverage**:
  - Caso exitoso con eventos de prueba
  - Validación de userID vacío
  - Validación de límite inválido

#### 4. DTO Tests
- **Ubicación**: `internal/application/dto/recommendation_dto_test.go`
- **Coverage**:
  - Conversión Domain → DTO
  - Conversión DTO → Domain
  - Conversión de slices

### Validaciones Realizadas

- [x] `internal/application` no depende de Gin
- [x] `internal/application` no depende de AWS SDK
- [x] `internal/application` no depende de infraestructura concreta
- [x] Todos los tests unitarios pasan
- [x] Tests con race detector pasan
- [x] Go vet no reporta errores
- [x] Formateo con gofmt aplicado

### Trade-offs Considerados

1. **Commands vs Functions**: Se usó el patrón Command/Handler en lugar de funciones simples para consistencia CQRS y extensibilidad futura (middleware, logging, etc.).

2. **DTOs vs Domain Entities**: Los handlers de commands/queries retornan entidades del dominio, no DTOs. La conversión a DTOs se delega a la capa HTTP para mantener la aplicación agnóstica del protocolo.

3. **GenerateRecommendations como Command**: Aunque es una operación de lectura, se implementó como Command para consistencia con el caso de uso de procesamiento por lotes offline.

### Estado del Módulo 3

- [x] Cerrado
