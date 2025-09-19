package domain

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type DataCategory string

const (
	DataCategoryCompanyName   DataCategory = "company_name"
	DataCategoryPersonName    DataCategory = "person_name"
	DataCategoryContributorID DataCategory = "contributor_id"
	DataCategoryAddress       DataCategory = "address"
)

// Generic interface for connectors that can process data of the same type
type GenericConnector[T any] interface {
	ProcessData(data T) (T, error)
	ValidateData(data T) error
	TransformData(data T) T
	GetDataByCategory(data T, category DataCategory) []string
	GetListOfSearchableCategory() []DataCategory
	GetListOfRetrievedCategory() []DataCategory
	GetDomainType() DomainType
}

// Extended interface that includes generic methods
type ExtendedConnectorInterface[T any] interface {
	GenericConnector[T]
}

// Generic utility function that works with any type implementing GenericConnector
func ProcessGenericData[T any](connector GenericConnector[T], data T) (T, error) {
	return connector.ProcessData(data)
}

// Generic utility function for batch processing
func ProcessBatchData[T any](connector GenericConnector[T], data []T) ([]T, error) {
	results := make([]T, 0, len(data))
	for _, item := range data {
		processed, err := connector.ProcessData(item)
		if err != nil {
			return nil, err
		}
		results = append(results, processed)
	}
	return results, nil
}

func GetCategoryByKeywords[T any](connector GenericConnector[T], data []T) map[DataCategory][]string {
	result := map[DataCategory][]string{}

	for _, d := range data {
		for _, rcv := range connector.GetListOfRetrievedCategory() {
			if result[rcv] == nil {
				result[rcv] = []string{}
			}

			result[rcv] = append(result[rcv], connector.GetDataByCategory(d, rcv)...)
		}
	}

	return result
}

// DomainType represents the different domain types available
type DomainType string

const (
	DomainTypeONAPI DomainType = "ONAPI"
	DomainTypeSCJ   DomainType = "SCJ"
	DomainTypeDGII  DomainType = "DGII"
	DomainTypePGR   DomainType = "PGR"
)

// DomainSearchParams holds the search parameters for different domains
type DomainSearchParams struct {
	Query string
	// Add more specific parameters as needed for different domains
	// For example:
	// PageSize int
	// PageIdx  int
	// Tipo     string
	// Subtipo  string
}

// DomainSearchResult represents the result of a domain search
type DomainSearchResult struct {
	Success             bool
	Error               error
	DomainType          DomainType
	SearchParameter     string
	KeywordsPerCategory map[DataCategory][]string
	Output              any
}

// SearchDomain performs a search using the specified domain type and parameters
func SearchDomain(domainType DomainType, params DomainSearchParams) (*DomainSearchResult, error) {
	connector, err := CreateDomainConnector(domainType)
	if err != nil {
		return &DomainSearchResult{
			Success:    false,
			Error:      err,
			DomainType: domainType,
		}, err
	}

	// Perform the search based on domain type
	var output any
	var searchErr error

	switch domainType {
	case DomainTypeONAPI:
		onapi := connector.(*Onapi)
		entities, err := onapi.SearchComercialName(params.Query)
		output = entities
		searchErr = err
	case DomainTypeSCJ:
		scj := connector.(*Scj)
		cases, err := scj.Search(params.Query)
		output = cases
		searchErr = err
	case DomainTypeDGII:
		dgii := connector.(*Dgii)
		registers, err := dgii.GetRegister(params.Query)
		output = registers
		searchErr = err
	case DomainTypePGR:
		pgr := connector.(*Pgr)
		registers, err := pgr.Search(params.Query)
		output = registers
		searchErr = err
	default:
		return &DomainSearchResult{
			Success:    false,
			Error:      fmt.Errorf("unsupported domain type: %s", domainType),
			DomainType: domainType,
		}, fmt.Errorf("unsupported domain type: %s", domainType)
	}

	// Extract keywords from the result
	var keywordsPerCategory map[DataCategory][]string
	if searchErr == nil && output != nil {
		switch domainType {
		case DomainTypeONAPI:
			if entities, ok := output.([]Entity); ok {
				spew.Dump("DomainTypePGR---entities", entities)
				keywordsPerCategory = GetCategoryByKeywords(&Onapi{}, entities)
			}
		case DomainTypeSCJ:
			if cases, ok := output.([]ScjCase); ok {
				keywordsPerCategory = GetCategoryByKeywords(&Scj{}, cases)
			}
		case DomainTypeDGII:
			if registers, ok := output.([]Register); ok {
				keywordsPerCategory = GetCategoryByKeywords(&Dgii{}, registers)
			}
		case DomainTypePGR:
			if registers, ok := output.([]PGRNews); ok {
				keywordsPerCategory = GetCategoryByKeywords(&Pgr{}, registers)
			}
		}
	}

	return &DomainSearchResult{
		Success:             searchErr == nil,
		Error:               searchErr,
		DomainType:          domainType,
		SearchParameter:     params.Query,
		KeywordsPerCategory: keywordsPerCategory,
		Output:              output,
	}, searchErr
}

