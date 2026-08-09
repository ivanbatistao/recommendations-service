# Kinesis - Decisiones de Implementación

## Módulo 6 — Event Processing con Kinesis

### Objetivo

Implementar procesamiento de eventos en tiempo real usando AWS Kinesis Data Streams.

### ¿Qué es Kinesis Data Streams?

**AWS Kinesis Data Streams** es un servicio de streaming de datos en tiempo real que:

- **Recibe** eventos de múltiples fuentes (productores)
- **Almacena** eventos temporalmente (hasta 365 días)
- **Permite** que múltiples consumidores lean los eventos
- **Escala** automáticamente según el volumen de datos

### Arquitectura del Sistema con Kinesis

```text
Ecommerce / Event Generator
            |
            v
       Kinesis Stream
            |
    Shard (Partition Key: UserID)
            |
 vvvvvvvvvvvvvvvvvvvvvv
 Recommendation Processor
            |
       Worker Pool
            |
            v
      DynamoDB
```

### Diseño del Stream de Kinesis

#### Configuración del Stream

**Nombre del stream**: `recommendation-events`

**Configuración inicial:**
- **Shards**: 1 (para desarrollo)
- **Retention period**: 24 horas
- **Partition key**: `UserID`

#### Decisión de Partition Key

**Elección:** UserID como partition key

**Razones:**
1. **Ordenamiento por usuario**: Todos los eventos de un usuario van al mismo shard
2. **Evita condiciones de carrera**: Procesamiento secuencial por usuario
3. **Distribución natural**: Si hay muchos usuarios, se distribuyen bien entre shards
4. **Caso de uso principal**: Procesamiento de eventos por usuario

**Trade-offs:**
- **Ventaja**: Procesamiento ordenado por usuario
- **Desventaja**: Si un usuario es muy activo, puede sobrecargar un shard
- **Mitigación**: Escalado automático de shards según throughput

### Implementación del Productor

#### Producer

**Ubicación**: `internal/infrastructure/streaming/kinesis/producer.go`

**Estructura:**
```go
type Producer struct {
    client    *kinesis.Client
    streamName string
}
```

**Métodos:**

##### 1. PublishEvent
```go
func (p *Producer) PublishEvent(
    ctx context.Context,
    ev event.Event,
) error
```

**Implementación:**
- Convierte evento a JSON
- Usa `PutRecord` operation
- Usa `UserID` como partition key

**Analogía:**
- Es como publicar un mensaje en una cola
- Pero con partition key para distribución
- Los eventos con mismo partition key van al mismo shard

##### 2. PublishBatch
```go
func (p *Producer) PublishBatch(
    ctx context.Context,
    events []event.Event,
) error
```

**Implementación:**
- Convierte múltiples eventos a JSON
- Usa `PutRecords` operation
- Soporta hasta 500 records por llamada

**Ventajas:**
- **Más eficiente**: Una llamada HTTP en lugar de muchas
- **Límite**: Máximo 500 records por llamada
- **Costo**: Menor número de API calls

### Implementación del Consumidor

#### Consumer

**Ubicación**: `internal/infrastructure/streaming/kinesis/consumer.go`

**Estructura:**
```go
type Consumer struct {
    client     *kinesis.Client
    streamName string
}
```

**Métodos:**

##### 1. GetShardIterator
```go
func (c *Consumer) GetShardIterator(
    ctx context.Context,
    shardID string,
    iteratorType types.ShardIteratorType,
) (string, error)
```

**Implementación:**
- Obtiene un "cursor" para leer eventos de un shard
- El tipo de iterator determina desde dónde empezar

**Iterator Types:**
- **TRIM_HORIZON**: Lee desde el evento más antiguo disponible
- **LATEST**: Lee solo eventos nuevos (después del iterator)
- **AT_SEQUENCE_NUMBER**: Lee desde un sequence number específico

##### 2. GetRecords
```go
func (c *Consumer) GetRecords(
    ctx context.Context,
    shardIterator string,
    limit int,
) ([]event.Event, string, error)
```

**Implementación:**
- Lee eventos desde el shard iterator
- Convierte JSON a eventos del dominio
- Retorna eventos y el siguiente iterator

**Features:**
- **Paginación**: Permite leer eventos secuencialmente
- **Límite**: Máximo 1000 records por llamada
- **Next iterator**: Para continuar reading

##### 3. ListShards
```go
func (c *Consumer) ListShards(
    ctx context.Context,
) ([]string, error)
```

