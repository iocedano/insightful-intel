package repositories

import (
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// DomainRepositoryUnion represents any of the domain repositories
type DomainRepositoryUnion interface {
	*OnapiRepository | *ScjRepository | *DgiiRepository | *PgrRepository | *DockingRepository
}

// DomainRepositoryHandler provides a way to work with any domain repository
type DomainRepositoryHandler struct {
	GetOnapiRepository   func() *OnapiRepository
	GetScjRepository     func() *ScjRepository
	GetDgiiRepository    func() *DgiiRepository
	GetPgrRepository     func() *PgrRepository
	GetDockingRepository func() *DockingRepository
}

func (h *DomainRepositoryHandler) GetDomainTypeRepo(domainType domain.DomainType) any {
	switch domainType {
	case domain.DomainTypeONAPI:
		return h.GetOnapiRepository()
	case domain.DomainTypeSCJ:
		return h.GetScjRepository()
	case domain.DomainTypeDGII:
		return h.GetDgiiRepository()
	case domain.DomainTypePGR:
		return h.GetPgrRepository()
	case domain.DomainTypeGoogleDocking:
		return h.GetDockingRepository()
	}
	return nil
}

// RepositoryFactory provides a centralized way to create repository instances
type RepositoryFactory struct {
	db database.Service
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db database.Service) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// GetOnapiRepository returns an ONAPI repository instance
func (f *RepositoryFactory) GetOnapiRepository() *OnapiRepository {
	return NewOnapiRepository(f.db)
}

// GetScjRepository returns an SCJ repository instance
func (f *RepositoryFactory) GetScjRepository() *ScjRepository {
	return NewScjRepository(f.db)
}

// GetDgiiRepository returns a DGII repository instance
func (f *RepositoryFactory) GetDgiiRepository() *DgiiRepository {
	return NewDgiiRepository(f.db)
}

// GetPgrRepository returns a PGR repository instance
func (f *RepositoryFactory) GetPgrRepository() *PgrRepository {
	return NewPgrRepository(f.db)
}

// GetDockingRepository returns a Google Docking repository instance
func (f *RepositoryFactory) GetDockingRepository() *DockingRepository {
	return NewDockingRepository(f.db)
}

// GetPipelineRepository returns a pipeline repository instance
func (f *RepositoryFactory) GetPipelineRepository() *PipelineRepository {
	return NewPipelineRepository(f.db)
}

// GetAllDomainRepositories returns all domain repositories
func (f *RepositoryFactory) GetAllDomainRepositories() *DomainRepositoryHandler {
	return &DomainRepositoryHandler{
		GetOnapiRepository:   f.GetOnapiRepository,
		GetScjRepository:     f.GetScjRepository,
		GetDgiiRepository:    f.GetDgiiRepository,
		GetPgrRepository:     f.GetPgrRepository,
		GetDockingRepository: f.GetDockingRepository,
	}
}

// GetRepositoryByDomainType returns a repository for a specific domain type
func (f *RepositoryFactory) GetRepositoryByDomainType(domainType domain.DomainType) any {
	return f.GetAllDomainRepositories().GetDomainTypeRepo(domainType)
}
