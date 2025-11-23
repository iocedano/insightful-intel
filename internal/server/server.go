package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"insightful-intel/internal/database"
	"insightful-intel/internal/interactor"
	"insightful-intel/internal/repositories"
)

type Server struct {
	Port int

	db           database.Service
	repositories *repositories.RepositoryFactory
	interactor   *interactor.DynamicPipelineInteractor
}

func NewServer(repoFactory *repositories.RepositoryFactory, interactor *interactor.DynamicPipelineInteractor) *http.Server {
	// Get port from environment variable, default to 8080 if not set or invalid
	portStr := os.Getenv("PORT")
	port := 8080 // Default port
	if portStr != "" {
		if parsedPort, err := strconv.Atoi(portStr); err == nil && parsedPort > 0 && parsedPort < 65536 {
			port = parsedPort
		} else {
			log.Printf("Warning: Invalid PORT environment variable '%s', using default port 8080", portStr)
		}
	} else {
		log.Printf("PORT environment variable not set, using default port 8080")
	}

	// Declare Server config
	srv := &Server{
		Port:         port,
		repositories: repoFactory,
		interactor:   interactor,
	}
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", srv.Port),
		Handler:     srv.RegisterRoutes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		// Increased WriteTimeout for long-running pipeline operations and streaming
		WriteTimeout: 5 * time.Minute,
	}

	log.Printf("Server configured to listen on port %d", port)

	return server
}

// GetRepositories returns the repository factory
func (s *Server) GetRepositories() *repositories.RepositoryFactory {
	return s.repositories
}
