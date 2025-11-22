package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// PgrRepository implements DomainRepository for PGR PGRNews domain type
type PgrRepository struct {
	db DatabaseAccessor
}

// NewPgrRepository creates a new PGR repository instance
func NewPgrRepository(db database.Service) *PgrRepository {
	return &PgrRepository{
		db: NewDatabaseAdapter(db),
	}
}

// Create inserts a new PGR news item
func (r *PgrRepository) Create(ctx context.Context, entity domain.PGRNews) error {
	return r.CreateWithDomainSearchResultID(ctx, entity)
}

// CreateWithDomainSearchResultID inserts a new PGR news item with a domain search result ID
func (r *PgrRepository) CreateWithDomainSearchResultID(ctx context.Context, entity domain.PGRNews) error {
	// Check if URL already exists
	exists, existingID, err := r.URLExists(ctx, entity.URL)
	if err != nil {
		return fmt.Errorf("error checking URL existence: %w", err)
	}

	if exists {
		// Update the existing entity instead of creating a new one
		entity.ID = existingID
		return r.Update(ctx, entity)
	}

	// Generate new ID for new entity
	entity.ID = domain.NewID()

	query := `
		INSERT INTO pgr_news (
			id, domain_search_result_id, url, title, created_at, updated_at
		) VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.DomainSearchResultID, entity.URL, entity.Title,
	)

	return err
}

// URLExists checks if a URL already exists in the database
func (r *PgrRepository) URLExists(ctx context.Context, url string) (bool, domain.ID, error) {
	query := `SELECT id FROM pgr_news WHERE url = ? LIMIT 1`

	var id domain.ID
	err := r.db.QueryRowContext(ctx, query, url).Scan(&id)

	if err == sql.ErrNoRows {
		return false, domain.ID{}, nil
	}
	if err != nil {
		return false, domain.ID{}, err
	}

	return true, id, nil
}

// GetByURL retrieves a PGR news item by URL
func (r *PgrRepository) GetByURL(ctx context.Context, url string) (domain.PGRNews, error) {
	query := `
		SELECT id, domain_search_result_id, url, title, created_at, updated_at
		FROM pgr_news 
		WHERE url = ?
	`

	var entity domain.PGRNews

	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&entity.ID, &entity.DomainSearchResultID, &entity.URL, &entity.Title,
	)

	if err != nil {
		return domain.PGRNews{}, err
	}

	return entity, nil
}

// GetByID retrieves a PGR news item by its URL
func (r *PgrRepository) GetByID(ctx context.Context, id string) (domain.PGRNews, error) {
	query := `
		SELECT id, domain_search_result_id, url, title, created_at, updated_at
		FROM pgr_news 
		WHERE id = ?
	`

	var entity domain.PGRNews

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID, &entity.DomainSearchResultID, &entity.URL, &entity.Title,
	)

	if err != nil {
		return domain.PGRNews{}, err
	}

	return entity, nil
}

// Update modifies an existing PGR news item
func (r *PgrRepository) Update(ctx context.Context, entity domain.PGRNews) error {
	query := `
		UPDATE pgr_news SET
			domain_search_result_id = ?, url = ?, title = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.DomainSearchResultID, entity.URL, entity.Title, entity.ID,
	)

	return err
}

// Delete removes a PGR news item by its URL
func (r *PgrRepository) Delete(ctx context.Context, entity domain.PGRNews) error {
	query := `DELETE FROM pgr_news WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, entity.ID)
	return err
}

// List retrieves multiple PGR news items with pagination
func (r *PgrRepository) List(ctx context.Context, offset, limit int) ([]domain.PGRNews, error) {
	query := `
		SELECT id, domain_search_result_id, url, title, created_at, updated_at
		FROM pgr_news 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.PGRNews
	for rows.Next() {
		var entity domain.PGRNews

		err := rows.Scan(
			&entity.ID,
			&entity.DomainSearchResultID,
			&entity.URL,
			&entity.Title,
			new(interface{}), // created_at (ignored)
			new(interface{}), // updated_at (ignored)
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// Count returns the total number of PGR news items
func (r *PgrRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM pgr_news`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Search performs a search query on PGR news items
func (r *PgrRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.PGRNews, error) {
	searchQuery := `
		SELECT id, domain_search_result_id, url, title, created_at, updated_at
		FROM pgr_news 
		WHERE title LIKE ? OR url LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery,
		searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.PGRNews
	for rows.Next() {
		var entity domain.PGRNews

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.URL, &entity.Title,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// SearchByCategory performs a search within a specific keyword category
func (r *PgrRepository) SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]domain.PGRNews, error) {
	var searchQuery string
	searchPattern := "%" + query + "%"

	switch category {
	case domain.KeywordCategoryCompanyName, domain.KeywordCategoryPersonName, domain.KeywordCategoryAddress:
		searchQuery = `
			SELECT id, domain_search_result_id, url, title, created_at, updated_at
			FROM pgr_news 
			WHERE title LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	default:
		return []domain.PGRNews{}, fmt.Errorf("unsupported category: %s", category)
	}

	rows, err := r.db.QueryContext(ctx, searchQuery,
		searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.PGRNews
	for rows.Next() {
		var entity domain.PGRNews

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.URL, &entity.Title,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// GetByDomainType retrieves PGR news items by domain type
func (r *PgrRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]domain.PGRNews, error) {
	// For PGR, all news items are of the same domain type, so we just return all
	return r.List(ctx, offset, limit)
}

// GetBySearchParameter retrieves PGR news items by search parameter
func (r *PgrRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]domain.PGRNews, error) {
	return r.Search(ctx, searchParam, offset, limit)
}

// GetKeywordsByCategory retrieves keywords grouped by category for a PGR news item
func (r *PgrRepository) GetKeywordsByCategory(ctx context.Context, entity domain.PGRNews) (map[domain.KeywordCategory][]string, error) {
	entity, err := r.GetByID(ctx, entity.ID.String())
	if err != nil {
		return nil, err
	}

	return map[domain.KeywordCategory][]string{}, nil
}
