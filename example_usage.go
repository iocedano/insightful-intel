package main

import (
	"fmt"
	"insightful-intel/internal/domain"
)

// This file demonstrates how to use the new domain search functions

func main() {
	// Example 1: Search a single domain
	fmt.Println("=== Example 1: Single Domain Search ===")
	searchParams := domain.DomainSearchParams{
		Query: "Novasco",
	}

	// Search ONAPI domain
	result, err := domain.SearchDomain(domain.DomainTypeONAPI, searchParams)
	if err != nil {
		fmt.Printf("Error searching ONAPI: %v\n", err)
	} else {
		fmt.Printf("ONAPI Search Result: Success=%v, Domain=%s, Query=%s\n",
			result.Success, result.DomainType, result.SearchParameter)
	}

	// Example 2: Search multiple domains at once
	fmt.Println("\n=== Example 2: Multiple Domain Search ===")
	domainTypes := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	results := domain.SearchMultipleDomains(domainTypes, searchParams)
	for _, result := range results {
		fmt.Printf("Domain: %s, Success: %v, Error: %v\n",
			result.DomainType, result.Success, result.Error)
	}

	// Example 3: Create a domain connector directly
	fmt.Println("\n=== Example 3: Direct Domain Connector Creation ===")
	connector, err := domain.CreateDomainConnector(domain.DomainTypeDGII)
	if err != nil {
		fmt.Printf("Error creating connector: %v\n", err)
	} else {
		fmt.Printf("Created connector: %T\n", connector)
	}

	// Example 4: Search with different queries
	fmt.Println("\n=== Example 4: Different Queries ===")
	queries := []string{"NOVASCO REAL ESTATE", "YVES ALEXANDRE GIROUX", "132-33710-7"}

	for _, query := range queries {
		fmt.Printf("\nSearching for: %s\n", query)
		params := domain.DomainSearchParams{Query: query}

		// Search DGII for company names and contributor IDs
		dgiiResult, err := domain.SearchDomain(domain.DomainTypeDGII, params)
		if err != nil {
			fmt.Printf("  DGII Error: %v\n", err)
		} else {
			fmt.Printf("  DGII Success: %v, Results: %d\n",
				dgiiResult.Success, len(dgiiResult.Output.([]domain.Register)))
		}
	}
}

// Example function that shows how to use the domain search in your own code
func searchCompany(companyName string) ([]domain.DomainSearchResult, error) {
	searchParams := domain.DomainSearchParams{
		Query: companyName,
	}

	// Search all available domains
	domainTypes := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	results := domain.SearchMultipleDomains(domainTypes, searchParams)
	return results, nil
}

// Example function for searching a specific domain
func searchSpecificDomain(domainType domain.DomainType, query string) (*domain.DomainSearchResult, error) {
	searchParams := domain.DomainSearchParams{
		Query: query,
	}

	return domain.SearchDomain(domainType, searchParams)
}
