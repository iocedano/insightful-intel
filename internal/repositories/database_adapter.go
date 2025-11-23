package repositories

import (
	"context"
	"database/sql"
	"insightful-intel/internal/database"
)

// databaseAdapter adapts database.Service to DatabaseAccessor
type databaseAdapter struct {
	db database.Service
}

func (da *databaseAdapter) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return da.db.GetDB().ExecContext(ctx, query, args...)
}

func (da *databaseAdapter) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return da.db.GetDB().QueryContext(ctx, query, args...)
}

func (da *databaseAdapter) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return da.db.GetDB().QueryRowContext(ctx, query, args...)
}

// NewDatabaseAdapter creates a new database adapter
func NewDatabaseAdapter(db database.Service) DatabaseAccessor {
	return &databaseAdapter{db: db}
}


