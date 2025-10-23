# Google Docking Builder

A comprehensive Go implementation of a Google Docking string search system with a fluent builder pattern API.

## Overview

The Google Docking Builder provides a powerful and flexible way to perform string searches with advanced relevance scoring, fuzzy matching, and keyword filtering. It implements the `GenericConnector` interface and integrates seamlessly with the existing domain system.

## Features

- **Fluent Builder Pattern**: Chain methods for intuitive API usage
- **Advanced Relevance Scoring**: Multi-factor scoring algorithm with configurable weights
- **Fuzzy String Matching**: Levenshtein distance-based similarity matching
- **Keyword Filtering**: Include/exclude keywords for precise results
- **Case Sensitivity**: Configurable case-sensitive and case-insensitive search
- **Exact Matching**: Support for exact string matching
- **Data Extraction**: Extract company names, person names, addresses, and social media
- **Search Statistics**: Comprehensive statistics about search results
- **Search Suggestions**: Generate search suggestions based on queries
- **Generic Connector**: Implements the `GenericConnector[GoogleDockingResult]` interface

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "insightful-intel/internal/domain"
)

func main() {
    // Simple search
    results, err := domain.NewGoogleDockingBuilder().
        Query("machine learning").
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    for _, result := range results {
        fmt.Printf("%s (relevance: %.2f)\n", result.Title, result.Relevance)
    }
}
```

### Advanced Usage

```go
// Advanced search with multiple parameters
results, err := domain.NewGoogleDockingBuilder().
    Query("artificial intelligence").
    MaxResults(20).
    MinRelevance(0.5).
    ExactMatch(false).
    CaseSensitive(false).
    IncludeKeywords("AI", "algorithm", "neural").
    ExcludeKeywords("spam", "advertisement").
    Build()
```

## API Reference

### GoogleDockingBuilder

The main builder class for constructing Google Docking searches.

#### Constructor

```go
func NewGoogleDockingBuilder() *GoogleDockingBuilder
```

Creates a new Google Docking builder with default parameters.

#### Methods

##### Query(query string) *GoogleDockingBuilder
Sets the search query string.

```go
builder.Query("machine learning")
```

##### MaxResults(max int) *GoogleDockingBuilder
Sets the maximum number of results to return.

```go
builder.MaxResults(50)
```

##### MinRelevance(min float64) *GoogleDockingBuilder
Sets the minimum relevance threshold (0.0-1.0).

```go
builder.MinRelevance(0.3)
```

##### ExactMatch(exact bool) *GoogleDockingBuilder
Enables or disables exact matching.

```go
builder.ExactMatch(true)
```

##### CaseSensitive(caseSensitive bool) *GoogleDockingBuilder
Enables or disables case-sensitive search.

```go
builder.CaseSensitive(true)
```

##### IncludeKeywords(keywords ...string) *GoogleDockingBuilder
Adds keywords that must be present in results.

```go
builder.IncludeKeywords("AI", "algorithm", "neural")
```

##### ExcludeKeywords(keywords ...string) *GoogleDockingBuilder
Adds keywords to exclude from results.

```go
builder.ExcludeKeywords("spam", "advertisement")
```

##### Build() ([]GoogleDockingResult, error)
Executes the search and returns results.

```go
results, err := builder.Build()
```

##### BuildWithStats() ([]GoogleDockingResult, map[string]interface{}, error)
Executes the search and returns results with statistics.

```go
results, stats, err := builder.BuildWithStats()
```

### GoogleDockingResult

Represents a search result from Google Docking.

```go
type GoogleDockingResult struct {
    URL         string   `json:"url"`          // Result URL
    Title       string   `json:"title"`        // Result title
    Description string   `json:"description"`  // Result description
    Relevance   float64  `json:"relevance"`    // Relevance score (0.0-1.0)
    Rank        int      `json:"rank"`         // Result rank
    Keywords    []string `json:"keywords"`     // Associated keywords
}
```

### GoogleDockingSearchParams

Holds parameters for Google Docking search.

```go
type GoogleDockingSearchParams struct {
    Query           string   `json:"query"`            // Search query
    MaxResults      int      `json:"max_results"`      // Maximum results
    MinRelevance    float64  `json:"min_relevance"`    // Minimum relevance threshold
    ExactMatch      bool     `json:"exact_match"`      // Enable exact matching
    CaseSensitive   bool     `json:"case_sensitive"`   // Case-sensitive search
    IncludeKeywords []string `json:"include_keywords"` // Required keywords
    ExcludeKeywords []string `json:"exclude_keywords"` // Excluded keywords
}
```

## Helper Functions

### Quick Search Functions

```go
// Simple one-liner search
results, err := domain.QuickSearch("query")

// Advanced search with parameters
results, err := domain.AdvancedSearch("query", maxResults, minRelevance)

// Exact match search
results, err := domain.ExactSearch("query")

// Case-sensitive search
results, err := domain.CaseSensitiveSearch("Query")

// Filtered search
results, err := domain.FilteredSearch("query", includeKeywords, excludeKeywords)
```

## Direct Usage

### Using GoogleDocking Struct

```go
gd := domain.NewGoogleDockingDomain()

// Basic search
results, err := gd.Search("query")

