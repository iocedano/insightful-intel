package domain

import "fmt"

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
		// pgr := NewPgrDomain()
		// return &pgr, nil
		return nil, fmt.Errorf("PGR domain not implemented yet")
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
