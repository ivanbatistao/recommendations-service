# Real-Time Recommendation Service

## Objetivo

Construir un microservicio de recomendaciones en tiempo real para un escenario de ecommerce, utilizando Go, Gin, AWS Lambda, DynamoDB y Kinesis, con MiniStack para simular los servicios de AWS localmente.

El proyecto debe seguir buenas prácticas de ingeniería, Clean Architecture, procesamiento concurrente, pruebas automatizadas, pruebas de carga y documentación suficiente para defenderlo en entrevistas técnicas.

---

# Reglas del proyecto

1. No implementar funcionalidades sin entender primero el problema y el diseño.
2. Mantener el dominio independiente de Gin, AWS y cualquier framework.
3. Usar Gin como framework HTTP.
4. Usar el SDK oficial de AWS para Go.
5. Usar MiniStack para el desarrollo local de los servicios AWS.
6. Mantener el sistema ejecutable localmente durante todo el desarrollo.
7. Aplicar concurrencia únicamente donde aporte valor real.
8. No inventar métricas de rendimiento.
9. Todas las cifras de rendimiento del CV deben provenir de pruebas reproducibles.
10. Cada etapa debe terminar con código funcional, pruebas y documentación cuando corresponda.
11. No avanzar a la siguiente etapa si la etapa actual tiene errores o conceptos fundamentales sin resolver.
12. Registrar las decisiones arquitectónicas importantes y sus trade-offs.
13. Mantener los resultados de tests y benchmarks como evidencia del comportamiento real; no asumir propiedades que el código todavía no garantiza.

---

# Stack tecnológico

| Área            | Tecnología                        |
| --------------- | --------------------------------- |
| Lenguaje        | Go 1.26.5                         |
| HTTP Framework  | Gin                               |
| Arquitectura    | Clean Architecture                |
| API             | REST                              |
| Compute         | AWS Lambda                        |
| Event Streaming | AWS Kinesis                       |
| Database        | AWS DynamoDB                      |
| AWS local       | MiniStack                         |
| Concurrencia    | Goroutines, Channels, Worker Pool |
| Load Testing    | k6                                |
| Containers      | Docker                            |
| Testing         | Go testing                        |
| Race Detection  | Go Race Detector                  |
| Profiling       | pprof                             |
| SDK             | AWS SDK for Go v2                 |

---

# Arquitectura objetivo

```text
                         Event Generator
                                |
                                v
                         Kinesis Stream
                                |
                                v
                 Recommendation Processor
                                |
                     Go Worker Pool
                                |
                                v
                           DynamoDB
                                ^
                                |
                     Recommendation API
                                |
                               Gin
                                |
                                v
                             Client
```

En el despliegue AWS:

```text
Client
  |
  v
API Gateway
  |
  v
AWS Lambda
  |
  v
Gin / HTTP Adapter
  |
  v
Application Layer
  |
  v
Domain
  |
  v
DynamoDB
```

El flujo de eventos será:

```text
Ecommerce / Event Generator
            |
            v
       Kinesis Stream
            |
            v
 Recommendation Processor
            |
       Worker Pool
            |
            v
      Recommendation
          Store
            |
            v
         DynamoDB
```

---

# Roadmap completo

## Módulo 0 — Problem Definition & System Design

### Objetivo

Definir completamente qué vamos a construir antes de comenzar la implementación.

### Tareas

- [x] Definir Problem Statement.
- [x] Identificar stakeholders.
- [x] Definir alcance.
- [x] Definir requisitos funcionales.
- [x] Definir requisitos no funcionales.
- [x] Definir casos de uso.
- [x] Identificar actores.
- [x] Identificar eventos del dominio.
- [x] Identificar comandos y operaciones.
- [x] Definir entidades del dominio.
- [x] Diseñar arquitectura de alto nivel.
- [x] Diseñar flujo de eventos.
- [x] Diseñar API inicial.
- [x] Diseñar modelo inicial de datos.
- [x] Definir objetivos de latencia.
- [x] Definir estrategia inicial de escalabilidad.
- [x] Definir estrategia de concurrencia.
- [x] Definir algoritmo inicial de recomendaciones.
- [x] Registrar decisiones arquitectónicas.
- [x] Documentar trade-offs.
- [x] Definir criterios de éxito.

