package repositories

import (
	"context"
	"database/sql"
	"insightful-intel/internal/domain"
)

// DatabaseAccessor provides access to the underlying database connection
type DatabaseAccessor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// BaseRepository defines common operations for all repositories
type BaseRepository[T any] interface {
	// Create inserts a new record
	Create(ctx context.Context, entity T) error

	// GetByID retrieves a record by its ID
	GetByID(ctx context.Context, id string) (T, error)

	// Update modifies an existing record
	Update(ctx context.Context, id string, entity T) error

	// Delete removes a record by its ID
	Delete(ctx context.Context, id string) error

	// List retrieves multiple records with pagination
	List(ctx context.Context, offset, limit int) ([]T, error)

	// Count returns the total number of records
	Count(ctx context.Context) (int64, error)
}

// SearchableRepository defines operations for repositories that support search
type SearchableRepository[T any] interface {
	BaseRepository[T]

	// Search performs a search query
	Search(ctx context.Context, query string, offset, limit int) ([]T, error)

	// SearchByCategory performs a search within a specific keyword category
	SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]T, error)
}

// DomainRepository defines operations specific to domain entities
type DomainRepository[T any] interface {
	SearchableRepository[T]

	// GetByDomainType retrieves records by domain type
	GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]T, error)

	// GetBySearchParameter retrieves records by search parameter
	GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]T, error)

	// GetKeywordsByCategory retrieves keywords grouped by category
	GetKeywordsByCategory(ctx context.Context, entityID string) (map[domain.KeywordCategory][]string, error)
}

// PipelineResultRepository defines operations for pipeline results
type PipelineResultRepository interface {
	BaseRepository[any]

	// GetByDomainType retrieves pipeline results by domain type
	GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]any, error)

	// GetBySuccessStatus retrieves pipeline results by success status
	GetBySuccessStatus(ctx context.Context, success bool, offset, limit int) ([]any, error)

	// GetBySearchParameter retrieves pipeline results by search parameter
	GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]any, error)

	// GetKeywordsByCategory retrieves keywords grouped by category for a pipeline result
	GetKeywordsByCategory(ctx context.Context, resultID string) (map[domain.KeywordCategory][]string, error)
}
