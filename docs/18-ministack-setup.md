# MiniStack Local Development Setup

## Overview

Este proyecto utiliza **LocalStack** para simular servicios de AWS localmente, permitiendo desarrollo y testing sin necesidad de una cuenta AWS real.

## Servicios Simulados

- **DynamoDB**: Base de datos NoSQL para almacenar recomendaciones
- **Kinesis**: Stream de eventos para procesamiento en tiempo real

## Requisitos Previos

- Docker y Docker Compose instalados
- AWS CLI instalado (opcional, para scripts de inicialización)

## Inicio Rápido

### 1. Iniciar LocalStack

```bash
docker-compose up -d localstack
```

Esto iniciará LocalStack en `http://localhost:4566` con DynamoDB y Kinesis habilitados.

### 2. Inicializar Recursos AWS

Ejecutar el script de inicialización para crear la tabla DynamoDB y el stream Kinesis:

```bash
./scripts/init-localstack.sh
```

Esto creará:
- Tabla DynamoDB: `Recommendations` (partition key: `UserID`)
- Stream Kinesis: `recommendations-events` (1 shard)

### 3. Verificar Funcionamiento

```bash
./scripts/test-localstack.sh
```

Esto ejecutará tests de conectividad con ambos servicios.

### 4. Iniciar la Aplicación

```bash
docker-compose up recommendation-service
```

La aplicación se conectará automáticamente a LocalStack usando las variables de entorno configuradas.

## Configuración

### Variables de Entorno

| Variable | Valor Default | Descripción |
|----------|---------------|-------------|
| `AWS_REGION` | `us-east-1` | Región AWS simulada |
| `AWS_ACCESS_KEY_ID` | `test` | Credencial falsa para LocalStack |
| `AWS_SECRET_ACCESS_KEY` | `test` | Credencial falsa para LocalStack |
| `DYNAMODB_ENDPOINT` | `http://localstack:4566` | Endpoint DynamoDB local |
| `KINESIS_ENDPOINT` | `http://localstack:4566` | Endpoint Kinesis local |
| `USE_LOCAL_AWS` | `true` | Usar servicios AWS locales |

### Docker Compose

```yaml
localstack:
  image: localstack/localstack:latest
  ports:
    - "4566:4566"
  environment:
    - SERVICES=dynamodb,kinesis
    - AWS_DEFAULT_REGION=us-east-1
    - AWS_ACCESS_KEY_ID=test
    - AWS_SECRET_ACCESS_KEY=test
```

## Uso

### Desarrollo Local

La aplicación detecta automáticamente cuando usar LocalStack basándose en las variables de entorno:

```go
if config.DynamoDBEndpoint != "" {
    // Usar DynamoDB local (LocalStack)
    client, err = dynamodb.NewLocalDynamoDBClient(
        context.Background(),
        config.DynamoDBEndpoint,
    )
} else {
    // Usar AWS DynamoDB real
    client, err = dynamodb.NewDynamoDBClient(
        context.Background(),
        config.AWSRegion,
    )
}
```

### Event Generator con LocalStack

Para usar el Event Generator con Kinesis local:

```bash
./event-generator --kinesis --stream-name recommendations-events --endpoint http://localhost:4566
```

## Scripts Disponibles

### `scripts/init-localstack.sh`
Inicializa recursos AWS en LocalStack:
- Crea tabla DynamoDB
- Crea stream Kinesis
- Espera que los recursos estén activos

### `scripts/test-localstack.sh`
Testea conectividad con LocalStack:
- Lista tablas DynamoDB
- Describe tabla Recommendations
- Lista streams Kinesis
- Escribe/lee datos de prueba

## Troubleshooting

### LocalStack no inicia
```bash
# Verificar logs
docker-compose logs localstack

# Reiniciar
docker-compose restart localstack
```

### Tabla o stream no existen
```bash
# Re-inicializar recursos
./scripts/init-localstack.sh
```

### Error de conexión
```bash
# Verificar que LocalStack está corriendo
curl http://localhost:4566

# Verificar que el puerto 4566 no está en uso
lsof -i :4566
```

### Permisos AWS
LocalStack ignora las credenciales reales, pero requiere valores válidos:
```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

## Datos Locales

Los datos de LocalStack se guardan en `./localstack_data/`:
```yaml
volumes:
  - "./localstack_data:/tmp/localstack/data"
```

Para limpiar datos:
```bash
rm -rf localstack_data/*
docker-compose restart localstack
```

## Limitaciones

LocalStack no es 100% compatible con AWS real. Limitaciones conocidas:

- Algunas características avanzadas de DynamoDB no están implementadas
- Kinesis tiene algunas diferencias en comportamiento
- No hay costos, pero tampoco garantías de disponibilidad

## Próximos Pasos

Una vez LocalStack está funcionando:
1. ✅ Event Generator puede enviar eventos a Kinesis local
2. ✅ La aplicación puede almacenar recomendaciones en DynamoDB local
3. ✅ Ready para load testing con k6

## Referencias

- [LocalStack Documentation](https://docs.localstack.cloud/)
- [DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)
- [Kinesis Developer Guide](https://docs.aws.amazon.com/streams/latest/dev/introduction.html)