### Entregables

- `docs/01-problem-statement.md`
- `docs/02-requirements.md`
- `docs/03-domain.md`
- `docs/04-events.md`
- `docs/05-api-design.md`
- `docs/06-architecture.md`
- `docs/07-data-model.md`
- `docs/08-decisions.md`
- `docs/09-non-functional.md`
- `docs/10-roadmap.md`

### Estado

- [x] Cerrado

---

# Módulo 1 — Project Bootstrap

## Objetivo

Crear una base limpia, ejecutable y mantenible para el proyecto.

### Tareas

- [x] Crear repositorio.
- [x] Inicializar Go module.
- [x] Definir versión de Go — Go 1.26.5.
- [x] Crear estructura de directorios.
- [x] Instalar Gin.
- [x] Configurar aplicación.
- [x] Crear `cmd/api`.
- [x] Crear servidor HTTP.
- [x] Crear router Gin.
- [x] Crear `GET /health`.
- [x] Crear configuración por variables de entorno.
- [x] Configurar logging inicial con `log/slog`.
- [x] Crear Request ID middleware.
- [x] Configurar Recovery middleware.
- [x] Crear `.gitignore`.
- [x] Crear `Dockerfile`.
- [x] Crear `docker-compose.yml` base.
- [x] Crear README inicial.
- [x] Crear tests HTTP.
- [x] Ejecutar `go test ./...`.
- [x] Ejecutar `go test -race ./...`.
- [x] Ejecutar `go vet ./...`.
- [x] Ejecutar formatter.
- [x] Verificar que el proyecto arranca correctamente.
- [x] Verificar ejecución mediante Docker.
- [x] Verificar `GET /health` dentro del contenedor.
- [x] Verificar `X-Request-ID`.
- [x] Verificar graceful shutdown.

### Decisiones importantes

- `main.go` actúa como Composition Root.
- Gin se utiliza explícitamente mediante `gin.New()` en lugar de `gin.Default()`.
- Los middleware se registran explícitamente.
- El servidor HTTP está encapsulado en una abstracción propia.
- El logging utiliza `log/slog` con JSON.
- El Dockerfile utiliza multi-stage build.
- El builder utiliza Go 1.26.5.
- La imagen final utiliza Distroless.
- El entorno local utiliza Docker Compose.
- La concurrencia se validará posteriormente con `go test -race`.

### Entregable

API Go + Gin ejecutándose localmente y mediante Docker.

### Estado

- [x] Cerrado

---

# Módulo 2 — Domain Layer

## Objetivo

Implementar el dominio sin depender de Gin, AWS o infraestructura.

### Tareas

- [x] Definir `User`.
- [x] Definir `Product`.
- [x] Definir `InteractionEvent`.
- [x] Definir `Recommendation`.
- [x] Definir tipos de eventos.
- [x] Definir reglas del dominio.
- [x] Definir interfaces de repositorio.
- [x] Definir servicio de recomendaciones.
- [x] Definir errores de dominio.
- [x] Crear tests unitarios.
- [x] Revisar dependencias arquitectónicas.
- [x] Verificar que `internal/domain` no dependa de Gin, AWS SDK ni infraestructura.
- [x] Documentar decisiones de modelado.

### Implementación realizada

- [x] Scoring por tipo de evento.
- [x] Agregación de interés por producto.
- [x] Ranking por score descendente.
- [x] Top N.
- [x] Integración de scoring + interest + ranking en `GenerateRecommendations`.
- [x] Filtrado de eventos que no pertenecen al `userId`.
- [x] Tests para scoring, interés, ranking y servicio.
- [x] Desempate determinístico cuando dos productos tienen el mismo score.
- [x] Validar nuevamente `go test ./...`.
- [x] Validar `go test -race ./...`.
- [x] Validar `go vet ./...`.

### Decisiones de modelado

