package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/module"

	"github.com/google/uuid"
)

// PipelineRepository implements PipelineResultRepository for pipeline results
type PipelineRepository struct {
	db DatabaseAccessor
}

// NewPipelineRepository creates a new pipeline repository instance
func NewPipelineRepository(db database.Service) *PipelineRepository {
	return &PipelineRepository{
		db: NewDatabaseAdapter(db),
	}
}

// createDomainSearchResult inserts a DomainSearchResult
func (r *PipelineRepository) CreateDomainSearchResult(ctx context.Context, result *domain.DomainSearchResult) (*domain.DomainSearchResult, error) {
	id := domain.NewID()
	result.ID = id

	query := `
		INSERT INTO domain_search_results (
			id, success, error_message, domain_type, search_parameter, keywords_per_category, output, pipeline_steps_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	keywordsJSON, _ := json.Marshal(result.KeywordsPerCategory)
	outputJSON, _ := json.Marshal(result.Output)

	var errorMessage string
	if result.Error != nil {
		errorMessage = result.Error.Error()
	}

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.Success, errorMessage, string(result.DomainType), result.SearchParameter, keywordsJSON, outputJSON, result.PipelineStepsID,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// createDynamicPipelineResult inserts a DynamicPipelineResult
func (r *PipelineRepository) CreateDynamicPipelineResult(ctx context.Context, result *module.DynamicPipelineResult) (*module.DynamicPipelineResult, error) {
	if result.ID == domain.ID(uuid.Nil) {
		result.ID = domain.NewID()
	}

	query := `
		INSERT INTO dynamic_pipeline_results (
			id, total_steps, successful_steps, failed_steps, max_depth_reached, config, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	configJSON, _ := json.Marshal(result.Config)

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.TotalSteps, result.SuccessfulSteps, result.FailedSteps, result.MaxDepthReached, configJSON,
	)
	if err != nil {
		return nil, err
	}

	// // Get the generated UUID
	// err = r.CreateDynamicPipelineSteps(ctx, result.ID.String(), result.Steps)
	// if err != nil {
	// 	return nil, err
	// }

	return result, nil
}

// CreateDynamicPipelineSteps inserts individual pipeline steps
func (r *PipelineRepository) CreateDynamicPipelineStep(ctx context.Context, step *module.DynamicPipelineStep) error {
	if step.ID == domain.ID(uuid.Nil) {
		step.ID = domain.NewID()
	}

	query := `
		INSERT INTO dynamic_pipeline_steps (
			id, pipeline_id, domain_type, search_parameter, category, keywords, success, error_message, 
			output, keywords_per_category, depth, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	keywordsJSON, _ := json.Marshal(step.Keywords)
	outputJSON, _ := json.Marshal(step.Output)
	keywordsPerCategoryJSON, _ := json.Marshal(step.KeywordsPerCategory)

	var errorMessage string
	if step.Error != nil {
		errorMessage = step.Error.Error()
	}

	_, err := r.db.ExecContext(ctx, query,
		step.ID, step.PipelineID, string(step.DomainType), step.SearchParameter, string(step.Category),
		keywordsJSON, step.Success, errorMessage, outputJSON, keywordsPerCategoryJSON, step.Depth,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a pipeline result by its ID
func (r *PipelineRepository) GetByID(ctx context.Context, id string) (any, error) {
	// Try DomainSearchResult first
	domainResult, err := r.getDomainSearchResultByID(ctx, id)
	if err == nil {
		return domainResult, nil
	}

	// Try DynamicPipelineResult
	dynamicResult, err := r.getDynamicPipelineResultByID(ctx, id)
	if err == nil {
		return dynamicResult, nil
	}

	return nil, fmt.Errorf("pipeline result not found with ID: %s", id)
}

// getDomainSearchResultByID retrieves a DomainSearchResult by ID
func (r *PipelineRepository) getDomainSearchResultByID(ctx context.Context, id string) (*domain.DomainSearchResult, error) {
	query := `
		SELECT success, error_message, domain_type, search_parameter, keywords_per_category, output, created_at, updated_at
		FROM domain_search_results 
		WHERE id = ?
	`

	var result domain.DomainSearchResult
	var errorMessage, domainType, keywordsJSON, outputJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.Success, &errorMessage, &domainType, &result.SearchParameter, &keywordsJSON, &outputJSON,
	)

	if err != nil {
		return nil, err
	}

	if errorMessage != "" {
		result.Error = fmt.Errorf(errorMessage)
	}
	result.DomainType = domain.DomainType(domainType)
	json.Unmarshal([]byte(keywordsJSON), &result.KeywordsPerCategory)
	json.Unmarshal([]byte(outputJSON), &result.Output)

	return &result, nil
}

// getDynamicPipelineResultByID retrieves a DynamicPipelineResult by ID
func (r *PipelineRepository) getDynamicPipelineResultByID(ctx context.Context, id string) (*module.DynamicPipelineResult, error) {
	query := `
		SELECT total_steps, successful_steps, failed_steps, max_depth_reached, config, created_at, updated_at
		FROM dynamic_pipeline_results 
		WHERE id = ?
	`

	var result module.DynamicPipelineResult
	var configJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.TotalSteps, &result.SuccessfulSteps, &result.FailedSteps, &result.MaxDepthReached, &configJSON,
	)

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(configJSON), &result.Config)

	// Get steps
	steps, err := r.getPipelineSteps(ctx, id)
	if err != nil {
		return nil, err
	}
	result.Steps = steps

	return &result, nil
}

// getPipelineSteps retrieves steps for a pipeline result
func (r *PipelineRepository) getPipelineSteps(ctx context.Context, pipelineID string) ([]module.DynamicPipelineStep, error) {
	query := `
		SELECT domain_type, search_parameter, category, keywords, success, error_message, 
			   output, keywords_per_category, depth, created_at, updated_at
		FROM dynamic_pipeline_steps 
		WHERE pipeline_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, pipelineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []module.DynamicPipelineStep
	for rows.Next() {
		var step module.DynamicPipelineStep
		var domainType, category, keywordsJSON, outputJSON, keywordsPerCategoryJSON, errorMessage string

		err := rows.Scan(
			&domainType, &step.SearchParameter, &category, &keywordsJSON, &step.Success, &errorMessage,
			&outputJSON, &keywordsPerCategoryJSON, &step.Depth,
		)
		if err != nil {
			return nil, err
		}

		step.DomainType = domain.DomainType(domainType)
		step.Category = domain.KeywordCategory(category)
		if errorMessage != "" {
			step.Error = fmt.Errorf(errorMessage)
		}
		json.Unmarshal([]byte(keywordsJSON), &step.Keywords)
		json.Unmarshal([]byte(outputJSON), &step.Output)
		json.Unmarshal([]byte(keywordsPerCategoryJSON), &step.KeywordsPerCategory)

		steps = append(steps, step)
	}

	return steps, nil
}

