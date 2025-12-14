# Implementando un Nuevo Tipo de Dominio

Esta guía explica cómo agregar un nuevo tipo de dominio a la plataforma Insightful Intel. Un tipo de dominio representa una fuente de datos (por ejemplo, ONAPI, SCJ, DGII) que puede ser buscada e integrada en el sistema de pipeline dinámico.

## Tabla de Contenidos

1. [Resumen](#resumen)
2. [Pasos de Implementación](#pasos-de-implementación)
3. [Paso 1: Definir Tipo de Dominio](#paso-1-definir-tipo-de-dominio)
4. [Paso 2: Crear Modelo de Dominio](#paso-2-crear-modelo-de-dominio)
5. [Paso 3: Implementar Conector de Dominio](#paso-3-implementar-conector-de-dominio)
6. [Paso 4: Crear Repositorio](#paso-4-crear-repositorio)
7. [Paso 5: Integrar con Capa de Módulo](#paso-5-integrar-con-capa-de-módulo)
8. [Paso 6: Actualizar Esquema de Base de Datos](#paso-6-actualizar-esquema-de-base-de-datos)
9. [Paso 7: Agregar Endpoints de API](#paso-7-agregar-endpoints-de-api)
10. [Paso 8: Actualizar Frontend (Opcional)](#paso-8-actualizar-frontend-opcional)
11. [Pruebas](#pruebas)
12. [Ejemplo Completo](#ejemplo-completo)

---

## Resumen

Un tipo de dominio en Insightful Intel consta de varios componentes:

1. **Constante de Tipo de Dominio** - Identificador tipo enum en `domain/connector.go`
2. **Modelo de Dominio** - Estructura de entidad en el paquete `domain/`
3. **Conector de Dominio** - Implementación en el paquete `module/` que implementa la interfaz `DomainConnector[T]`
4. **Repositorio** - Capa de acceso a datos en el paquete `repositories/`
5. **Esquema de Base de Datos** - Definición de tabla en `database/schema.sql`
6. **Integración de Módulo** - Registro en `module/dynamic.go`
7. **Manejador de API** - Endpoint HTTP en `server/routes.go`

---

## Pasos de Implementación

### Paso 1: Definir Tipo de Dominio

**Archivo**: `internal/domain/connector.go`

Agregar la constante del nuevo tipo de dominio y actualizar los mapeos:

```go
const (
    // ... tipos existentes
    DomainTypeNewDomain DomainType = "NEW_DOMAIN"
)

// Actualizar AllDomainTypes()
func AllDomainTypes() []DomainType {
    return []DomainType{
        // ... tipos existentes
        DomainTypeNewDomain,
    }
}

// Actualizar mapa StringToDomainType
var StringToDomainType = map[string]DomainType{
    // ... mapeos existentes
    "new_domain": DomainTypeNewDomain,
}

// Actualizar mapa DomainTypeToString
var DomainTypeToString = map[DomainType]string{
    // ... mapeos existentes
    DomainTypeNewDomain: "new_domain",
}
```

---

### Paso 2: Crear Modelo de Dominio

**Archivo**: `internal/domain/newdomain.go`

Crear un nuevo archivo para el modelo de dominio:

```go
package domain

import "time"

// NewDomainEntity representa una entidad del nuevo dominio
type NewDomainEntity struct {
    ID                   ID        `json:"id"`
    DomainSearchResultID ID        `json:"domain_search_result_id"`
    
    // Agregar campos específicos del dominio
    Name                 string    `json:"name"`
    Description          string    `json:"description"`
    Identifier           string    `json:"identifier"`
    
    // Campos comunes (heredados del patrón Common struct)
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}
```

**Puntos Clave**:
- Siempre incluir campos `ID` y `DomainSearchResultID`
- Incluir `CreatedAt` y `UpdatedAt` para auditoría
- Usar etiquetas JSON apropiadas para serialización
- Elegir nombres de campos significativos que reflejen el dominio

---

### Paso 3: Implementar Conector de Dominio

**Archivo**: `internal/module/newdomain.go`

Crear un nuevo archivo implementando la interfaz `DomainConnector[T]`:

```go
package module

import (
    "fmt"
    "insightful-intel/internal/custom"
    "insightful-intel/internal/domain"
    "strings"
)

// Verificar implementación de interfaz en tiempo de compilación
var _ domain.DomainConnector[domain.NewDomainEntity] = &NewDomain{}

// NewDomain implementa DomainConnector para el nuevo dominio
type NewDomain struct {
    Stuff    custom.Client
    BasePath string
    PathMap  custom.CustomPathMap
}

// NewNewDomainDomain crea una nueva instancia de conector de dominio
func NewNewDomainDomain() domain.DomainConnector[domain.NewDomainEntity] {
    return &NewDomain{
        BasePath: "https://api.example.com/endpoint",
        Stuff:    *custom.NewClient(),
    }
}

// GetDomainType retorna el identificador del tipo de dominio
func (n *NewDomain) GetDomainType() domain.DomainType {
    return domain.DomainTypeNewDomain
}

// Search realiza una consulta de búsqueda y retorna resultados
func (n *NewDomain) Search(query string) ([]domain.NewDomainEntity, error) {
    // Implementar lógica de búsqueda aquí
    // Esto podría ser:
    // - Llamada a API HTTP
    // - Web scraping
    // - Consulta a base de datos
    // - Búsqueda en sistema de archivos
    
    // Ejemplo: Llamada a API HTTP
    resp, err := n.Stuff.Get(n.BasePath+"?q="+query, map[string]string{
        "Content-Type": "application/json",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()
    
    // Parsear respuesta y convertir a entidades de dominio
    // var results []domain.NewDomainEntity
    // ... lógica de parsing ...
    
    return results, nil
}

// ProcessData procesa y valida datos de entidad
func (n *NewDomain) ProcessData(data domain.NewDomainEntity) (domain.NewDomainEntity, error) {
    if err := n.ValidateData(data); err != nil {
        return domain.NewDomainEntity{}, err
    }
    return n.TransformData(data), nil
}

// ValidateData valida datos de entidad
func (n *NewDomain) ValidateData(data domain.NewDomainEntity) error {
    // Agregar lógica de validación
    if data.Identifier == "" {
        return fmt.Errorf("identifier is required")
    }
    return nil
}

// TransformData transforma/limpia datos de entidad
func (n *NewDomain) TransformData(data domain.NewDomainEntity) domain.NewDomainEntity {
    transformed := data
    transformed.Name = strings.TrimSpace(data.Name)
    transformed.Description = strings.TrimSpace(data.Description)
    // Agregar otras transformaciones
    return transformed
}

// GetDataByCategory extrae datos por categoría de palabra clave
func (n *NewDomain) GetDataByCategory(data domain.NewDomainEntity, category domain.KeywordCategory) []string {
    result := []string{}
    
    switch category {
    case domain.KeywordCategoryCompanyName:
        result = append(result, data.Name)
    case domain.KeywordCategoryPersonName:
        // Extraer nombres de personas si aplica
        // result = append(result, data.PersonName)
    case domain.KeywordCategoryAddress:
        // Extraer direcciones si aplica
        // result = append(result, data.Address)
    case domain.KeywordCategoryContributorID:
        result = append(result, data.Identifier)
    }
    
    return result
}

// GetSearchableKeywordCategories retorna categorías que este dominio puede buscar
func (n *NewDomain) GetSearchableKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
        domain.KeywordCategoryContributorID,
        // Agregar categorías que este dominio puede buscar
    }
}

// GetFoundKeywordCategories retorna categorías que este dominio puede extraer de resultados
func (n *NewDomain) GetFoundKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
        domain.KeywordCategoryPersonName,
        // Agregar categorías que este dominio puede extraer
    }
}
```

**Métodos Clave a Implementar**:

1. **`Search(query string)`** - Realiza la operación de búsqueda real
2. **`GetDataByCategory(data, category)`** - Extrae palabras clave por categoría de resultados
3. **`GetSearchableKeywordCategories()`** - Define qué categorías puede buscar este dominio
4. **`GetFoundKeywordCategories()`** - Define qué categorías puede extraer este dominio
5. **`ValidateData()`** - Valida datos de entidad
6. **`TransformData()`** - Limpia/transforma datos de entidad

---

### Paso 4: Crear Repositorio

**Archivo**: `internal/repositories/newdomain.go`

Crear un repositorio para operaciones de base de datos:

```go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "insightful-intel/internal/database"
    "insightful-intel/internal/domain"
)

// NewDomainRepository implementa DomainRepository para NewDomainEntity
type NewDomainRepository struct {
    db DatabaseAccessor
}

// NewNewDomainRepository crea una nueva instancia de repositorio
func NewNewDomainRepository(db database.Service) *NewDomainRepository {
    return &NewDomainRepository{
        db: NewDatabaseAdapter(db),
    }
}

// Create inserta una nueva entidad
func (r *NewDomainRepository) Create(ctx context.Context, entity domain.NewDomainEntity) error {
    entity.ID = domain.NewID()
    
    query := `
        INSERT INTO new_domain_entities (
            id, domain_search_result_id, name, description, identifier,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, NOW(), NOW())
    `
    
    _, err := r.db.ExecContext(ctx, query,
        entity.ID, entity.DomainSearchResultID, entity.Name,
        entity.Description, entity.Identifier,
    )
    
    return err
}

// GetByID recupera una entidad por ID
func (r *NewDomainRepository) GetByID(ctx context.Context, id string) (domain.NewDomainEntity, error) {
    query := `
        SELECT id, domain_search_result_id, name, description, identifier,
               created_at, updated_at
        FROM new_domain_entities
        WHERE id = ?
    `
    
    var entity domain.NewDomainEntity
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &entity.ID, &entity.DomainSearchResultID, &entity.Name,
        &entity.Description, &entity.Identifier,
        &entity.CreatedAt, &entity.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return domain.NewDomainEntity{}, fmt.Errorf("entity not found")
    }
    if err != nil {
        return domain.NewDomainEntity{}, err
    }
    
    return entity, nil
}

// Search realiza una consulta de búsqueda
func (r *NewDomainRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.NewDomainEntity, error) {
    searchQuery := `
        SELECT id, domain_search_result_id, name, description, identifier,
               created_at, updated_at
        FROM new_domain_entities
        WHERE name LIKE ? OR description LIKE ? OR identifier LIKE ?
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    
    pattern := "%" + query + "%"
    rows, err := r.db.QueryContext(ctx, searchQuery, pattern, pattern, pattern, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var entities []domain.NewDomainEntity
    for rows.Next() {
        var entity domain.NewDomainEntity
        err := rows.Scan(
            &entity.ID, &entity.DomainSearchResultID, &entity.Name,
            &entity.Description, &entity.Identifier,
            &entity.CreatedAt, &entity.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        entities = append(entities, entity)
    }
    
    return entities, nil
}

// Implementar otros métodos requeridos:
// - Update(ctx, id, entity)
// - Delete(ctx, id)
// - List(ctx, offset, limit)
// - Count(ctx)
// - SearchByCategory(ctx, category, query, offset, limit)
// - GetByDomainType(ctx, domainType, offset, limit)
// - GetBySearchParameter(ctx, searchParam, offset, limit)
// - GetKeywordsByCategory(ctx, entityID)
```

**Actualizar Factory de Repositorio** (`internal/repositories/factory.go`):

```go
// GetNewDomainRepository retorna una nueva instancia de repositorio de dominio
func (f *RepositoryFactory) GetNewDomainRepository() *NewDomainRepository {
    return NewNewDomainRepository(f.db)
}

// Actualizar GetAllDomainRepositories() si es necesario
func (f *RepositoryFactory) GetAllDomainRepositories() *DomainRepositoryHandler {
    return &DomainRepositoryHandler{
        // ... repositorios existentes
        GetNewDomainRepository: f.GetNewDomainRepository,
    }
}

// Actualizar switch statement de GetRepositoryByDomainType()
func (f *RepositoryFactory) GetRepositoryByDomainType(domainType domain.DomainType) any {
    switch domainType {
    // ... casos existentes
    case domain.DomainTypeNewDomain:
        return f.GetNewDomainRepository()
    }
    return nil
}
```

---

### Paso 5: Integrar con Capa de Módulo

**Archivo**: `internal/module/dynamic.go`

Actualizar la función `SearchDomain` para incluir el nuevo dominio:

```go
func SearchDomain(domainType domain.DomainType, params domain.DomainSearchParams) (*domain.DomainSearchResult, error) {
    // ... validación existente
    
    switch domainType {
    // ... casos existentes
    case domain.DomainTypeNewDomain:
        newDomain := NewNewDomainDomain()
        output, searchErr = newDomain.Search(params.Query)
    default:
        return &domain.DomainSearchResult{
            Success:    false,
            Error:      fmt.Errorf("unsupported domain type: %s", domainType),
            DomainType: domainType,
        }, fmt.Errorf("unsupported domain type: %s", domainType)
    }
    
    // Extraer palabras clave
    var keywordsPerCategory map[domain.KeywordCategory][]string
    if searchErr == nil && output != nil {
        switch domainType {
        // ... casos existentes
        case domain.DomainTypeNewDomain:
            if entities, ok := output.([]domain.NewDomainEntity); ok {
                keywordsPerCategory = domain.GetCategoryByKeywords(NewNewDomainDomain(), entities)
            }
        }
    }
    
    return &domain.DomainSearchResult{
        Success:             searchErr == nil,
        Error:               searchErr,
        DomainType:          domainType,
        SearchParameter:     params.Query,
        KeywordsPerCategory: keywordsPerCategory,
        Output:              output,
    }, searchErr
}
```

**Actualizar función `CreateDomainConnector`**:

```go
func CreateDomainConnector(domainType domain.DomainType) (any, error) {
    switch domainType {
    // ... casos existentes
    case domain.DomainTypeNewDomain:
        newDomain := NewNewDomainDomain()
        return &newDomain, nil
    default:
        return nil, fmt.Errorf("unsupported domain type: %s", domainType)
    }
}
```

**Actualizar función `CreateDynamicPipeline`** para establecer mapeo de categoría inicial:

```go
initialDomainCategories := map[domain.DomainType]domain.KeywordCategory{
    // ... mapeos existentes
    domain.DomainTypeNewDomain: domain.KeywordCategoryCompanyName, // o categoría apropiada
}
```

---

### Paso 6: Actualizar Esquema de Base de Datos

**Archivo**: `internal/database/schema.sql`

Agregar definición de tabla para el dominio:

```sql
-- Tabla de entidades del nuevo dominio
CREATE TABLE IF NOT EXISTS new_domain_entities (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    name VARCHAR(255),
    description TEXT,
    identifier VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_name (name),
    INDEX idx_identifier (identifier),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Puntos Clave**:
- Siempre incluir `id` como CHAR(36) para UUID
- Incluir `domain_search_result_id` con clave foránea
- Agregar índices apropiados para rendimiento de búsqueda
- Incluir timestamps `created_at` y `updated_at`

---

### Paso 7: Agregar Endpoints de API

**Archivo**: `internal/server/routes.go`

Agregar registro de ruta:

```go
func (s *Server) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()
    
    // ... rutas existentes
    mux.HandleFunc("/api/newdomain", s.newDomainHandler)
    
    return s.corsMiddleware(mux)
}
```

**Agregar función manejadora**:

```go
func (s *Server) newDomainHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Obtener parámetros de consulta
    query := r.URL.Query().Get("q")
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    
    if limit == 0 {
        limit = 10
    }
    
    repo := s.repositories.GetNewDomainRepository()
    
    var results []domain.NewDomainEntity
    var err error
    
    if query != "" {
        results, err = repo.Search(r.Context(), query, offset, limit)
    } else {
        results, err = repo.List(r.Context(), offset, limit)
    }
    
    if err != nil {
        http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "data":    results,
        "count":   len(results),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

**Actualizar interactor** (`internal/interactor/dymanic.go`) para manejar nuevo dominio en pipeline:

```go
switch step.DomainType {
// ... casos existentes
case domain.DomainTypeNewDomain:
    results, ok := created.Output.([]domain.NewDomainEntity)
    if !ok {
        log.Println("Error casting result output to []domain.NewDomainEntity")
        return nil, err
    }
    for _, result := range results {
        result.DomainSearchResultID = created.ID
        if err := d.repositories.GetNewDomainRepository().Create(ctx, result); err != nil {
            log.Println("Error creating new domain repository", err)
            return nil, err
        }
    }
}
```

---

### Paso 8: Actualizar Frontend (Opcional)

Si deseas agregar soporte frontend:

1. **Agregar mapeo de tipo de dominio** (`frontend/src/types.ts` o similar):
```typescript
export const DOMAIN_TYPE_MAP = {
  // ... mapeos existentes
  NEW_DOMAIN: 'NEW_DOMAIN',
};
```

2. **Crear componente** (`frontend/src/components/NewDomainRow.tsx`):
```typescript
interface NewDomainRowProps {
  entity: NewDomainEntity;
}

export function NewDomainRow({ entity }: NewDomainRowProps) {
  return (
    <div>
      <h3>{entity.name}</h3>
      <p>{entity.description}</p>
      <p>ID: {entity.identifier}</p>
    </div>
  );
}
```

3. **Actualizar PipelineDetails** para incluir nuevo tipo de dominio en filtrado

---

## Pruebas

### Pruebas Unitarias

Crear pruebas para el conector de dominio:

```go
// internal/module/newdomain_test.go
package module

import (
    "testing"
    "insightful-intel/internal/domain"
)

func TestNewDomain_GetDomainType(t *testing.T) {
    connector := NewNewDomainDomain()
    if connector.GetDomainType() != domain.DomainTypeNewDomain {
        t.Errorf("Expected DomainTypeNewDomain, got %v", connector.GetDomainType())
    }
}

func TestNewDomain_GetSearchableKeywordCategories(t *testing.T) {
    connector := NewNewDomainDomain()
    categories := connector.GetSearchableKeywordCategories()
    
    expected := []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
    }
    
    if len(categories) != len(expected) {
        t.Errorf("Expected %d categories, got %d", len(expected), len(categories))
    }
}

// Agregar más pruebas...
```

### Pruebas de Integración

Probar operaciones de repositorio:

```go
// internal/repositories/newdomain_test.go
package repositories

import (
    "context"
    "testing"
    "insightful-intel/internal/domain"
)

func TestNewDomainRepository_Create(t *testing.T) {
    // Configurar base de datos de prueba
    // Crear repositorio
    // Probar operación Create
}
```

---

## Ejemplo Completo

Aquí hay un ejemplo mínimo para un dominio hipotético "Registro de Negocios":

### 1. Definición de Tipo de Dominio
```go
// domain/connector.go
DomainTypeBusinessRegistry DomainType = "BUSINESS_REGISTRY"
```

### 2. Modelo de Dominio
```go
// domain/businessregistry.go
type BusinessRegistry struct {
    ID                   ID
    DomainSearchResultID ID
    BusinessName         string
    RegistrationNumber   string
    OwnerName           string
    Address             string
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

### 3. Conector de Dominio
```go
// module/businessregistry.go
type BusinessRegistry struct {
    Stuff    custom.Client
    BasePath string
}

func (b *BusinessRegistry) Search(query string) ([]domain.BusinessRegistry, error) {
    // Implementación
}

func (b *BusinessRegistry) GetDataByCategory(data domain.BusinessRegistry, category domain.KeywordCategory) []string {
    switch category {
    case domain.KeywordCategoryCompanyName:
        return []string{data.BusinessName}
    case domain.KeywordCategoryPersonName:
        return []string{data.OwnerName}
    case domain.KeywordCategoryAddress:
        return []string{data.Address}
    }
    return []string{}
}
```

### 4. Repositorio
```go
// repositories/businessregistry.go
type BusinessRegistryRepository struct {
    db DatabaseAccessor
}

func (r *BusinessRegistryRepository) Create(ctx context.Context, entity domain.BusinessRegistry) error {
    // Implementación
}
```

---

## Lista de Verificación

Usar esta lista de verificación al implementar un nuevo tipo de dominio:

- [ ] Agregar constante de tipo de dominio a `domain/connector.go`
- [ ] Actualizar `AllDomainTypes()`, `StringToDomainType`, y `DomainTypeToString`
- [ ] Crear estructura de modelo de dominio en paquete `domain/`
- [ ] Implementar interfaz `DomainConnector[T]` en paquete `module/`
- [ ] Crear repositorio en paquete `repositories/`
- [ ] Actualizar factory de repositorio
- [ ] Agregar tabla de base de datos en `schema.sql`
- [ ] Integrar con `module/dynamic.go` (SearchDomain, CreateDomainConnector, CreateDynamicPipeline)
- [ ] Agregar endpoint de API en `server/routes.go`
- [ ] Actualizar interactor para manejar nuevo dominio en ejecución de pipeline
- [ ] Escribir pruebas unitarias
- [ ] Escribir pruebas de integración
- [ ] Actualizar frontend (si es necesario)
- [ ] Actualizar documentación

---

## Patrones Comunes

### Patrón 1: Integración con API HTTP

```go
func (d *Domain) Search(query string) ([]domain.Entity, error) {
    resp, err := d.Stuff.Get(d.BasePath+"?q="+query, headers)
    // Parsear respuesta JSON
    // Convertir a entidades de dominio
    return entities, nil
}
```

### Patrón 2: Web Scraping

```go
func (d *Domain) Search(query string) ([]domain.Entity, error) {
    // Usar Colly o librería similar
    // Scrapear HTML
    // Extraer datos
    // Convertir a entidades de dominio
    return entities, nil
}
```

### Patrón 3: Extracción de Palabras Clave con Regex

```go
func (d *Domain) GetDataByCategory(data domain.Entity, category domain.KeywordCategory) []string {
    switch category {
    case domain.KeywordCategoryPersonName:
        // Usar regex para extraer nombres
        re := regexp.MustCompile(`(?i)\s*,\s*|\s+vs\.?\s*`)
        names := re.Split(data.Involucrados, -1)
        // Filtrar y retornar
    }
}
```

---

## Solución de Problemas

### Problema: Dominio no aparece en búsquedas

**Solución**: Verificar que:
1. El tipo de dominio está agregado a `AllDomainTypes()`
2. La función `SearchDomain()` incluye el dominio en el switch statement
3. El conector de dominio está correctamente instanciado

### Problema: Palabras clave no se extraen

**Solución**: Verificar:
1. `GetFoundKeywordCategories()` retorna las categorías correctas
2. `GetDataByCategory()` extrae correctamente datos para cada categoría
3. Las categorías coinciden entre categorías buscables y encontradas

### Problema: Errores de base de datos

**Solución**: Asegurar:
1. La tabla existe en la base de datos (ejecutar migraciones)
2. Las consultas SQL del repositorio coinciden con el esquema de tabla
3. Las relaciones de clave foránea son correctas

---

## Mejores Prácticas

1. **Seguir Convenciones de Nomenclatura**: Usar nomenclatura consistente (por ejemplo, `NewDomainDomain`, `NewDomainRepository`)
2. **Manejo de Errores**: Siempre retornar errores descriptivos con contexto
3. **Validación**: Validar datos en múltiples niveles (conector, repositorio, base de datos)
4. **Logging**: Agregar logging apropiado para depuración
5. **Documentación**: Documentar cualquier lógica o particularidad específica del dominio
6. **Pruebas**: Escribir pruebas para rutas críticas
7. **Rendimiento**: Considerar caché para datos accedidos frecuentemente
8. **Seguridad**: Sanitizar entradas y usar consultas parametrizadas

---

## Recursos Adicionales

- [Documentación de Diseño Dirigido por Dominio](PROJECT_DOCUMENTATION.md#ddd-implementation-in-the-project)
- [Guía de Pipeline Dinámico](DYNAMIC_PIPELINE_GUIDE.md)
- [README de Capa de Repositorio](../internal/repositories/README.md)

---

**¿Necesitas Ayuda?** Revisa implementaciones existentes:
- ONAPI: `internal/module/onapi.go`
- SCJ: `internal/module/scj.go`
- DGII: `internal/module/dgii.go`
- PGR: `internal/module/pgr.go`