// Search with parameters
params := domain.GoogleDockingSearchParams{
    Query:        "query",
    MaxResults:   10,
    MinRelevance: 0.3,
    ExactMatch:   false,
    CaseSensitive: false,
}
results, err := gd.SearchWithParams(params)

// Search with filters
filters := map[string]interface{}{
    "max_results":   20,
    "min_relevance": 0.5,
    "exact_match":   true,
    "case_sensitive": false,
    "include_keywords": []string{"AI", "algorithm"},
    "exclude_keywords": []string{"spam"},
}
results, err := gd.SearchWithFilters("query", filters)
```

### Utility Methods

```go
// Get search suggestions
suggestions, err := gd.GetSearchSuggestions("partial query")

// Get search statistics
stats := gd.GetSearchStatistics(results)

// Extract data by category
companies := gd.GetDataByCategory(result, domain.KeywordCategoryCompanyName)
persons := gd.GetDataByCategory(result, domain.KeywordCategoryPersonName)
addresses := gd.GetDataByCategory(result, domain.KeywordCategoryAddress)
social := gd.GetDataByCategory(result, domain.KeywordCategorySocialMedia)
```

## Relevance Scoring Algorithm

The Google Docking system uses a sophisticated relevance scoring algorithm:

### Scoring Factors

1. **Title Match** (weight: 3.0)
   - Exact matches get the highest score
   - Partial matches with position and frequency bonuses
   - Fuzzy matches using Levenshtein distance

2. **Description Match** (weight: 2.0)
   - Similar scoring to title matches
   - Medium weight for content relevance

3. **URL Match** (weight: 1.0)
   - Lower weight for URL relevance
   - Useful for domain-specific searches

4. **Keywords Match** (weight: 1.5)
   - Matches against associated keywords
   - Medium-high weight for keyword relevance

5. **Exact Match Bonus** (+2.0)
   - Additional bonus for exact matches when enabled

6. **Include Keywords Bonus**
   - Bonus for including required keywords

7. **Exclude Keywords Penalty**
   - Penalty for including excluded keywords

### String Matching

- **Exact Match**: Direct string comparison
- **Partial Match**: Substring matching with position and frequency bonuses
- **Fuzzy Match**: Levenshtein distance with similarity threshold of 0.6

### Score Normalization

Final scores are normalized to a 0.0-1.0 range for consistency.

## Data Extraction

The system can extract various types of data from search results:

### Company Names
- Identifies capitalized words that might be company names
- Recognizes common company suffixes (Inc, Corp, LLC, Ltd, Co, Company)

### Person Names
- Identifies patterns like "First Last" or "Mr. Last"
- Recognizes title prefixes (Mr., Mrs., Ms., Dr.)

### Addresses
- Identifies street numbers and address patterns
- Collects following words that might be part of addresses

### Social Media
- Identifies social media handles (@)
- Recognizes social media URLs (twitter.com, facebook.com, etc.)

## Generic Connector Interface

The Google Docking system implements the `GenericConnector[GoogleDockingResult]` interface:

```go
// Process data
processed, err := domain.ProcessGenericData(&gd, result)

// Validate data
err := gd.ValidateData(result)

// Transform data
transformed := gd.TransformData(result)

// Get searchable categories
categories := gd.GetSearchableKeywordCategories()

// Get found categories
categories := gd.GetFoundKeywordCategories()
```

## Integration with Domain System

The Google Docking system integrates with the existing domain system:

```go
// Search using domain search
result, err := domain.SearchDomain(domain.DomainTypeGoogleDocking, searchParams)

// Search multiple domains
domains := []domain.DomainType{
    domain.DomainTypeONAPI,
    domain.DomainTypeSCJ,
    domain.DomainTypeDGII,
    domain.DomainTypePGR,
    domain.DomainTypeGoogleDocking,
}
results := domain.SearchMultipleDomains(domains, searchParams)
```

## Error Handling

All functions return proper error handling:

```go
results, err := builder.Build()
if err != nil {
    // Handle error appropriately
    log.Printf("Search failed: %v", err)
    return
}
```

Common error scenarios:
- Empty query strings
- Invalid relevance scores
- Missing required fields in results
- Network or API errors (in production)

## Performance Considerations

- Results are limited by `MaxResults` parameter
- Relevance filtering reduces processing overhead
- Fuzzy matching has O(m*n) complexity for Levenshtein distance
- Consider caching for frequently searched terms

## Examples

See `google_docking_example.go` for comprehensive usage examples covering:

1. Basic search using builder pattern
2. Advanced search with parameters
3. Search with keyword filtering
4. Search with statistics
5. Using helper functions
6. Using GoogleDocking struct directly
7. Data extraction by category
8. Generic connector interface
9. Keyword categories
10. Domain information

## Running Examples

```bash
go run google_docking_example.go
```

This will demonstrate all the features and capabilities of the Google Docking Builder system.

## Testing

The system includes comprehensive test coverage for:
- Basic search functionality
- Advanced parameter handling
- Relevance scoring accuracy
- String matching algorithms
- Data validation and transformation
- Error handling
- Generic connector interface compliance

## Future Enhancements

Potential improvements for production use:
- Real Google API integration
- Caching mechanisms
- Performance optimizations
- Additional string matching algorithms
- Machine learning-based relevance scoring
- Real-time search suggestions
- Advanced filtering options

