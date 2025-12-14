package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"insightful-intel/internal/domain"
)

// onapiHandler handles ONAPI repository operations
func (s *Server) onapiHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	onapiRepo := repos.GetOnapiRepository()

	switch r.Method {
	case http.MethodGet:
		// List ONAPI entities with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		entities, err := onapiRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list entities: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    entities,
			"count":   len(entities),
		})

	case http.MethodPost:
		// Create new ONAPI entity
		var entity domain.Entity
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if err := onapiRepo.Create(r.Context(), entity); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create entity: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Entity created successfully",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// scjHandler handles SCJ repository operations
func (s *Server) scjHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	scjRepo := repos.GetScjRepository()

	switch r.Method {
	case http.MethodGet:
		// List SCJ cases with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		cases, err := scjRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list cases: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    cases,
			"count":   len(cases),
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// dgiiHandler handles DGII repository operations
func (s *Server) dgiiHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	dgiiRepo := repos.GetDgiiRepository()

	switch r.Method {
	case http.MethodGet:
		// List DGII registers with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		registers, err := dgiiRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list registers: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    registers,
			"count":   len(registers),
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// pgrHandler handles PGR repository operations
func (s *Server) pgrHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	pgrRepo := repos.GetPgrRepository()

	switch r.Method {
	case http.MethodGet:
		// List PGR news with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		news, err := pgrRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list news: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    news,
			"count":   len(news),
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// dockingHandler handles Google Docking repository operations
func (s *Server) dockingHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	dockingRepo := repos.GetDockingRepository()

	switch r.Method {
	case http.MethodGet:
		// List Google Docking results with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		results, err := dockingRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list results: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    results,
			"count":   len(results),
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// pipelineHandler handles Pipeline repository operations
func (s *Server) pipelineHandler(w http.ResponseWriter, r *http.Request) {
	repos := s.GetRepositories()
	pipelineRepo := repos.GetPipelineRepository()

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			dynamicResult, err := pipelineRepo.GetPipelineByID(r.Context(), id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get pipeline result: %v", err), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"data":    dynamicResult,
			})
			return
		}

		// List pipeline results with pagination
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}

		results, err := pipelineRepo.List(r.Context(), offset, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list pipeline results: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    results,
			"count":   len(results),
		})

	case http.MethodPost:
		// Save pipeline result to database
		var pipelineData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&pipelineData); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Try to parse as DynamicPipelineResult first
		if dynamicResult, ok := s.parseDynamicPipelineResult(pipelineData); ok {
			_, err := pipelineRepo.CreateDynamicPipelineResult(r.Context(), dynamicResult)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to save dynamic pipeline result: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Dynamic pipeline result saved successfully",
				"type":    "dynamic_pipeline",
			})
			return
		}

		// Try to parse as DomainSearchResult
		if domainResult, ok := s.parseDomainSearchResult(pipelineData); ok {
			_, err := pipelineRepo.CreateDomainSearchResult(r.Context(), domainResult)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to save domain search result: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Domain search result saved successfully",
				"type":    "domain_search",
			})
			return
		}

		http.Error(w, "Unsupported pipeline result format", http.StatusBadRequest)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// parseDynamicPipelineResult attempts to parse the input data as a DynamicPipelineResult
func (s *Server) parseDynamicPipelineResult(data map[string]interface{}) (*domain.DynamicPipelineResult, bool) {
	// Check if this looks like a dynamic pipeline result
	if _, hasSteps := data["steps"]; !hasSteps {
		return nil, false
	}

	// Convert to JSON and back to ensure proper type conversion
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var result domain.DynamicPipelineResult
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, false
	}

	return &result, true
}

// parseDomainSearchResult attempts to parse the input data as a DomainSearchResult
func (s *Server) parseDomainSearchResult(data map[string]interface{}) (*domain.DomainSearchResult, bool) {
	// Check if this looks like a domain search result
	if _, hasDomainType := data["domain_type"]; !hasDomainType {
		return nil, false
	}

	// Convert to JSON and back to ensure proper type conversion
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var result domain.DomainSearchResult
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, false
	}

	return &result, true
}

// savePipelineHandler handles saving pipeline responses to the database
func (s *Server) savePipelineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repos := s.GetRepositories()
	pipelineRepo := repos.GetPipelineRepository()

	// Parse the request body
	var pipelineData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pipelineData); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Try to parse as DynamicPipelineResult first
	if dynamicResult, ok := s.parseDynamicPipelineResult(pipelineData); ok {
		_, err := pipelineRepo.CreateDynamicPipelineResult(r.Context(), dynamicResult)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save dynamic pipeline result: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Dynamic pipeline result saved successfully",
			"type":    "dynamic_pipeline",
		})
		return
	}

	// Try to parse as DomainSearchResult
	if domainResult, ok := s.parseDomainSearchResult(pipelineData); ok {
		_, err := pipelineRepo.CreateDomainSearchResult(r.Context(), domainResult)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save domain search result: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Domain search result saved successfully",
			"type":    "domain_search",
		})
		return
	}

	http.Error(w, "Unsupported pipeline result format", http.StatusBadRequest)
}

// pipelineStepsHandler handles retrieving steps for a pipeline result
func (s *Server) pipelineStepsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get pipeline ID from query parameters
	pipelineID := r.URL.Query().Get("pipeline_id")
	if pipelineID == "" {
		http.Error(w, "Query parameter 'pipeline_id' is required", http.StatusBadRequest)
		return
	}

	repos := s.GetRepositories()
	pipelineRepo := repos.GetPipelineRepository()

	// Get the pipeline result by ID
	result, err := pipelineRepo.GetPipelineStepsByID(r.Context(), pipelineID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get pipeline result: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"pipeline_id": pipelineID,
		"steps":       result,
		"count":       len(result),
	})
}
