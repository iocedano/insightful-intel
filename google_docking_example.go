package main

import (
	"fmt"
	"insightful-intel/internal/domain"
	"log"
)

func main() {
	fmt.Println("=== Google Docking Builder Examples ===\n")

	// Example 1: Basic search using the builder pattern
	fmt.Println("1. Basic Search:")
	results, err := domain.NewGoogleDockingBuilder().
		Query("machine learning").
		Build()

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d results:\n", len(results))
	for i, result := range results {
		fmt.Printf("  %d. %s (relevance: %.2f)\n", i+1, result.Title, result.Relevance)
	}
	fmt.Println()

	// Example 2: Advanced search with parameters
	fmt.Println("2. Advanced Search with Parameters:")
	results, err = domain.NewGoogleDockingBuilder().
		Query("artificial intelligence").
		MaxResults(5).
		MinRelevance(0.3).
		ExactMatch(false).
		CaseSensitive(false).
		Build()

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d results:\n", len(results))
	for i, result := range results {
		fmt.Printf("  %d. %s (relevance: %.2f, rank: %d)\n", i+1, result.Title, result.Relevance, result.Rank)
	}
	fmt.Println()

	// Example 3: Search with keyword filtering
	fmt.Println("3. Search with Keyword Filtering:")
	results, err = domain.NewGoogleDockingBuilder().
		Query("data science").
		IncludeKeywords("python", "algorithm").
		ExcludeKeywords("spam", "advertisement").
		MaxResults(3).
		Build()

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Found %d filtered results:\n", len(results))
	for i, result := range results {
		fmt.Printf("  %d. %s (relevance: %.2f)\n", i+1, result.Title, result.Relevance)
	}
	fmt.Println()

	// Example 4: Search with statistics
	fmt.Println("4. Search with Statistics:")
	results, stats, err := domain.NewGoogleDockingBuilder().
		Query("deep learning").
		MaxResults(4).
		BuildWithStats()

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Search Statistics:\n")
	fmt.Printf("  Total results: %v\n", stats["total_results"])
	fmt.Printf("  Average relevance: %.2f\n", stats["average_relevance"])
	fmt.Printf("  Max relevance: %.2f\n", stats["max_relevance"])
	fmt.Printf("  Min relevance: %.2f\n", stats["min_relevance"])
	fmt.Println()

	// Example 5: Using helper functions
	fmt.Println("5. Using Helper Functions:")

	// Quick search
	quickResults, err := domain.QuickSearch("neural networks")
	if err != nil {
		log.Printf("Quick search error: %v", err)
	} else {
		fmt.Printf("Quick search found %d results\n", len(quickResults))
	}

	// Advanced search
	advancedResults, err := domain.AdvancedSearch("computer vision", 3, 0.5)
	if err != nil {
		log.Printf("Advanced search error: %v", err)
	} else {
		fmt.Printf("Advanced search found %d results\n", len(advancedResults))
	}

	// Exact search
	exactResults, err := domain.ExactSearch("exact match")
	if err != nil {
		log.Printf("Exact search error: %v", err)
	} else {
		fmt.Printf("Exact search found %d results\n", len(exactResults))
	}

	// Case sensitive search
	caseResults, err := domain.CaseSensitiveSearch("Case Sensitive")
	if err != nil {
		log.Printf("Case sensitive search error: %v", err)
	} else {
		fmt.Printf("Case sensitive search found %d results\n", len(caseResults))
	}

	// Filtered search
	filteredResults, err := domain.FilteredSearch("machine learning", []string{"AI", "algorithm"}, []string{"spam"})
	if err != nil {
		log.Printf("Filtered search error: %v", err)
	} else {
		fmt.Printf("Filtered search found %d results\n", len(filteredResults))
	}
	fmt.Println()

	// Example 6: Using the GoogleDocking struct directly
	fmt.Println("6. Using GoogleDocking Struct Directly:")
	gd := domain.NewGoogleDockingDomain()

	// Basic search
	directResults, err := gd.Search("natural language processing")
	if err != nil {
		log.Printf("Direct search error: %v", err)
	} else {
		fmt.Printf("Direct search found %d results\n", len(directResults))
	}

	// Search with filters
	filters := map[string]interface{}{
		"max_results":      2,
		"min_relevance":    0.4,
		"exact_match":      false,
		"case_sensitive":   false,
		"include_keywords": []string{"NLP", "text"},
		"exclude_keywords": []string{"spam"},
	}

	filterResults, err := gd.SearchWithFilters("text processing", filters)
	if err != nil {
		log.Printf("Filter search error: %v", err)
	} else {
		fmt.Printf("Filter search found %d results\n", len(filterResults))
	}

	// Get suggestions
	suggestions, err := gd.GetSearchSuggestions("machine")
	if err != nil {
		log.Printf("Suggestions error: %v", err)
	} else {
		fmt.Printf("Search suggestions: %v\n", suggestions)
	}
	fmt.Println()

	// Example 7: Data extraction by category
	fmt.Println("7. Data Extraction by Category:")
	if len(results) > 0 {
		result := results[0]
		fmt.Printf("Analyzing result: %s\n", result.Title)

		// Extract company names
		companies := gd.GetDataByCategory(result, domain.KeywordCategoryCompanyName)
		fmt.Printf("  Company names: %v\n", companies)

		// Extract person names
		persons := gd.GetDataByCategory(result, domain.KeywordCategoryPersonName)
		fmt.Printf("  Person names: %v\n", persons)

		// Extract addresses
		addresses := gd.GetDataByCategory(result, domain.KeywordCategoryAddress)
		fmt.Printf("  Addresses: %v\n", addresses)

		// Extract social media
		social := gd.GetDataByCategory(result, domain.KeywordCategorySocialMedia)
		fmt.Printf("  Social media: %v\n", social)
	}
	fmt.Println()

	// Example 8: Generic connector interface
	fmt.Println("8. Generic Connector Interface:")
	if len(results) > 0 {
		result := results[0]

		// Process data using the generic interface
		processed, err := domain.ProcessGenericData(&gd, result)
		if err != nil {
			log.Printf("Processing error: %v", err)
		} else {
			fmt.Printf("Processed result: %s\n", processed.Title)
			fmt.Printf("  URL: %s\n", processed.URL)
			fmt.Printf("  Keywords: %v\n", processed.Keywords)
		}

		// Validate data
		err = gd.ValidateData(result)
		if err != nil {
			fmt.Printf("Validation error: %v\n", err)
		} else {
			fmt.Println("Data validation passed")
		}
	}
	fmt.Println()

	// Example 9: Keyword categories
	fmt.Println("9. Keyword Categories:")
	searchableCategories := gd.GetSearchableKeywordCategories()
	foundCategories := gd.GetFoundKeywordCategories()

	fmt.Printf("Searchable categories: %v\n", searchableCategories)
	fmt.Printf("Found categories: %v\n", foundCategories)
	fmt.Println()

	// Example 10: Domain type
	fmt.Println("10. Domain Information:")
	fmt.Printf("Domain type: %s\n", gd.GetDomainType())
	fmt.Printf("Base path: %s\n", gd.BasePath)
	fmt.Println()

	fmt.Println("=== All Examples Completed Successfully! ===")
}