// CreateDomainConnector creates a domain connector instance based on the domain type
func CreateDomainConnector(domainType DomainType) (any, error) {
	switch domainType {
	case DomainTypeONAPI:
		onapi := NewOnapiDomain()
		return &onapi, nil
	case DomainTypeSCJ:
		scj := NewScjDomain()
		return &scj, nil
	case DomainTypeDGII:
		dgii := NewDgiiDomain()
		return &dgii, nil
	case DomainTypePGR:
		pgr := NewPgrDomain()
		return &pgr, nil
	default:
		return nil, fmt.Errorf("unsupported domain type: %s", domainType)
	}
}

// SearchMultipleDomains performs searches across multiple domains
func SearchMultipleDomains(domainTypes []DomainType, params DomainSearchParams) []*DomainSearchResult {
	results := make([]*DomainSearchResult, 0, len(domainTypes))

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
	DomainType          DomainType
	SearchParameter     string
	Category            DataCategory
	Keywords            []string
	Success             bool
	Error               error
	Output              any
	KeywordsPerCategory map[DataCategory][]string
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
	availableDomains []DomainType,
	config DynamicPipelineConfig,
) (*DynamicPipelineResult, error) {

	// Initialize the pipeline
	pipeline := &DynamicPipelineResult{
		Steps:  make([]DynamicPipelineStep, 0),
		Config: config,
	}

	// Track searched keywords per domain to avoid duplicates
	searchedKeywordsPerDomain := make(map[DomainType]map[string]bool)
	for _, domainType := range availableDomains {
		searchedKeywordsPerDomain[domainType] = make(map[string]bool)
	}

	// Start with the initial query
	initialStep := DynamicPipelineStep{
		DomainType:      DomainTypeONAPI, // Start with ONAPI as it's most comprehensive
		SearchParameter: initialQuery,
		Category:        DataCategoryCompanyName,
		Keywords:        []string{initialQuery},
		Depth:           0,
	}

	// Add initial step
	pipeline.Steps = append(pipeline.Steps,
		initialStep,
		DynamicPipelineStep{
			DomainType:      DomainTypeDGII, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        DataCategoryContributorID,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      DomainTypePGR, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        DataCategoryPersonName,
			Keywords:        []string{initialQuery},
			Depth:           0,
		},
		DynamicPipelineStep{
			DomainType:      DomainTypeSCJ, // Start with ONAPI as it's most comprehensive
			SearchParameter: initialQuery,
			Category:        DataCategoryContributorID,
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
		result, err := SearchDomain(step.DomainType, DomainSearchParams{Query: step.SearchParameter})
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
	keywordsPerCategory map[DataCategory][]string,
	availableDomains []DomainType,
	searchedKeywordsPerDomain map[DomainType]map[string]bool,
	depth int,
	config DynamicPipelineConfig,
) []DynamicPipelineStep {

	var newSteps []DynamicPipelineStep

	// Get all available domain connectors
	domainConnectors := make(map[DomainType]any)
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
func getSearchableCategories(connector any) []DataCategory {
	switch c := connector.(type) {
	case *Onapi:
		return c.GetListOfSearchableCategory()
	case *Scj:
		return c.GetListOfSearchableCategory()
	case *Dgii:
		return c.GetListOfSearchableCategory()
	case *Pgr:
		return c.GetListOfSearchableCategory()
	default:
		return []DataCategory{}
	}
}

// contains checks if a slice contains a specific element
func contains(slice []DataCategory, item DataCategory) bool {
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
	availableDomains []DomainType,
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

		result, err := SearchDomain(step.DomainType, DomainSearchParams{Query: step.SearchParameter})
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
