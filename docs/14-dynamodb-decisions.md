# DynamoDB - Decisiones de Implementación

## Módulo 5 — Persistence Layer con DynamoDB

### Objetivo

Implementar persistencia real usando DynamoDB para producción y desarrollo local con MiniStack.

### Diseño de la Tabla DynamoDB

#### Esquema de la Tabla

```text
Table Name: Recommendations
Partition Key: UserID (String)
Sort Key: ProductID (String)
Attributes:
  - Score (Number)
```

#### Decisión de Partition Key

**Elección:** UserID como Partition Key, ProductID como Sort Key

**Razones:**
1. **Caso de uso principal**: `GET /recommendations/:userId` - queries por usuario
2. **Performance**: Query por partition key es O(1)
3. **Simplicidad**: Implementación directa del Repository interface
4. **Costo**: Sin GSI, más económico

**Trade-offs:**
- **Ventaja**: Queries por usuario son muy rápidos
- **Desventaja**: Un usuario con miles de productos puede tener partición grande
- **Mitigación**: DynamoDB maneja particiones grandes con paginación automática

### Implementación del Repository

#### DynamoDBRepository

**Ubicación**: `internal/infrastructure/persistence/dynamodb/repository.go`

**Estructura:**
```go
type DynamoDBRepository struct {
    client    *dynamodb.Client
    tableName string
}
```

**Métodos:**

##### 1. GetByUserID
```go
func (r *DynamoDBRepository) GetByUserID(
    ctx context.Context,
    userID string,
) ([]recommendation.Recommendation, error)
```

**Implementación:**
- Usa `Query` operation con KeyConditionExpression
- Expresión: `UserID = :uid`
- Convierte AttributeValues a structs Go con `UnmarshalMap`

**Analogía SQL:**
```sql
SELECT * FROM Recommendations WHERE UserID = :uid
```

##### 2. Save
```go
func (r *DynamoDBRepository) Save(
    ctx context.Context,
    rec recommendation.Recommendation,
) error
```

**Implementación:**
- Usa `PutItem` operation
- Convierte struct Go a AttributeValues con `MarshalMap`
- **Idempotente**: Si existe (mismo UserID + ProductID), lo reemplaza

**Comportamiento:**
- **Create**: Item no existe → crea nuevo
- **Update**: Item existe → reemplaza completamente
- **Upsert**: Operación de upsert automática

### Mapeo DynamoDB ↔ Go

#### Tags dynamodbav

```go
type Recommendation struct {
    UserID    string `dynamodbav:"UserID"`
    ProductID string `dynamodbav:"ProductID"`
    Score     float64 `dynamodbav:"Score"`
}
```

**Propósito:**
- Control explícito del mapeo entre structs Go y atributos DynamoDB
- Evita inconsistencias de nombres
- Permite nombres de atributos diferentes a nombres de campos Go

**Sin tags vs Con tags:**
- **Sin tags**: SDK intenta mapeo automático (propenso a errores)
- **Con tags**: Mapeo explícito y robusto

### Configuración del Cliente DynamoDB

#### NewDynamoDBClient (AWS Real)

```go
func NewDynamoDBClient(ctx context.Context, region string) (*dynamodb.Client, error)
```

