package domain

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
	GetName() string
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
