package module

import (
	"fmt"
	"insightful-intel/internal/domain"
)

// DomainSearchParams holds the search parameters for different domains

// SearchDomain performs a search using the specified domain type and parameters
func SearchDomain(domainType domain.DomainType, params domain.DomainSearchParams) (*domain.DomainSearchResult, error) {
	connector, err := CreateDomainConnector(domainType)
	if err != nil {
		return &domain.DomainSearchResult{
			Success:    false,
			Error:      err,
			DomainType: domainType,
		}, err
	}

	// Perform the search based on domain type
	var output any
	var searchErr error

	switch domainType {
	case domain.DomainTypeONAPI:
		onapi := connector.(*Onapi)
		entities, err := onapi.SearchComercialName(params.Query)
		output = entities
		searchErr = err
	case domain.DomainTypeSCJ:
		scj := connector.(*Scj)
		cases, err := scj.Search(params.Query)
		output = cases
		searchErr = err
	case domain.DomainTypeDGII:
		dgii := connector.(*Dgii)
		registers, err := dgii.GetRegister(params.Query)
		output = registers
		searchErr = err
	case domain.DomainTypePGR:
		pgr := connector.(*Pgr)
		registers, err := pgr.Search(params.Query)
		output = registers
		searchErr = err
	case domain.DomainTypeGoogleDocking:
		registers, err := NewGoogleDockingBuilder().
			Query(params.Query).
			IncludeKeywords(FRAUD_KEYWORDS...).
			Build()

		output = registers
		searchErr = err
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
				keywordsPerCategory = domain.GetCategoryByKeywords(&Onapi{}, entities)
			}
		case domain.DomainTypeSCJ:
			if cases, ok := output.([]domain.ScjCase); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(&Scj{}, cases)
			}
		case domain.DomainTypeDGII:
			if registers, ok := output.([]domain.Register); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(&Dgii{}, registers)
			}
		case domain.DomainTypePGR:
			if registers, ok := output.([]domain.PGRNews); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(&domain.Pgr{}, registers)
			}
		case domain.DomainTypeGoogleDocking:
			if registers, ok := output.([]domain.GoogleDockingResult); ok {
				keywordsPerCategory = domain.GetCategoryByKeywords(&domain.GoogleDocking{}, registers)
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
	case domain.DomainTypeGoogleDocking:
		docking := NewGoogleDockingDomain()
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

// DynamicPipelineConfig holds configuration for the dynamic pipeline
type DynamicPipelineConfig struct {
	MaxDepth           int
	MaxConcurrentSteps int
	DelayBetweenSteps  int // seconds
	SkipDuplicates     bool
}

// DefaultDynamicPipelineConfig returns a default configuration
func DefaultDynamicPipelineConfig() DynamicPipelineConfig {
	return DynamicPipelineConfig{
		MaxDepth:           5,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     true,
	}
}

// DynamicPipelineStep represents a single step in the pipeline
type DynamicPipelineStep struct {
	DomainType          domain.DomainType
	SearchParameter     string
	Category            domain.KeywordCategory
	Keywords            []string
	Success             bool
	Error               error
	Output              any
	KeywordsPerCategory map[domain.KeywordCategory][]string
	Depth               int
}

// DynamicPipelineResult represents the complete pipeline result
type DynamicPipelineResult struct {
	Steps           []DynamicPipelineStep
	TotalSteps      int
	SuccessfulSteps int
	FailedSteps     int
	MaxDepthReached int
	Config          DynamicPipelineConfig
}

// CreateDynamicPipeline creates a dynamic pipeline based on searchable categories
func CreateDynamicPipeline(
	initialQuery string,
	availableDomains []domain.DomainType,
	config DynamicPipelineConfig,
) (*DynamicPipelineResult, error) {

	// Initialize the pipeline
	pipeline := &DynamicPipelineResult{
		Steps:  make([]DynamicPipelineStep, 0),
		Config: config,
	}

	// Track searched keywords per domain to avoid duplicates
	searchedKeywordsPerDomain := make(map[domain.DomainType]map[string]bool)
	for _, domainType := range availableDomains {
		searchedKeywordsPerDomain[domainType] = make(map[string]bool)
	}

	// Add initial step
	pipeline.Steps = append(pipeline.Steps,
		DynamicPipelineStep{
			DomainType:      domain.DomainTypeONAPI, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        domain.KeywordCategoryCompanyName,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      domain.DomainTypeDGII, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        domain.KeywordCategoryContributorID,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      domain.DomainTypePGR, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        domain.KeywordCategoryPersonName,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      domain.DomainTypeSCJ, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        domain.KeywordCategoryContributorID,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      domain.DomainTypeGoogleDocking, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        domain.KeywordCategoryCompanyName,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
	)

	// Process the pipeline dynamically
	currentStep := 0
	for currentStep < len(pipeline.Steps) && currentStep < config.MaxDepth {
		step := pipeline.Steps[currentStep]

		// Skip if we've already processed this step
		if step.Success || step.Error != nil {
			currentStep++
			continue
		}

		// Execute the step
		result, err := SearchDomain(step.DomainType, domain.DomainSearchParams{Query: step.SearchParameter})
		if err != nil {
			step.Error = err
			step.Success = false
			pipeline.Steps[currentStep] = step
			currentStep++
			continue
		}

		// Update step with results
		step.Success = result.Success
		step.Output = result.Output
		step.KeywordsPerCategory = result.KeywordsPerCategory
		pipeline.Steps[currentStep] = step

		// Generate new steps based on keywords
		newSteps := generateStepsFromKeywords(
			step.KeywordsPerCategory,
			availableDomains,
			searchedKeywordsPerDomain,
			step.Depth+1,
			config,
		)

		// Add new steps to pipeline
		pipeline.Steps = append(pipeline.Steps, newSteps...)

		currentStep++
	}

	// Calculate statistics
	pipeline.TotalSteps = len(pipeline.Steps)
	pipeline.SuccessfulSteps = 0
	pipeline.FailedSteps = 0
	maxDepth := 0

	for _, step := range pipeline.Steps {
		if step.Success {
			pipeline.SuccessfulSteps++
		} else {
			pipeline.FailedSteps++
		}
		if step.Depth > maxDepth {
			maxDepth = step.Depth
		}
	}

	pipeline.MaxDepthReached = maxDepth

	return pipeline, nil
}

// generateStepsFromKeywords creates new pipeline steps based on extracted keywords
func generateStepsFromKeywords(
	keywordsPerCategory map[domain.KeywordCategory][]string,
	availableDomains []domain.DomainType,
	searchedKeywordsPerDomain map[domain.DomainType]map[string]bool,
	depth int,
	config DynamicPipelineConfig,
) []DynamicPipelineStep {

	var newSteps []DynamicPipelineStep

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
			searchableCategories := getSearchableCategories(connector)
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

				newStep := DynamicPipelineStep{
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

// getSearchableCategories extracts searchable categories from a connector
func getSearchableCategories(connector any) []domain.KeywordCategory {
	switch c := connector.(type) {
	case *Onapi:
		return c.GetSearchableKeywordCategories()
	case *Scj:
		return c.GetSearchableKeywordCategories()
	case *Dgii:
		return c.GetSearchableKeywordCategories()
	case *Pgr:
		return c.GetSearchableKeywordCategories()
	case *GoogleDocking:
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
	initialQuery string,
	availableDomains []domain.DomainType,
	config DynamicPipelineConfig,
) (*DynamicPipelineResult, error) {

	// Create the pipeline structure
	pipeline, err := CreateDynamicPipeline(initialQuery, availableDomains, config)
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