1. `userId` se recibe explícitamente en `GenerateRecommendations`.
2. `CalculateInterest` ignora eventos de otros usuarios.
3. El ranking usa score descendente.
4. Los scores iguales se resuelven determinísticamente por `productId` ascendente.
5. El algoritmo permanece en el dominio; la persistencia queda detrás de `Repository`.

### Entregable

Dominio independiente y testeable.

### Estado

- [x] Cerrado

---

# Módulo 3 — Application Layer

## Objetivo

Implementar los casos de uso de la aplicación.

### Casos de uso

- [x] Get Recommendations.
- [x] Process Event.
- [x] Generate Recommendations.
- [x] Update Recommendations (integrado en Process Event).

### Tareas

- [x] Crear comandos.
- [x] Crear queries.
- [x] Crear DTOs.
- [x] Implementar use cases.
- [x] Inyección de dependencias.
- [x] Manejo de errores.
- [x] Tests unitarios.
- [x] Tests con mocks/fakes cuando sea apropiado.

### Entregable

Casos de uso funcionando sin depender de infraestructura concreta.

### Estado

- [x] Cerrado

---

# Módulo 4 — HTTP API con Gin

## Objetivo

Exponer los casos de uso mediante una API REST usando Gin.

### Endpoints

```text
GET  /health
GET  /recommendations/:userId
POST /events
POST /recommendations/generate
```
```

### Tareas

- [x] Crear handlers.
- [x] Crear routing.
- [x] Request validation.
- [x] Response DTOs.
- [x] HTTP status codes.
- [x] Error handling.
- [x] Middleware adicional cuando corresponda.
- [x] Request ID / correlation ID.
- [x] Tests HTTP.
- [x] Integration tests.
- [x] Documentar API.

### Entregable

API REST funcional.

### Estado

- [x] Cerrado

---

# Módulo 5 — DynamoDB

## Objetivo

Implementar persistencia usando DynamoDB y MiniStack.

### Tareas

- [ ] Levantar MiniStack.
- [ ] Configurar DynamoDB local.
- [ ] Definir tablas.
- [ ] Definir partition keys.
- [ ] Definir sort keys.
- [ ] Definir índices si son necesarios.
- [ ] Implementar repository.
- [ ] Implementar reads.
- [ ] Implementar writes.
- [ ] Implementar batch operations cuando corresponda.
- [ ] Configurar AWS SDK v2.
- [ ] Tests de integración.
- [ ] Verificar funcionamiento contra MiniStack.

### Entregable

Persistencia DynamoDB funcional localmente.

### Estado

- [ ] En progreso

---

# Módulo 6 — Kinesis & Event-Driven Architecture

## Objetivo

Introducir procesamiento de eventos.

### Tareas

- [ ] Levantar Kinesis en MiniStack.
- [ ] Crear stream.
- [ ] Definir schema de eventos.
- [ ] Implementar event producer.
- [ ] Implementar event generator.
- [ ] Implementar consumer.
- [ ] Implementar deserialización.
- [ ] Validar eventos.
- [ ] Manejar eventos inválidos.
- [ ] Implementar retry strategy.
- [ ] Analizar ordering.
- [ ] Analizar partition keys.
- [ ] Documentar event flow.

### Entregable

Flujo completo:

```text
Event Generator
      |
      v
Kinesis
      |
      v
Consumer
      |
      v
Application
      |
      v
