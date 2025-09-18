package server

import (
	"encoding/json"
	"insightful-intel/internal/domain"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/search", s.searchHandler)
	mux.HandleFunc("/health", s.healthHandler)

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
	Success             bool
	Error               error
	Name                string
	SearchParameter     string
	keywordsPerCategory map[domain.DataCategory][]string
	Output              any
}

func Step[T any](
	domainConnector domain.GenericConnector[T],
	searchableCategory []domain.DataCategory,
	category domain.DataCategory,
	keywords []string,
	seachedKeywordsPerDomain map[domain.DomainType][]string,
) []ConnectorPipeline {
	pipeline := []ConnectorPipeline{}

	// Check if the given category is searchable by the provided domain connector
	searchableCategories := domainConnector.GetListOfSearchableCategory()
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

		result, err := domain.SearchDomain(domainType, domain.DomainSearchParams{Query: keyword})
		if err != nil {
			continue
		}

		pipeline = append(pipeline, ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			keywordsPerCategory: result.KeywordsPerCategory,
		})

		seachedKeywordsPerDomain[domainType] = append(seachedKeywordsPerDomain[domainType], keyword)
	}

	return pipeline
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	searchParams := domain.DomainSearchParams{
		Query: query,
	}

	// // Convert to the existing ConnectorPipeline format for compatibility
	pipeline := []ConnectorPipeline{}

	// Example 2: Search multiple domains at once
	domainTypes := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	multiResults := domain.SearchMultipleDomains(domainTypes, searchParams)
	seachedKeywordsPerDomain := map[domain.DomainType][]string{}

	// Add results to pipeline
	for _, result := range multiResults {
		pipeline = append(pipeline, ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			keywordsPerCategory: result.KeywordsPerCategory,
		})

		seachedKeywordsPerDomain[result.DomainType] = append(seachedKeywordsPerDomain[result.DomainType], result.SearchParameter)
	}

	// Example 3: Dynamic pipeline based on keywords (keeping the original logic)
	scj := domain.NewScjDomain()
	scjSearchableCategory := scj.GetListOfSearchableCategory()
	dgii := domain.NewDgiiDomain()
	dgiiSearchableCategory := dgii.GetListOfSearchableCategory()

	// Use goroutines and channels to parallelize Step calls for SCJ and DGII
	nextStep := 0
	type stepResult struct {
		results []ConnectorPipeline
	}
	for nextStep < len(pipeline) {
		collector := pipeline[nextStep]

		type stepCall struct {
			connector          any
			searchableCategory []domain.DataCategory
			category           domain.DataCategory
			keywords           []string
		}

		var calls []stepCall
		for category, keywords := range collector.keywordsPerCategory {
			calls = append(calls, stepCall{&scj, scjSearchableCategory, category, keywords})
			calls = append(calls, stepCall{&dgii, dgiiSearchableCategory, category, keywords})
		}

		resultsCh := make(chan []ConnectorPipeline, len(calls))
		doneCh := make(chan struct{})
		var wg sync.WaitGroup

		for _, call := range calls {
			wg.Add(1)
			go func(call stepCall) {
				defer wg.Done()
				switch c := call.connector.(type) {
				case *domain.Scj:
					resultsCh <- Step(c, call.searchableCategory, call.category, call.keywords, seachedKeywordsPerDomain)
				case *domain.Dgii:
					resultsCh <- Step(c, call.searchableCategory, call.category, call.keywords, seachedKeywordsPerDomain)
				default:
					resultsCh <- nil
				}
			}(call)
		}

		// Wait for all goroutines to finish, then close the results channel
		go func() {
			wg.Wait()
			close(resultsCh)
			close(doneCh)
		}()

		// Collect all results from the channel
		for p := range resultsCh {
			if len(p) > 0 {
				pipeline = append(pipeline, p...)
			}
		}
		<-doneCh // Ensure all goroutines are finished

		nextStep++
	}

	spew.Dump("-----Finish----")

	jsonResp, err := json.Marshal(pipeline)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

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
		switch domainType {
		case "onapi":
			result, err = domain.SearchDomain(domain.DomainTypeONAPI, searchParams)
		case "scj":
			result, err = domain.SearchDomain(domain.DomainTypeSCJ, searchParams)
		case "dgii":
			result, err = domain.SearchDomain(domain.DomainTypeDGII, searchParams)
		default:
			http.Error(w, "Invalid domain type. Use: onapi, scj, or dgii", http.StatusBadRequest)
			return
		}

		if err != nil {
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
			keywordsPerCategory: result.KeywordsPerCategory,
		}

		jsonResp, err := json.Marshal(pipeline)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
		return
	}

	// If no specific domain, search all domains
	domainTypes := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
	}

	results := domain.SearchMultipleDomains(domainTypes, searchParams)

	// Convert to ConnectorPipeline format
	pipeline := make([]ConnectorPipeline, 0, len(results))
	for _, result := range results {
		pipeline = append(pipeline, ConnectorPipeline{
			Success:             result.Success,
			Error:               result.Error,
			Name:                string(result.DomainType),
			SearchParameter:     result.SearchParameter,
			Output:              result.Output,
			keywordsPerCategory: result.KeywordsPerCategory,
		})
	}

	jsonResp, err := json.Marshal(pipeline)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
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
