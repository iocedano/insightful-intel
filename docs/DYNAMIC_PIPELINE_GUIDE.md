# Dynamic Pipeline Guide

This guide explains the new dynamic pipeline system that automatically creates search pipelines based on `GetSearchableKeywordCategories()` from domain connectors.

## Overview

The dynamic pipeline system automatically:
1. **Discovers searchable categories** from each domain connector
2. **Extracts keywords** from search results
3. **Creates new search steps** based on extracted keywords
4. **Executes searches** across multiple domains in parallel
5. **Tracks progress** and prevents duplicate searches

## Key Components

### 1. DynamicPipelineConfig
```go
type DynamicPipelineConfig struct {
    MaxDepth           int  // Maximum pipeline depth (default: 5)
    MaxConcurrentSteps int  // Maximum concurrent steps (default: 10)
    DelayBetweenSteps  int  // Delay between steps in seconds (default: 2)
    SkipDuplicates     bool // Skip duplicate keyword searches (default: true)
}
```

### 2. DynamicPipelineStep
```go
type DynamicPipelineStep struct {
    DomainType          DomainType
    SearchParameter     string
    Category            KeywordCategory
    Keywords            []string
    Success             bool
    Error               error
    Output              any
    KeywordsPerCategory map[KeywordCategory][]string
    Depth               int
}
```

### 3. DynamicPipelineResult
```go
type DynamicPipelineResult struct {
    Steps           []DynamicPipelineStep
    TotalSteps      int
    SuccessfulSteps int
    FailedSteps     int
    MaxDepthReached int
    Config          DynamicPipelineConfig
}
```

## How It Works

### 1. Initial Search
The pipeline starts with an initial search query, typically using ONAPI as it provides the most comprehensive data.

### 2. Keyword Extraction
After each successful search, the system extracts keywords from the results using `GetDataByCategory()` for each `KeywordCategory` returned by `GetFoundKeywordCategories()`.

### 3. Category Matching
For each extracted keyword, the system finds domains that can search that category using `GetSearchableKeywordCategories()`.

### 4. Step Generation
New pipeline steps are created for each valid keyword-domain combination.

### 5. Execution
Steps are executed in parallel where possible, with configurable delays and concurrency limits.

## Usage Examples

### Basic Usage
```go
query := "Novasco"
availableDomains := []domain.DomainType{
    domain.DomainTypeONAPI,
    domain.DomainTypeSCJ,
    domain.DomainTypeDGII,
}

config := domain.DefaultDynamicPipelineConfig()
result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

### Custom Configuration
```go
config := domain.DynamicPipelineConfig{
    MaxDepth:           3,
    MaxConcurrentSteps: 5,
    DelayBetweenSteps:  1,
    SkipDuplicates:     true,
}

result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

### Step-by-Step Creation
```go
// Create pipeline structure first
pipeline, err := domain.CreateDynamicPipeline(query, availableDomains, config)

// Then execute it
result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
```

## HTTP API Endpoints

### Dynamic Pipeline Endpoint
```
GET /dynamic?q=Novasco&depth=3&skip_duplicates=true
```

**Parameters:**
- `q`: Search query (required)
- `depth`: Maximum pipeline depth (optional, default: 3, max: 10)
- `skip_duplicates`: Skip duplicate searches (optional, default: true)

**Response:**
```json
{
  "dynamic_result": {
    "steps": [...],
    "total_steps": 15,
    "successful_steps": 12,
    "failed_steps": 3,
    "max_depth_reached": 3,
    "config": {...}
  },
  "pipeline": [...],
  "summary": {
    "total_steps": 15,
    "successful_steps": 12,
    "failed_steps": 3,
    "max_depth_reached": 3
  }
}
```

## Domain Integration

### Required Methods
Each domain connector must implement:
```go
type GenericConnector[T any] interface {
    GetDomainType() DomainType
    GetSearchableKeywordCategories() []KeywordCategory
    GetFoundKeywordCategories() []KeywordCategory
    GetDataByCategory(data T, category KeywordCategory) []string
    // ... other methods
}
```

### Category Flow
1. **ONAPI** searches `company_name` → retrieves `company_name`, `person_name`, `address`
2. **SCJ** can search `person_name`, `company_name` → retrieves `person_name`
3. **DGII** can search `company_name` → retrieves `company_name`, `contributor_id`

## Configuration Options

### MaxDepth
Controls how deep the pipeline can go. Higher values mean more comprehensive searches but longer execution times.

### MaxConcurrentSteps
Limits the number of concurrent searches to prevent overwhelming external APIs.

### DelayBetweenSteps
Adds delays between steps to be respectful to external APIs.

### SkipDuplicates
Prevents searching the same keyword multiple times across domains.

## Best Practices

### 1. Start with ONAPI
ONAPI typically provides the most comprehensive initial data, making it ideal as the starting point.

### 2. Configure Depth Appropriately
- **Depth 1-2**: Quick searches, good for simple queries
- **Depth 3-4**: Comprehensive searches, good for complex investigations
- **Depth 5+**: Very thorough searches, use sparingly

### 3. Use Appropriate Delays
- **1-2 seconds**: For development/testing
- **2-5 seconds**: For production use
- **5+ seconds**: For very respectful API usage

### 4. Monitor Results
Check the `Summary` section to understand pipeline performance and adjust configuration accordingly.

## Error Handling

The system handles errors gracefully:
- Failed steps are marked but don't stop the pipeline
- Network errors are logged but don't crash the system
- Invalid domains are skipped automatically

## Performance Considerations

### Memory Usage
- Each step stores its full output in memory
- Consider `MaxDepth` and `MaxConcurrentSteps` for memory usage
- Large result sets may require pagination

### API Rate Limits
- Use `DelayBetweenSteps` to respect rate limits
- Monitor `FailedSteps` for rate limit indicators
- Consider implementing exponential backoff for production

### Parallel Execution
- Steps at the same depth can run in parallel
- Steps at different depths run sequentially
- `MaxConcurrentSteps` controls parallel execution

## Troubleshooting

### Common Issues

1. **No results from dynamic pipeline**
   - Check if initial query returns results
   - Verify domain connectors are working
   - Check `GetSearchableKeywordCategories()` implementation

2. **Too many duplicate searches**
   - Enable `SkipDuplicates`
   - Check keyword extraction logic
   - Verify `GetDataByCategory()` implementation

3. **Pipeline stops early**
   - Check `MaxDepth` setting
   - Verify all domains are available
   - Check for errors in step execution

4. **Performance issues**
   - Reduce `MaxDepth` or `MaxConcurrentSteps`
   - Increase `DelayBetweenSteps`
   - Check network connectivity

### Debugging
Enable debug logging to see:
- Step creation process
- Keyword extraction results
- Domain matching decisions
- Execution progress

## Migration from Static Pipeline

### Before (Static)
```go
// Manual step creation
onapi := domain.NewOnapiDomain()
entities, err := onapi.SearchComercialName("Novasco")

// Manual keyword extraction
keywords := domain.GetCategoryByKeywords(&domain.Onapi{}, entities)

// Manual step execution
for _, keyword := range keywords["person_name"] {
    scj := domain.NewScjDomain()
    cases, err := scj.Search(keyword)
    // ...
}
```

### After (Dynamic)
```go
// Automatic pipeline creation and execution
config := domain.DefaultDynamicPipelineConfig()
result, err := domain.ExecuteDynamicPipeline("Novasco", availableDomains, config)
```

The dynamic pipeline system automatically handles all the manual steps, making your code cleaner and more maintainable.
