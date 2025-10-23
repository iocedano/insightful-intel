package main

import (
	"fmt"
	"insightful-intel/internal/domain"
	"log"
)

// This file demonstrates the new dynamic pipeline functionality

func main() {
	fmt.Println("=== Dynamic Pipeline Example ===")

	// Example 1: Basic dynamic pipeline
	fmt.Println("\n1. Basic Dynamic Pipeline")
	basicExample()

	// Example 2: Custom configuration
	fmt.Println("\n2. Custom Configuration")
	customConfigExample()

	// Example 3: Step-by-step pipeline creation
	fmt.Println("\n3. Step-by-step Pipeline Creation")
	stepByStepExample()
}

func basicExample() {
	query := "Novasco"
	availableDomains := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	config := domain.DefaultDynamicPipelineConfig()
	config.MaxDepth = 2 // Limit depth for demo

	fmt.Printf("Searching for: %s\n", query)
	fmt.Printf("Available domains: %v\n", availableDomains)
	fmt.Printf("Max depth: %d\n", config.MaxDepth)

	result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Pipeline executed successfully!\n")
	fmt.Printf("Total steps: %d\n", result.TotalSteps)
	fmt.Printf("Successful steps: %d\n", result.SuccessfulSteps)
	fmt.Printf("Failed steps: %d\n", result.FailedSteps)
	fmt.Printf("Max depth reached: %d\n", result.MaxDepthReached)

	// Show some steps
	fmt.Println("\nFirst few steps:")
	for i, step := range result.Steps {
		if i >= 3 { // Show only first 3 steps
			break
		}
		fmt.Printf("  Step %d: %s searching '%s' (Success: %v, Depth: %d)\n",
			i+1, step.DomainType, step.SearchParameter, step.Success, step.Depth)
	}
}

func customConfigExample() {
	query := "NOVASCO REAL ESTATE"
	availableDomains := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeDGII,
	}

	// Custom configuration
	config := domain.DynamicPipelineConfig{
		MaxDepth:           4,
		MaxConcurrentSteps: 5,
		DelayBetweenSteps:  1,
		SkipDuplicates:     true,
	}

	fmt.Printf("Searching for: %s\n", query)
	fmt.Printf("Custom config: MaxDepth=%d, SkipDuplicates=%v\n",
		config.MaxDepth, config.SkipDuplicates)

	result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Pipeline completed with %d total steps\n", result.TotalSteps)

	// Show steps by domain
	domainSteps := make(map[domain.DomainType]int)
	for _, step := range result.Steps {
		domainSteps[step.DomainType]++
	}

	fmt.Println("Steps per domain:")
	for domainType, count := range domainSteps {
		fmt.Printf("  %s: %d steps\n", domainType, count)
	}
}

func stepByStepExample() {
	query := "YVES ALEXANDRE GIROUX"

	// Step 1: Create the pipeline structure
	fmt.Println("Step 1: Creating pipeline structure...")
	availableDomains := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	config := domain.DynamicPipelineConfig{
		MaxDepth:       3,
		SkipDuplicates: true,
	}

	pipeline, err := domain.CreateDynamicPipeline(query, availableDomains, config)
	if err != nil {
		log.Printf("Error creating pipeline: %v", err)
		return
	}

	fmt.Printf("Pipeline created with %d initial steps\n", len(pipeline.Steps))

	// Step 2: Execute the pipeline
	fmt.Println("Step 2: Executing pipeline...")
	result, err := domain.ExecuteDynamicPipeline(query, availableDomains, config)
	if err != nil {
		log.Printf("Error executing pipeline: %v", err)
		return
	}

	// Step 3: Analyze results
	fmt.Println("Step 3: Analyzing results...")
	analyzeResults(result)
}

func analyzeResults(result *domain.DynamicPipelineResult) {
	fmt.Printf("\n=== Pipeline Analysis ===\n")
	fmt.Printf("Total Steps: %d\n", result.TotalSteps)
	fmt.Printf("Successful: %d (%.1f%%)\n",
		result.SuccessfulSteps,
		float64(result.SuccessfulSteps)/float64(result.TotalSteps)*100)
	fmt.Printf("Failed: %d (%.1f%%)\n",
		result.FailedSteps,
		float64(result.FailedSteps)/float64(result.TotalSteps)*100)
	fmt.Printf("Max Depth Reached: %d\n", result.MaxDepthReached)

	// Group by domain
	domainStats := make(map[domain.DomainType]struct {
		Total    int
		Success  int
		Failed   int
		MaxDepth int
	})

	for _, step := range result.Steps {
		stats := domainStats[step.DomainType]
		stats.Total++
		if step.Success {
			stats.Success++
		} else {
			stats.Failed++
		}
		if step.Depth > stats.MaxDepth {
			stats.MaxDepth = step.Depth
		}
		domainStats[step.DomainType] = stats
	}

	fmt.Println("\nPer-Domain Statistics:")
	for domainType, stats := range domainStats {
		fmt.Printf("  %s: %d total, %d success, %d failed, max depth %d\n",
			domainType, stats.Total, stats.Success, stats.Failed, stats.MaxDepth)
	}

	// Show keyword extraction
	fmt.Println("\nKeyword Categories Found:")
	categoryCounts := make(map[domain.KeywordCategory]int)
	for _, step := range result.Steps {
		for category := range step.KeywordsPerCategory {
			categoryCounts[category]++
		}
	}

	for category, count := range categoryCounts {
		fmt.Printf("  %s: %d occurrences\n", category, count)
	}
}

// Example function showing how to use the dynamic pipeline in your own code
func searchWithDynamicPipeline(companyName string) (*domain.DynamicPipelineResult, error) {
	availableDomains := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
		domain.DomainTypePGR,
	}

	config := domain.DynamicPipelineConfig{
		MaxDepth:           5,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     true,
	}

	return domain.ExecuteDynamicPipeline(companyName, availableDomains, config)
}

// Example function for creating a custom pipeline configuration
func createCustomPipelineConfig(maxDepth int, skipDuplicates bool) domain.DynamicPipelineConfig {
	return domain.DynamicPipelineConfig{
		MaxDepth:           maxDepth,
		MaxConcurrentSteps: 15,
		DelayBetweenSteps:  1,
		SkipDuplicates:     skipDuplicates,
	}
}
