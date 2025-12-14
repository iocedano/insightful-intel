# Guía de Pipeline Dinámico

Esta guía explica el nuevo sistema de pipeline dinámico que crea automáticamente pipelines de búsqueda basados en `GetSearchableKeywordCategories()` de los conectores de dominio.

## Resumen

El sistema de pipeline dinámico automáticamente:
1. **Descubre categorías buscables** desde cada conector de dominio
2. **Extrae palabras clave** de los resultados de búsqueda
3. **Crea nuevos pasos de búsqueda** basados en palabras clave extraídas
4. **Ejecuta búsquedas** a través de múltiples dominios en paralelo
5. **Rastrea el progreso** y previene búsquedas duplicadas

## Componentes Clave

### 1. DynamicPipelineConfig
```go
type DynamicPipelineConfig struct {
    MaxDepth           int  // Profundidad máxima del pipeline (por defecto: 5)
    MaxConcurrentSteps int  // Máximo de pasos concurrentes (por defecto: 10)
    DelayBetweenSteps  int  // Retraso entre pasos en segundos (por defecto: 2)
    SkipDuplicates     bool // Omitir búsquedas de palabras clave duplicadas (por defecto: true)
}
```

### 2. DynamicPipelineStep
```go
type DynamicPipelineStep struct {
    DomainType          DomainType
    SearchParameter     string
    Category            KeywordCategory
    Keywords            []string
    Success             bool
    Error               error
    Output              any
    KeywordsPerCategory map[KeywordCategory][]string
    Depth               int
}
```

### 3. DynamicPipelineResult
```go
type DynamicPipelineResult struct {
    Steps           []DynamicPipelineStep
    TotalSteps      int
    SuccessfulSteps int
    FailedSteps     int
    MaxDepthReached int
    Config          DynamicPipelineConfig
}
```

## Cómo Funciona

### 1. Búsqueda Inicial
El pipeline comienza con una consulta de búsqueda inicial, típicamente usando ONAPI ya que proporciona los datos más completos.

### 2. Extracción de Palabras Clave
Después de cada búsqueda exitosa, el sistema extrae palabras clave de los resultados usando `GetDataByCategory()` para cada `KeywordCategory` retornado por `GetFoundKeywordCategories()`.

### 3. Coincidencia de Categorías
Para cada palabra clave extraída, el sistema encuentra dominios que pueden buscar esa categoría usando `GetSearchableKeywordCategories()`.

### 4. Generación de Pasos
Se crean nuevos pasos de pipeline para cada combinación válida de palabra clave-dominio.

### 5. Ejecución
Los pasos se ejecutan en paralelo cuando es posible, con retrasos configurables y límites de concurrencia.

## Ejemplos de Uso

### Uso Básico
```go
query := "Novasco"
availableDomains := []domain.DomainType{
    domain.DomainTypeONAPI,
    domain.DomainTypeSCJ,
    domain.DomainTypeDGII,
}

config := domain.DefaultDynamicPipelineConfig()
result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

### Configuración Personalizada
```go
config := domain.DynamicPipelineConfig{
    MaxDepth:           3,
    MaxConcurrentSteps: 5,
    DelayBetweenSteps:  1,
    SkipDuplicates:     true,
}

result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