**Características:**
- Usa `config.LoadDefaultConfig` para cargar credenciales automáticamente
- Busca credenciales en:
  1. Variables de entorno (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
  2. Archivo ~/.aws/credentials
  3. IAM roles (si está en EC2/Lambda)
- **Uso**: Producción en AWS

#### NewLocalDynamoDBClient (Desarrollo Local)

```go
func NewLocalDynamoDBClient(ctx context.Context, endpoint string) (*dynamodb.Client, error)
```

**Características:**
- Override del endpoint con `BaseEndpoint`
- Apunta a local: `http://localhost:8000`
- **Uso**: Desarrollo con MiniStack o DynamoDB Local

### Variables de Entorno

#### Configuración

```go
type Config struct {
    Port              string
    UseDynamoDB       bool
    DynamoDBTable     string
    DynamoDBEndpoint  string
    AWSRegion         string
}
```

#### Variables de Entorno

| Variable            | Default           | Descripción                              |
| ------------------- | ----------------- | ---------------------------------------- |
| `PORT`              | `8080`            | Puerto del servidor HTTP                 |
| `USE_DYNAMODB`      | `false`           | Activar DynamoDB                         |
| `DYNAMODB_TABLE`    | `Recommendations`  | Nombre de la tabla DynamoDB               |
| `DYNAMODB_ENDPOINT` | `""`              | Endpoint para DynamoDB Local             |
| `AWS_REGION`        | `us-east-1`       | Región AWS                              |

#### Ejemplos de Uso

**Desarrollo con Memory Repository:**
```bash
# Sin variables de entorno adicionales
go run cmd/api/main.go
```

**Desarrollo con DynamoDB Local:**
```bash
export USE_DYNAMODB=true
export DYNAMODB_ENDPOINT=http://localhost:8000
go run cmd/api/main.go
```

**Producción en AWS:**
```bash
export USE_DYNAMODB=true
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
go run cmd/api/main.go
```

### Selección Dinámica del Repository

#### Lógica en main.go

```go
if config.UseDynamoDB {
    // Configurar cliente DynamoDB
    if config.DynamoDBEndpoint != "" {
        client = dynamodb.NewLocalDynamoDBClient(...)
    } else {
        client = dynamodb.NewDynamoDBClient(...)
    }
    repository = dynamodb.NewDynamoDBRepository(client, config.DynamoDBTable)
} else {
    repository = memory.NewMemoryRepository()
}
```

**Beneficios:**
- **Desarrollo rápido**: Memory repository por defecto
- **Testing local**: DynamoDB Local con endpoint override
- **Producción**: DynamoDB con credenciales AWS
- **Mismo código**: Cambio transparente sin modificar lógica de negocio

### Paquetes AWS SDK Instalados

1. **`github.com/aws/aws-sdk-go-v2`**: SDK principal de AWS
2. **`github.com/aws/aws-sdk-go-v2/config`**: Configuración de credenciales y regiones
3. **`github.com/aws/aws-sdk-go-v2/service/dynamodb`**: Cliente específico de DynamoDB
4. **`github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue`**: Conversión entre Go y DynamoDB

### Tests

#### Estado Actual

Los tests de DynamoDB están en `Skip` porque:

1. **SDK no usa interfaces**: El cliente DynamoDB es un struct concreto, no una interfaz
2. **Mocking difícil**: No se puede mockear fácilmente sin wrappers
3. **Tests de integración**: Requieren DynamoDB Local o LocalStack

#### Plan para Tests

Los tests de integración se implementarán en el **Módulo 9 - MiniStack** cuando configuremos:
- DynamoDB Local para testing
- Tests end-to-end con DynamoDB real
- Validación de schema y operaciones

### Trade-offs Considerados

1. **Partition Key Strategy**
   - **Elección**: UserID como PK
   - **Trade-off**: Queries rápidos vs particiones grandes
   - **Decisión**: Queries por usuario son el caso de uso principal

2. **Memory vs DynamoDB Repository**
   - **Elección**: Ambos con selección dinámica
   - **Trade-off**: Complejidad vs flexibilidad
   - **Decisión**: Flexibilidad para diferentes entornos

3. **Tags dynamodbav**
   - **Elección**: Tags explícitos en todos los campos
   - **Trade-off**: Verbosidad vs robustez
   - **Decisión**: Robustez y control explícito

4. **Tests Unitarios vs Integración**
   - **Elección**: Tests de integración en Módulo 9
   - **Trade-off**: Cobertura temprana vs infraestructura real
   - **Decisión**: Tests de integración con DynamoDB Local son más valiosos

### Validaciones Realizadas

- [x] Código compila sin errores
- [x] Todos los tests existentes pasan
- [x] Configuración de variables de entorno
- [x] Selección dinámica de repository
- [x] Mapeo DynamoDB ↔ Go con tags
- [x] Cliente DynamoDB para AWS y local
- [x] Documentación de decisiones

### Estado del Módulo 5

- [x] Cerrado (tests de integración pendientes en Módulo 9)
