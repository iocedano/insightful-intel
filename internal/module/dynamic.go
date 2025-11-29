package module

import (
	"context"
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/infra"
)

// SearchDomain performs a search using the specified domain type and parameters
func SearchDomain(domainType domain.DomainType, params domain.DomainSearchParams) (*domain.DomainSearchResult, error) {
	// Validate domain type
	if !domain.IsValidDomainType(domainType) {
		return &domain.DomainSearchResult{
			Success:    false,
			Error:      fmt.Errorf("unsupported domain type: %s", domainType),
			DomainType: domainType,
		}, fmt.Errorf("unsupported domain type: %s", domainType)
	}

	// Perform the search based on domain type
	var output any
	var searchErr error

	switch domainType {
	case domain.DomainTypeONAPI:
		onapi := NewOnapiDomain()
		output, searchErr = onapi.Search(params.Query)
	case domain.DomainTypeSCJ:
		scj := NewScjDomain()
		output, searchErr = scj.Search(params.Query)
	case domain.DomainTypeDGII:
		dgii := NewDgiiDomain()
		output, searchErr = dgii.Search(params.Query)
	case domain.DomainTypePGR:
		pgr := NewPgrDomain()
		output, searchErr = pgr.Search(params.Query)
	case domain.DomainTypeGoogleDorking:
		output, searchErr = NewGoogleDorkingBuilder().
			Query(params.Query).
			IncludeKeywords(domain.FRAUD_KEYWORDS...).
			Build()
	case domain.DomainTypeSocialMedia:
		output, searchErr = NewGoogleDorkingBuilder().
			Query(params.Query).
			SitesKeywords(domain.SOCIAL_MEDIA_SITES_KEYWORDS...).
			Build()
	case domain.DomainTypeFileType:
		output, searchErr = NewGoogleDorkingBuilder().
			Query(params.Query).
			FileTypeKeywords(domain.FILE_TYPE_KEYWORDS...).
			Build()
	case domain.DomainTypeXSocialMedia:
		output, searchErr = NewGoogleDorkingBuilder().
			Query(params.Query).
			InURLKeywords(domain.X_IN_URL_KEYWORDS...).
			SitesKeywords([]string{"x.com"}...).
			Build()
	default:
		return &domain.DomainSearchResult{
			Success:    false,
			Error:      fmt.Errorf("unsupported domain type: %s", domainType),
			DomainType: domainType,
		}, fmt.Errorf("unsupported domain type: %s", domainType)
	}

	// Extract keywords from the result
	var keywordsPerCategory map[domain.KeywordCategory][]string
	if searchErr == nil && output != nil {
		switch domainType {
		case domain.DomainTypeONAPI:
			if entities, ok := output.([]domain.Entity); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(NewOnapiDomain(), entities)
			}
		case domain.DomainTypeSCJ:
			if cases, ok := output.([]domain.ScjCase); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(NewScjDomain(), cases)
			}
		case domain.DomainTypeDGII:
			if registers, ok := output.([]domain.Register); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(NewDgiiDomain(), registers)
			}
		case domain.DomainTypePGR:
			if registers, ok := output.([]domain.PGRNews); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(NewPgrDomain(), registers)
			}
		case domain.DomainTypeGoogleDorking, domain.DomainTypeSocialMedia, domain.DomainTypeFileType, domain.DomainTypeXSocialMedia:
			if registers, ok := output.([]domain.GoogleDorkingResult); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(&GoogleDorking{}, registers)
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

// CreateDomainConnector creates a domain connector instance based on the domain type
func CreateDomainConnector(domainType domain.DomainType) (any, error) {
	// Validate domain type
	if !domain.IsValidDomainType(domainType) {
		return nil, fmt.Errorf("unsupported domain type: %s", domainType)
	}

	switch domainType {
	case domain.DomainTypeONAPI:
		onapi := NewOnapiDomain()
		return &onapi, nil
	case domain.DomainTypeSCJ:
		scj := NewScjDomain()
		return &scj, nil
	case domain.DomainTypeDGII:
		dgii := NewDgiiDomain()
		return &dgii, nil
	case domain.DomainTypePGR:
		pgr := NewPgrDomain()
		return &pgr, nil
	case domain.DomainTypeGoogleDorking, domain.DomainTypeSocialMedia, domain.DomainTypeFileType, domain.DomainTypeXSocialMedia:
		docking := NewGoogleDorkingDomain()
		return &docking, nil
	default:
		return nil, fmt.Errorf("unsupported domain type: %s", domainType)
	}
}

// SearchMultipleDomains performs searches across multiple domains
func SearchMultipleDomains(domainTypes []domain.DomainType, params domain.DomainSearchParams) []*domain.DomainSearchResult {
	results := make([]*domain.DomainSearchResult, 0, len(domainTypes))

	for _, domainType := range domainTypes {
		result, err := SearchDomain(domainType, params)
		if err != nil {
			result.Error = err
			result.Success = false
		}
		results = append(results, result)
	}

	return results
}

