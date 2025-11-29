# Insightful Intel - Complete Project Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Technologies Used](#technologies-used)
3. [Development Process](#development-process)
4. [Use Cases](#use-cases)
5. [Domain-Driven Design (DDD) Overview](#domain-driven-design-ddd-overview)
6. [DDD Implementation in the Project](#ddd-implementation-in-the-project)
7. [Architecture](#architecture)
8. [Key Features](#key-features)

---

## Project Overview

**Insightful Intel** is an intelligent data aggregation and investigation platform designed to perform comprehensive searches across multiple Dominican Republic government and public data sources. The system enables automated cross-domain intelligence gathering by dynamically creating search pipelines that extract keywords from one domain and use them to search other related domains, creating a comprehensive investigative profile.

### Core Purpose
The platform automates the process of gathering intelligence from multiple public data sources, including:
- **ONAPI** (Oficina Nacional de la Propiedad Industrial) - Trademark and patent registrations
- **SCJ** (Suprema Corte de Justicia) - Supreme Court case records
- **DGII** (Dirección General de Impuestos Internos) - Tax authority registrations
- **PGR** (Procuraduría General de la República) - Attorney General's Office news
- **Google Docking** - Web search results with relevance scoring
- **Social Media** - Social media platform searches
- **File Type Searches** - Document and file searches

Traditional investigative work requires manually searching multiple government databases and websites, which is time-consuming, error-prone, and often misses connections between entities. This platform automates this process by:
1. Starting with a single search query (e.g., a company name or RNC)
2. Extracting relevant keywords from search results
3. Automatically creating new searches across other domains using those keywords
4. Building a comprehensive intelligence profile through iterative, depth-based exploration

---

## Technologies Used

### Backend Technologies

#### **Go (Golang) 1.24.2**
- **Purpose**: Core backend language
- **Why**: High performance, excellent concurrency support, strong typing
- **Key Features Used**:
  - Goroutines for concurrent pipeline execution
  - Channels for streaming results
  - Context for cancellation and timeout handling
  - Interfaces for domain abstraction

#### **Key Go Libraries**:
- **github.com/go-sql-driver/mysql**: MySQL database driver
- **github.com/gocolly/colly**: Web scraping framework for extracting data from government websites
- **github.com/google/uuid**: UUID generation for entity identification
- **github.com/spf13/cobra**: CLI command framework
- **github.com/testcontainers/testcontainers-go**: Integration testing with Docker containers
- **golang.org/x/net**: Network utilities

### Frontend Technologies

#### **React 19.1.0**
- **Purpose**: User interface framework
- **Features**: Component-based architecture, hooks for state management

#### **TypeScript 5.8.3**
- **Purpose**: Type-safe frontend development
- **Benefits**: Compile-time error checking, better IDE support

#### **Vite 6.3.5**
- **Purpose**: Build tool and development server
- **Benefits**: Fast hot module replacement, optimized production builds

#### **Tailwind CSS 4.1.10**
- **Purpose**: Utility-first CSS framework
- **Benefits**: Rapid UI development, responsive design

#### **React Router DOM 7.9.6**
- **Purpose**: Client-side routing
- **Features**: Navigation between pages (Dashboard, Pipeline, Search)

### Database

#### **MySQL**
- **Purpose**: Relational database for persistent storage
- **Schema Features**:
  - JSON columns for complex nested data
  - Foreign key constraints for referential integrity
  - Comprehensive indexing for search performance
  - Timestamp tracking for audit trails

### Infrastructure & DevOps

#### **Docker & Docker Compose**
- **Purpose**: Containerization and local development environment
- **Files**: `docker-compose.yml`, `docker-compose.api.yml`, `docker-compose.cli.yml`
- **Benefits**: Consistent development environment, easy deployment

#### **Makefile**
- **Purpose**: Build automation and task management
- **Commands**: `make build`, `make run`, `make test`, `make watch`, `make docker-run`

### Development Tools

#### **Air (Live Reload)**
- **Purpose**: Automatic application reloading during development
- **Integration**: Configured in Makefile for Go backend

#### **ESLint**
- **Purpose**: JavaScript/TypeScript code quality and linting
- **Configuration**: Custom ESLint config for React/TypeScript

---

## Development Process

### Project Structure

The project follows a clean, layered architecture:

```
insightful-intel/
├── docs/                 # Project documentation, guides, technical references
├── cmd/                  # Application entry points
│   ├── api/              # HTTP API server
│   └── cli/              # Command-line interface
├── internal/             # Private application code
│   ├── domain/           # Domain models and business logic
│   ├── repositories/     # Data access layer
│   ├── interactor/       # Application use cases
│   ├── module/           # Domain services
│   ├── server/           # HTTP handlers and routes
│   ├── database/         # Database connection and migrations
│   └── infra/            # Infrastructure concerns
├── frontend/             # React frontend application
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── api.ts       # API client
│   │   └── types.ts     # TypeScript type definitions
├── config/              # Configuration management
├── doc/                # Documentation
└── vendor/             # Go dependencies
```

### Development Workflow

1. **Local Development Setup**:
   ```bash
   make docker-run      # Start MySQL database container
   make watch          # Start backend with live reload
   npm run dev         # Start frontend development server
   ```

2. **Database Migrations**:
   - Migrations run automatically on application startup
   - Schema defined in `internal/database/schema.sql`
   - Migration service handles versioning

3. **Testing**:
   ```bash
   make test           # Run unit tests
   make itest          # Run integration tests with testcontainers
   ```

4. **Building**:
   ```bash
   make build          # Build API server
   make build-cli     # Build CLI tool
   ```

### Code Organization Principles

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Dependency Inversion**: High-level modules depend on abstractions (interfaces)
3. **Single Responsibility**: Each package/module has one clear purpose
4. **Interface-Based Design**: Domain logic defined through interfaces

---

## Use Cases

### 1. **Single Domain Search**
**Description**: Search a specific domain (e.g., ONAPI) with a query string.

**Flow**:
- User provides query (e.g., "Novasco") and domain type
- System performs search using domain-specific connector
- Results are returned with extracted keywords categorized by type

**Example**:
```
GET /search?q=Novasco&domain=onapi
```

**Output**: List of trademark/patent entities with extracted keywords (company names, person names, addresses)

### 2. **Multi-Domain Search**
**Description**: Search across multiple domains simultaneously with a single query.

**Flow**:
- User provides query without specifying domain
- System searches all default domains (ONAPI, SCJ, DGII) in parallel
- Aggregated results returned with keywords from each domain

**Example**:
```
GET /search?q=Novasco
```

**Output**: Combined results from multiple domains

### 3. **Dynamic Pipeline Execution** (Primary Use Case)
**Description**: Automated, iterative search pipeline that discovers new search targets from previous results.

**Flow**:
1. **Initialization**: User provides initial query and configuration (max depth, skip duplicates)
2. **Step Creation**: System creates initial search steps for all available domains
3. **Execution**: Each step is executed, results stored in database
4. **Keyword Extraction**: Keywords are extracted from results and categorized
5. **Step Generation**: New steps are created for each keyword in compatible domains
6. **Iteration**: Process repeats up to maximum depth
7. **Streaming**: Results streamed to client in real-time via Server-Sent Events (SSE)

**Example**:
```
GET /dynamic?q=Novasco&depth=5&skip_duplicates=true&stream=true
```

**Output**: Real-time stream of pipeline steps with:
- Domain type
- Search parameter
- Success/failure status
- Extracted keywords
- Results
- Depth level

**Use Case Scenarios**:
- **Due Diligence**: Investigate a company's legal standing, trademarks, tax status, and court cases
- **Fraud Detection**: Cross-reference entities across multiple databases to identify inconsistencies
- **Background Checks**: Comprehensive profile of individuals or companies
- **Compliance Verification**: Verify business registrations and legal status

### 4. **Pipeline Result Retrieval**
**Description**: Retrieve previously executed pipeline results from database.

**Endpoints**:
- `GET /api/pipeline` - List all pipelines
- `GET /api/pipeline/steps?pipeline_id={id}` - Get steps for a pipeline
- `GET /api/pipeline/save` - Save pipeline execution

### 5. **Domain-Specific Data Access**
**Description**: Query stored data from specific domains with filtering and pagination.

**Endpoints**:
- `GET /api/onapi` - ONAPI entities
- `GET /api/scj` - SCJ cases
- `GET /api/dgii` - DGII registers
- `GET /api/pgr` - PGR news
- `GET /api/docking` - Google Docking results

**Features**:
- Search by keyword
- Filter by category
- Pagination support
- Category-based keyword extraction

---

## Domain-Driven Design (DDD) Overview

### What is Domain-Driven Design?

Domain-Driven Design (DDD) is a software development approach introduced by Eric Evans that focuses on:
1. **Ubiquitous Language**: A common vocabulary shared by developers and domain experts
2. **Domain Models**: Rich, behavior-rich models that represent business concepts
3. **Bounded Contexts**: Explicit boundaries where a particular domain model applies
4. **Layered Architecture**: Separation between domain logic, application logic, and infrastructure

### Core DDD Concepts

#### **Entities**
Objects with unique identity that persist over time. In this project:
- `Entity` (ONAPI trademark/patent)
- `ScjCase` (Court case)
- `Register` (DGII tax registration)
- `PGRNews` (News article)
- `GoogleDockingResult` (Web search result)

#### **Value Objects**
Objects defined by their attributes rather than identity:
- `ID` (UUID wrapper)
- `KeywordCategory` (enum-like string type)
- `DomainType` (enum-like string type)

#### **Aggregates**
Clusters of entities and value objects treated as a single unit:
- `DynamicPipelineResult` (aggregate root containing multiple `DynamicPipelineStep` entities)

#### **Domain Services**
Operations that don't naturally fit within a single entity:
- `SearchDomain()` - Coordinates search across different domain types
- `CreateDynamicPipeline()` - Orchestrates pipeline creation

#### **Repositories**
Abstractions for accessing aggregates:
- `OnapiRepository`, `ScjRepository`, `DgiiRepository`, etc.
- `PipelineRepository` for pipeline aggregates

#### **Application Services (Interactors)**
Orchestrate domain objects to fulfill use cases:
- `DynamicPipelineInteractor` - Coordinates pipeline execution

---

## DDD Implementation in the Project

### 1. **Domain Layer** (`internal/domain/`)

The domain layer contains the core business logic and domain models, completely independent of infrastructure concerns.

#### **Domain Entities**

Each domain entity represents a real-world concept from the Dominican Republic's public data sources:

**ONAPI Entity** (`domain/onapi.go`):
```go
type Entity struct {
    ID                   ID
    SerieExpediente      int32
    NumeroExpediente     int32
    Titular              string      // Company/Person name
    Gestor               string      // Manager name
    Domicilio            string      // Address
    // ... other trademark/patent fields
}
```

**SCJ Case** (`domain/scj.go`):
```go
type ScjCase struct {
    ID              ID
    NoExpediente    string      // Case number
    NoSentencia     string      // Sentence number
    Involucrados    string      // Involved parties
    DescTribunal    string      // Court description
    // ... other case fields
}
```

**DGII Register** (`domain/dgii.go`):
```go
type Register struct {
    ID                ID
    RNC               string      // Tax ID
    RazonSocial       string      // Legal name
    NombreComercial   string      // Commercial name
    Estado            string      // Status
    // ... other tax registration fields
}
```

#### **Domain Connector Interface**

The `DomainConnector[T]` interface (`domain/connector.go`) is the heart of the DDD implementation:

```go
type DomainConnector[T any] interface {
    ProcessData(data T) (T, error)
    ValidateData(data T) error
    TransformData(data T) T
    GetDataByCategory(data T, category KeywordCategory) []string
    GetSearchableKeywordCategories() []KeywordCategory
    GetFoundKeywordCategories() []KeywordCategory
    GetDomainType() DomainType
    Search(query string) ([]T, error)
}
```

**Key DDD Principles Applied**:
- **Polymorphism**: Each domain (ONAPI, SCJ, DGII, etc.) implements this interface
- **Encapsulation**: Domain logic is encapsulated within each connector
- **Ubiquitous Language**: Method names reflect business concepts (Search, GetDataByCategory)

**Example Implementation** (`module/onapi.go`):
```go
type Onapi struct {
    // Domain-specific fields
}

func (o *Onapi) GetSearchableKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,
        domain.KeywordCategoryPersonName,
        domain.KeywordCategoryAddress,
    }
}

func (o *Onapi) GetFoundKeywordCategories() []domain.KeywordCategory {
    return []domain.KeywordCategory{
        domain.KeywordCategoryCompanyName,  // From Titular
        domain.KeywordCategoryPersonName,   // From Gestor
        domain.KeywordCategoryAddress,      // From Domicilio
    }
}

func (o *Onapi) GetDataByCategory(entity domain.Entity, category domain.KeywordCategory) []string {
    switch category {
    case domain.KeywordCategoryCompanyName:
        return []string{entity.Titular}
    case domain.KeywordCategoryPersonName:
        return []string{entity.Gestor}
    case domain.KeywordCategoryAddress:
        return []string{entity.Domicilio}
    }
    return []string{}
}
```

#### **Keyword Categories** (Value Objects)

```go
type KeywordCategory string

const (
    KeywordCategoryAddress       KeywordCategory = "address"
    KeywordCategoryCompanyName   KeywordCategory = "company_name"
    KeywordCategoryContributorID KeywordCategory = "contributor_id"
    KeywordCategoryPersonName    KeywordCategory = "person_name"
    KeywordCategorySocialMedia   KeywordCategory = "social_media"
    KeywordCategoryFileType      KeywordCategory = "file_type"
    KeywordCategoryXSocialMedia  KeywordCategory = "x_social_media"
)
```

These categories represent the **Ubiquitous Language** - terms understood by both developers and domain experts (investigators, legal professionals).

#### **Dynamic Pipeline Aggregate**

The `DynamicPipelineResult` is an aggregate root:

```go
type DynamicPipelineResult struct {
    ID              ID
    Steps           []DynamicPipelineStep  // Child entities
    TotalSteps      int
    SuccessfulSteps int
    FailedSteps     int
    MaxDepthReached int
    Config          DynamicPipelineConfig  // Value object
}
```

**Aggregate Rules**:
- Only `DynamicPipelineResult` can be accessed from outside the aggregate
- `DynamicPipelineStep` entities are managed within the aggregate
- The aggregate ensures consistency (e.g., updating counters when steps complete)

### 2. **Repository Layer** (`internal/repositories/`)

Repositories provide a clean abstraction over data persistence, following the Repository Pattern from DDD.

#### **Repository Interfaces** (`repositories/interfaces.go`)

```go
// Base repository with CRUD operations
type BaseRepository[T any] interface {
    Create(ctx context.Context, entity T) error
    GetByID(ctx context.Context, id string) (T, error)
    Update(ctx context.Context, id string, entity T) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, offset, limit int) ([]T, error)
    Count(ctx context.Context) (int64, error)
}

// Searchable repository extends base with search capabilities
type SearchableRepository[T any] interface {
    BaseRepository[T]
    Search(ctx context.Context, query string, offset, limit int) ([]T, error)
    SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]T, error)
}

// Domain-specific repository with domain-aware operations
type DomainRepository[T any] interface {
    SearchableRepository[T]
    GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]T, error)
    GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]T, error)
    GetKeywordsByCategory(ctx context.Context, entityID string) (map[domain.KeywordCategory][]string, error)
}
```

**DDD Benefits**:
- **Persistence Ignorance**: Domain entities don't know about databases
- **Testability**: Easy to mock repositories for unit testing
- **Flexibility**: Can swap database implementations without changing domain code

#### **Repository Factory** (`repositories/factory.go`)

The factory pattern centralizes repository creation:

```go
type RepositoryFactory struct {
    db database.Service
}

func (f *RepositoryFactory) GetOnapiRepository() *OnapiRepository {
    return NewOnapiRepository(f.db)
}
// ... other repositories
```

**DDD Principle**: Dependency Injection - repositories are created and injected, not instantiated directly in domain code.

### 3. **Application Layer** (`internal/interactor/`)

The interactor (Application Service) orchestrates domain objects to fulfill use cases.

#### **DynamicPipelineInteractor** (`interactor/dymanic.go`)

```go
type DynamicPipelineInteractor struct {
    repositories *repositories.RepositoryFactory
}

func (d *DynamicPipelineInteractor) ExecuteDynamicPipeline(
    ctx context.Context,
    query string,
    maxDepth int,
    skipDuplicates bool,
) (*domain.DynamicPipelineResult, error) {
    // 1. Create pipeline configuration (value object)
    config := domain.DynamicPipelineConfig{
        Query:              query,
        MaxDepth:           maxDepth,
        MaxConcurrentSteps: 10,
        DelayBetweenSteps:  2,
        SkipDuplicates:     skipDuplicates,
        AvailableDomains:   domain.AllDomainTypes(),
    }

    // 2. Execute pipeline using domain services
    return d.executeStreamingPipeline(ctx, query, availableDomains, config, stepChan)
}
```

**DDD Responsibilities**:
- **Orchestration**: Coordinates multiple domain objects
- **Transaction Management**: Ensures data consistency
- **Use Case Implementation**: Encapsulates business workflows

**Key Workflow**:
1. Create pipeline aggregate using domain service (`module.CreateDynamicPipeline`)
2. Execute steps using domain service (`module.SearchDomain`)
3. Extract keywords using domain function (`domain.GetCategoryByKeywords`)
4. Generate new steps based on domain rules
5. Persist results through repositories

### 4. **Module Layer** (`internal/module/`)

The module layer contains domain services - operations that don't belong to a single entity.

#### **SearchDomain** (`module/dynamic.go`)

```go
func SearchDomain(domainType domain.DomainType, params domain.DomainSearchParams) (*domain.DomainSearchResult, error) {
    // Factory pattern to create domain connector
    switch domainType {
    case domain.DomainTypeONAPI:
        onapi := NewOnapiDomain()
        output, searchErr = onapi.Search(params.Query)
    case domain.DomainTypeSCJ:
        scj := NewScjDomain()
        output, searchErr = scj.Search(params.Query)
    // ... other domains
    }

    // Extract keywords using domain function
    keywordsPerCategory := domain.GetCategoryByKeywords(connector, output)

    return &domain.DomainSearchResult{
        Success:             searchErr == nil,
        DomainType:          domainType,
        SearchParameter:     params.Query,
        KeywordsPerCategory: keywordsPerCategory,
        Output:              output,
    }, searchErr
}
```

**DDD Role**: Domain Service - coordinates multiple domain objects to perform a complex operation.

#### **CreateDynamicPipeline** (`module/dynamic.go`)

```go
func CreateDynamicPipeline(
    ctx context.Context,
    initialQuery string,
    availableDomains []domain.DomainType,
    config domain.DynamicPipelineConfig,
) (*domain.DynamicPipelineResult, error) {
    // Create aggregate root
    pipeline := &domain.DynamicPipelineResult{
        ID:     domain.NewID(),
        Steps:  make([]domain.DynamicPipelineStep, 0),
        Config: config,
    }

    // Domain logic: Map domains to initial categories
    initialDomainCategories := map[domain.DomainType]domain.KeywordCategory{
        domain.DomainTypeONAPI:         domain.KeywordCategoryCompanyName,
        domain.DomainTypeDGII:          domain.KeywordCategoryContributorID,
        domain.DomainTypePGR:           domain.KeywordCategoryPersonName,
        domain.DomainTypeSCJ:           domain.KeywordCategoryContributorID,
        domain.DomainTypeGoogleDocking: domain.KeywordCategoryCompanyName,
    }

    // Create initial steps based on domain knowledge
    for _, domainType := range availableDomains {
        if category, ok := initialDomainCategories[domainType]; ok {
            pipeline.Steps = append(pipeline.Steps, domain.DynamicPipelineStep{
                DomainType:      domainType,
                SearchParameter: initialQuery,
                Category:        category,
                Keywords:        []string{initialQuery},
                Depth:           0,
            })
        }
    }

    return pipeline, nil
}
```

**DDD Principle**: Domain Knowledge - the mapping of domains to initial categories is business logic, not infrastructure.

### 5. **Infrastructure Layer** (`internal/server/`, `internal/database/`)

The infrastructure layer handles technical concerns without polluting domain logic.

#### **HTTP Handlers** (`server/routes.go`)

HTTP handlers are thin adapters that:
1. Parse HTTP requests
2. Call application services (interactors)
3. Format responses

```go
func (s *Server) dynamicPipelineHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    maxDepth := parseMaxDepth(r)
    
    // Delegate to application service
    _, err := s.interactor.ExecuteDynamicPipeline(ctx, query, maxDepth, skipDuplicates)
    
    // Format response
    json.NewEncoder(w).Encode(response)
}
```

**DDD Principle**: **Hexagonal Architecture** - the application core (domain + application layers) doesn't depend on HTTP. HTTP is just one way to access the application.

#### **Database Adapter** (`repositories/database_adapter.go`)

The database adapter implements repository interfaces:

```go
type DatabaseAdapter struct {
    db *sql.DB
}

func (d *DatabaseAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    return d.db.ExecContext(ctx, query, args...)
}
```

**DDD Benefit**: The domain layer doesn't know about SQL - it only knows about repositories (interfaces).

### 6. **DDD Patterns Used**

#### **Factory Pattern**
- `RepositoryFactory`: Creates repository instances
- `CreateDomainConnector()`: Creates domain connector instances
- `NewOnapiDomain()`, `NewScjDomain()`, etc.: Domain object factories

#### **Strategy Pattern**
- `DomainConnector[T]` interface: Different domains implement different search strategies
- Each domain connector has its own search algorithm

#### **Template Method Pattern**
- `GetCategoryByKeywords()`: Generic algorithm that delegates to domain-specific `GetDataByCategory()`

#### **Aggregate Pattern**
- `DynamicPipelineResult`: Aggregate root
- `DynamicPipelineStep`: Child entities within the aggregate
- Repository only exposes the aggregate root

#### **Repository Pattern**
- All data access through repository interfaces
- Domain entities never directly access the database

#### **Application Service Pattern**
- `DynamicPipelineInteractor`: Orchestrates domain objects for use cases
- Thin layer that coordinates repositories and domain services

### 7. **Bounded Contexts**

While this is a single application, it effectively has multiple bounded contexts:

1. **Search Context**: Domain connectors, search operations
2. **Pipeline Context**: Dynamic pipeline execution, step management
3. **Data Access Context**: Repository implementations, database schema

Each context has its own models and language, though they share some common concepts (DomainType, KeywordCategory).

---

## Architecture

### Layered Architecture

```
┌─────────────────────────────────────────┐
│         Presentation Layer               │
│  (HTTP Handlers, React Frontend)        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Application Layer                  │
│  (Interactors, Use Cases)              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Domain Layer                     │
│  (Entities, Value Objects, Services)   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Infrastructure Layer               │
│  (Repositories, Database, HTTP)        │
└─────────────────────────────────────────┘
```

### Data Flow

1. **Request Flow**:
   ```
   HTTP Request → Handler → Interactor → Domain Service → Domain Entity
                                                          ↓
   Response ← Handler ← Interactor ← Repository ← Database
   ```

2. **Pipeline Execution Flow**:
   ```
   User Query → Create Pipeline → Execute Step → Extract Keywords
                                                    ↓
   Generate New Steps ← Domain Logic ← Store Results
   ```

### Key Design Decisions

1. **Generic Domain Connector Interface**: Allows adding new domains without changing core pipeline logic
2. **Keyword Category System**: Enables cross-domain keyword propagation
3. **Streaming Architecture**: Real-time results via Server-Sent Events for better UX
4. **Repository Factory**: Centralized repository management with dependency injection
5. **Context-Based Execution**: Supports cancellation and timeout via Go contexts

---

## Key Features

### 1. **Dynamic Pipeline System**
- Automatically generates search steps based on extracted keywords
- Configurable depth and duplicate prevention
- Real-time streaming of results

### 2. **Multi-Domain Search**
- Unified interface for searching across 7+ data sources
- Parallel execution for performance
- Consistent result format

### 3. **Keyword Extraction & Categorization**
- Automatic extraction of relevant keywords from search results
- Categorization (company name, person name, address, etc.)
- Cross-domain keyword propagation

### 4. **Persistent Storage**
- All search results and pipeline executions stored in MySQL
- Full audit trail with timestamps
- Queryable history of investigations

### 5. **RESTful API**
- Clean REST endpoints for all operations
- JSON responses
- CORS support for frontend integration

### 6. **Real-Time Streaming**
- Server-Sent Events (SSE) for live pipeline updates
- Progressive result delivery
- Better user experience for long-running operations

### 7. **CLI Tool**
- Command-line interface for automated/scripted usage
- Same core logic as API
- Useful for batch processing

### 8. **Type Safety**
- Go's strong typing prevents runtime errors
- TypeScript on frontend ensures UI consistency
- Generic interfaces maintain type safety across domains

---

## Conclusion

Insightful Intel demonstrates a well-structured Domain-Driven Design implementation with:
- **Clear separation of concerns** across layers
- **Rich domain models** that encapsulate business logic
- **Repository pattern** for data access abstraction
- **Application services** for use case orchestration
- **Domain services** for cross-entity operations
- **Ubiquitous language** reflected in code (KeywordCategory, DomainType, etc.)

The architecture is maintainable, testable, and extensible - new domains can be added by implementing the `DomainConnector` interface without modifying core pipeline logic.

