package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// DgiiRepository implements DomainRepository for DGII Register domain type
type DgiiRepository struct {
	db DatabaseAccessor
}

// NewDgiiRepository creates a new DGII repository instance
func NewDgiiRepository(db database.Service) *DgiiRepository {
	return &DgiiRepository{
		db: NewDatabaseAdapter(db),
	}
}

// Create inserts a new DGII register
func (r *DgiiRepository) Create(ctx context.Context, entity domain.Register) error {
	return r.CreateWithDomainSearchResultID(ctx, entity)
}

// CreateWithDomainSearchResultID inserts a new DGII register with a domain search result ID
func (r *DgiiRepository) CreateWithDomainSearchResultID(ctx context.Context, entity domain.Register) error {
	entity.ID = domain.NewID()

	query := `
		INSERT INTO dgii_registers (
			id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
			facturador_electronico, licencia_comercial, estado, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.DomainSearchResultID, entity.RNC, entity.RazonSocial, entity.NombreComercial, entity.Categoria,
		entity.RegimenPagos, entity.FacturadorElectronico, entity.LicenciaComercial, entity.Estado,
	)

	return err
}

// GetByID retrieves a DGII register by its RNC
func (r *DgiiRepository) GetByID(ctx context.Context, id string) (domain.Register, error) {
	query := `
		SELECT id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
			   facturador_electronico, licencia_comercial, estado, created_at, updated_at
		FROM dgii_registers 
		WHERE id = ?
	`

	var entity domain.Register

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID, &entity.DomainSearchResultID, &entity.RNC, &entity.RazonSocial, &entity.NombreComercial, &entity.Categoria,
		&entity.RegimenPagos, &entity.FacturadorElectronico, &entity.LicenciaComercial, &entity.Estado,
	)

	if err != nil {
		return domain.Register{}, err
	}

	return entity, nil
}

// Update modifies an existing DGII register
func (r *DgiiRepository) Update(ctx context.Context, id string, entity domain.Register) error {
	query := `
		UPDATE dgii_registers SET
			 domain_search_result_id = ?, razon_social = ?, nombre_comercial = ?, categoria = ?, regimen_pagos = ?,
			facturador_electronico = ?, licencia_comercial = ?, estado = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.DomainSearchResultID, entity.RazonSocial, entity.NombreComercial, entity.Categoria, entity.RegimenPagos,
		entity.FacturadorElectronico, entity.LicenciaComercial, entity.Estado, id,
	)

	return err
}

// Delete removes a DGII register by its RNC
func (r *DgiiRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM dgii_registers WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves multiple DGII registers with pagination
func (r *DgiiRepository) List(ctx context.Context, offset, limit int) ([]domain.Register, error) {
	query := `
		SELECT id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
			   facturador_electronico, licencia_comercial, estado, created_at, updated_at
		FROM dgii_registers 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.Register
	for rows.Next() {
		var entity domain.Register

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.RNC, &entity.RazonSocial, &entity.NombreComercial, &entity.Categoria,
			&entity.RegimenPagos, &entity.FacturadorElectronico, &entity.LicenciaComercial, &entity.Estado,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// Count returns the total number of DGII registers
func (r *DgiiRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM dgii_registers`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Search performs a search query on DGII registers
func (r *DgiiRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.Register, error) {
	searchQuery := `
		SELECT id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
			   facturador_electronico, licencia_comercial, estado, created_at, updated_at
		FROM dgii_registers 
		WHERE id LIKE ? OR domain_search_result_id LIKE ? OR rnc LIKE ? OR razon_social LIKE ? OR nombre_comercial LIKE ? 
		   OR categoria LIKE ? OR estado LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery,
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.Register
	for rows.Next() {
		var entity domain.Register

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.RNC, &entity.RazonSocial, &entity.NombreComercial, &entity.Categoria,
			&entity.RegimenPagos, &entity.FacturadorElectronico, &entity.LicenciaComercial, &entity.Estado,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// SearchByCategory performs a search within a specific keyword category
func (r *DgiiRepository) SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]domain.Register, error) {
	var searchQuery string
	searchPattern := "%" + query + "%"

	switch category {
	case domain.KeywordCategoryCompanyName:
		searchQuery = `
			SELECT id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
				   facturador_electronico, licencia_comercial, estado, created_at, updated_at
			FROM dgii_registers 
			WHERE id LIKE ? OR domain_search_result_id LIKE ? OR razon_social LIKE ? OR nombre_comercial LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryContributorID:
		searchQuery = `
			SELECT id, domain_search_result_id, rnc, razon_social, nombre_comercial, categoria, regimen_pagos,
				   facturador_electronico, licencia_comercial, estado, created_at, updated_at
			FROM dgii_registers 
			WHERE id LIKE ? OR domain_search_result_id LIKE ? OR rnc LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	default:
		return []domain.Register{}, fmt.Errorf("unsupported category: %s", category)
	}

	var rows *sql.Rows
	var err error

	if category == domain.KeywordCategoryCompanyName {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, searchPattern, searchPattern, searchPattern, limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, searchPattern, searchPattern, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.Register
	for rows.Next() {
		var entity domain.Register

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.RNC, &entity.RazonSocial, &entity.NombreComercial, &entity.Categoria,
			&entity.RegimenPagos, &entity.FacturadorElectronico, &entity.LicenciaComercial, &entity.Estado,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// GetByDomainType retrieves DGII registers by domain type
func (r *DgiiRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]domain.Register, error) {
	// For DGII, all registers are of the same domain type, so we just return all
	return r.List(ctx, offset, limit)
}

// GetBySearchParameter retrieves DGII registers by search parameter
func (r *DgiiRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]domain.Register, error) {
	return r.Search(ctx, searchParam, offset, limit)
}

// GetKeywordsByCategory retrieves keywords grouped by category for a DGII register
func (r *DgiiRepository) GetKeywordsByCategory(ctx context.Context, entityID string) (map[domain.KeywordCategory][]string, error) {
	entity, err := r.GetByID(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// DGII registers don't implement GenericConnector, so we manually extract keywords
	return map[domain.KeywordCategory][]string{
		domain.KeywordCategoryCompanyName:   []string{entity.RazonSocial, entity.NombreComercial},
		domain.KeywordCategoryContributorID: []string{entity.RNC},
	}, nil
}