DynamoDB
```

### Estado

- [ ] En progreso

---

# Módulo 7 — Recommendation Engine

## Objetivo

Implementar el algoritmo de recomendaciones.

### Primera versión

Utilizar un algoritmo sencillo basado en co-ocurrencia/interacciones entre productos.

### Tareas

- [ ] Definir scoring.
- [ ] Definir reglas.
- [ ] Procesar `product_viewed`.
- [ ] Procesar `add_to_cart`.
- [ ] Procesar `purchase`.
- [ ] Procesar `search`.
- [ ] Actualizar relaciones entre productos.
- [ ] Generar recomendaciones.
- [ ] Persistir recomendaciones.
- [ ] Recuperar recomendaciones.
- [ ] Tests del algoritmo.
- [ ] Tests con diferentes escenarios.

### Entregable

Recommendation Engine funcional.

### Estado

- [ ] En progreso

---

# Módulo 8 — Concurrency in Go

## Objetivo

Aplicar concurrencia de manera justificada y medible.

### Conceptos

- [ ] Goroutines.
- [ ] Channels.
- [ ] Buffered channels.
- [ ] Worker Pool.
- [ ] Fan-out.
- [ ] Fan-in.
- [ ] Context.
- [ ] Cancellation.
- [ ] Timeouts.
- [ ] Synchronization.
- [ ] Mutex cuando sea necesario.
- [ ] Atomic operations cuando sean apropiadas.
- [ ] Backpressure.
- [ ] Graceful shutdown.

### Tareas

- [ ] Diseñar worker pool.
- [ ] Implementar workers.
- [ ] Distribuir eventos.
- [ ] Manejar errores concurrentes.
- [ ] Evitar data races.
- [ ] Ejecutar `go test -race`.
- [ ] Comparar procesamiento secuencial vs concurrente.
- [ ] Crear benchmarks.

### Entregable

Procesador concurrente medido y justificado.

### Estado

- [ ] En progreso

---

# Módulo 9 — AWS Lambda

## Objetivo

Preparar el servicio para ejecución serverless.

### Tareas

- [ ] Diseñar Lambda handler.
- [ ] Separar transporte HTTP de application layer.
- [ ] Adaptar Gin cuando corresponda.
- [ ] Analizar lifecycle de Lambda.
- [ ] Analizar cold starts.
- [ ] Analizar concurrency.
- [ ] Configurar variables de entorno.
- [ ] Ejecutar Lambda localmente.
- [ ] Integrar con MiniStack.
- [ ] Probar API completa.

### Entregable

Servicio ejecutándose como Lambda localmente.

### Estado

- [ ] En progreso

---

# Módulo 10 — Testing

## Objetivo

Construir una estrategia de testing profesional.

### Unit Tests

- [ ] Domain.
- [ ] Application.
- [ ] Recommendation algorithm.
- [ ] Validation.
- [ ] Error handling.

### Integration Tests

- [ ] DynamoDB.
- [ ] Kinesis.
- [ ] API.

### Concurrency Tests

- [ ] Race detector.
- [ ] Concurrent requests.
- [ ] Worker pool.

### Tareas adicionales

- [ ] Definir test fixtures.
- [ ] Definir test helpers.
- [ ] Analizar cobertura.
- [ ] Evitar tests frágiles.

### Entregable

Suite automatizada de tests.

### Estado

- [ ] En progreso

---

# Módulo 11 — Performance & Benchmarking

## Objetivo

Medir y optimizar antes de realizar pruebas de carga distribuidas.

### Tareas

- [ ] Crear Go benchmarks.
- [ ] Medir recommendation algorithm.
- [ ] Medir serialization.
- [ ] Medir persistence operations.
- [ ] Analizar allocations.
- [ ] Usar `pprof`.
- [ ] Identificar bottlenecks.
- [ ] Optimizar.
- [ ] Repetir benchmarks.
- [ ] Documentar resultados.

### Entregable

Benchmark report.

### Estado

- [ ] En progreso

---

# Módulo 12 — Load Testing

## Objetivo

Determinar el rendimiento real del sistema.

### Herramienta

k6.

### Métricas

- [ ] Requests/sec.
- [ ] Throughput.
- [ ] P50.
- [ ] P90.
- [ ] P95.
- [ ] P99.
- [ ] Error rate.
- [ ] Response time.
- [ ] Saturation.
- [ ] Resource utilization.

### Escenarios

- [ ] Low load.
- [ ] Normal load.
- [ ] High load.
- [ ] Stress test.
- [ ] Spike test.
- [ ] Soak test si el tiempo lo permite.

### Tareas

- [ ] Crear scripts k6.
- [ ] Definir escenarios reproducibles.
- [ ] Ejecutar baseline.
- [ ] Identificar bottlenecks.
- [ ] Optimizar.
- [ ] Repetir pruebas.
- [ ] Documentar resultados.
- [ ] Guardar resultados en `benchmarks/results/`.

### Regla

Nunca escribir una cifra de rendimiento en el CV sin evidencia de una prueba reproducible.

### Entregable

Load Test Report.

### Estado

- [ ] En progreso

---

# Módulo 13 — Observability

## Objetivo

Poder entender qué está ocurriendo dentro del sistema.

### Logging

- [ ] Structured logging.
- [ ] Request ID.
- [ ] Event ID.
- [ ] User ID cuando sea apropiado.
- [ ] Error context.

### Metrics

- [ ] Requests.
- [ ] Errors.
- [ ] Latency.
- [ ] Events processed.
- [ ] Events failed.
- [ ] Queue/stream lag si es medible.
- [ ] Worker utilization.

### Profiling

- [ ] pprof.
- [ ] CPU profiling.
- [ ] Memory profiling.

### Entregable

Sistema observable y documentado.

### Estado

- [ ] En progreso

---

# Módulo 14 — Resilience & Failure Handling

## Objetivo

Probar qué ocurre cuando las cosas fallan.

### Escenarios

- [ ] DynamoDB unavailable.
- [ ] Kinesis unavailable.
- [ ] Invalid event.
- [ ] Duplicate event.
- [ ] Timeout.
- [ ] Slow dependency.
- [ ] Worker failure.
- [ ] Partial processing failure.

### Tareas

- [ ] Retry.
- [ ] Backoff.
- [ ] Timeout.
- [ ] Context cancellation.
- [ ] Idempotency.
- [ ] Graceful degradation.
- [ ] Dead-letter strategy conceptual.

### Entregable

Failure handling documentado y probado.

### Estado

- [ ] En progreso

---

# Módulo 15 — Docker & Local Developer Experience

## Objetivo

Permitir que cualquier desarrollador levante el proyecto rápidamente.

### Estado actual

Parte del trabajo base ya fue realizado en el Módulo 1:

- Dockerfile base.
- Docker Compose base.
- API ejecutándose en contenedor.

### Tareas restantes

- [ ] Integrar MiniStack.
- [ ] Integrar API.
- [ ] Integrar Event Generator.
- [ ] Configuración completa.
- [ ] Scripts de inicialización.
- [ ] Seed data.
- [ ] Makefile o scripts equivalentes.
- [ ] Documentar comandos.

### Objetivo final

```bash
docker compose up
```

debe dejar el entorno listo para desarrollar.

### Entregable

Local development environment.

### Estado

- [ ] En progreso

---

# Módulo 16 — Infrastructure as Code

## Objetivo

Definir infraestructura reproducible.

### Posible tecnología

Terraform.

### Recursos

- [ ] Lambda.
- [ ] API Gateway.
- [ ] Kinesis.
- [ ] DynamoDB.
- [ ] IAM.
- [ ] CloudWatch.

### Nota

Este módulo es secundario respecto al MVP de dos semanas. Se implementará si el tiempo disponible lo permite.

### Estado

- [ ] Post-MVP

---

# Módulo 17 — AWS Deployment

## Objetivo

Desplegar una versión real en AWS.

### Tareas

- [ ] Configurar AWS account.
- [ ] IAM.
- [ ] Provisionar infraestructura.
- [ ] Deploy Lambda.
- [ ] Deploy API Gateway.
- [ ] Deploy Kinesis.
- [ ] Deploy DynamoDB.
- [ ] Configurar environment variables.
- [ ] Probar endpoint.
- [ ] Ejecutar smoke tests.
- [ ] Medir latencia real.

### Entregable

Versión desplegada en AWS.

### Estado

- [ ] Post-MVP

---

# Módulo 18 — CI/CD

## Objetivo

Automatizar calidad y despliegue.

### Tareas

- [ ] GitHub Actions.
- [ ] Formatting.
- [ ] Linting.
- [ ] Unit tests.
- [ ] Integration tests.
- [ ] Race detector.
- [ ] Build.
- [ ] Docker build.
- [ ] Deployment pipeline.

### Entregable

CI/CD pipeline.

### Estado

- [ ] Post-MVP

---

# Módulo 19 — Documentation & Architecture Review

## Objetivo

Convertir el proyecto en un repositorio profesional.

### Tareas

- [ ] README final.
- [ ] Architecture diagram.
- [ ] Sequence diagrams.
- [ ] Event flow diagram.
- [ ] Data model.
- [ ] API documentation.
- [ ] ADRs.
- [ ] Performance report.
- [ ] Failure scenarios.
- [ ] Local setup.
- [ ] AWS setup.
- [ ] Trade-offs.
- [ ] Future improvements.

### Entregable

Repositorio presentable profesionalmente.

### Estado

- [ ] En progreso

---

# Módulo 20 — Interview Preparation

## Objetivo

Ser capaz de explicar y defender el proyecto.

### System Design

- [ ] Explicar arquitectura.
- [ ] Explicar flujo de eventos.
- [ ] Explicar DynamoDB.
- [ ] Explicar Kinesis.
- [ ] Explicar Lambda.
- [ ] Explicar escalabilidad.
- [ ] Explicar concurrencia.
- [ ] Explicar trade-offs.
- [ ] Explicar failure handling.
- [ ] Explicar performance.

### Go

- [ ] Goroutines.
- [ ] Channels.
- [ ] Worker pools.
- [ ] Context.
- [ ] Interfaces.
- [ ] Error handling.
- [ ] Memory.
- [ ] Race conditions.

### Performance

- [ ] Explicar benchmarks.
- [ ] Explicar load tests.
- [ ] Explicar bottlenecks.
- [ ] Explicar optimizaciones.

### Entregables

- [ ] Pitch de 30 segundos.
- [ ] Pitch de 60 segundos.
- [ ] Explicación técnica de 5 minutos.
- [ ] Preguntas y respuestas de entrevista.
- [ ] STAR stories.

### Estado

- [ ] En progreso

---

# Módulo 21 — CV & LinkedIn

## Objetivo

Convertir el trabajo realizado en evidencia profesional.

### Tareas

- [ ] Redactar proyecto para CV.
- [ ] Redactar bullets orientados a impacto.
- [ ] Incorporar métricas reales.
- [ ] Redactar descripción para LinkedIn.
- [ ] Preparar GitHub README.
- [ ] Preparar portfolio description.

### Regla

No utilizar cifras ficticias.

Las métricas deberán corresponder a resultados medidos y documentados.

### Estado

- [ ] En progreso

---

# Definición de Done

El MVP de dos semanas estará terminado cuando:

- [ ] API Go + Gin funcionando.
- [ ] Clean Architecture implementada.
- [ ] Dominio desacoplado de infraestructura.
- [ ] DynamoDB funcionando mediante MiniStack.
- [ ] Kinesis funcionando mediante MiniStack.
- [ ] Event Generator funcionando.
- [ ] Recommendation Engine funcionando.
- [ ] Procesamiento concurrente funcionando.
- [ ] Lambda funcionando localmente.
- [ ] Tests unitarios funcionando.
- [ ] Integration tests críticos funcionando.
- [ ] Race detector ejecutándose correctamente.
- [ ] Benchmarks disponibles.
- [ ] Load tests con k6 disponibles.
- [ ] Métricas reales documentadas.
- [ ] README completo.
- [ ] Arquitectura documentada.
- [ ] Proyecto explicable en una entrevista.

---

# MVP de 2 semanas

## Prioridad P0 — Obligatorio

- Módulo 0 — Diseño.
- Módulo 1 — Bootstrap.
- Módulo 2 — Domain.
- Módulo 3 — Application.
- Módulo 4 — API + Gin.
- Módulo 5 — DynamoDB.
- Módulo 6 — Kinesis.
- Módulo 7 — Recommendation Engine.
- Módulo 8 — Concurrency.
- Módulo 9 — Lambda.
- Módulo 10 — Testing.
- Módulo 11 — Benchmarking.
- Módulo 12 — Load Testing.
- Módulo 19 — Documentation.
- Módulo 20 — Interview Preparation.
- Módulo 21 — CV.

## Prioridad P1 — Si hay tiempo

- Módulo 13 — Observability.
- Módulo 14 — Resilience.
- Módulo 15 — Developer Experience.

## Prioridad P2 — Post-MVP

- Módulo 16 — Terraform.
- Módulo 17 — AWS Deployment.
- Módulo 18 — CI/CD.

---

# Métricas objetivo

Estas son metas iniciales, no resultados garantizados.

## Latencia

Target:

```text
P95 < 200 ms
```

## Throughput

Target:

```text
High concurrency
```

El throughput máximo será determinado mediante pruebas de carga.

No asumir 500K RPS antes de medirlo.

## Disponibilidad

Target conceptual:

```text
99.9%
```

La disponibilidad real se evaluará únicamente cuando exista un despliegue AWS apropiado.

---

# Performance Validation

El ciclo de optimización será:

```text
Implement
   |
   v