// Update modifies an existing pipeline result
func (r *PipelineRepository) Update(ctx context.Context, id string, entity any) error {
	switch v := entity.(type) {
	case *domain.DomainSearchResult:
		return r.updateDomainSearchResult(ctx, v)
	case *module.DynamicPipelineResult:
		return r.UpdateDynamicPipelineResult(ctx, v)
	default:
		return fmt.Errorf("unsupported pipeline result type: %T", entity)
	}
}

// updateDomainSearchResult updates a DomainSearchResult
func (r *PipelineRepository) updateDomainSearchResult(ctx context.Context, result *domain.DomainSearchResult) error {
	query := `
		UPDATE domain_search_results SET
			success = ?, error_message = ?, domain_type = ?, search_parameter = ?, 
			keywords_per_category = ?, output = ?, updated_at = NOW()
		WHERE id = ?
	`

	keywordsJSON, _ := json.Marshal(result.KeywordsPerCategory)
	outputJSON, _ := json.Marshal(result.Output)

	var errorMessage string
	if result.Error != nil {
		errorMessage = result.Error.Error()
	}

	_, err := r.db.ExecContext(ctx, query,
		result.Success, errorMessage, string(result.DomainType), result.SearchParameter,
		keywordsJSON, outputJSON, result.ID,
	)

	return err
}

