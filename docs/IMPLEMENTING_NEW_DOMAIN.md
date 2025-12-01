# Implementing a New Domain Type

This guide explains how to add a new domain type to the Insightful Intel platform. A domain type represents a data source (e.g., ONAPI, SCJ, DGII) that can be searched and integrated into the dynamic pipeline system.

## Table of Contents

1. [Overview](#overview)
2. [Implementation Steps](#implementation-steps)
3. [Step 1: Define Domain Type](#step-1-define-domain-type)
4. [Step 2: Create Domain Model](#step-2-create-domain-model)
5. [Step 3: Implement Domain Connector](#step-3-implement-domain-connector)
6. [Step 4: Create Repository](#step-4-create-repository)
7. [Step 5: Integrate with Module Layer](#step-5-integrate-with-module-layer)
8. [Step 6: Update Database Schema](#step-6-update-database-schema)
9. [Step 7: Add API Endpoints](#step-7-add-api-endpoints)
10. [Step 8: Update Frontend (Optional)](#step-8-update-frontend-optional)
11. [Testing](#testing)
12. [Complete Example](#complete-example)

---

## Overview

A domain type in Insightful Intel consists of several components:

1. **Domain Type Constant** - Enum-like identifier in `domain/connector.go`
2. **Domain Model** - Entity struct in `domain/` package
3. **Domain Connector** - Implementation in `module/` package that implements `DomainConnector[T]` interface
4. **Repository** - Data access layer in `repositories/` package
5. **Database Schema** - Table definition in `database/schema.sql`
6. **Module Integration** - Registration in `module/dynamic.go`
7. **API Handler** - HTTP endpoint in `server/routes.go`

---

## Implementation Steps

### Step 1: Define Domain Type

**File**: `internal/domain/connector.go`

Add your new domain type constant and update the mappings:

```go
const (
    // ... existing types
    DomainTypeNewDomain DomainType = "NEW_DOMAIN"
)

// Update AllDomainTypes()
func AllDomainTypes() []DomainType {
    return []DomainType{
        // ... existing types
        DomainTypeNewDomain,
    }
}

// Update StringToDomainType map
var StringToDomainType = map[string]DomainType{
    // ... existing mappings
    "new_domain": DomainTypeNewDomain,
}

// Update DomainTypeToString map
var DomainTypeToString = map[DomainType]string{
    // ... existing mappings
    DomainTypeNewDomain: "new_domain",
}
```

---

### Step 2: Create Domain Model

**File**: `internal/domain/newdomain.go`

Create a new file for your domain model:

```go
package domain

import "time"

// NewDomainEntity represents an entity from the new domain
type NewDomainEntity struct {
    ID                   ID        `json:"id"`
    DomainSearchResultID ID        `json:"domain_search_result_id"`
    
    // Add your domain-specific fields
    Name                 string    `json:"name"`
    Description          string    `json:"description"`
    Identifier           string    `json:"identifier"`
    
    // Common fields (inherited from Common struct pattern)
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}
```

**Key Points**:
- Always include `ID` and `DomainSearchResultID` fields
- Include `CreatedAt` and `UpdatedAt` for audit trails
- Use appropriate JSON tags for serialization
- Choose meaningful field names that reflect the domain

---

### Step 3: Implement Domain Connector

**File**: `internal/module/newdomain.go`

Create a new file implementing the `DomainConnector[T]` interface:

```go
package module

import (
    "fmt"
    "insightful-intel/internal/custom"
    "insightful-intel/internal/domain"
    "strings"
)

// Verify interface implementation at compile time
var _ domain.DomainConnector[domain.NewDomainEntity] = &NewDomain{}

// NewDomain implements DomainConnector for the new domain
type NewDomain struct {
    Stuff    custom.Client
    BasePath string
    PathMap  custom.CustomPathMap
}

// NewNewDomainDomain creates a new domain connector instance
func NewNewDomainDomain() domain.DomainConnector[domain.NewDomainEntity] {
    return &NewDomain{
        BasePath: "https://api.example.com/endpoint",
        Stuff:    *custom.NewClient(),
    }
}

// GetDomainType returns the domain type identifier
func (n *NewDomain) GetDomainType() domain.DomainType {
    return domain.DomainTypeNewDomain
}

// Search performs a search query and returns results
func (n *NewDomain) Search(query string) ([]domain.NewDomainEntity, error) {
    // Implement your search logic here
    // This could be:
    // - HTTP API call
    // - Web scraping
    // - Database query
    // - File system search
    
    // Example: HTTP API call
    resp, err := n.Stuff.Get(n.BasePath+"?q="+query, map[string]string{
        "Content-Type": "application/json",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()
    
    // Parse response and convert to domain entities
    // var results []domain.NewDomainEntity
    // ... parsing logic ...
    
    return results, nil
}

// ProcessData processes and validates entity data
func (n *NewDomain) ProcessData(data domain.NewDomainEntity) (domain.NewDomainEntity, error) {
    if err := n.ValidateData(data); err != nil {
        return domain.NewDomainEntity{}, err
    }
    return n.TransformData(data), nil
}

// ValidateData validates entity data
func (n *NewDomain) ValidateData(data domain.NewDomainEntity) error {
    // Add validation logic
    if data.Identifier == "" {
        return fmt.Errorf("identifier is required")
    }
    return nil
}

// TransformData transforms/cleans entity data
func (n *NewDomain) TransformData(data domain.NewDomainEntity) domain.NewDomainEntity {
    transformed := data
    transformed.Name = strings.TrimSpace(data.Name)
    transformed.Description = strings.TrimSpace(data.Description)
    // Add any other transformations
    return transformed
}

// GetDataByCategory extracts data by keyword category
func (n *NewDomain) GetDataByCategory(data domain.NewDomainEntity, category domain.KeywordCategory) []string {
    result := []string{}
    
    switch category {
    case domain.KeywordCategoryCompanyName:
        result = append(result, data.Name)
    case domain.KeywordCategoryPersonName:
        // Extract person names if applicable
        // result = append(result, data.PersonName)
    case domain.KeywordCategoryAddress:
        // Extract addresses if applicable
        // result = append(result, data.Address)
    case domain.KeywordCategoryContributorID:
        result = append(result, data.Identifier)
    }
    
    return result
}

// GetSearchableKeywordCategories returns categories this domain can search
func (n *NewDomain) GetSearchableKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
        domain.KeywordCategoryContributorID,
        // Add categories this domain can search
    }
}

// GetFoundKeywordCategories returns categories this domain can extract from results
func (n *NewDomain) GetFoundKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
        domain.KeywordCategoryPersonName,
        // Add categories this domain can extract
    }
}
```

**Key Methods to Implement**:

1. **`Search(query string)`** - Performs the actual search operation
2. **`GetDataByCategory(data, category)`** - Extracts keywords by category from results
3. **`GetSearchableKeywordCategories()`** - Defines what categories this domain can search
4. **`GetFoundKeywordCategories()`** - Defines what categories this domain can extract
5. **`ValidateData()`** - Validates entity data
6. **`TransformData()`** - Cleans/transforms entity data

---

### Step 4: Create Repository

**File**: `internal/repositories/newdomain.go`

Create a repository for database operations:

```go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "insightful-intel/internal/database"
    "insightful-intel/internal/domain"
)

// NewDomainRepository implements DomainRepository for NewDomainEntity
type NewDomainRepository struct {
    db DatabaseAccessor
}

// NewNewDomainRepository creates a new repository instance
func NewNewDomainRepository(db database.Service) *NewDomainRepository {
    return &NewDomainRepository{
        db: NewDatabaseAdapter(db),
    }
}

// Create inserts a new entity
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

// GetByID retrieves an entity by ID
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

// Search performs a search query
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

// Implement other required methods:
// - Update(ctx, id, entity)
// - Delete(ctx, id)
// - List(ctx, offset, limit)
// - Count(ctx)
// - SearchByCategory(ctx, category, query, offset, limit)
// - GetByDomainType(ctx, domainType, offset, limit)
// - GetBySearchParameter(ctx, searchParam, offset, limit)
// - GetKeywordsByCategory(ctx, entityID)
```

**Update Repository Factory** (`internal/repositories/factory.go`):

```go
// GetNewDomainRepository returns a new domain repository instance
func (f *RepositoryFactory) GetNewDomainRepository() *NewDomainRepository {
    return NewNewDomainRepository(f.db)
}

// Update GetAllDomainRepositories() if needed
func (f *RepositoryFactory) GetAllDomainRepositories() *DomainRepositoryHandler {
    return &DomainRepositoryHandler{
        // ... existing repositories
        GetNewDomainRepository: f.GetNewDomainRepository,
    }
}

// Update GetRepositoryByDomainType() switch statement
func (f *RepositoryFactory) GetRepositoryByDomainType(domainType domain.DomainType) any {
    switch domainType {
    // ... existing cases
    case domain.DomainTypeNewDomain:
        return f.GetNewDomainRepository()
    }
    return nil
}
```

---

### Step 5: Integrate with Module Layer

**File**: `internal/module/dynamic.go`

Update the `SearchDomain` function to include your new domain:

```go
func SearchDomain(domainType domain.DomainType, params domain.DomainSearchParams) (*domain.DomainSearchResult, error) {
    // ... existing validation
    
    switch domainType {
    // ... existing cases
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
    
    // Extract keywords
    var keywordsPerCategory map[domain.KeywordCategory][]string
    if searchErr == nil && output != nil {
        switch domainType {
        // ... existing cases
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

**Update `CreateDomainConnector` function**:

```go
func CreateDomainConnector(domainType domain.DomainType) (any, error) {
    switch domainType {
    // ... existing cases
    case domain.DomainTypeNewDomain:
        newDomain := NewNewDomainDomain()
        return &newDomain, nil
    default:
        return nil, fmt.Errorf("unsupported domain type: %s", domainType)
    }
}
```

**Update `CreateDynamicPipeline` function** to set initial category mapping:

```go
initialDomainCategories := map[domain.DomainType]domain.KeywordCategory{
    // ... existing mappings
    domain.DomainTypeNewDomain: domain.KeywordCategoryCompanyName, // or appropriate category
}
```

---

### Step 6: Update Database Schema

**File**: `internal/database/schema.sql`

Add table definition for your domain:

```sql
-- New Domain entities table
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

**Key Points**:
- Always include `id` as CHAR(36) for UUID
- Include `domain_search_result_id` with foreign key
- Add appropriate indexes for search performance
- Include `created_at` and `updated_at` timestamps

---

### Step 7: Add API Endpoints

**File**: `internal/server/routes.go`

Add route registration:

```go
func (s *Server) RegisterRoutes() http.Handler {
    mux := http.NewServeMux()
    
    // ... existing routes
    mux.HandleFunc("/api/newdomain", s.newDomainHandler)
    
    return s.corsMiddleware(mux)
}
```

**Add handler function**:

```go
func (s *Server) newDomainHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get query parameters
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

**Update interactor** (`internal/interactor/dymanic.go`) to handle new domain in pipeline:

```go
switch step.DomainType {
// ... existing cases
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

### Step 8: Update Frontend (Optional)

If you want to add frontend support:

1. **Add domain type mapping** (`frontend/src/types.ts` or similar):
```typescript
export const DOMAIN_TYPE_MAP = {
  // ... existing mappings
  NEW_DOMAIN: 'NEW_DOMAIN',
};
```

2. **Create component** (`frontend/src/components/NewDomainRow.tsx`):
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

3. **Update PipelineDetails** to include new domain type in filtering

---

## Testing

### Unit Tests

Create tests for your domain connector:

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

// Add more tests...
```

### Integration Tests

Test repository operations:

```go
// internal/repositories/newdomain_test.go
package repositories

import (
    "context"
    "testing"
    "insightful-intel/internal/domain"
)

func TestNewDomainRepository_Create(t *testing.T) {
    // Setup test database
    // Create repository
    // Test Create operation
}
```

---

## Complete Example

Here's a minimal example for a hypothetical "Business Registry" domain:

### 1. Domain Type Definition
```go
// domain/connector.go
DomainTypeBusinessRegistry DomainType = "BUSINESS_REGISTRY"
```

### 2. Domain Model
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

### 3. Domain Connector
```go
// module/businessregistry.go
type BusinessRegistry struct {
    Stuff    custom.Client
    BasePath string
}

func (b *BusinessRegistry) Search(query string) ([]domain.BusinessRegistry, error) {
    // Implementation
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

### 4. Repository
```go
// repositories/businessregistry.go
type BusinessRegistryRepository struct {
    db DatabaseAccessor
}

func (r *BusinessRegistryRepository) Create(ctx context.Context, entity domain.BusinessRegistry) error {
    // Implementation
}
```

---

## Checklist

Use this checklist when implementing a new domain type:

- [ ] Add domain type constant to `domain/connector.go`
- [ ] Update `AllDomainTypes()`, `StringToDomainType`, and `DomainTypeToString`
- [ ] Create domain model struct in `domain/` package
- [ ] Implement `DomainConnector[T]` interface in `module/` package
- [ ] Create repository in `repositories/` package
- [ ] Update repository factory
- [ ] Add database table in `schema.sql`
- [ ] Integrate with `module/dynamic.go` (SearchDomain, CreateDomainConnector, CreateDynamicPipeline)
- [ ] Add API endpoint in `server/routes.go`
- [ ] Update interactor to handle new domain in pipeline execution
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Update frontend (if needed)
- [ ] Update documentation

---

## Common Patterns

### Pattern 1: HTTP API Integration

```go
func (d *Domain) Search(query string) ([]domain.Entity, error) {
    resp, err := d.Stuff.Get(d.BasePath+"?q="+query, headers)
    // Parse JSON response
    // Convert to domain entities
    return entities, nil
}
```

### Pattern 2: Web Scraping

```go
func (d *Domain) Search(query string) ([]domain.Entity, error) {
    // Use Colly or similar library
    // Scrape HTML
    // Extract data
    // Convert to domain entities
    return entities, nil
}
```

### Pattern 3: Keyword Extraction with Regex

```go
func (d *Domain) GetDataByCategory(data domain.Entity, category domain.KeywordCategory) []string {
    switch category {
    case domain.KeywordCategoryPersonName:
        // Use regex to extract names
        re := regexp.MustCompile(`(?i)\s*,\s*|\s+vs\.?\s*`)
        names := re.Split(data.Involucrados, -1)
        // Filter and return
    }
}
```

---

## Troubleshooting

### Issue: Domain not appearing in searches

**Solution**: Check that:
1. Domain type is added to `AllDomainTypes()`
2. `SearchDomain()` function includes your domain in the switch statement
3. Domain connector is properly instantiated

### Issue: Keywords not being extracted

**Solution**: Verify:
1. `GetFoundKeywordCategories()` returns the correct categories
2. `GetDataByCategory()` properly extracts data for each category
3. Categories match between searchable and found categories

### Issue: Database errors

**Solution**: Ensure:
1. Table exists in database (run migrations)
2. Repository SQL queries match table schema
3. Foreign key relationships are correct

---

## Best Practices

1. **Follow Naming Conventions**: Use consistent naming (e.g., `NewDomainDomain`, `NewDomainRepository`)
2. **Error Handling**: Always return descriptive errors with context
3. **Validation**: Validate data at multiple levels (connector, repository, database)
4. **Logging**: Add appropriate logging for debugging
5. **Documentation**: Document any domain-specific logic or quirks
6. **Testing**: Write tests for critical paths
7. **Performance**: Consider caching for frequently accessed data
8. **Security**: Sanitize inputs and use parameterized queries

---

## Additional Resources

- [Domain-Driven Design Documentation](PROJECT_DOCUMENTATION.md#ddd-implementation-in-the-project)
- [Dynamic Pipeline Guide](DYNAMIC_PIPELINE_GUIDE.md)
- [Repository Layer README](../internal/repositories/README.md)

---

**Need Help?** Review existing implementations:
- ONAPI: `internal/module/onapi.go`
- SCJ: `internal/module/scj.go`
- DGII: `internal/module/dgii.go`
- PGR: `internal/module/pgr.go`

