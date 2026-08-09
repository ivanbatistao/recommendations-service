# Worker Pool - Decisiones de Implementación

## Módulo 7 — Worker Pool con Go Concurrency

### Objetivo

Implementar procesamiento concurrente de eventos usando un worker pool con goroutines y channels.

### ¿Qué es un Worker Pool?

Un **Worker Pool** es un patrón de concurrencia que:

- **Crea un número fijo de workers** (goroutines)
- **Usa un channel** como cola de tareas
- **Los workers procesan tareas** en paralelo
- **Controla el uso de recursos** limitando la concurrencia

### ¿Por qué necesitamos un Worker Pool?

**Problema:**
- Kinesis puede enviar muchos eventos rápidamente
- Si procesamos cada evento en una goroutine nueva, podemos:
  - Agotar memoria
  - Sobrecargar la base de datos
  - Causar rate limiting

**Solución:**
- Worker pool con número fijo de workers
- Control de concurrencia
- Procesamiento eficiente sin sobrecarga

### Arquitectura del Worker Pool

```text
Kinesis Consumer
      |
      v
  Event Channel (Buffered)
      |
      v
  Worker Pool (N workers)
      |
  +---+---+---+---+
  |   |   |   |   |  (Goroutines)
  +---+---+---+---+
      |
      v
Recommendation Service
      |
      v
   DynamoDB
```

### Implementación del Worker Pool

#### WorkerPool

**Ubicación**: `internal/infrastructure/processing/workerpool/worker_pool.go`

**Estructura:**
```go
type WorkerPool struct {
    numWorkers   int                    // Número de workers
    eventChan    chan event.Event       // Channel de eventos (buffered)
    processor    EventProcessor         // Interfaz para procesar eventos
    wg           sync.WaitGroup         // Para esperar que terminen los workers
    logger       *slog.Logger          // Logger
    ctx          context.Context       // Contexto para cancelación
    cancel       context.CancelFunc     // Función para cancelar
}
```

**Componentes clave:**
- **Channel buffered**: Cola de eventos con capacidad limitada
- **NumWorkers**: Número fijo de goroutines workers
- **EventProcessor**: Interfaz para desacoplar lógica de procesamiento
- **Context**: Para graceful shutdown

#### EventProcessor Interface

```go
type EventProcessor interface {
    ProcessEvent(ctx context.Context, ev event.Event) error
}
```

**¿Por qué una interfaz?**
- **Desacoplamiento**: El worker pool no sabe cómo procesar eventos
- **Testabilidad**: Podemos usar mocks en tests
- **Flexibilidad**: Podemos cambiar la implementación fácilmente

#### Métodos del Worker Pool

##### 1. Start - Iniciar Workers
```go
func (wp *WorkerPool) Start() {
    for i := 0; i < wp.numWorkers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}
```

**¿Qué hace?**
- Crea `numWorkers` goroutines
- Cada goroutine ejecuta `worker(i)`
- Usa `WaitGroup` para trackear workers activos

##### 2. Worker - Goroutine Individual
```go
func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case <-wp.ctx.Done():
            return  // Cerrar worker
        case ev := <-wp.eventChan:
            wp.processor.ProcessEvent(wp.ctx, ev)  // Procesar evento
        }
    }
}
```

**Patrón clásico de worker:**
- **Loop infinito**: El worker siempre está escuchando
- **Select**: Espera either contexto cancelado o nuevo evento
- **Context.Done()**: Graceful shutdown
- **Channel receive**: Procesa evento cuando llega

##### 3. Submit - Enviar Eventos
```go
func (wp *WorkerPool) Submit(ev event.Event) {
    select {
    case wp.eventChan <- ev:
        // Evento enviado exitosamente
    default:
        // Channel lleno, drop evento
    }
}
```

**¿Por qué `select` con `default`?**
- **Non-blocking**: No espera si el channel está lleno
- **Backpressure**: Si el channel está lleno, dropea eventos
- **Trade-off**: Mejor dropear eventos que bloquear al productor

**Alternativa:** Podríamos usar blocking send, pero eso podría saturar al productor.

##### 4. Stop - Graceful Shutdown
```go
func (wp *WorkerPool) Stop() {
    wp.cancel()           // Cancelar contexto
    close(wp.eventChan)   // Cerrar channel
    wp.wg.Wait()         // Esperar que terminen workers
}
```

