package server

import (
	"context"
	"encoding/json"
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/infra"
	"insightful-intel/internal/module"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	// mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/search", s.searchHandler)
	mux.HandleFunc("/dynamic", s.dynamicPipelineHandler)
	mux.HandleFunc("/health", s.healthHandler)

	// Repository-based routes
	mux.HandleFunc("/api/onapi", s.onapiHandler)
	mux.HandleFunc("/api/scj", s.scjHandler)
	mux.HandleFunc("/api/dgii", s.dgiiHandler)
	mux.HandleFunc("/api/pgr", s.pgrHandler)
	mux.HandleFunc("/api/docking", s.dockingHandler)
	mux.HandleFunc("/api/pipeline", s.pipelineHandler)
	mux.HandleFunc("/api/pipeline/steps", s.pipelineStepsHandler)
	mux.HandleFunc("/api/pipeline/save", s.savePipelineHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

type ConnectorPipeline struct {
	Success             bool                                `json:"success"`
	Error               error                               `json:"error"`
	Name                string                              `json:"name"`
	SearchParameter     string                              `json:"search_parameter"`
	Output              any                                 `json:"output"`
	KeywordsPerCategory map[domain.KeywordCategory][]string `json:"keywords_per_category"`
}

// Step executes a single step in the pipeline for a specific domain connector
func Step[T any](
	domainConnector domain.DomainConnector[T],
	searchableCategory []domain.KeywordCategory,
	category domain.KeywordCategory,
	keywords []string,
	seachedKeywordsPerDomain map[domain.DomainType][]string,
) []ConnectorPipeline {
	pipeline := []ConnectorPipeline{}

	// Check if the given category is searchable by the provided domain connector
	searchableCategories := domainConnector.GetSearchableKeywordCategories()
	if !slices.Contains(searchableCategories, category) {
		return pipeline
	}

	domainType := domainConnector.GetDomainType()
	if seachedKeywordsPerDomain[domainType] == nil {
		seachedKeywordsPerDomain[domainType] = []string{}
	}

	for _, keyword := range keywords {
		if slices.Contains(seachedKeywordsPerDomain[domainType], keyword) || keyword == "" {
			continue
		}

		result, err := module.SearchDomain(domainType, domain.DomainSearchParams{Query: keyword})
		if err != nil {
			continue
		}

		pipeline = append(pipeline, ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			KeywordsPerCategory: result.KeywordsPerCategory,
		})

		seachedKeywordsPerDomain[domainType] = append(seachedKeywordsPerDomain[domainType], keyword)
	}

	return pipeline
}

// convertDynamicPipelineToConnectorPipeline converts DynamicPipelineResult to ConnectorPipeline format
func convertDynamicPipelineToConnectorPipeline(dynamicResult *domain.DynamicPipelineResult) []ConnectorPipeline {
	pipeline := make([]ConnectorPipeline, 0, len(dynamicResult.Steps))

	for _, step := range dynamicResult.Steps {
		if step.DomainType != "SUMMARY" {
			pipeline = append(pipeline, ConnectorPipeline{
				Success:             step.Success,
				Error:               step.Error,
				Name:                string(step.DomainType),
				SearchParameter:     step.SearchParameter,
				Output:              step.Output,
				KeywordsPerCategory: step.KeywordsPerCategory,
			})
		}
	}

	return pipeline
}

const Seconds = 2

// searchHandler demonstrates how to use the new domain search function
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query().Get("q")
	domainType := r.URL.Query().Get("domain")

	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	searchParams := domain.DomainSearchParams{
		Query: query,
	}

	var result *domain.DomainSearchResult
	var err error

	// If specific domain is requested, search that domain
	if domainType != "" {
		dt, err := domain.GetDomainTypeFromString(domainType)
		if err != nil {
			// Build error message with available domain types
			availableTypes := make([]string, 0, len(domain.StringToDomainType))
			for k := range domain.StringToDomainType {
				availableTypes = append(availableTypes, k)
			}
			http.Error(w, fmt.Sprintf("Invalid domain type. Use: %v", availableTypes), http.StatusBadRequest)
			return
		}
		result, err = module.SearchDomain(dt, searchParams)

		if err != nil {
			spew.Dump("error", err)
			http.Error(w, "Search failed", http.StatusInternalServerError)
			return
		}

		// Convert to ConnectorPipeline format
		pipeline := ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			KeywordsPerCategory: result.KeywordsPerCategory,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pipeline)
		return
	}

	// If no specific domain, search default domains
	domainTypes := domain.DefaultDomainTypes()

	results := module.SearchMultipleDomains(domainTypes, searchParams)

	// Convert to ConnectorPipeline format
	pipeline := make([]ConnectorPipeline, 0, len(results))
	for _, result := range results {
		pipeline = append(pipeline, ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			KeywordsPerCategory: result.KeywordsPerCategory,
		})
	}

	response := map[string]interface{}{
		"data":    pipeline,
		"success": true,
		"message": "Search completed successfully",
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// dynamicPipelineHandler demonstrates the new dynamic pipeline functionality
func (s *Server) dynamicPipelineHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		query = "Novasco" // Default query
	}
	// User the dymanic interactor

	// Get configuration parameters
	maxDepth := 53
	if depth := r.URL.Query().Get("depth"); depth != "" {
		if d, err := strconv.Atoi(depth); err == nil && d > 0 && d <= 10 {
			maxDepth = d
		}
	}

	skipDuplicates := true
	if skip := r.URL.Query().Get("skip_duplicates"); skip == "false" {
		skipDuplicates = false
	}

	executionID := r.URL.Query().Get("execution_id")
	if executionID == "" {
		executionID = domain.NewID().String()
	}

	// Check if streaming is requested
	stream := r.URL.Query().Get("stream") == "true"
	if stream {
		s.dynamicPipelineStreamHandler(w, r, query, maxDepth, skipDuplicates)
		return
	}

	spew.Dump("before set execution id")
	ctx := infra.SetExecutionID(context.Background(), executionID)
	spew.Dump("ctx.Err() == context.Canceled", ctx.Err() == context.Canceled)
	spew.Dump("ctx.Err() == context.DeadlineExceeded", ctx.Err() == context.DeadlineExceeded)
	spew.Dump("ctx.Err()", ctx.Err())
	// Start pipeline execution in the background
	go func() {
		log.Printf("[%s] Starting background pipeline execution with query: %s, max depth: %d, skip duplicates: %v",
			executionID, query, maxDepth, skipDuplicates)

		spew.Dump("before execute dynamic pipeline")
		spew.Dump(ctx.Err() == context.Canceled)
		spew.Dump(ctx.Err() == context.DeadlineExceeded)
		spew.Dump(ctx.Err())

		_, err := s.interactor.ExecuteDynamicPipeline(ctx, query, maxDepth, skipDuplicates)
		if err != nil {
			log.Printf("[%s] Background pipeline execution failed: %v", executionID, err)
		} else {
			log.Printf("[%s] Background pipeline execution completed successfully", executionID)
		}
	}()

	spew.Dump("go func")
	spew.Dump("ctx.Err() == context.Canceled", ctx.Err() == context.Canceled)
	spew.Dump("ctx.Err() == context.DeadlineExceeded", ctx.Err() == context.DeadlineExceeded)
	spew.Dump("ctx.Err()", ctx.Err())
	// _, err := s.interactor.ExecuteDynamicPipeline(ctx, query, maxDepth, skipDuplicates)
	// if err != nil {
	// 	http.Error(w, "Failed to execute dynamic pipeline", http.StatusInternalServerError)
	// 	return
	// }

	// Return immediately with executionID
	response := struct {
		ExecutionID string `json:"execution_id"`
		Message     string `json:"message"`
		Status      string `json:"status"`
	}{
		ExecutionID: executionID,
		Message:     "Pipeline execution started in background",
		Status:      "processing",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// dynamicPipelineStreamHandler handles streaming pipeline results
func (s *Server) dynamicPipelineStreamHandler(w http.ResponseWriter, r *http.Request, query string, maxDepth int, skipDuplicates bool) {
	// Set headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Create a channel to receive pipeline steps
	stepChan := make(chan domain.DynamicPipelineStep, 100)
	done := make(chan bool)

	// Configure the dynamic pipeline
	config := domain.DynamicPipelineConfig{
		MaxDepth:           maxDepth,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     skipDuplicates,
	}

	// Available domains
	availableDomains := domain.AllDomainTypes()

	// Start pipeline execution in a goroutine
	go func() {
		defer close(stepChan)
		defer close(done)

		// Execute the dynamic pipeline with step callback
		dynamicResult, err := s.executeDynamicPipelineWithCallback(r.Context(), query, availableDomains, config, stepChan)
		if err != nil {
			// Send error as a step
			errorStep := domain.DynamicPipelineStep{
				DomainType:      "ERROR",
				SearchParameter: query,
				Success:         false,
				Error:           err,
				Output:          nil,
				Depth:           0,
			}
			stepChan <- errorStep
			return
		}

		// Send final summary
		summaryStep := domain.DynamicPipelineStep{
			DomainType:      "SUMMARY",
			SearchParameter: query,
			Success:         true,
			Error:           nil,
			Output: map[string]interface{}{
				"total_steps":       dynamicResult.TotalSteps,
				"successful_steps":  dynamicResult.SuccessfulSteps,
				"failed_steps":      dynamicResult.FailedSteps,
				"max_depth_reached": dynamicResult.MaxDepthReached,
			},
			Depth: dynamicResult.MaxDepthReached,
		}
		stepChan <- summaryStep
	}()
	// Flush the response to ensure immediate delivery
	flusher, ok := w.(http.Flusher)
	if ok {
		flusher.Flush()
	}
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Stream the steps as they come
	stepCount := 0
	for {
		select {
		case step, ok := <-stepChan:
			if !ok {
				// Channel closed, send completion event
				s.writeSSEEvent(w, "complete", map[string]interface{}{
					"message":     "Pipeline execution completed",
					"total_steps": stepCount,
				}, flusher)
				return
			}

			stepCount++

			// Convert step to ConnectorPipeline format
			pipelineStep := ConnectorPipeline{
				Success:             step.Success,
				Error:               step.Error,
				Name:                string(step.DomainType),
				SearchParameter:     step.SearchParameter,
				Output:              step.Output,
				KeywordsPerCategory: step.KeywordsPerCategory,
			}

			// Send step as SSE event
			eventData := map[string]interface{}{
				"step_number": stepCount,
				"step":        pipelineStep,
				"depth":       step.Depth,
				"category":    string(step.Category),
				"keywords":    step.Keywords,
			}

			eventType := "step"
			switch step.DomainType {
			case "error":
				eventType = "error"
			case "SUMMARY":
				eventType = "sumary"
			}

			s.writeSSEEvent(w, eventType, eventData, flusher)

		case <-r.Context().Done():
			// Client disconnected
			return
		}
	}
}

// executeDynamicPipelineWithCallback executes the dynamic pipeline and sends steps to a channel
func (s *Server) executeDynamicPipelineWithCallback(ctx context.Context, query string, availableDomains []domain.DomainType, config domain.DynamicPipelineConfig, stepChan chan<- domain.DynamicPipelineStep) (*domain.DynamicPipelineResult, error) {
	// Create a custom pipeline executor that streams steps
	return s.executeStreamingPipeline(ctx, query, availableDomains, config, stepChan)
}

// executeStreamingPipeline executes the pipeline with real-time streaming
func (s *Server) executeStreamingPipeline(ctx context.Context, query string, availableDomains []domain.DomainType, config domain.DynamicPipelineConfig, stepChan chan<- domain.DynamicPipelineStep) (*domain.DynamicPipelineResult, error) {
	// Create the initial pipeline steps
	initialResult, err := module.CreateDynamicPipeline(ctx, query, availableDomains, config)
	if err != nil {
		return nil, err
	}

	// Get initial steps from the result
	initialSteps := initialResult.Steps

	totalSteps := 0
	successfulSteps := 0
	failedSteps := 0
	maxDepthReached := 0

	// Track searched keywords per domain to avoid duplicates
	searchedKeywordsPerDomain := make(map[domain.DomainType]map[string]bool)
	for _, domainType := range availableDomains {
		searchedKeywordsPerDomain[domainType] = make(map[string]bool)
	}

	// Process steps with streaming
	processedSteps := make([]domain.DynamicPipelineStep, 0)

	// Create a queue for steps to process
	stepQueue := make([]domain.DynamicPipelineStep, len(initialSteps))
	copy(stepQueue, initialSteps)

	for len(stepQueue) > 0 {
		// Get next step from queue
		step := stepQueue[0]
		stepQueue = stepQueue[1:]

		// Send step start event
		startStep := step
		startStep.Success = false
		startStep.Output = nil
		stepChan <- startStep

		// Execute the step
		result, err := module.SearchDomain(step.DomainType, domain.DomainSearchParams{Query: step.SearchParameter})

		// Update step with results
		step.Success = err == nil
		step.Error = err
		if result != nil {
			step.Output = result.Output
			step.KeywordsPerCategory = result.KeywordsPerCategory
		}

		// Update counters
		totalSteps++
		if step.Success {
			successfulSteps++
		} else {
			failedSteps++
		}

		if step.Depth > maxDepthReached {
			maxDepthReached = step.Depth
		}

		// Send completed step
		stepChan <- step
		processedSteps = append(processedSteps, step)

		// Add delay between steps for better streaming experience
		time.Sleep(time.Duration(config.DelayBetweenSteps) * time.Second)

		// Generate new steps from keywords if not at max depth
		if step.Depth < config.MaxDepth && step.Success && step.Output != nil {
			newSteps := s.generateNextSteps(step, availableDomains, searchedKeywordsPerDomain, config)
			stepQueue = append(stepQueue, newSteps...)
		}
	}

	// Create final result
	dynamicResult := &domain.DynamicPipelineResult{
		Steps:           processedSteps,
		TotalSteps:      totalSteps,
		SuccessfulSteps: successfulSteps,
		FailedSteps:     failedSteps,
		MaxDepthReached: maxDepthReached,
		Config:          config,
	}

	return dynamicResult, nil
}

// generateNextSteps generates new pipeline steps from a completed step
func (s *Server) generateNextSteps(completedStep domain.DynamicPipelineStep, availableDomains []domain.DomainType, searchedKeywordsPerDomain map[domain.DomainType]map[string]bool, config domain.DynamicPipelineConfig) []domain.DynamicPipelineStep {
	var newSteps []domain.DynamicPipelineStep

	// Extract keywords from the completed step
	keywordsPerCategory := completedStep.KeywordsPerCategory
	if keywordsPerCategory == nil {
		return newSteps
	}

	// Generate new steps for each keyword category
	for category, keywords := range keywordsPerCategory {
		for _, keyword := range keywords {
			// Skip if already searched or if keyword is too short
			if len(keyword) < 3 {
				continue
			}

			// Generate steps for each available domain
			for _, domainType := range availableDomains {
				// Skip if already searched this keyword for this domain
				if searchedKeywordsPerDomain[domainType][keyword] {
					continue
				}

				// Skip if same domain as current step
				if domainType == completedStep.DomainType {
					continue
				}

				// Mark as searched
				searchedKeywordsPerDomain[domainType][keyword] = true

				// Create new step
				newStep := domain.DynamicPipelineStep{
					DomainType:          domainType,
					SearchParameter:     keyword,
					Category:            category,
					Keywords:            []string{keyword},
					Success:             false,
					Error:               nil,
					Output:              nil,
					KeywordsPerCategory: nil,
					Depth:               completedStep.Depth + 1,
				}

				newSteps = append(newSteps, newStep)
			}
		}
	}

	return newSteps
}

// writeSSEEvent writes a Server-Sent Event
func (s *Server) writeSSEEvent(w http.ResponseWriter, eventType string, data interface{}, flusher http.Flusher) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling SSE data: %v", err)
		return
	}

	// Write SSE format: event: type\ndata: json\n\n
	fmt.Fprintf(w, "event: %s\n", eventType)
	fmt.Fprintf(w, "data: %s\n\n", string(jsonData))

	flusher.Flush()
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// Google Docking handlers