// updateDynamicPipelineResult updates a DynamicPipelineResult
func (r *PipelineRepository) UpdateDynamicPipelineResult(ctx context.Context, result *module.DynamicPipelineResult) error {
	query := `
		UPDATE dynamic_pipeline_results SET
			total_steps = ?, successful_steps = ?, failed_steps = ?, max_depth_reached = ?, 
			config = ?, updated_at = NOW()
		WHERE id = ?
	`

	configJSON, _ := json.Marshal(result.Config)

	_, err := r.db.ExecContext(ctx, query,
		result.TotalSteps, result.SuccessfulSteps, result.FailedSteps, result.MaxDepthReached,
		configJSON, result.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a pipeline result by its ID
func (r *PipelineRepository) Delete(ctx context.Context, id string) error {
	// Delete from both tables (cascade should handle steps)
	queries := []string{
		`DELETE FROM dynamic_pipeline_steps WHERE pipeline_id = ?`,
		`DELETE FROM dynamic_pipeline_results WHERE id = ?`,
		`DELETE FROM domain_search_results WHERE id = ?`,
	}

	for _, query := range queries {
		_, err := r.db.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// List retrieves multiple pipeline results with pagination
func (r *PipelineRepository) List(ctx context.Context, offset, limit int) ([]any, error) {
	// Get DomainSearchResults
	domainQuery := `
		SELECT id, success, error_message, domain_type, search_parameter, keywords_per_category, output, created_at, updated_at
		FROM domain_search_results 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, domainQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var result domain.DomainSearchResult
		var id, errorMessage, domainType, keywordsJSON, outputJSON string

		err := rows.Scan(
			&id, &result.Success, &errorMessage, &domainType, &result.SearchParameter, &keywordsJSON, &outputJSON,
		)
		if err != nil {
			return nil, err
		}

		if errorMessage != "" {
			result.Error = fmt.Errorf(errorMessage)
		}
		result.DomainType = domain.DomainType(domainType)
		json.Unmarshal([]byte(keywordsJSON), &result.KeywordsPerCategory)
		json.Unmarshal([]byte(outputJSON), &result.Output)

		results = append(results, &result)
	}

	return results, nil
}

// Count returns the total number of pipeline results
func (r *PipelineRepository) Count(ctx context.Context) (int64, error) {
	query := `
		SELECT (
			(SELECT COUNT(*) FROM domain_search_results) + 
			(SELECT COUNT(*) FROM dynamic_pipeline_results)
		) as total
	`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// GetByDomainType retrieves pipeline results by domain type
func (r *PipelineRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]any, error) {
	query := `
		SELECT id, success, error_message, domain_type, search_parameter, keywords_per_category, output, created_at, updated_at
		FROM domain_search_results 
		WHERE domain_type = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, string(domainType), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var result domain.DomainSearchResult
		var id, errorMessage, domainTypeStr, keywordsJSON, outputJSON string

		err := rows.Scan(
			&id, &result.Success, &errorMessage, &domainTypeStr, &result.SearchParameter, &keywordsJSON, &outputJSON,
		)
		if err != nil {
			return nil, err
		}

		if errorMessage != "" {
			result.Error = fmt.Errorf(errorMessage)
		}
		result.DomainType = domain.DomainType(domainTypeStr)
		json.Unmarshal([]byte(keywordsJSON), &result.KeywordsPerCategory)
		json.Unmarshal([]byte(outputJSON), &result.Output)

		results = append(results, &result)
	}

	return results, nil
}

// GetBySuccessStatus retrieves pipeline results by success status
func (r *PipelineRepository) GetBySuccessStatus(ctx context.Context, success bool, offset, limit int) ([]any, error) {
	query := `
		SELECT id, success, error_message, domain_type, search_parameter, keywords_per_category, output, created_at, updated_at
		FROM domain_search_results 
		WHERE success = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, success, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var result domain.DomainSearchResult
		var id, errorMessage, domainType, keywordsJSON, outputJSON string

		err := rows.Scan(
			&id, &result.Success, &errorMessage, &domainType, &result.SearchParameter, &keywordsJSON, &outputJSON,
		)
		if err != nil {
			return nil, err
		}

		if errorMessage != "" {
			result.Error = fmt.Errorf(errorMessage)
		}
		result.DomainType = domain.DomainType(domainType)
		json.Unmarshal([]byte(keywordsJSON), &result.KeywordsPerCategory)
		json.Unmarshal([]byte(outputJSON), &result.Output)

		results = append(results, &result)
	}

	return results, nil
}

// GetBySearchParameter retrieves pipeline results by search parameter
func (r *PipelineRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]any, error) {
	query := `
		SELECT id, success, error_message, domain_type, search_parameter, keywords_per_category, output, created_at, updated_at
		FROM domain_search_results 
		WHERE search_parameter LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + searchParam + "%"
	rows, err := r.db.QueryContext(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var result domain.DomainSearchResult
		var id, errorMessage, domainType, keywordsJSON, outputJSON string

		err := rows.Scan(
			&id, &result.Success, &errorMessage, &domainType, &result.SearchParameter, &keywordsJSON, &outputJSON,
		)
		if err != nil {
			return nil, err
		}

		if errorMessage != "" {
			result.Error = fmt.Errorf(errorMessage)
		}
		result.DomainType = domain.DomainType(domainType)
		json.Unmarshal([]byte(keywordsJSON), &result.KeywordsPerCategory)
		json.Unmarshal([]byte(outputJSON), &result.Output)

		results = append(results, &result)
	}

	return results, nil
}

// GetKeywordsByCategory retrieves keywords grouped by category for a pipeline result
func (r *PipelineRepository) GetKeywordsByCategory(ctx context.Context, resultID string) (map[domain.KeywordCategory][]string, error) {
	result, err := r.GetByID(ctx, resultID)
	if err != nil {
		return nil, err
	}

	switch v := result.(type) {
	case *domain.DomainSearchResult:
		return v.KeywordsPerCategory, nil
	case *module.DynamicPipelineResult:
		// Aggregate keywords from all steps
		aggregated := make(map[domain.KeywordCategory][]string)
		for _, step := range v.Steps {
			for category, keywords := range step.KeywordsPerCategory {
				if aggregated[category] == nil {
					aggregated[category] = []string{}
				}
				aggregated[category] = append(aggregated[category], keywords...)
			}
		}
		return aggregated, nil
	default:
		return nil, fmt.Errorf("unsupported result type: %T", result)
	}
}