**Graceful shutdown:**
1. **Cancel context**: Workers reciben señal de parar
2. **Close channel**: Desbloquea workers esperando en channel
3. **WaitGroup**: Espera que todos los workers terminen

##### 5. Stats - Métricas del Pool
```go
func (wp *WorkerPool) Stats() map[string]interface{} {
    return map[string]interface{}{
        "num_workers":  wp.numWorkers,
        "buffer_size":  cap(wp.eventChan),
        "buffer_used":  len(wp.eventChan),
    }
}
```

**Métricas disponibles:**
- **num_workers**: Número de workers activos
- **buffer_size**: Capacidad del channel
- **buffer_used**: Eventos actualmente en el channel

### RecommendationProcessorAdapter

**Ubicación**: `internal/infrastructure/processing/workerpool/processor_adapter.go`

**Propósito:**
- Conectar el `recommendation.Service` con el `WorkerPool`
- Implementar la interfaz `EventProcessor`

**Implementación:**
```go
type RecommendationProcessorAdapter struct {
    service *recommendation.Service
}

func (a *RecommendationProcessorAdapter) ProcessEvent(
    ctx context.Context,
    ev event.Event,
) error {
    return a.service.ProcessEvent(ctx, ev)
}
```

**¿Por qué un adapter?**
- **Interface segregation**: El worker pool solo conoce `EventProcessor`
- **Single responsibility**: El adapter conecta dos capas
- **Testabilidad**: Podemos mockear `EventProcessor` fácilmente

### Configuración del Worker Pool

#### Parámetros Configurables

```go
NewWorkerPool(
    numWorkers int,           // Número de workers (goroutines)
    bufferSize int,            // Tamaño del channel buffer
    processor EventProcessor,  // Procesador de eventos
    logger *slog.Logger,       // Logger
)
```

**Recomendaciones:**
- **numWorkers**: Basado en número de CPUs o límites de DB
- **bufferSize**: 2-3x numWorkers para evitar bloqueos
- **Ejemplo**: 4 workers, buffer de 10-20 eventos

### Tests Implementados

#### 1. TestWorkerPool_StartAndStop
- Verifica que el pool inicia y detiene workers correctamente
- Verifica graceful shutdown

#### 2. TestWorkerPool_SubmitAndProcess
- Verifica que los eventos se procesan correctamente
- Verifica que todos los eventos sean procesados

#### 3. TestWorkerPool_ConcurrentProcessing
- Verifica procesamiento concurrente con múltiples workers
- Verifica que más eventos que workers se procesen correctamente

#### 4. TestWorkerPool_Stats
- Verifica que las estadísticas sean correctas
- Verifica num_workers y buffer_size

#### 5. TestWorkerPool_Backpressure
- Verifica comportamiento cuando el channel está lleno
- Verifica que eventos sean dropeados cuando no hay capacidad

### Trade-offs Considerados

1. **Non-blocking vs Blocking Submit**
   - **Elección**: Non-blocking con backpressure
   - **Trade-off**: Pérdida de eventos vs saturación del productor
   - **Decisión**: Mejor dropear eventos que bloquear al productor

2. **Fixed vs Dynamic Workers**
   - **Elección**: Número fijo de workers
   - **Trade-off**: Simplicidad vs auto-scaling
   - **Decisión**: Número fijo es más simple y predecible

3. **Buffer Size**
   - **Elección**: 2-3x numWorkers
   - **Trade-off**: Memoria vs throughput
   - **Decisión**: Balance razonable para la mayoría de casos

4. **Graceful Shutdown vs Immediate**
   - **Elección**: Graceful shutdown con context
   - **Trade-off**: Latencia vs completitud
   - **Decisión**: Completitud es más importante para eventos

### Validaciones Realizadas

- [x] Código compila sin errores
- [x] Todos los tests del worker pool pasan
- [x] Tests de concurrencia pasan
- [x] Mock processor para testing
- [x] Graceful shutdown implementado
- [x] Backpressure handling implementado
- [x] Stats function implementada
- [x] Documentación de decisiones

### Estado del Módulo 7

- [x] Cerrado
