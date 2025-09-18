# Domain Search Functions

This document explains how to use the new generic domain search functions that allow you to pass values based on the domain type.

## Overview

The new domain search functionality provides a unified interface for searching across different domains (ONAPI, SCJ, DGII, PGR) without having to manually instantiate and call each domain's specific methods.

## Key Components

### 1. Domain Types
```go
type DomainType string

const (
    DomainTypeONAPI DomainType = "ONAPI"
    DomainTypeSCJ   DomainType = "SCJ"
    DomainTypeDGII  DomainType = "DGII"
    DomainTypePGR   DomainType = "PGR"
)
```

### 2. Search Parameters
```go
type DomainSearchParams struct {
    Query string
    // Additional parameters can be added as needed
}
```

### 3. Search Results
```go
type DomainSearchResult struct {
    Success             bool
    Error               error
    DomainType          DomainType
    SearchParameter     string
    KeywordsPerCategory map[DataCategory][]string
    Output              any
}
```

## Usage Examples

### 1. Search a Single Domain

```go
searchParams := domain.DomainSearchParams{
    Query: "Novasco",
}

// Search ONAPI domain
result, err := domain.SearchDomain(domain.DomainTypeONAPI, searchParams)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

fmt.Printf("Success: %v, Domain: %s", result.Success, result.DomainType)
```

### 2. Search Multiple Domains

```go
domainTypes := []domain.DomainType{
    domain.DomainTypeONAPI,
    domain.DomainTypeSCJ,
    domain.DomainTypeDGII,
}

results := domain.SearchMultipleDomains(domainTypes, searchParams)
for _, result := range results {
    fmt.Printf("Domain: %s, Success: %v", result.DomainType, result.Success)
}
```

### 3. Create Domain Connector Directly

```go
connector, err := domain.CreateDomainConnector(domain.DomainTypeDGII)
if err != nil {
    log.Printf("Error creating connector: %v", err)
    return
}

// Use the connector as needed
dgii := connector.(*domain.Dgii)
registers, err := dgii.GetRegister("132-33710-7")
```

## HTTP API Usage

### Search All Domains
```
GET /search?q=Novasco
```

### Search Specific Domain
```
GET /search?q=Novasco&domain=onapi
GET /search?q=Novasco&domain=scj
GET /search?q=Novasco&domain=dgii
```

## Benefits

1. **Unified Interface**: Single function to search any domain
2. **Type Safety**: Compile-time checking of domain types
3. **Extensibility**: Easy to add new domains
4. **Consistency**: Standardized result format across all domains
5. **Error Handling**: Centralized error handling
6. **Keyword Extraction**: Automatic extraction of searchable keywords

## Adding New Domains

To add a new domain:

1. Add the domain type to the `DomainType` constants
2. Implement the `GenericConnector[T]` interface for your domain
3. Add a case in `CreateDomainConnector` function
4. Add a case in `SearchDomain` function

## Migration from Old Code

### Before (Old Way)
```go
onapi := domain.NewOnapiDomain()
entities, err := onapi.SearchComercialName("Novasco")

scj := domain.NewScjDomain()
cases, err := scj.Search("Novasco")

dgii := domain.NewDgiiDomain()
registers, err := dgii.GetRegister("Novasco")
```

### After (New Way)
```go
searchParams := domain.DomainSearchParams{Query: "Novasco"}

// Single domain
result, err := domain.SearchDomain(domain.DomainTypeONAPI, searchParams)

// Multiple domains
domainTypes := []domain.DomainType{
    domain.DomainTypeONAPI,
    domain.DomainTypeSCJ,
    domain.DomainTypeDGII,
}
results := domain.SearchMultipleDomains(domainTypes, searchParams)
```

## Error Handling

All functions return proper error handling:
- `SearchDomain` returns `(*DomainSearchResult, error)`
- `SearchMultipleDomains` returns `[]*DomainSearchResult` (errors are contained within each result)
- `CreateDomainConnector` returns `(any, error)`

Check the `Success` field in `DomainSearchResult` to determine if the search was successful, and the `Error` field for specific error details.
