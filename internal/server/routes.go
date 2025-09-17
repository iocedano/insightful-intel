package server

import (
	"encoding/json"
	"insightful-intel/internal/domain"
	"log"
	"net/http"
	"slices"

	"github.com/davecgh/go-spew/spew"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.HelloWorldHandler)

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

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	pipeline := []ConnectorPipeline{}

	onapi := domain.NewOnapiDomain()
	scj := domain.NewScjDomain()
	scjSearchableCategory := scj.GetListOfSearchableCategory()
	dgii := domain.NewDgiiDomain()
	dgiiSearchableCategory := scj.GetListOfSearchableCategory()
	// pgr := domain.NewPgrDomain()
	// pgrSearchableCategory := scj.GetListOfSearchableCategory()

	onapiResp, err := onapi.SearchComercialName("Novasco")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	startPoint := ConnectorPipeline{
		Success:             true,
		Error:               err,
		Name:                onapi.GetName(),
		SearchParameter:     "Novasco",
		Output:              onapiResp,
		keywordsPerCategory: domain.GetCategoryByKeywords(&domain.Onapi{}, onapiResp),
	}

	pipeline = append(pipeline, startPoint)

	nextStep := 0

	for nextStep <= len(pipeline) {
		collector := pipeline[nextStep]

		// Add condition to end the cycle

		for category, keywords := range collector.keywordsPerCategory {
			if slices.Contains(scjSearchableCategory, category) {
				for _, keyword := range keywords {

					scjResp, err := scj.Search(keyword)
					if err != nil {
						http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
						continue
					}

					pipeline = append(pipeline, ConnectorPipeline{
						Success:             err != nil,
						Error:               err,
						Name:                scj.GetName(),
						SearchParameter:     keyword,
						Output:              scjResp,
						keywordsPerCategory: domain.GetCategoryByKeywords(&domain.Scj{}, scjResp),
					})

					spew.Dump(scj.GetName(), keyword, scjResp)

				}
			}

			if slices.Contains(dgiiSearchableCategory, category) {
				for _, keyword := range keywords {

					dgiiResp, err := dgii.GetRegister(keyword)
					if err != nil {
						http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
						continue
					}

					pipeline = append(pipeline, ConnectorPipeline{
						Success:             err != nil,
						Error:               err,
						Name:                dgii.GetName(),
						SearchParameter:     keyword,
						Output:              dgiiResp,
						keywordsPerCategory: domain.GetCategoryByKeywords(&domain.Dgii{}, dgiiResp),
					})
					spew.Dump(dgii.GetName(), keyword, dgiiResp)

				}
			}

		}

		nextStep++
	}

	spew.Dump(pipeline)

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
