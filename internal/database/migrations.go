package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed schema.sql
var schemaFS embed.FS

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationService handles database migrations
type MigrationService struct {
	db *sql.DB
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *sql.DB) *MigrationService {
	return &MigrationService{
		db: db,
	}
}

// CreateMigrationsTable creates the migrations tracking table
func (m *MigrationService) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			version INT NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_version (version)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`

	_, err := m.db.Exec(query)
	return err
}

// GetExecutedMigrations returns a list of executed migration versions
func (m *MigrationService) GetExecutedMigrations() (map[int]bool, error) {
	query := `SELECT version FROM migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		executed[version] = true
	}

	return executed, nil
}

// ExecuteMigration executes a single migration
func (m *MigrationService) ExecuteMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the migration SQL statements
	statements := strings.Split(migration.UpSQL, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// Execute each statement individually
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	// Record the migration as executed
	recordQuery := `INSERT INTO migrations (version, name) VALUES (?, ?)`
	if _, err := tx.Exec(recordQuery, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

// RollbackMigration rolls back a single migration
func (m *MigrationService) RollbackMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the rollback SQL
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return fmt.Errorf("failed to rollback migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	// Remove the migration record
	deleteQuery := `DELETE FROM migrations WHERE version = ?`
	if _, err := tx.Exec(deleteQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d (%s): %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

// RunMigrations executes all pending migrations
func (m *MigrationService) RunMigrations(migrations []Migration) error {
	// Create migrations table if it doesn't exist
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get executed migrations
	executed, err := m.GetExecutedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Execute pending migrations
	for _, migration := range migrations {
		if !executed[migration.Version] {
			log.Printf("Executing migration %d: %s", migration.Version, migration.Name)
			if err := m.ExecuteMigration(migration); err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}
			log.Printf("Migration %d completed successfully", migration.Version)
		} else {
			log.Printf("Migration %d already executed, skipping", migration.Version)
		}
	}

	return nil
}

// getEmbeddedSchemaSQL returns the embedded schema.sql content
func getEmbeddedSchemaSQL() (string, error) {
	content, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetInitialMigrations returns the initial set of migrations
func GetInitialMigrations() []Migration {
	// Load schema from embedded file
	schemaSQL, err := getEmbeddedSchemaSQL()
	if err != nil {
		log.Printf("Warning: Could not load embedded schema, using fallback schema: %v", err)
		// Fallback to embedded schema
		return []Migration{}
	}

	// Clean up the SQL content
	cleanSQL := cleanSQLContent(schemaSQL)

	return []Migration{
		{
			Version: 1,
			Name:    "create_repository_tables",
			UpSQL:   cleanSQL,
			DownSQL: getRepositorySchemaRollbackSQL(),
		},
	}
}

// getRepositorySchemaRollbackSQL returns the SQL for rolling back repository tables
func getRepositorySchemaRollbackSQL() string {
	return `
		DROP TABLE IF EXISTS dynamic_pipeline_steps;
		DROP TABLE IF EXISTS dynamic_pipeline_results;
		DROP TABLE IF EXISTS google_docking_results;
		DROP TABLE IF EXISTS pgr_news;
		DROP TABLE IF EXISTS dgii_registers;
		DROP TABLE IF EXISTS scj_cases;
		DROP TABLE IF EXISTS onapi_entities;
		DROP TABLE IF EXISTS domain_search_results;
	`
}

// cleanSQLContent cleans up SQL content by removing comments and empty lines
func cleanSQLContent(sqlContent string) string {
	// Clean up the SQL content
	sqlContent = strings.TrimSpace(sqlContent)

	// Remove comments and empty lines while preserving statement boundaries
	lines := strings.Split(sqlContent, "\n")
	var cleanLines []string
	inMultiLineComment := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle multi-line comments
		if strings.Contains(line, "/*") {
			inMultiLineComment = true
		}
		if strings.Contains(line, "*/") {
			inMultiLineComment = false
			continue
		}
		if inMultiLineComment {
			continue
		}

		// Skip single-line comments
		if strings.HasPrefix(line, "--") {
			continue
		}

		cleanLines = append(cleanLines, line)
	}

	// Join the clean lines
	return strings.Join(cleanLines, "\n")
}

// LoadMigrationsFromFile loads migrations from a SQL file
func LoadMigrationsFromFile(filePath string) ([]Migration, error) {
	// This function is kept for backward compatibility but now uses embedded files
	// For now, we'll return an error to indicate this method is deprecated
	return nil, fmt.Errorf("LoadMigrationsFromFile is deprecated, use embedded schema instead")
}
