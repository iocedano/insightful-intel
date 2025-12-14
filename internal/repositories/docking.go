package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// DockingRepository implements DomainRepository for Google Docking GoogleDorkingResult domain type
type DockingRepository struct {
	db DatabaseAccessor
}

// NewDockingRepository creates a new Google Docking repository instance
func NewDockingRepository(db database.Service) *DockingRepository {
	return &DockingRepository{
		db: NewDatabaseAdapter(db),
	}
}

// Create inserts a new Google Docking result
func (r *DockingRepository) Create(ctx context.Context, entity domain.GoogleDorkingResult) error {
	return r.CreateWithDomainSearchResultID(ctx, entity)
}

// CreateWithDomainSearchResultID inserts a new Google Docking result with a domain search result ID
func (r *DockingRepository) CreateWithDomainSearchResultID(ctx context.Context, entity domain.GoogleDorkingResult) error {
	entity.ID = domain.NewID()

	query := `
		INSERT INTO google_docking_results (
			id, domain_search_result_id, search_parameter, url, title, description, relevance, search_rank, keywords, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	keywordsJSON, _ := json.Marshal(entity.Keywords)

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.DomainSearchResultID, entity.SearchParameter, entity.URL, entity.Title, entity.Description, entity.Relevance, entity.Rank, keywordsJSON,
	)

	return err
}

// GetByID retrieves a Google Docking result by its URL
func (r *DockingRepository) GetByID(ctx context.Context, id string) (domain.GoogleDorkingResult, error) {
	query := `
		SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
		FROM google_docking_results 
		WHERE id = ?
	`

	var entity domain.GoogleDorkingResult
	var keywordsJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID, &entity.DomainSearchResultID, &entity.URL, &entity.Title, &entity.Description, &entity.Relevance, &entity.Rank, &keywordsJSON,
	)

	if err != nil {
		return domain.GoogleDorkingResult{}, err
	}

	// Parse JSON field
	json.Unmarshal([]byte(keywordsJSON), &entity.Keywords)

	return entity, nil
}

// Update modifies an existing Google Docking result
func (r *DockingRepository) Update(ctx context.Context, id string, entity domain.GoogleDorkingResult) error {
	query := `
		UPDATE google_docking_results SET
			domain_search_result_id = ?, title = ?, description = ?, relevance = ?, search_rank = ?, keywords = ?, updated_at = NOW()
		WHERE id = ?
	`

	keywordsJSON, _ := json.Marshal(entity.Keywords)

	_, err := r.db.ExecContext(ctx, query,
		entity.DomainSearchResultID, entity.Title, entity.Description, entity.Relevance, entity.Rank, keywordsJSON, id,
	)

	return err
}

// Delete removes a Google Docking result by its URL
func (r *DockingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM google_docking_results WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves multiple Google Docking results with pagination
func (r *DockingRepository) List(ctx context.Context, offset, limit int) ([]domain.GoogleDorkingResult, error) {
	query := `
		SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
		FROM google_docking_results 
		ORDER BY relevance DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.GoogleDorkingResult
	for rows.Next() {
		var entity domain.GoogleDorkingResult
		var keywordsJSON string

		err := rows.Scan(
			&entity.ID,                   // id
			&entity.DomainSearchResultID, // domain_search_result_id
			&entity.URL,                  // url
			&entity.Title,                // title
			&entity.Description,          // description
			&entity.Relevance,            // relevance
			&entity.Rank,                 // search_rank
			&keywordsJSON,                // keywords
			new(interface{}),             // created_at (ignored)
			new(interface{}),             // updated_at (ignored)
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON field
		json.Unmarshal([]byte(keywordsJSON), &entity.Keywords)

		entities = append(entities, entity)
	}

	return entities, nil
}

// Count returns the total number of Google Docking results
func (r *DockingRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM google_docking_results`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Search performs a search query on Google Docking results
func (r *DockingRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.GoogleDorkingResult, error) {
	searchQuery := `
		SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
		FROM google_docking_results 
		WHERE id LIKE ? OR domain_search_result_id LIKE ? OR title LIKE ? OR description LIKE ? OR url LIKE ?
		ORDER BY relevance DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery,
		searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.GoogleDorkingResult
	for rows.Next() {
		var entity domain.GoogleDorkingResult
		var keywordsJSON string

		err := rows.Scan(
			&entity.ID,                   // id
			&entity.DomainSearchResultID, // domain_search_result_id
			&entity.URL,                  // url
			&entity.Title,                // title
			&entity.Description,          // description
			&entity.Relevance,            // relevance
			&entity.Rank,                 // search_rank
			&keywordsJSON,                // keywords
			new(interface{}),             // created_at (ignored)
			new(interface{}),             // updated_at (ignored)
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON field
		json.Unmarshal([]byte(keywordsJSON), &entity.Keywords)

		entities = append(entities, entity)
	}

	return entities, nil
}

// SearchByCategory performs a search within a specific keyword category
func (r *DockingRepository) SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]domain.GoogleDorkingResult, error) {
	var searchQuery string
	searchPattern := "%" + query + "%"

	switch category {
	case domain.KeywordCategoryPersonName:
		searchQuery = `
			SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
			FROM google_docking_results 
			WHERE title LIKE ? OR description LIKE ?
			ORDER BY relevance DESC, created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryCompanyName:
		searchQuery = `
			SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
			FROM google_docking_results 
			WHERE title LIKE ? OR description LIKE ?
			ORDER BY relevance DESC, created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryAddress:
		searchQuery = `
			SELECT id, domain_search_result_id, description url, title, description, relevance, search_rank, keywords, created_at, updated_at
			FROM google_docking_results 
			WHERE description LIKE ?
			ORDER BY relevance DESC, created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategorySocialMedia:
		searchQuery = `
			SELECT id, domain_search_result_id, url, title, description, relevance, search_rank, keywords, created_at, updated_at
			FROM google_docking_results 
			WHERE url LIKE ? OR title LIKE ? OR description LIKE ?
			ORDER BY relevance DESC, created_at DESC
			LIMIT ? OFFSET ?
		`
	default:
		return []domain.GoogleDorkingResult{}, fmt.Errorf("unsupported category: %s", category)
	}

	var rows *sql.Rows
	var err error

	if category == domain.KeywordCategorySocialMedia {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, searchPattern, searchPattern, limit, offset)
	} else if category == domain.KeywordCategoryAddress {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, searchPattern, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.GoogleDorkingResult
	for rows.Next() {
		var entity domain.GoogleDorkingResult
		var keywordsJSON string

		err := rows.Scan(
			&entity.ID,                   // id
			&entity.DomainSearchResultID, // domain_search_result_id
			&entity.URL,                  // url
			&entity.Title,                // title
			&entity.Description,          // description
			&entity.Relevance,            // relevance
			&entity.Rank,                 // search_rank
			&keywordsJSON,                // keywords
			new(interface{}),             // created_at (ignored)
			new(interface{}),             // updated_at (ignored)
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON field
		json.Unmarshal([]byte(keywordsJSON), &entity.Keywords)

		entities = append(entities, entity)
	}

	return entities, nil
}

// GetByDomainType retrieves Google Docking results by domain type
func (r *DockingRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]domain.GoogleDorkingResult, error) {
	// For Google Docking, all results are of the same domain type, so we just return all
	return r.List(ctx, offset, limit)
}

// GetBySearchParameter retrieves Google Docking results by search parameter
func (r *DockingRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]domain.GoogleDorkingResult, error) {
	return r.Search(ctx, searchParam, offset, limit)
}

// GetKeywordsByCategory retrieves keywords grouped by category for a Google Docking result
func (r *DockingRepository) GetKeywordsByCategory(ctx context.Context, entityID string) (map[domain.KeywordCategory][]string, error) {
	// entity, err := r.GetByID(ctx, entityID)
	// if err != nil {
	// 	return nil, err
	// }

	// docking := domain.NewGoogleDorkingDomain()
	return map[domain.KeywordCategory][]string{
		// domain.KeywordCategoryCompanyName: docking.GetDataByCategory(entity, domain.KeywordCategoryCompanyName),
		// domain.KeywordCategoryPersonName:  docking.GetDataByCategory(entity, domain.KeywordCategoryPersonName),
		// domain.KeywordCategoryAddress:     docking.GetDataByCategory(entity, domain.KeywordCategoryAddress),
		// domain.KeywordCategorySocialMedia: docking.GetDataByCategory(entity, domain.KeywordCategorySocialMedia),
	}, nil
}