func (s *Server) googleDorkingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Parse optional parameters
	maxResults := 10
	if mr := r.URL.Query().Get("max_results"); mr != "" {
		if parsed, err := strconv.Atoi(mr); err == nil && parsed > 0 {
			maxResults = parsed
		}
	}

	minRelevance := 0.1
	if mr := r.URL.Query().Get("min_relevance"); mr != "" {
		if parsed, err := strconv.ParseFloat(mr, 64); err == nil && parsed >= 0 && parsed <= 1 {
			minRelevance = parsed
		}
	}

	exactMatch := r.URL.Query().Get("exact_match") == "true"
	caseSensitive := r.URL.Query().Get("case_sensitive") == "true"

	// Parse include/exclude keywords
	var includeKeywords, excludeKeywords []string
	if ik := r.URL.Query().Get("include_keywords"); ik != "" {
		includeKeywords = strings.Split(ik, ",")
	}
	if ek := r.URL.Query().Get("exclude_keywords"); ek != "" {
		excludeKeywords = strings.Split(ek, ",")
	}

	// Create Google Docking connector
	googleDorking := module.NewGoogleDorkingDomain()

	// Create search parameters
	params := domain.GoogleDorkingSearchParams{
		Query:           query,
		MaxResults:      maxResults,
		MinRelevance:    minRelevance,
		ExactMatch:      exactMatch,
		CaseSensitive:   caseSensitive,
		IncludeKeywords: includeKeywords,
		ExcludeKeywords: excludeKeywords,
	}

	// Perform search
	results, err := googleDorking.SearchWithParams(params)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get statistics
	stats := googleDorking.GetSearchStatistics(results)

	// Create response
	response := map[string]interface{}{
		"success":    true,
		"query":      query,
		"results":    results,
		"statistics": stats,
		"parameters": map[string]interface{}{
			"max_results":      maxResults,
			"min_relevance":    minRelevance,
			"exact_match":      exactMatch,
			"case_sensitive":   caseSensitive,
			"include_keywords": includeKeywords,
			"exclude_keywords": excludeKeywords,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) googleDorkingSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	googleDorking := module.NewGoogleDorkingDomain()
	suggestions, err := googleDorking.GetSearchSuggestions(query)
	if err != nil {
		http.Error(w, "Failed to get suggestions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"query":       query,
		"suggestions": suggestions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) googleDorkingStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Results []domain.GoogleDorkingResult `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	googleDorking := module.NewGoogleDorkingDomain()
	stats := googleDorking.GetSearchStatistics(request.Results)

	response := map[string]interface{}{
		"success":       true,
		"statistics":    stats,
		"total_results": len(request.Results),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