Benchmark
   |
   v
Profile
   |
   v
Identify bottleneck
   |
   v
Optimize
   |
   v
Benchmark again
   |
   v
Load test
   |
   v
Document
```

---

# Estructura final esperada

```text
recommendation-service/
│
├── cmd/
│   ├── api/
│   └── event-generator/
│
├── internal/
│   ├── domain/
│   │   ├── event/
│   │   ├── product/
│   │   ├── recommendation/
│   │   └── user/
│   │
│   ├── application/
│   │   ├── commands/
│   │   ├── queries/
│   │   └── dto/
│   │
│   └── infrastructure/
│       ├── http/
│       │   └── gin/
│       ├── dynamodb/
│       ├── kinesis/
│       └── logger/
│
├── configs/
├── deploy/
├── scripts/
├── tests/
├── benchmarks/
├── docs/
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

# Estado global del proyecto

Actualiza esta sección al terminar cada etapa.

## Progreso

```text
Módulo actual: 3
Módulos completados: 3 / 21
```

### Estado actual

El Módulo 2 — Domain Layer está **cerrado**.

Las validaciones finales fueron satisfactorias:

```text
go vet ./...       → pasa
go test ./...      → pasa
go test -race ./... → pasa
```

El ranking ahora es determinístico:

```text
1. Score descendente.
2. En caso de empate, ProductID ascendente.
```

