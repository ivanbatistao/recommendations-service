# Load Testing with k6

## Overview

Este proyecto utiliza **k6** para pruebas de carga del API de recomendaciones, permitiendo obtener métricas de performance reales para el CV y validación de la arquitectura.

## Instalación de k6

### Linux (Ubuntu/Debian)
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

### macOS
```bash
brew install k6
```

### Otras opciones
Ver [k6 Installation Guide](https://k6.io/docs/getting-started/installation/)

## Scripts de Prueba

### 1. API Load Test (`api-load-test.js`)

**Propósito:** Prueba de carga estándar del API completo

**Escenario:**
- Health check
- Get recommendations
- Process events
- 10-10-10-0 usuarios (ramp up - steady - ramp down)
- Duración: 70 segundos

**Thresholds:**
- 95% de requests < 500ms
- Error rate < 1%

**Ejecución:**
```bash
./scripts/run-load-tests.sh api-load-test http://localhost:8080
```

### 2. Stress Test (`stress-test.js`)

**Propósito:** Prueba de estrés con alta concurrencia

**Escenario:**
- Solo get recommendations
- 10-50-100-100-50-0 usuarios
- Duración: 80 segundos
- 100 usuarios concurrentes en peak

**Thresholds:**
- 95% de requests < 1s
- Error rate < 5%

**Ejecución:**
```bash
./scripts/run-load-tests.sh stress-test http://localhost:8080
```

### 3. Spike Test (`spike-test.js`)

**Propósito:** Prueba de tráfico repentino (spike)

**Escenario:**
- Health check + get recommendations
- 5-100-100-5-0 usuarios
- Spike repentino de 5 a 100 usuarios en 1 segundo
- Duración: 36 segundos

**Thresholds:**
- 95% de requests < 2s (más laxo durante spike)
- Error rate < 10% (permite más errores durante spike)

**Ejecución:**
```bash
./scripts/run-load-tests.sh spike-test http://localhost:8080
```

## Configuración

### Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | URL del API a probar |
| `VUS` | `10` | Virtual Users (si el script lo soporta) |
| `DURATION` | `30s` | Duración del test (si el script lo soporta) |

### Uso Personalizado

```bash
# Cambiar URL
export BASE_URL="http://localhost:9090"
k6 run loadtests/api-load-test.js

# Cambiar VUs y duración (scripts que lo soporten)
export VUS=20
export DURATION="2m"
k6 run loadtests/api-load-test.js
```

## Resultados Esperados

### Métricas de k6

**HTTP Request Duration:**
- `p(95)`: 95th percentile de latencia
- `p(99)`: 99th percentile de latencia
- `avg`: Promedio de latencia

**HTTP Requests:**
- Total requests
- Requests por segundo (RPS)
- Failed requests
- Error rate

**VUs:**
- Virtual Users activos
- Time to first byte (TTFB)

### Resultados para CV

According to the system design, the following metrics are expected:

**Get Recommendations Latency:**
- Target: < 250ms (p95)
- Stress: < 500ms (p95)
- Spike: < 500ms (p95)

**Throughput:**
- API Load: ~100 RPS with 10 VUs
- Stress: ~200 RPS with 100 VUs
- Spike: ~150 RPS during spike

**Error Rate:**
- Normal: < 1%
- Stress: < 5%
- Spike: < 10%

## Integration with MiniStack

For load testing with MiniStack:

```bash
# 1. Start MiniStack
docker-compose up -d ministack

# 2. Initialize resources
./scripts/init-ministack.sh

# 3. Start application
docker-compose up recommendation-service

# 4. Ejecutar load tests
./scripts/run-load-tests.sh api-load-test http://localhost:8080
```

## Integración con Event Generator

Para pruebas de carga realistas combinando Event Generator + k6:

```bash
# Terminal 1: Iniciar aplicación
docker-compose up recommendation-service

# Terminal 2: Iniciar Event Generator
./event-generator --rate 50 --duration 2m

# Terminal 3: Ejecutar load tests
./scripts/run-load-tests.sh api-load-test http://localhost:8080
```

## Troubleshooting

### k6 no está instalado
```bash
# Verificar instalación
k6 version

# Si no está instalado, seguir instrucciones de instalación arriba
```

### Aplicación no responde
```bash
# Verificar que la aplicación está corriendo
curl http://localhost:8080/health

# Verificar logs
docker-compose logs recommendation-service
```

### Error rate alto
- Verificar que DynamoDB/Kinesis están disponibles
- Revisar logs de la aplicación
- Aumentar thresholds si el sistema está funcionando correctamente pero más lento de lo esperado

### LocalStack es más lento que AWS real
- Normal: LocalStack es más lento que AWS
- Considerar usar AWS real para métricas de producción
- Usar LocalStack solo para desarrollo local

## Próximos Pasos

Una vez obtenidas las métricas de performance:

1. **Documentar resultados**: Guardar outputs de k6
2. **Optimizar si necesario**: Si no se cumplen los thresholds
3. **Comparar vs diseño**: Validar que cumple requisitos no funcionales
4. **Actualizar CV**: Incluir métricas reales de performance

## Referencias

- [k6 Documentation](https://k6.io/docs/)
- [k6 Examples](https://k6.io/docs/examples/)
- [Performance Testing Best Practices](https://k6.io/docs/test-types/load-testing/)
