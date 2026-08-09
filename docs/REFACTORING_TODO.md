# Refactorizaciones Pendientes

Este documento contiene refactorizaciones identificadas pero no implementadas todavía, priorizadas por impacto y complejidad.

## Prioridad Alta 🔴

### 1. Mover Composition Root a `internal/app/composition/`

**Estado:** Identificado pero no implementado

**Ubicación actual:** `internal/infrastructure/composition/root.go`

**Ubicación sugerida:** `internal/app/composition/root.go`

**Razón:**
- Composition root es *application assembly*, no infraestructura AWS
- Más semánticamente correcto que "infrastructure"
- Sigue mejor el patrón de Clean Architecture

**Impacto:** 
- Cambio de ubicación, sin cambio funcional
- Mejor claridad arquitectónica
- Requiere actualizar imports en cmd/api y cmd/lambda

**Complejidad:** Baja (mover archivo y actualizar imports)

---

## Prioridad Media 🟡

### 2. Mejorar Routing de Lambda Handler

**Estado:** Funcional pero puede mejorarse

**Ubicación:** `cmd/lambda/main.go` - HandleRequest method

**Problema actual:**
- Switch/case verboso para routing
- Código duplicado de respuesta HTTP
- No escalable fácilmente

**Sugerencia 1 - Router Pattern:**
```go
type Route struct {
    Path    string
    Method  string
    Handler func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

routes := []Route{
    {"/health", "GET", h.handleHealth},
    {"/recommendations/{userId}", "GET", h.handleGetRecommendations},
    // ...
}
```

**Sugerencia 2 - Map de Handlers:**
```go
handlers := map[string]func(ctx context.Context, body string) (events.APIGatewayProxyResponse, error){
    "GET_/health": h.handleHealth,
    "GET_/recommendations/{userId}": h.handleGetRecommendations,
    // ...
}
```

**Impacto:** Mejor maintainability si crece el número de rutas

**Complejidad:** Media (refactorizar estructura de routing)

---

### 3. Agregar Helper para Respuestas HTTP

**Estado:** Código duplicado en Lambda handlers

**Ubicación:** `cmd/lambda/main.go` - múltiples métodos

**Problema actual:**
```go
return events.APIGatewayProxyResponse{
    StatusCode: 200,
    Headers: map[string]string{"Content-Type": "application/json"},
    Body: body,
}, nil
```

**Sugerencia:**
```go
func jsonResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {
    data, _ := json.Marshal(body)
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: string(data),
    }, nil
}

func jsonError(statusCode int, message string) events.APIGatewayProxyResponse {
    return jsonResponse(statusCode, map[string]string{"error": message})
}
```

**Impacto:** Reducción de código duplicado

**Complejidad:** Baja (agregar helpers)

---

## Prioridad Baja 🟢

### 4. Extraer Validación de Request

**Estado:** Validación inline en cada handler

**Ubicación:** `cmd/lambda/main.go` - handleGetRecommendations, handleProcessEvent, etc.

**Problema actual:**
```go
if userID == "" {
    return events.APIGatewayProxyResponse{
        StatusCode: 400,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: `{"error":"user_id is required"}`,
    }, nil
}
```

**Sugerencia:**
```go
func validateUserID(userID string) error {
    if userID == "" {
        return errors.New("user_id is required")
    }
    return nil
}
```

**Impacto:** Mejor testability y reutilización

**Complejidad:** Baja (extraer funciones de validación)

---

### 5. Implementar Error Handling Centralizado

**Estado:** Error handling disperso en cada handler

**Ubicación:** Múltiples handlers en cmd/lambda/main.go

**Problema actual:**
```go
if err != nil {
    h.Logger.Error("failed to get recommendations", slog.String("error", err.Error()))
    return events.APIGatewayProxyResponse{
        StatusCode: 500,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: `{"error":"internal server error"}`,
    }, nil
}
```

**Sugerencia:**
```go
func (h *LambdaHandler) handleError(err error, context string) events.APIGatewayProxyResponse {
    h.Logger.Error(context, slog.String("error", err.Error()))
    
    var statusCode int
    var message string
    
    switch {
    case errors.Is(err, recommendation.ErrInvalidUserID):
        statusCode = 400
        message = "invalid user_id"
    default:
        statusCode = 500
        message = "internal server error"
    }
    
    return jsonError(statusCode, message)
}
```

**Impacto:** Mejor consistencia de respuestas de error

**Complejidad:** Media (requiere análisis de tipos de error)

---

### 6. Agregar Context Propagation

**Estado:** Context no se propaga completamente

**Ubicación:** Varios handlers en cmd/lambda/main.go

**Problema actual:**
- Request ID no se propaga desde API Gateway
- No hay tracing distribuido

**Sugerencia:**
```go
func (h *LambdaHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Extract request ID from headers
    requestID := req.Headers["X-Request-ID"]
    if requestID == "" {
        requestID = generateRequestID()
    }
    
    ctx = context.WithValue(ctx, "requestID", requestID)
    ctx = h.Logger.With("request_id", requestID)
    
    // Resto del código con ctx enriquecido
}
```

**Impacto:** Mejor observabilidad y debugging

**Complejidad:** Media (requiere agregar tracing)

---

## Detalles de Implementación Menores ⚪

### 7. Eliminar Código Dead

**Estado:** Posible código no utilizado después de refactorizaciones

**Lugares a revisar:**
- Verificar si hay imports no utilizados
- Verificar si hay funciones no utilizadas después de refactorizaciones
- Limpiar tests que ya no aplican

**Impacto:** Limpieza de código

**Complejidad:** Baja (análisis y limpieza)

---

### 8. Mejorar Naming de Variables

**Estado:** Algunos nombres podrían ser más descriptivos

**Ejemplos:**
- `eventList` podría ser `events` (ya se usa en domain)
- `app` podría ser `application` en algunos contextos

**Impacto:** Mejor legibilidad

**Complejidad:** Baja (renombrar variables)

---

## Decisiones Pendientes

### 9. Usar gin-gonic/lambdaadapter

**Estado:** No investigado

**Descripción:** Evaluar si usar `github.com/awslabs/aws-lambda-go-api-proxy` para integrar Gin con Lambda

**Ventajas:**
- Reutilizar código Gin exacto
- No necesitar adaptador manual
- Proven solution

**Desventajas:**
- Dependencia adicional
- Possible overhead
- Less control sobre request/response

**Impacto:** Simplificar código Lambda

**Complejidad:** Media (evaluar librería y refactorizar)

---

## Criterios para Priorización

**Alta Prioridad 🔴:**
- Impacto arquitectónico significativo
- Mejora maintainability a largo plazo
- Baja complejidad de implementación

**Media Prioridad 🟡:**
- Mejora code quality significativamente
- Reducción de duplicación
- Complejidad media

**Baja Prioridad 🟢:**
- Mejoras cosméticas
- Optimizaciones menores
- Nice-to-have features

---

## Proceso para Implementar Refactorizaciones

1. **Crear rama:** `refactor/nombre-refactorizacion`
2. **Implementar:** Hacer los cambios propuestos
3. **Tests:** Asegurar que todos los tests pasan
4. **Documentar:** Actualizar documentación relevante
5. **Commit:** Con mensaje descriptivo
6. **Merge:** A main después de revisión

---

## Notas

- Este documento es vivo - agregar nuevas refactorizaciones según se identifiquen
- Antes de implementar, evaluar si realmente agrega valor
- Considerar el costo/beneficio de cada refactorización
- Priorizar funcionalidad sobre perfección del código
