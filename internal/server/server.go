package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"insightful-intel/internal/database"
	"insightful-intel/internal/repositories"
)

type Server struct {
	Port int

	db           database.Service
	repositories *repositories.RepositoryFactory
}

func NewServer(repoFactory *repositories.RepositoryFactory) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Declare Server config
	srv := &Server{
		Port:         port,
		repositories: repoFactory,
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.Port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// GetRepositories returns the repository factory
func (s *Server) GetRepositories() *repositories.RepositoryFactory {
	return s.repositories
}
