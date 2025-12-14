package domain

import "fmt"

type KeywordCategory string

const (
	KeywordCategoryAddress       KeywordCategory = "address"
	KeywordCategoryCompanyName   KeywordCategory = "company_name"
	KeywordCategoryContributorID KeywordCategory = "contributor_id"
	KeywordCategoryPersonName    KeywordCategory = "person_name"
	KeywordCategorySocialMedia   KeywordCategory = "social_media"
	KeywordCategoryFileType      KeywordCategory = "file_type"
	KeywordCategoryXSocialMedia  KeywordCategory = "x_social_media"
)

// DomainConnector is an interface for domain-specific connectors that can process,
// search, and extract keywords from domain data of a specific type.
type DomainConnector[T any] interface {
	ProcessData(data T) (T, error)
	ValidateData(data T) error
	TransformData(data T) T
	GetDataByCategory(data T, category KeywordCategory) []string
	GetSearchableKeywordCategories() []KeywordCategory
	GetFoundKeywordCategories() []KeywordCategory
	GetDomainType() DomainType
	Search(query string) ([]T, error)
}

// Extended interface that includes domain connector methods
type ExtendedConnectorInterface[T any] interface {
	DomainConnector[T]
}

// ProcessDomainData processes data using a domain connector
func ProcessDomainData[T any](connector DomainConnector[T], data T) (T, error) {
	return connector.ProcessData(data)
}

// ProcessBatchDomainData processes a batch of data using a domain connector
func ProcessBatchDomainData[T any](connector DomainConnector[T], data []T) ([]T, error) {
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

// GetCategoryByKeywords extracts keywords by category from domain data using a connector
func GetCategoryByKeywords[T any](connector DomainConnector[T], data []T) map[KeywordCategory][]string {
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
	DomainTypeGoogleDorking DomainType = "GOOGLE_DOCKING"
	DomainTypeSocialMedia   DomainType = "SOCIAL_MEDIA"
	DomainTypeXSocialMedia  DomainType = "X_SOCIAL_MEDIA"
	DomainTypeFileType      DomainType = "FILE_TYPE"
)

// AllDomainTypes returns a list of all available domain types (excluding ERROR)
func AllDomainTypes() []DomainType {
	return []DomainType{
		DomainTypeONAPI,
		DomainTypeSCJ,
		DomainTypeDGII,
		DomainTypePGR,
		DomainTypeGoogleDorking,
		DomainTypeSocialMedia,
		DomainTypeXSocialMedia,
		DomainTypeFileType,
	}
}

// DefaultDomainTypes returns a list of commonly used domain types for default searches
func DefaultDomainTypes() []DomainType {
	return []DomainType{
		DomainTypeONAPI,
		DomainTypeSCJ,
		DomainTypeDGII,
	}
}

// StringToDomainType maps string identifiers (typically from URL parameters) to DomainType
var StringToDomainType = map[string]DomainType{
	"onapi":          DomainTypeONAPI,
	"scj":            DomainTypeSCJ,
	"dgii":           DomainTypeDGII,
	"pgr":            DomainTypePGR,
	"docking":        DomainTypeGoogleDorking,
	"social_media":   DomainTypeSocialMedia,
	"x_social_media": DomainTypeXSocialMedia,
	"file_type":      DomainTypeFileType,
}

// DomainTypeToString maps DomainType to string identifiers (for URL parameters)
var DomainTypeToString = map[DomainType]string{
	DomainTypeONAPI:         "onapi",
	DomainTypeSCJ:           "scj",
	DomainTypeDGII:          "dgii",
	DomainTypePGR:           "pgr",
	DomainTypeGoogleDorking: "docking",
	DomainTypeSocialMedia:   "social_media",
	DomainTypeXSocialMedia:  "x_social_media",
	DomainTypeFileType:      "file_type",
}

// GetDomainTypeFromString converts a string to DomainType, returns error if not found
func GetDomainTypeFromString(s string) (DomainType, error) {
	if dt, ok := StringToDomainType[s]; ok {
		return dt, nil
	}
	return DomainTypeERROR, fmt.Errorf("unknown domain type: %s", s)
}

// IsValidDomainType checks if a DomainType is valid (not ERROR)
func IsValidDomainType(dt DomainType) bool {
	for _, validType := range AllDomainTypes() {
		if dt == validType {
			return true
		}
	}
	return false
}

// DomainSearchParams holds the search parameters for different domains
type DomainSearchParams struct {
	Query string
}

// DomainSearchResult represents the result of a domain search
type DomainSearchResult struct {
	ID                  ID
	PipelineStepsID     ID
	DomainType          DomainType
	SearchParameter     string
	KeywordsPerCategory map[KeywordCategory][]string `json:"keywordsPerCategory"`
	Output              any                          `json:"output"`
	Success             bool                         `json:"success"`
	Error               error                        `json:"error"`
}