### Creación Paso a Paso
```go
// Crear estructura de pipeline primero
pipeline, err := domain.CreateDynamicPipeline(query, availableDomains, config)

// Luego ejecutarlo
result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

## Endpoints de API HTTP

### Endpoint de Pipeline Dinámico
```
GET /dynamic?q=Novasco&depth=3&skip_duplicates=true
```

**Parámetros:**
- `q`: Consulta de búsqueda (requerido)
- `depth`: Profundidad máxima del pipeline (opcional, por defecto: 3, máximo: 10)
- `skip_duplicates`: Omitir búsquedas duplicadas (opcional, por defecto: true)

**Respuesta:**
```json
{
  "dynamic_result": {
    "steps": [...],
    "total_steps": 15,
    "successful_steps": 12,
    "failed_steps": 3,
    "max_depth_reached": 3,
    "config": {...}
  },
  "pipeline": [...],
  "summary": {
    "total_steps": 15,
    "successful_steps": 12,
    "failed_steps": 3,
    "max_depth_reached": 3
  }
}
```

## Integración de Dominio

### Métodos Requeridos
Cada conector de dominio debe implementar:
```go
type GenericConnector[T any] interface {
    GetDomainType() DomainType
    GetSearchableKeywordCategories() []KeywordCategory
    GetFoundKeywordCategories() []KeywordCategory
    GetDataByCategory(data T, category KeywordCategory) []string
    // ... otros métodos
}
```

### Flujo de Categorías
1. **ONAPI** busca `company_name` → recupera `company_name`, `person_name`, `address`
2. **SCJ** puede buscar `person_name`, `company_name` → recupera `person_name`
3. **DGII** puede buscar `company_name` → recupera `company_name`, `contributor_id`

## Opciones de Configuración

### MaxDepth
Controla qué tan profundo puede ir el pipeline. Valores más altos significan búsquedas más completas pero tiempos de ejecución más largos.

### MaxConcurrentSteps
Limita el número de búsquedas concurrentes para evitar sobrecargar APIs externas.

### DelayBetweenSteps
Agrega retrasos entre pasos para ser respetuoso con las APIs externas.

### SkipDuplicates
Previene buscar la misma palabra clave múltiples veces a través de dominios.

## Mejores Prácticas

### 1. Comenzar con ONAPI
ONAPI típicamente proporciona los datos iniciales más completos, haciéndolo ideal como punto de partida.

### 2. Configurar Profundidad Apropiadamente
- **Profundidad 1-2**: Búsquedas rápidas, buenas para consultas simples
- **Profundidad 3-4**: Búsquedas completas, buenas para investigaciones complejas
- **Profundidad 5+**: Búsquedas muy exhaustivas, usar con moderación

### 3. Usar Retrasos Apropiados
- **1-2 segundos**: Para desarrollo/pruebas
- **2-5 segundos**: Para uso en producción
- **5+ segundos**: Para uso muy respetuoso de API

### 4. Monitorear Resultados
Revisar la sección `Summary` para entender el rendimiento del pipeline y ajustar la configuración en consecuencia.

## Manejo de Errores

El sistema maneja errores de manera elegante:
- Los pasos fallidos se marcan pero no detienen el pipeline
- Los errores de red se registran pero no crashean el sistema
- Los dominios inválidos se omiten automáticamente

## Consideraciones de Rendimiento

### Uso de Memoria
- Cada paso almacena su salida completa en memoria
- Considerar `MaxDepth` y `MaxConcurrentSteps` para uso de memoria
- Los conjuntos de resultados grandes pueden requerir paginación

### Límites de Tasa de API
- Usar `DelayBetweenSteps` para respetar límites de tasa
- Monitorear `FailedSteps` para indicadores de límite de tasa
- Considerar implementar backoff exponencial para producción

### Ejecución Paralela
- Los pasos en la misma profundidad pueden ejecutarse en paralelo
- Los pasos en diferentes profundidades se ejecutan secuencialmente
- `MaxConcurrentSteps` controla la ejecución paralela

## Solución de Problemas

### Problemas Comunes

1. **Sin resultados del pipeline dinámico**
   - Verificar si la consulta inicial retorna resultados
   - Verificar que los conectores de dominio estén funcionando
   - Revisar implementación de `GetSearchableKeywordCategories()`

2. **Demasiadas búsquedas duplicadas**
   - Habilitar `SkipDuplicates`
   - Revisar lógica de extracción de palabras clave
   - Verificar implementación de `GetDataByCategory()`

3. **Pipeline se detiene temprano**
   - Revisar configuración de `MaxDepth`
   - Verificar que todos los dominios estén disponibles
   - Revisar errores en ejecución de pasos

4. **Problemas de rendimiento**
   - Reducir `MaxDepth` o `MaxConcurrentSteps`
   - Aumentar `DelayBetweenSteps`
   - Revisar conectividad de red

### Depuración
Habilitar logging de depuración para ver:
- Proceso de creación de pasos
- Resultados de extracción de palabras clave
- Decisiones de coincidencia de dominio
- Progreso de ejecución

## Migración desde Pipeline Estático

### Antes (Estático)
```go
// Creación manual de pasos
onapi := domain.NewOnapiDomain()
entities, err := onapi.SearchComercialName("Novasco")

// Extracción manual de palabras clave
keywords := domain.GetCategoryByKeywords(&domain.Onapi{}, entities)

// Ejecución manual de pasos
for _, keyword := range keywords["person_name"] {
    scj := domain.NewScjDomain()
    cases, err := scj.Search(keyword)
    // ...
}
```

### Después (Dinámico)
```go
// Creación y ejecución automática de pipeline
config := domain.DefaultDynamicPipelineConfig()
result, err := domain.ExecuteDynamicPipeline("Novasco", availableDomains, config)
```

El sistema de pipeline dinámico maneja automáticamente todos los pasos manuales, haciendo tu código más limpio y mantenible.
