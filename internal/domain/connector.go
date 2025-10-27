package domain

type KeywordCategory string

const (
	KeywordCategoryAddress       KeywordCategory = "address"
	KeywordCategoryCompanyName   KeywordCategory = "company_name"
	KeywordCategoryContributorID KeywordCategory = "contributor_id"
	KeywordCategoryPersonName    KeywordCategory = "person_name"
	KeywordCategorySocialMedia   KeywordCategory = "social_media"
)

// Generic interface for connectors that can process data of the same type
type GenericConnector[T any] interface {
	ProcessData(data T) (T, error)
	ValidateData(data T) error
	TransformData(data T) T
	GetDataByCategory(data T, category KeywordCategory) []string
	GetSearchableKeywordCategories() []KeywordCategory
	GetFoundKeywordCategories() []KeywordCategory
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

func GetCategoryByKeywords[T any](connector GenericConnector[T], data []T) map[KeywordCategory][]string {
	result := map[KeywordCategory][]string{}

	for _, d := range data {
		for _, rcv := range connector.GetFoundKeywordCategories() {
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
	DomainTypeERROR         DomainType = "error"
	DomainTypeONAPI         DomainType = "ONAPI"
	DomainTypeSCJ           DomainType = "SCJ"
	DomainTypeDGII          DomainType = "DGII"
	DomainTypePGR           DomainType = "PGR"
	DomainTypeGoogleDocking DomainType = "GOOGLE_DOCKING"
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
	ID                  ID
	DomainType          DomainType
	SearchParameter     string
	KeywordsPerCategory map[KeywordCategory][]string `json:"keywordsPerCategory"`
	Output              any                          `json:"output"`
	Success             bool                         `json:"success"`
	Error               error                        `json:"error"`
}
