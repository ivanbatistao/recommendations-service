# Decisiones de Modelado del Dominio

## Módulo 2 — Domain Layer

### Entidades del Dominio

#### 1. Event
- **Ubicación**: `internal/domain/event/entity.go`
- **Propósito**: Representa un evento de interacción del usuario con productos
- **Campos**:
  - `EventID`: Identificador único del evento
  - `EventType`: Tipo de evento (viewed, search, added_to_cart, purchased)
  - `UserID`: Identificador del usuario que generó el evento
  - `ProductID`: Identificador del producto involucrado
  - `ProductCategory`: Categoría del producto (opcional para filtrado futuro)
  - `ProductBrand`: Marca del producto (opcional para filtrado futuro)
  - `Metadata`: Información adicional (device, country)
  - `OccurredAt`: Timestamp del evento

**Decisiones**:
- Se usa un string para `EventID` para compatibilidad con AWS Kinesis
- `EventType` es un tipo enum para seguridad de tipos
- `Metadata` es un struct opcional para extender información sin cambiar la entidad principal

#### 2. Recommendation
- **Ubicación**: `internal/domain/recommendation/entity.go`
- **Propósito**: Representa una recomendación calculada para un usuario
- **Campos**:
  - `UserID`: Identificador del usuario destinatario
  - `ProductID`: Identificador del producto recomendado
  - `Score`: Puntuación de relevancia calculada

**Decisiones**:
- El score es un `float64` para permitir cálculos de suma y promedio
- No incluye timestamps porque el score refleta el estado actual
- Estructura minimalista para optimizar almacenamiento en DynamoDB

### Algoritmos del Dominio

#### 1. ScoreEvent
- **Ubicación**: `internal/domain/recommendation/scoring.go`
- **Propósito**: Asignar peso a cada tipo de evento
- **Scores**:
  - `ProductViewed`: 1.0
  - `SearchPerformed`: 2.0
  - `ProductAddedCart`: 3.0
  - `ProductPurchased`: 5.0

**Decisiones**:
- Los scores son constantes simples pero configuras en futuro
- El peso incrementa con la intención de compra
- Eventos desconocidos retornan 0 para ignorarlos

#### 2. CalculateInterest
- **Ubicación**: `internal/domain/recommendation/interest.go`
- **Propósito**: Agregar scores por producto para un usuario específico
- **Lógica**:
  - Filtra eventos que no pertenecen al usuario
  - Suma scores por producto
  - Ignora eventos con score 0

**Decisiones**:
- Filtrado explícito por `userID` para evitar contaminación entre usuarios
- Retorna un mapa para O(1) lookup en ranking
- Función pura sin efectos secundarios

#### 3. Rank
- **Ubicación**: `internal/domain/recommendation/ranking.go`
- **Propósito**: Ordenar productos por score y aplicar límite
- **Lógica**:
  - Convierte mapa de interés a slice de recomendaciones
  - Ordena por score descendente
  - Desempate por `ProductID` ascendente (determinístico)
  - Aplica límite N (Top N)

**Decisiones**:
- Ordenamiento determinístico para resultados reproducibles
- Desempate por ID alfabético para consistencia
- Retorna slice vacío si no hay suficientes productos
- Ajusta límite si hay menos productos que el solicitado

### Errores del Dominio

#### Recommendation Errors
- **Ubicación**: `internal/domain/recommendation/errors.go`
- **Errores definidos**:
  - `ErrInvalidUserID`: UserID vacío o inválido
  - `ErrInvalidProductID`: ProductID vacío o inválido
  - `ErrInvalidEventType`: Tipo de evento no reconocido
  - `ErrInvalidLimit`: Límite <= 0
  - `ErrUserNotFound`: Usuario no existe en repositorio
  - `ErrProductNotFound`: Producto no existe en repositorio
  - `ErrRepositoryError`: Error genérico de repositorio
  - `ErrRecommendationNotFound`: Recomendación no encontrada

**Decisiones**:
- Errores simples sin wrapping para permitir comparación directa
- Separación entre errores de validación y errores de infraestructura
- Errores específicos para cada caso de uso

#### Event Errors
- **Ubicación**: `internal/domain/event/errors.go`
- **Errores definidos**:
  - `ErrInvalidEventID`: EventID vacío o inválido
  - `ErrInvalidUserID`: UserID vacío o inválido
  - `ErrInvalidProductID`: ProductID vacío o inválido
  - `ErrInvalidEventType`: Tipo de evento no reconocido
  - `ErrInvalidTimestamp`: Timestamp inválido
  - `ErrInvalidMetadata`: Metadata inválida
  - `ErrEventNotFound`: Evento no encontrado

**Decisiones**:
- Consistentes con errores de recommendation
- Preparados para validación futura de eventos

### Interfaces del Dominio

#### Repository
- **Ubicación**: `internal/domain/recommendation/repository.go`
- **Métodos**:
  - `GetByUserID(ctx, userID)`: Obtener recomendaciones por usuario
  - `Save(ctx, recommendation)`: Guardar o actualizar recomendación

**Decisiones**:
- Interfaz minimalista para facilitar múltiples implementaciones
- Usa `context.Context` para cancelación y timeouts
- No incluye métodos de eliminación (las recomendaciones se actualizan)

### Servicio del Dominio

#### Service
- **Ubicación**: `internal/domain/recommendation/service.go`
- **Métodos**:
  - `GetByUserID`: Obtiene recomendaciones con validación
  - `GenerateRecommendations`: Genera recomendaciones desde eventos
  - `ProcessEvent`: Procesa un evento y actualiza recomendaciones

**Decisiones**:
- `GenerateRecommendations` es funcional (no depende de repositorio)
- `ProcessEvent` es imperativo (usa repositorio para persistencia)
- Validación de parámetros en nivel de servicio
- `ProcessEvent` ignora eventos con score 0

### Trade-offs Considerados

1. **Score simple vs modelo complejo**: Se eligió score simple sumatorio por simplicidad y velocidad. Futuras mejoras podrían incluir decaimiento temporal o factores contextuales.

2. **UserID explícito vs implícito**: Se pasa `userID` explícitamente en `GenerateRecommendations` para evitar ambigüedad y permitir procesamiento por lotes.

3. **Ranking determinístico vs aleatorio**: Se usa desempate por `ProductID` para resultados reproducibles en tests y debugging.

4. **Errores simples vs error wrapping**: Se usan errores simples sin wrapping para comparación directa (`errors.Is`). Futuros errores complejos podrían usar `fmt.Errorf` con wrapping.

### Validaciones Realizadas

- [x] `internal/domain` no depende de Gin
- [x] `internal/domain` no depende de AWS SDK
- [x] `internal/domain` no depende de infraestructura
- [x] Todos los tests unitarios pasan
- [x] Tests con race detector pasan
- [x] Go vet no reporta errores
- [x] Formateo con gofmt aplicado

### Estado del Módulo 2

- [x] Cerrado