// DefaultDynamicPipelineConfig returns a default configuration
func DefaultDynamicPipelineConfig() domain.DynamicPipelineConfig {
	return domain.DynamicPipelineConfig{
		MaxDepth:           5,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     false,
	}
}

// CreateDynamicPipeline creates a dynamic pipeline based on searchable categories
func CreateDynamicPipeline(
	ctx context.Context,
	initialQuery string,
	availableDomains []domain.DomainType,
	config domain.DynamicPipelineConfig,
) (*domain.DynamicPipelineResult, error) {
	pipelineID := domain.NewID()

	executionID, _ := infra.GetExecutionID(ctx)

	if executionID != "" {
		pipelineID = domain.NewIDFromString(executionID)
	}

	// Initialize the pipeline
	pipeline := &domain.DynamicPipelineResult{
		ID:     pipelineID,
		Steps:  make([]domain.DynamicPipelineStep, 0),
		Config: config,
	}

	// Track searched keywords per domain to avoid duplicates
	searchedKeywordsPerDomain := make(map[domain.DomainType]map[string]bool)
	for _, domainType := range availableDomains {
		searchedKeywordsPerDomain[domainType] = make(map[string]bool)
	}

	// Add initial steps for all available domains
	// Map each domain to its default category for initial search
	initialDomainCategories := map[domain.DomainType]domain.KeywordCategory{
		domain.DomainTypeONAPI:         domain.KeywordCategoryCompanyName,
		domain.DomainTypeDGII:          domain.KeywordCategoryContributorID,
		domain.DomainTypePGR:           domain.KeywordCategoryPersonName,
		domain.DomainTypeSCJ:           domain.KeywordCategoryContributorID,
		domain.DomainTypeGoogleDorking: domain.KeywordCategoryCompanyName,
	}

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

// generateStepsFromKeywords creates new pipeline steps based on extracted keywords
func generateStepsFromKeywords(
	keywordsPerCategory map[domain.KeywordCategory][]string,
	availableDomains []domain.DomainType,
	searchedKeywordsPerDomain map[domain.DomainType]map[string]bool,
	depth int,
	config domain.DynamicPipelineConfig,
) []domain.DynamicPipelineStep {

	var newSteps []domain.DynamicPipelineStep

	// Get all available domain connectors
	domainConnectors := make(map[domain.DomainType]any)
	for _, domainType := range availableDomains {
		connector, err := CreateDomainConnector(domainType)
		if err != nil {
			continue
		}
		domainConnectors[domainType] = connector
	}

	// For each category and its keywords
	for category, keywords := range keywordsPerCategory {
		// Find domains that can search this category
		for domainType, connector := range domainConnectors {
			searchableCategories := GetSearchableKeywordCategories(connector)
			if !contains(searchableCategories, category) {
				continue
			}

			// Create steps for each keyword
			for _, keyword := range keywords {
				if keyword == "" {
					continue
				}

				// Skip duplicates if configured
				if config.SkipDuplicates {
					if searchedKeywordsPerDomain[domainType][keyword] {
						continue
					}
					searchedKeywordsPerDomain[domainType][keyword] = true
				}

				newStep := domain.DynamicPipelineStep{
					DomainType:      domainType,
					SearchParameter: keyword,
					Category:        category,
					Keywords:        []string{keyword},
					Depth:           depth,
				}

				newSteps = append(newSteps, newStep)
			}
		}
	}

	return newSteps
}

// GetSearchableKeywordCategories extracts searchable categories from a connector
func GetSearchableKeywordCategories(connector any) []domain.KeywordCategory {
	switch c := connector.(type) {
	case *Onapi:
		return c.GetSearchableKeywordCategories()
	case *Scj:
		return c.GetSearchableKeywordCategories()
	case *Dgii:
		return c.GetSearchableKeywordCategories()
	case *Pgr:
		return c.GetSearchableKeywordCategories()
	case *GoogleDorking:
		return c.GetSearchableKeywordCategories()
	default:
		return []domain.KeywordCategory{}
	}
}

// contains checks if a slice contains a specific element
func contains(slice []domain.KeywordCategory, item domain.KeywordCategory) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ExecuteDynamicPipeline executes the dynamic pipeline with parallel processing
func ExecuteDynamicPipeline(
	ctx context.Context,
	initialQuery string,
	availableDomains []domain.DomainType,
	config domain.DynamicPipelineConfig,
) (*domain.DynamicPipelineResult, error) {

	// Create the pipeline structure
	pipeline, err := CreateDynamicPipeline(ctx, initialQuery, availableDomains, config)
	if err != nil {
		return nil, err
	}

	// Execute steps in parallel where possible
	// This is a simplified version - in practice you might want more sophisticated parallel execution
	for i := range pipeline.Steps {
		step := &pipeline.Steps[i]
		if step.Success || step.Error != nil {
			continue // Already processed
		}

		result, err := SearchDomain(step.DomainType, domain.DomainSearchParams{Query: step.SearchParameter})
		if err != nil {
			step.Error = err
			step.Success = false
		} else {
			step.Success = result.Success
			step.Output = result.Output
			step.KeywordsPerCategory = result.KeywordsPerCategory
		}
	}

	return pipeline, nil
}