La lógica de recomendaciones está desacoplada de Gin, AWS SDK e infraestructura.

## Módulos cerrados

- [x] Módulo 0
- [x] Módulo 1
- [x] Módulo 2
- [ ] Módulo 3
- [ ] Módulo 4
- [ ] Módulo 5
- [ ] Módulo 6
- [ ] Módulo 7
- [ ] Módulo 8
- [ ] Módulo 9
- [ ] Módulo 10
- [ ] Módulo 11
- [ ] Módulo 12
- [ ] Módulo 13
- [ ] Módulo 14
- [ ] Módulo 15
- [ ] Módulo 16
- [ ] Módulo 17
- [ ] Módulo 18
- [ ] Módulo 19
- [ ] Módulo 20
- [ ] Módulo 21

## Último módulo completado

```text
Módulo 2 — Domain Layer
```

## Próximo módulo

```text
Módulo 3 — Application Layer
```

### Objetivo inmediato

Implementar los casos de uso:

```text
ProcessEvent
GenerateRecommendations
GetRecommendations
UpdateRecommendations
```

manteniendo los casos de uso independientes de Gin, DynamoDB, Kinesis y AWS.

---

# Cómo usar este archivo durante el proyecto

En cada sesión:

1. Compartir este archivo actualizado.
2. Indicar qué módulo(s) se cerraron.
3. Revisar el estado actual.
4. Continuar únicamente desde el siguiente punto pendiente.
5. No repetir módulos ya cerrados salvo que sea necesario corregir un problema.
6. Mantener este roadmap como fuente de verdad del progreso.

Cuando una etapa esté completamente terminada, marcar:

```text
- [x] Cerrado
```

y actualizar:

```text
Módulo actual: X
Módulos completados: X / 21
```
