# Repository Layer Documentation

This directory contains the repository layer implementation for the Insightful Intel application. The repository layer provides a clean abstraction over data persistence operations for all domain types and pipeline responses.

## Architecture

The repository layer follows the Repository pattern and provides:

- **BaseRepository[T]**: Common CRUD operations for any entity type
- **SearchableRepository[T]**: Extends BaseRepository with search capabilities
- **DomainRepository[T]**: Domain-specific operations for domain entities
- **PipelineRepository[T]**: Specialized operations for pipeline results

## Domain Repositories

### ONAPI Repository (`onapi_repository.go`)
- **Entity Type**: `domain.Entity`
- **Database Table**: `onapi_entities`
- **Key Features**:
  - Stores trademark and patent information
  - Supports search by company name, person name, and address
  - Handles JSON fields for images and class lists

### SCJ Repository (`scj_repository.go`)
- **Entity Type**: `domain.ScjCase`
- **Database Table**: `scj_cases`
- **Key Features**:
  - Stores Supreme Court of Justice case information
  - Supports search by case number, sentence number, and involved parties
  - Tracks case metadata and document URLs

### DGII Repository (`dgii_repository.go`)
- **Entity Type**: `domain.Register`
- **Database Table**: `dgii_registers`
- **Key Features**:
  - Stores Dominican tax authority register information
  - Supports search by RNC, company name, and commercial name
  - Tracks tax compliance status

### PGR Repository (`pgr_repository.go`)
- **Entity Type**: `domain.PGRNews`
- **Database Table**: `pgr_news`
- **Key Features**:
  - Stores Attorney General's Office news items
  - Supports search by title and URL
  - Minimal structure for news articles

### Google Docking Repository (`docking_repository.go`)
- **Entity Type**: `domain.GoogleDorkingResult`
- **Database Table**: `google_docking_results`
- **Key Features**:
  - Stores Google search results with relevance scoring
  - Supports search by title, description, and URL
  - Handles JSON fields for keywords
  - Includes ranking and relevance information

## Pipeline Repository

### Pipeline Repository (`pipeline_repository.go`)
- **Entity Types**: `domain.DomainSearchResult`, `module.DynamicPipelineResult`
- **Database Tables**: `domain_search_results`, `dynamic_pipeline_results`, `dynamic_pipeline_steps`
- **Key Features**:
  - Stores both individual domain search results and complete pipeline results
  - Supports complex pipeline step tracking
  - Handles JSON serialization for complex data structures
  - Provides aggregation capabilities for pipeline statistics

## Repository Factory

### RepositoryFactory (`factory.go`)
Provides centralized repository instantiation:

```go
factory := repositories.NewRepositoryFactory(db)

// Get specific repositories
onapiRepo := factory.GetOnapiRepository()
scjRepo := factory.GetScjRepository()
pipelineRepo := factory.GetPipelineRepository()

// Get all domain repositories
allRepos := factory.GetAllDomainRepositories()
```

## Database Schema

The `schema.sql` file contains all necessary table definitions with:
- Proper indexing for search performance
- Foreign key constraints where appropriate
- JSON columns for complex data structures
- Timestamp tracking for audit purposes

## Usage Examples

### Basic CRUD Operations
```go
// Create
entity := domain.Entity{...}
err := onapiRepo.Create(ctx, entity)

// Read
entity, err := onapiRepo.GetByID(ctx, "123")

// Update
err := onapiRepo.Update(ctx, "123", updatedEntity)

// Delete
err := onapiRepo.Delete(ctx, "123")
```

### Search Operations
```go
// General search
results, err := onapiRepo.Search(ctx, "Novasco", 0, 10)

// Category-specific search
results, err := onapiRepo.SearchByCategory(ctx, domain.KeywordCategoryCompanyName, "Novasco", 0, 10)

// Domain type search
results, err := onapiRepo.GetByDomainType(ctx, domain.DomainTypeONAPI, 0, 10)
```

### Pipeline Operations
```go
// Store search result
searchResult := &domain.DomainSearchResult{...}
err := pipelineRepo.Create(ctx, searchResult)

// Store pipeline result
pipelineResult := &module.DynamicPipelineResult{...}
err := pipelineRepo.Create(ctx, pipelineResult)

// Get keywords by category
keywords, err := pipelineRepo.GetKeywordsByCategory(ctx, "result-id")
```

## Error Handling

All repositories return standard Go errors and handle:
- Database connection issues
- SQL execution errors
- JSON marshaling/unmarshaling errors
- Data validation errors

## Performance Considerations

- All tables include appropriate indexes for common search patterns
- JSON fields are used sparingly and only for complex data structures
- Pagination is supported for all list operations
- Database connections are managed through the database service layer

## Testing

Repository implementations should be tested with:
- Unit tests for individual methods
- Integration tests with actual database
- Mock database for isolated testing
- Performance tests for search operations

## Future Enhancements

- Add caching layer for frequently accessed data
- Implement repository-level transaction support
- Add bulk operations for batch processing
- Implement soft delete functionality
- Add audit logging for data changes
