# Documentación de la Capa de Repositorio

Este directorio contiene la implementación de la capa de repositorio para la aplicación Insightful Intel. La capa de repositorio proporciona una abstracción limpia sobre las operaciones de persistencia de datos para todos los tipos de dominio y respuestas de pipeline.

## Arquitectura

La capa de repositorio sigue el patrón Repository y proporciona:

- **BaseRepository[T]**: Operaciones CRUD comunes para cualquier tipo de entidad
- **SearchableRepository[T]**: Extiende BaseRepository con capacidades de búsqueda
- **DomainRepository[T]**: Operaciones específicas de dominio para entidades de dominio
- **PipelineRepository[T]**: Operaciones especializadas para resultados de pipeline

## Repositorios de Dominio

### Repositorio ONAPI (`onapi_repository.go`)
- **Tipo de Entidad**: `domain.Entity`
- **Tabla de Base de Datos**: `onapi_entities`
- **Características Clave**:
  - Almacena información de marcas y patentes
  - Soporta búsqueda por nombre de empresa, nombre de persona y dirección
  - Maneja campos JSON para imágenes y listas de clases

### Repositorio SCJ (`scj_repository.go`)
- **Tipo de Entidad**: `domain.ScjCase`
- **Tabla de Base de Datos**: `scj_cases`
- **Características Clave**:
  - Almacena información de casos de la Suprema Corte de Justicia
  - Soporta búsqueda por número de caso, número de sentencia y partes involucradas
  - Rastrea metadatos de casos y URLs de documentos

### Repositorio DGII (`dgii_repository.go`)
- **Tipo de Entidad**: `domain.Register`
- **Tabla de Base de Datos**: `dgii_registers`
- **Características Clave**:
  - Almacena información de registro de la autoridad fiscal dominicana
  - Soporta búsqueda por RNC, nombre de empresa y nombre comercial
  - Rastrea estado de cumplimiento fiscal

### Repositorio PGR (`pgr_repository.go`)
- **Tipo de Entidad**: `domain.PGRNews`
- **Tabla de Base de Datos**: `pgr_news`
- **Características Clave**:
  - Almacena elementos de noticias de la Procuraduría General de la República
  - Soporta búsqueda por título y URL
  - Estructura mínima para artículos de noticias

### Repositorio Google Docking (`docking_repository.go`)
- **Tipo de Entidad**: `domain.GoogleDorkingResult`
- **Tabla de Base de Datos**: `google_docking_results`
- **Características Clave**:
  - Almacena resultados de búsqueda de Google con puntuación de relevancia
  - Soporta búsqueda por título, descripción y URL
  - Maneja campos JSON para palabras clave
  - Incluye información de clasificación y relevancia

## Repositorio de Pipeline

### Repositorio de Pipeline (`pipeline_repository.go`)
- **Tipos de Entidad**: `domain.DomainSearchResult`, `module.DynamicPipelineResult`
- **Tablas de Base de Datos**: `domain_search_results`, `dynamic_pipeline_results`, `dynamic_pipeline_steps`
- **Características Clave**:
  - Almacena tanto resultados de búsqueda de dominio individuales como resultados completos de pipeline
  - Soporta seguimiento complejo de pasos de pipeline
  - Maneja serialización JSON para estructuras de datos complejas
  - Proporciona capacidades de agregación para estadísticas de pipeline

## Factory de Repositorio

### RepositoryFactory (`factory.go`)
Proporciona instanciación centralizada de repositorios:

```go
factory := repositories.NewRepositoryFactory(db)

// Obtener repositorios específicos
onapiRepo := factory.GetOnapiRepository()
scjRepo := factory.GetScjRepository()
pipelineRepo := factory.GetPipelineRepository()

// Obtener todos los repositorios de dominio
allRepos := factory.GetAllDomainRepositories()
```

## Esquema de Base de Datos

El archivo `schema.sql` contiene todas las definiciones de tabla necesarias con:
- Indexación apropiada para rendimiento de búsqueda
- Restricciones de clave foránea donde corresponda
- Columnas JSON para estructuras de datos complejas
- Seguimiento de timestamps para propósitos de auditoría

## Ejemplos de Uso

### Operaciones CRUD Básicas
```go
// Crear
entity := domain.Entity{...}
err := onapiRepo.Create(ctx, entity)

// Leer
entity, err := onapiRepo.GetByID(ctx, "123")

// Actualizar
err := onapiRepo.Update(ctx, "123", updatedEntity)

// Eliminar
err := onapiRepo.Delete(ctx, "123")
```

### Operaciones de Búsqueda
```go
// Búsqueda general
results, err := onapiRepo.Search(ctx, "Novasco", 0, 10)

// Búsqueda específica por categoría
results, err := onapiRepo.SearchByCategory(ctx, domain.KeywordCategoryCompanyName, "Novasco", 0, 10)

// Búsqueda por tipo de dominio
results, err := onapiRepo.GetByDomainType(ctx, domain.DomainTypeONAPI, 0, 10)
```

### Operaciones de Pipeline
```go
// Almacenar resultado de búsqueda
searchResult := &domain.DomainSearchResult{...}
err := pipelineRepo.Create(ctx, searchResult)

// Almacenar resultado de pipeline
pipelineResult := &module.DynamicPipelineResult{...}
err := pipelineRepo.Create(ctx, pipelineResult)

// Obtener palabras clave por categoría
keywords, err := pipelineRepo.GetKeywordsByCategory(ctx, "result-id")
```

## Manejo de Errores

Todos los repositorios retornan errores estándar de Go y manejan:
- Problemas de conexión a base de datos
- Errores de ejecución SQL
- Errores de marshaling/unmarshaling JSON
- Errores de validación de datos

## Consideraciones de Rendimiento

- Todas las tablas incluyen índices apropiados para patrones de búsqueda comunes
- Los campos JSON se usan con moderación y solo para estructuras de datos complejas
- La paginación está soportada para todas las operaciones de lista
- Las conexiones a base de datos se gestionan a través de la capa de servicio de base de datos

## Mejoras Futuras

- Agregar capa de caché para datos accedidos frecuentemente
- Implementar soporte de transacciones a nivel de repositorio
- Agregar operaciones en lote para procesamiento por lotes
- Implementar funcionalidad de eliminación suave
- Agregar registro de auditoría para cambios de datos