**Implementación:**
- Lista todos los shards del stream
- Retorna IDs de shards

**Uso:**
- Necesario para saber qué shards consumir
- Útil para monitoreo y scaling

### Configuración del Cliente Kinesis

#### NewKinesisClient (AWS Real)

```go
func NewKinesisClient(ctx context.Context, region string) (*kinesis.Client, error)
```

**Características:**
- Usa `config.LoadDefaultConfig` para cargar credenciales automáticamente
- Busca credenciales en:
  1. Variables de entorno (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  2. Archivo ~/.aws/credentials
  3. IAM roles (si está en EC2/Lambda)
- **Uso**: Producción en AWS

#### NewLocalKinesisClient (Desarrollo Local)

```go
func NewLocalKinesisClient(ctx context.Context, endpoint string) (*kinesis.Client, error)
```

**Características:**
- Override del endpoint con `BaseEndpoint`
- Apunta a local: `http://localhost:4566` (LocalStack)
- **Uso**: Desarrollo con LocalStack

### Serialización JSON

#### Tags en Event Entity

```go
type Event struct {
    EventID         string    `json:"event_id"`
    EventType       Type      `json:"event_type"`
    UserID          string    `json:"user_id"`
    ProductID       string    `json:"product_id"`
    ProductCategory string    `json:"product_category,omitempty"`
    ProductBrand    string    `json:"product_brand,omitempty"`
    Metadata        Metadata  `json:"metadata"`
    OccurredAt      time.Time `json:"occurred_at"`
}
```

**Propósito:**
- Control explícito del mapeo JSON
- Nombres snake_case para JSON (convención HTTP)
- Nombres PascalCase en Go (convención Go)

**Sin tags vs Con tags:**
- **Sin tags**: JSON tags sería eventID, eventType (no convención HTTP)
- **Con tags**: event_id, event_type (convención HTTP estándar)

### Paquetes AWS SDK Instalados

1. **`github.com/aws/aws-sdk-go-v2/service/kinesis`**: Cliente específico de Kinesis
2. **`github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream`**: Protocolo de streaming

### Tests

#### Estado Actual

Los tests de Kinesis están en `Skip` porque:

1. **SDK no usa interfaces**: El cliente Kinesis es un struct concreto
2. **Mocking difícil**: No se puede mockear fácilmente sin wrappers
3. **Tests de integración**: Requieren Kinesis Local o LocalStack

#### Plan para Tests

Los tests de integración se implementarán en el **Módulo 9 - MiniStack** cuando configuremos:
- Kinesis Local para testing
- Tests end-to-end con Kinesis real
- Validación de producer y consumer

### Flujo de Trabajo del Sistema

#### 1. Publicación de Eventos
```text
Ecommerce API
    |
    v
POST /events
    |
    v
Process Event Command
    |
    v
Kinesis Producer
    |
    v
Kinesis Stream (recommendation-events)
```

#### 2. Consumo de Eventos
```text
Kinesis Stream
    |
    v
Kinesis Consumer
    |
    v
Worker Pool (Módulo 7)
    |
    v
Recommendation Service
    |
    v
DynamoDB
```

### Trade-offs Considerados

1. **Partition Key Strategy**
   - **Elección**: UserID como partition key
   - **Trade-off**: Ordenamiento por usuario vs distribución uniforme
   - **Decisión**: Ordenamiento por usuario es más importante para el caso de uso

2. **Single Record vs Batch**
   - **Elección**: Ambos métodos disponibles
   - **Trade-off**: Simplicidad vs eficiencia
   - **Decisión**: Batch para alta throughput, single para baja latencia

3. **JSON vs Binary Serialization**
   - **Elección**: JSON para Kinesis
   - **Trade-off**: Tamaño vs legibilidad
   - **Decisión**: JSON es más legible y suficiente para este caso de uso

4. **Tests Unitarios vs Integración**
   - **Elección**: Tests de integración en Módulo 9
   - **Trade-off**: Cobertura temprana vs infraestructura real
   - **Decisión**: Tests de integración con Kinesis Local son más valiosos

### Validaciones Realizadas

- [x] Código compila sin errores
- [x] Todos los tests existentes pasan
- [x] JSON tags en Event entity
- [x] Producer implementado (single y batch)
- [x] Consumer implementado (GetRecords, GetShardIterator, ListShards)
- [x] Cliente Kinesis para AWS y local
- [x] Documentación de decisiones

### Estado del Módulo 6

- [x] Cerrado (tests de integración pendientes en Módulo 9)
