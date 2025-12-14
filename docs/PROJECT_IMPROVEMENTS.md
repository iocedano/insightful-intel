# Project Improvements Guide

This document outlines recommended improvements for the Insightful Intel project, covering database migrations, structure, domain implementation patterns, testing strategies, and more.

## Table of Contents

1. [Database Migrations](#database-migrations)
2. [Database Structure Improvements](#database-structure-improvements)
3. [Improved Domain Implementation](#improved-domain-implementation)
4. [Testing Improvements](#testing-improvements)
5. [Code Generation & Templates](#code-generation--templates)
6. [Error Handling & Logging](#error-handling--logging)
7. [Configuration Management](#configuration-management)
8. [Performance Optimizations](#performance-optimizations)
9. [Security Enhancements](#security-enhancements)
10. [Documentation & Developer Experience](#documentation--developer-experience)

---

## Database Migrations

### Current State

The project currently uses a basic migration system with:
- Single migration file (`schema.sql`)
- Embedded migrations in Go code
- Basic up/down SQL support
- Version tracking in `migrations` table

### Recommended Improvements

#### 1. File-Based Migration System

**Structure**:
```
internal/database/migrations/
├── 0001_initial_schema.up.sql
├── 0001_initial_schema.down.sql
├── 0002_add_indexes.up.sql
├── 0002_add_indexes.down.sql
├── 0003_add_new_domain_table.up.sql
├── 0003_add_new_domain_table.down.sql
└── ...
```

**Implementation**:

```go
// internal/database/migrations.go

package database

import (
    "database/sql"
    "embed"
    "fmt"
    "io/fs"
    "log"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type MigrationFile struct {
    Version int
    Name    string
    UpSQL   string
    DownSQL string
}

// LoadMigrationsFromFS loads migrations from embedded filesystem
func LoadMigrationsFromFS() ([]MigrationFile, error) {
    migrations := make(map[int]*MigrationFile)
    
    err := fs.WalkDir(migrationsFS, "migrations", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        
        if d.IsDir() {
            return nil
        }
        
        // Parse filename: 0001_name.up.sql or 0001_name.down.sql
        baseName := filepath.Base(path)
        parts := strings.Split(baseName, "_")
        if len(parts) < 3 {
            return nil // Skip invalid files
        }
        
        version, err := strconv.Atoi(parts[0])
        if err != nil {
            return nil // Skip invalid version
        }
        
        direction := strings.TrimSuffix(parts[len(parts)-1], ".sql")
        name := strings.Join(parts[1:len(parts)-1], "_")
        
        if migrations[version] == nil {
            migrations[version] = &MigrationFile{
                Version: version,
                Name:    name,
            }
        }
        
        content, err := migrationsFS.ReadFile(path)
        if err != nil {
            return err
        }
        
        if direction == "up" {
            migrations[version].UpSQL = string(content)
        } else if direction == "down" {
            migrations[version].DownSQL = string(content)
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Convert map to sorted slice
    result := make([]MigrationFile, 0, len(migrations))
    for _, m := range migrations {
        result = append(result, *m)
    }
    
    sort.Slice(result, func(i, j int) bool {
        return result[i].Version < result[j].Version
    })
    
    return result, nil
}

// Enhanced MigrationService with better error handling
type EnhancedMigrationService struct {
    db *sql.DB
}

func (m *EnhancedMigrationService) RunMigrations() error {
    migrations, err := LoadMigrationsFromFS()
    if err != nil {
        return fmt.Errorf("failed to load migrations: %w", err)
    }
    
    if err := m.CreateMigrationsTable(); err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }
    
    executed, err := m.GetExecutedMigrations()
    if err != nil {
        return fmt.Errorf("failed to get executed migrations: %w", err)
    }
    
    for _, migration := range migrations {
        if executed[migration.Version] {
            log.Printf("Migration %d (%s) already executed, skipping", migration.Version, migration.Name)
            continue
        }
        
        log.Printf("Executing migration %d: %s", migration.Version, migration.Name)
        if err := m.ExecuteMigration(migration); err != nil {
            return fmt.Errorf("migration %d failed: %w", migration.Version, err)
        }
        log.Printf("Migration %d completed successfully", migration.Version)
    }
    
    return nil
}

// RollbackLastMigration rolls back the last executed migration
func (m *EnhancedMigrationService) RollbackLastMigration() error {
    query := `SELECT version, name FROM migrations ORDER BY version DESC LIMIT 1`
    var version int
    var name string
    
    err := m.db.QueryRow(query).Scan(&version, &name)
    if err == sql.ErrNoRows {
        return fmt.Errorf("no migrations to rollback")
    }
    if err != nil {
        return err
    }
    
    migrations, err := LoadMigrationsFromFS()
    if err != nil {
        return err
    }
    
    var migration MigrationFile
    for _, m := range migrations {
        if m.Version == version {
            migration = m
            break
        }
    }
    
    if migration.Version == 0 {
        return fmt.Errorf("migration %d not found", version)
    }
    
    log.Printf("Rolling back migration %d: %s", version, name)
    return m.RollbackMigration(migration)
}

// RollbackToVersion rolls back migrations down to a specific version
func (m *EnhancedMigrationService) RollbackToVersion(targetVersion int) error {
    executed, err := m.GetExecutedMigrations()
    if err != nil {
        return err
    }
    
    migrations, err := LoadMigrationsFromFS()
    if err != nil {
        return err
    }
    
    // Get migrations to rollback (in reverse order)
    toRollback := make([]MigrationFile, 0)
    for _, migration := range migrations {
        if migration.Version > targetVersion && executed[migration.Version] {
            toRollback = append(toRollback, migration)
        }
    }
    
    // Sort in reverse order
    sort.Slice(toRollback, func(i, j int) bool {
        return toRollback[i].Version > toRollback[j].Version
    })
    
    for _, migration := range toRollback {
        log.Printf("Rolling back migration %d: %s", migration.Version, migration.Name)
        if err := m.RollbackMigration(migration); err != nil {
            return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
        }
    }
    
    return nil
}
```

**Example Migration Files**:

```sql
-- migrations/0003_add_new_domain_table.up.sql
CREATE TABLE IF NOT EXISTS new_domain_entities (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    name VARCHAR(255) NOT NULL,
    identifier VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_name (name),
    INDEX idx_identifier (identifier)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

```sql
-- migrations/0003_add_new_domain_table.down.sql
DROP TABLE IF EXISTS new_domain_entities;
```

#### 2. Migration CLI Commands

Add to `cmd/cli/main.go`:

```go
var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Database migration commands",
    Long:  "Manage database migrations",
}

var migrateUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Run pending migrations",
    Run: func(cmd *cobra.Command, args []string) {
        db := database.New()
        migrationService := database.NewEnhancedMigrationService(db.GetDB())
        if err := migrationService.RunMigrations(); err != nil {
            log.Fatalf("Migration failed: %v", err)
        }
        log.Println("Migrations completed successfully")
    },
}

var migrateDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Rollback last migration",
    Run: func(cmd *cobra.Command, args []string) {
        db := database.New()
        migrationService := database.NewEnhancedMigrationService(db.GetDB())
        if err := migrationService.RollbackLastMigration(); err != nil {
            log.Fatalf("Rollback failed: %v", err)
        }
        log.Println("Rollback completed successfully")
    },
}

var migrateStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show migration status",
    Run: func(cmd *cobra.Command, args []string) {
        db := database.New()
        migrationService := database.NewEnhancedMigrationService(db.GetDB())
        status, err := migrationService.GetMigrationStatus()
        if err != nil {
            log.Fatalf("Failed to get status: %v", err)
        }
        // Print status table
    },
}
```

---

## Database Structure Improvements

### 1. Normalization Improvements

**Current Issues**:
- Some tables may have redundant data
- JSON columns used where normalized tables would be better
- Missing relationships between entities

**Recommendations**:

```sql
-- Separate keywords table instead of JSON
CREATE TABLE IF NOT EXISTS keyword_extractions (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36) NOT NULL,
    category VARCHAR(50) NOT NULL,
    keyword VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_category (category),
    INDEX idx_keyword (keyword),
    UNIQUE KEY unique_extraction (domain_search_result_id, category, keyword)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Separate pipeline configuration table
CREATE TABLE IF NOT EXISTS pipeline_configurations (
    id CHAR(36) PRIMARY KEY,
    pipeline_id CHAR(36) NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (pipeline_id) REFERENCES dynamic_pipeline_results(id) ON DELETE CASCADE,
    INDEX idx_pipeline_id (pipeline_id),
    INDEX idx_config_key (config_key),
    UNIQUE KEY unique_config (pipeline_id, config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2. Audit Trail Enhancement

```sql
-- Add audit log table
CREATE TABLE IF NOT EXISTS audit_logs (
    id CHAR(36) PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id CHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL, -- CREATE, UPDATE, DELETE
    old_values JSON,
    new_values JSON,
    user_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 3. Soft Deletes

```sql
-- Add deleted_at column to all entity tables
ALTER TABLE onapi_entities ADD COLUMN deleted_at TIMESTAMP NULL;
ALTER TABLE scj_cases ADD COLUMN deleted_at TIMESTAMP NULL;
ALTER TABLE dgii_registers ADD COLUMN deleted_at TIMESTAMP NULL;
-- ... etc

-- Add indexes
CREATE INDEX idx_deleted_at ON onapi_entities(deleted_at);
```

### 4. Database Indexing Strategy

```sql
-- Composite indexes for common queries
CREATE INDEX idx_domain_type_success ON domain_search_results(domain_type, success);
CREATE INDEX idx_pipeline_domain_depth ON dynamic_pipeline_steps(pipeline_id, domain_type, depth);
CREATE INDEX idx_search_parameter_domain ON domain_search_results(search_parameter, domain_type);

-- Full-text search indexes (MySQL 5.6+)
ALTER TABLE onapi_entities ADD FULLTEXT INDEX ft_texto (texto);
ALTER TABLE scj_cases ADD FULLTEXT INDEX ft_involucrados (involucrados);
```

---

## Improved Domain Implementation

### 1. Domain Generator Tool

Create a code generator to scaffold new domains:

```go
// cmd/tools/generate_domain/main.go
package main

import (
    "fmt"
    "os"
    "strings"
    "text/template"
)

type DomainTemplate struct {
    DomainName      string
    DomainType      string
    DomainTypeConst string
    EntityName      string
    TableName       string
    Fields          []Field
}

type Field struct {
    Name       string
    Type       string
    JSONTag    string
    DBTag      string
    IsRequired bool
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: generate_domain <domain_name>")
        os.Exit(1)
    }
    
    domainName := strings.ToLower(os.Args[1])
    domainType := strings.ToUpper(domainName)
    entityName := toPascalCase(domainName)
    
    template := DomainTemplate{
        DomainName:      domainName,
        DomainType:      domainType,
        DomainTypeConst: "DomainType" + entityName,
        EntityName:      entityName,
        TableName:       domainName + "_entities",
        Fields: []Field{
            {Name: "Name", Type: "string", JSONTag: "name", IsRequired: true},
            {Name: "Identifier", Type: "string", JSONTag: "identifier", IsRequired: true},
        },
    }
    
    generateFiles(template)
}

func generateFiles(tmpl DomainTemplate) {
    // Generate domain model
    generateDomainModel(tmpl)
    // Generate connector
    generateConnector(tmpl)
    // Generate repository
    generateRepository(tmpl)
    // Generate migration
    generateMigration(tmpl)
    // Update factory
    updateFactory(tmpl)
}
```

### 2. Domain Interface Enhancements

```go
// Enhanced domain connector with additional capabilities
type EnhancedDomainConnector[T any] interface {
    DomainConnector[T]
    
    // Batch operations
    SearchBatch(queries []string) ([]T, error)
    
    // Caching support
    GetCacheKey(query string) string
    GetCacheTTL() time.Duration
    
    // Rate limiting
    GetRateLimit() RateLimit
    
    // Health check
    HealthCheck() error
}

type RateLimit struct {
    RequestsPerSecond int
    Burst             int
}

// Base implementation
type BaseDomainConnector[T any] struct {
    CacheEnabled bool
    CacheTTL     time.Duration
    RateLimiter  *rate.Limiter
}

func (b *BaseDomainConnector[T]) GetCacheKey(query string) string {
    return fmt.Sprintf("domain:%s:query:%s", b.GetDomainType(), query)
}

func (b *BaseDomainConnector[T]) GetCacheTTL() time.Duration {
    return b.CacheTTL
}
```

### 3. Domain Registry Pattern

```go
// internal/domain/registry.go
package domain

import "sync"

type DomainRegistry struct {
    connectors map[DomainType]any
    mu         sync.RWMutex
}

var globalRegistry = &DomainRegistry{
    connectors: make(map[DomainType]any),
}

func RegisterDomain[T any](domainType DomainType, connector DomainConnector[T]) {
    globalRegistry.mu.Lock()
    defer globalRegistry.mu.Unlock()
    globalRegistry.connectors[domainType] = connector
}

func GetDomain[T any](domainType DomainType) (DomainConnector[T], error) {
    globalRegistry.mu.RLock()
    defer globalRegistry.mu.RUnlock()
    
    connector, ok := globalRegistry.connectors[domainType]
    if !ok {
        return nil, fmt.Errorf("domain %s not registered", domainType)
    }
    
    typedConnector, ok := connector.(DomainConnector[T])
    if !ok {
        return nil, fmt.Errorf("domain %s type mismatch", domainType)
    }
    
    return typedConnector, nil
}

// Auto-registration on init
func init() {
    RegisterDomain(DomainTypeONAPI, NewOnapiDomain())
    RegisterDomain(DomainTypeSCJ, NewScjDomain())
    // ... etc
}
```

---

## Testing Improvements

### 1. Test Structure

```
internal/
├── domain/
│   ├── onapi.go
│   ├── onapi_test.go
│   └── ...
├── module/
│   ├── onapi.go
│   ├── onapi_test.go
│   └── ...
└── repositories/
    ├── onapi.go
    ├── onapi_test.go
    └── ...
```

### 2. Test Utilities

```go
// internal/testing/testutil.go
package testing

import (
    "context"
    "database/sql"
    "insightful-intel/internal/database"
    "testing"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// SetupTestDB creates a test database using testcontainers
func SetupTestDB(t *testing.T) (database.Service, func()) {
    ctx := context.Background()
    
    mysqlContainer, err := mysql.RunContainer(ctx,
        testcontainers.WithImage("mysql:8.0"),
        mysql.WithDatabase("test_db"),
        mysql.WithUsername("test"),
        mysql.WithPassword("test"),
    )
    if err != nil {
        t.Fatalf("failed to start mysql container: %v", err)
    }
    
    connStr, err := mysqlContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatalf("failed to get connection string: %v", err)
    }
    
    db := database.NewWithConnectionString(connStr)
    
    // Run migrations
    migrationService := database.NewMigrationService(db.GetDB())
    migrations := database.GetInitialMigrations()
    if err := migrationService.RunMigrations(migrations); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }
    
    cleanup := func() {
        if err := mysqlContainer.Terminate(ctx); err != nil {
            t.Logf("failed to terminate container: %v", err)
        }
    }
    
    return db, cleanup
}

// CreateTestEntity creates a test entity
func CreateTestEntity(t *testing.T, repo interface{}, entity interface{}) {
    // Generic entity creation helper
}

// AssertEntityEqual compares two entities
func AssertEntityEqual(t *testing.T, expected, actual interface{}) {
    // Deep comparison helper
}
```

### 3. Integration Test Example

```go
// internal/repositories/onapi_test.go
package repositories

import (
    "context"
    "testing"
    "insightful-intel/internal/domain"
    "insightful-intel/internal/testing"
)

func TestOnapiRepository_Create(t *testing.T) {
    db, cleanup := testing.SetupTestDB(t)
    defer cleanup()
    
    repo := NewOnapiRepository(db)
    
    entity := domain.Entity{
        SerieExpediente:  1,
        NumeroExpediente: 123,
        Titular:         "Test Company",
    }
    
    err := repo.Create(context.Background(), entity)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    
    retrieved, err := repo.GetByID(context.Background(), entity.ID.String())
    if err != nil {
        t.Fatalf("GetByID failed: %v", err)
    }
    
    if retrieved.Titular != entity.Titular {
        t.Errorf("Expected Titular %s, got %s", entity.Titular, retrieved.Titular)
    }
}
```

### 4. Mock Generation

```go
// Use gomock or similar for interface mocking
//go:generate mockgen -source=connector.go -destination=mocks/connector_mock.go

// Usage in tests
func TestSearchDomain_WithMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockConnector := mocks.NewMockDomainConnector(ctrl)
    mockConnector.EXPECT().
        Search("test").
        Return([]domain.Entity{{Titular: "Test"}}, nil)
    
    // Test with mock
}
```

### 5. Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Coverage threshold in CI
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}'
```

---

## Code Generation & Templates

### 1. Domain Template Files

Create template files for code generation:

```go
// templates/domain_model.tmpl
package domain

import "time"

// {{.EntityName}} represents an entity from the {{.DomainName}} domain
type {{.EntityName}} struct {
    ID                   ID        `json:"id"`
    DomainSearchResultID ID        `json:"domain_search_result_id"`
    {{range .Fields}}
    {{.Name}} {{.Type}} `json:"{{.JSONTag}}"`{{if .IsRequired}} // Required{{end}}
    {{end}}
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}
```

### 2. Makefile Targets

```makefile
# Generate new domain
generate-domain:
	@read -p "Enter domain name: " domain; \
	go run cmd/tools/generate_domain/main.go $$domain

# Generate mocks
generate-mocks:
	@go generate ./...

# Generate API docs
generate-docs:
	@swag init -g cmd/api/main.go -o docs/api
```

---

## Error Handling & Logging

### 1. Structured Logging

```go
// internal/logger/logger.go
package logger

import (
    "log/slog"
    "os"
)

var Logger *slog.Logger

func InitLogger(level slog.Level) {
    opts := &slog.HandlerOptions{
        Level: level,
        AddSource: true,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    Logger = slog.New(handler)
}

// Usage
logger.Logger.Info("Domain search started",
    "domain", domainType,
    "query", query,
    "user_id", userID,
)
```

### 2. Error Wrapping

```go
// internal/errors/errors.go
package errors

import "fmt"

type DomainError struct {
    Domain    string
    Operation string
    Err       error
}

func (e *DomainError) Error() string {
    return fmt.Sprintf("%s.%s: %v", e.Domain, e.Operation, e.Err)
}

func (e *DomainError) Unwrap() error {
    return e.Err
}

// Usage
func (d *Domain) Search(query string) ([]Entity, error) {
    results, err := d.performSearch(query)
    if err != nil {
        return nil, &DomainError{
            Domain:    string(d.GetDomainType()),
            Operation: "Search",
            Err:       fmt.Errorf("failed to search: %w", err),
        }
    }
    return results, nil
}
```

### 3. Error Codes

```go
// internal/errors/codes.go
package errors

type ErrorCode string

const (
    ErrCodeDomainNotFound    ErrorCode = "DOMAIN_NOT_FOUND"
    ErrCodeInvalidQuery     ErrorCode = "INVALID_QUERY"
    ErrCodeDatabaseError    ErrorCode = "DATABASE_ERROR"
    ErrCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
)

type CodedError struct {
    Code    ErrorCode
    Message string
    Err     error
}

func (e *CodedError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
}
```

---

## Configuration Management

### 1. Environment-Based Config

```go
// config/config.go
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    Domains  DomainsConfig
    Cache    CacheConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
    MaxConns int
}

func Load() (*Config, error) {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvInt("DB_PORT", 3306),
            User:     getEnv("DB_USER", "root"),
            Password: getEnv("DB_PASSWORD", ""),
            Name:     getEnv("DB_NAME", "insightful_intel"),
            MaxConns: getEnvInt("DB_MAX_CONNS", 10),
        },
        Server: ServerConfig{
            Port:         getEnvInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
            WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
        },
        // ... etc
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 2. Config Validation

```go
func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 {
        return fmt.Errorf("invalid database port")
    }
    // ... more validation
    return nil
}
```

---

## Performance Optimizations

### 1. Connection Pooling

```go
// internal/database/pool.go
func NewWithConfig(config DatabaseConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        config.User, config.Password, config.Host, config.Port, config.Name)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    db.SetMaxOpenConns(config.MaxConns)
    db.SetMaxIdleConns(config.MaxConns / 2)
    db.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}
```

### 2. Caching Layer

```go
// internal/cache/cache.go
package cache

import (
    "context"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type RedisCache struct {
    client *redis.Client
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
    return r.client.Get(ctx, key).Bytes()
}

// Usage in domain connector
func (d *Domain) Search(query string) ([]Entity, error) {
    cacheKey := d.GetCacheKey(query)
    
    // Try cache first
    if cached, err := cache.Get(context.Background(), cacheKey); err == nil {
        var entities []Entity
        if err := json.Unmarshal(cached, &entities); err == nil {
            return entities, nil
        }
    }
    
    // Perform search
    entities, err := d.performSearch(query)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    if data, err := json.Marshal(entities); err == nil {
        cache.Set(context.Background(), cacheKey, data, d.GetCacheTTL())
    }
    
    return entities, nil
}
```

### 3. Batch Operations

```go
// Batch insert for better performance
func (r *Repository) CreateBatch(ctx context.Context, entities []Entity) error {
    if len(entities) == 0 {
        return nil
    }
    
    query := `INSERT INTO table_name (id, field1, field2) VALUES `
    values := make([]interface{}, 0, len(entities)*3)
    
    for i, entity := range entities {
        if i > 0 {
            query += ", "
        }
        query += "(?, ?, ?)"
        values = append(values, entity.ID, entity.Field1, entity.Field2)
    }
    
    _, err := r.db.ExecContext(ctx, query, values...)
    return err
}
```

---

## Security Enhancements

### 1. SQL Injection Prevention

```go
// Always use parameterized queries
func (r *Repository) Search(ctx context.Context, query string) ([]Entity, error) {
    // ❌ BAD: String concatenation
    // sql := "SELECT * FROM table WHERE name = '" + query + "'"
    
    // ✅ GOOD: Parameterized query
    sql := "SELECT * FROM table WHERE name = ?"
    rows, err := r.db.QueryContext(ctx, sql, query)
    // ...
}
```

### 2. Input Validation

```go
// internal/validation/validator.go
package validation

import (
    "regexp"
    "strings"
)

func ValidateQuery(query string) error {
    if len(query) < 2 {
        return fmt.Errorf("query must be at least 2 characters")
    }
    if len(query) > 255 {
        return fmt.Errorf("query too long")
    }
    
    // Remove potentially dangerous characters
    dangerous := regexp.MustCompile(`[<>'"&]`)
    if dangerous.MatchString(query) {
        return fmt.Errorf("query contains invalid characters")
    }
    
    return nil
}
```

### 3. Rate Limiting

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "golang.org/x/time/rate"
    "net/http"
)

func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Documentation & Developer Experience

### 1. API Documentation

```go
// Use swaggo/swag for API docs
// @Summary Search domain
// @Description Search a specific domain with a query
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param domain query string false "Domain type"
// @Success 200 {object} ConnectorPipeline
// @Router /search [get]
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

### 2. Development Scripts

```makefile
# Development helpers
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	@docker-compose up -d
	@go mod download
	@npm install --prefix ./frontend
	@make migrate-up

.PHONY: dev-reset
dev-reset:
	@echo "Resetting development environment..."
	@docker-compose down -v
	@make dev-setup

.PHONY: lint
lint:
	@golangci-lint run
	@npm run lint --prefix ./frontend

.PHONY: format
format:
	@gofmt -w .
	@goimports -w .
	@npm run format --prefix ./frontend
```

### 3. Pre-commit Hooks

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests
make test
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

# Run linter
make lint
if [ $? -ne 0 ]; then
    echo "Linting failed. Commit aborted."
    exit 1
fi

# Format code
make format
git add -u
```

---

## Summary

This document outlines comprehensive improvements for the Insightful Intel project:

1. **Migration System**: File-based migrations with up/down support
2. **Database Structure**: Better normalization, indexing, audit trails
3. **Domain Implementation**: Code generation, registry pattern, enhanced interfaces
4. **Testing**: Test utilities, integration tests, mocks, coverage
5. **Error Handling**: Structured logging, error wrapping, error codes
6. **Configuration**: Environment-based config with validation
7. **Performance**: Connection pooling, caching, batch operations
8. **Security**: SQL injection prevention, input validation, rate limiting
9. **Developer Experience**: API docs, scripts, pre-commit hooks

These improvements will make the codebase more maintainable, testable, and scalable.

