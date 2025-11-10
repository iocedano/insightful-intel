package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// ScjRepository implements DomainRepository for SCJ ScjCase domain type
type ScjRepository struct {
	db DatabaseAccessor
}

// NewScjRepository creates a new SCJ repository instance
func NewScjRepository(db database.Service) *ScjRepository {
	return &ScjRepository{
		db: NewDatabaseAdapter(db),
	}
}

// Create inserts a new SCJ case
func (r *ScjRepository) Create(ctx context.Context, entity domain.ScjCase) error {
	return r.CreateWithDomainSearchResultID(ctx, entity)
}

// CreateWithDomainSearchResultID inserts a new SCJ case with a domain search result ID
func (r *ScjRepository) CreateWithDomainSearchResultID(ctx context.Context, entity domain.ScjCase) error {
	entity.ID = domain.NewID()
	// Check if insert columns and parameters match -- compare to ScjCase struct fields
	query := `
		INSERT INTO scj_cases 
		(id, domain_search_result_id, linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo, id_expediente, 
		no_expediente, no_sentencia, no_unico, no_interno, id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo, 
		involucrados, guid_blob, tipo_documento_adjunto, total_filas, url_blob, extension, origen, activo, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID,
		entity.DomainSearchResultID,
		entity.Linea,
		entity.AgnoCabecera,
		entity.MesCabecera,
		entity.URLCabecera,
		entity.URLCuerpo,
		entity.IDExpediente,
		entity.NoExpediente,
		entity.NoSentencia,
		entity.NoUnico,
		entity.NoInterno,
		entity.IDTribunal,
		entity.DescTribunal,
		entity.IDMateria,
		entity.DescMateria,
		entity.FechaFallo,
		entity.Involucrados,
		entity.GuidBlob,
		entity.TipoDocumentoAdjunto,
		entity.TotalFilas,
		entity.URLBlob,
		entity.Extension,
		entity.Origen,
		entity.Activo,
	)

	if err != nil {
		return fmt.Errorf("error creating scj case: %w", err)
	}

	return nil
}

// GetByID retrieves an SCJ case by its ID
func (r *ScjRepository) GetByID(ctx context.Context, id int) (domain.ScjCase, error) {
	query := `
		SELECT linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo,
			   id_expediente, no_expediente, no_sentencia, no_unico, no_interno,
			   id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo,
			   involucrados, guid_blob, tipo_documento_adjunto, total_filas,
			   url_blob, extension, origen, activo, created_at, updated_at
		FROM scj_cases 
		WHERE id_expediente = ?
	`

	var entity domain.ScjCase

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.Linea, &entity.AgnoCabecera, &entity.MesCabecera, &entity.URLCabecera, &entity.URLCuerpo,
		&entity.IDExpediente, &entity.NoExpediente, &entity.NoSentencia, &entity.NoUnico, &entity.NoInterno,
		&entity.IDTribunal, &entity.DescTribunal, &entity.IDMateria, &entity.DescMateria, &entity.FechaFallo,
		&entity.Involucrados, &entity.GuidBlob, &entity.TipoDocumentoAdjunto, &entity.TotalFilas,
		&entity.URLBlob, &entity.Extension, &entity.Origen, &entity.Activo,
	)

	if err != nil {
		return domain.ScjCase{}, err
	}

	return entity, nil
}

// Update modifies an existing SCJ case
func (r *ScjRepository) Update(ctx context.Context, idExpediente int, entity domain.ScjCase) error {
	query := `
		UPDATE scj_cases SET
			linea = ?, agno_cabecera = ?, mes_cabecera = ?, url_cabecera = ?, url_cuerpo = ?,
			no_expediente = ?, no_sentencia = ?, no_unico = ?, no_interno = ?,
			id_tribunal = ?, desc_tribunal = ?, id_materia = ?, desc_materia = ?, fecha_fallo = ?,
			involucrados = ?, guid_blob = ?, tipo_documento_adjunto = ?, total_filas = ?,
			url_blob = ?, extension = ?, origen = ?, activo = ?, updated_at = NOW()
		WHERE id_expediente = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.Linea, entity.AgnoCabecera, entity.MesCabecera, entity.URLCabecera, entity.URLCuerpo,
		entity.NoExpediente, entity.NoSentencia, entity.NoUnico, entity.NoInterno,
		entity.IDTribunal, entity.DescTribunal, entity.IDMateria, entity.DescMateria, entity.FechaFallo,
		entity.Involucrados, entity.GuidBlob, entity.TipoDocumentoAdjunto, entity.TotalFilas,
		entity.URLBlob, entity.Extension, entity.Origen, entity.Activo, idExpediente,
	)

	return err
}

// Delete removes an SCJ case by its ID
func (r *ScjRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM scj_cases WHERE id_expediente = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves multiple SCJ cases with pagination
func (r *ScjRepository) List(ctx context.Context, offset, limit int) ([]domain.ScjCase, error) {
	query := `
		SELECT id, domain_search_result_id, linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo,
			   id_expediente, no_expediente, no_sentencia, no_unico, no_interno,
			   id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo,
			   involucrados, guid_blob, tipo_documento_adjunto, total_filas,
			   url_blob, extension, origen, activo, created_at, updated_at
		FROM scj_cases 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.ScjCase
	for rows.Next() {
		var entity domain.ScjCase

		err := rows.Scan(
			&entity.Linea, &entity.AgnoCabecera, &entity.MesCabecera, &entity.URLCabecera, &entity.URLCuerpo,
			&entity.IDExpediente, &entity.NoExpediente, &entity.NoSentencia, &entity.NoUnico, &entity.NoInterno,
			&entity.IDTribunal, &entity.DescTribunal, &entity.IDMateria, &entity.DescMateria, &entity.FechaFallo,
			&entity.Involucrados, &entity.GuidBlob, &entity.TipoDocumentoAdjunto, &entity.TotalFilas,
			&entity.URLBlob, &entity.Extension, &entity.Origen, &entity.Activo,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// Count returns the total number of SCJ cases
func (r *ScjRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM scj_cases`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Search performs a search query on SCJ cases
func (r *ScjRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.ScjCase, error) {
	searchQuery := `
		SELECT id, domain_search_result_id, linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo,
			   id_expediente, no_expediente, no_sentencia, no_unico, no_interno,
			   id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo,
			   involucrados, guid_blob, tipo_documento_adjunto, total_filas,
			   url_blob, extension, origen, activo, created_at, updated_at
		FROM scj_cases 
		WHERE no_expediente LIKE ? OR no_sentencia LIKE ? OR involucrados LIKE ? 
		   OR desc_tribunal LIKE ? OR desc_materia LIKE ?
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

	var entities []domain.ScjCase
	for rows.Next() {
		var entity domain.ScjCase

		err := rows.Scan(
			&entity.Linea, &entity.AgnoCabecera, &entity.MesCabecera, &entity.URLCabecera, &entity.URLCuerpo,
			&entity.IDExpediente, &entity.NoExpediente, &entity.NoSentencia, &entity.NoUnico, &entity.NoInterno,
			&entity.IDTribunal, &entity.DescTribunal, &entity.IDMateria, &entity.DescMateria, &entity.FechaFallo,
			&entity.Involucrados, &entity.GuidBlob, &entity.TipoDocumentoAdjunto, &entity.TotalFilas,
			&entity.URLBlob, &entity.Extension, &entity.Origen, &entity.Activo,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// SearchByCategory performs a search within a specific keyword category
func (r *ScjRepository) SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]domain.ScjCase, error) {
	var searchQuery string
	searchPattern := "%" + query + "%"

	switch category {
	case domain.KeywordCategoryPersonName:
		searchQuery = `
			SELECT id, domain_search_result_id, linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo,
				   id_expediente, no_expediente, no_sentencia, no_unico, no_interno,
				   id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo,
				   involucrados, guid_blob, tipo_documento_adjunto, total_filas,
				   url_blob, extension, origen, activo, created_at, updated_at
			FROM scj_cases 
			WHERE involucrados LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryCompanyName:
		searchQuery = `
			SELECT id, domain_search_result_id, linea, agno_cabecera, mes_cabecera, url_cabecera, url_cuerpo,
				   id_expediente, no_expediente, no_sentencia, no_unico, no_interno,
				   id_tribunal, desc_tribunal, id_materia, desc_materia, fecha_fallo,
				   involucrados, guid_blob, tipo_documento_adjunto, total_filas,
				   url_blob, extension, origen, activo, created_at, updated_at
			FROM scj_cases 
			WHERE desc_tribunal LIKE ? OR desc_materia LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	default:
		return []domain.ScjCase{}, fmt.Errorf("unsupported category: %s", category)
	}

	var rows *sql.Rows
	var err error

	if category == domain.KeywordCategoryCompanyName {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, searchPattern, limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, searchQuery,
			searchPattern, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.ScjCase
	for rows.Next() {
		var entity domain.ScjCase

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.Linea, &entity.AgnoCabecera, &entity.MesCabecera, &entity.URLCabecera, &entity.URLCuerpo,
			&entity.IDExpediente, &entity.NoExpediente, &entity.NoSentencia, &entity.NoUnico, &entity.NoInterno,
			&entity.IDTribunal, &entity.DescTribunal, &entity.IDMateria, &entity.DescMateria, &entity.FechaFallo,
			&entity.Involucrados, &entity.GuidBlob, &entity.TipoDocumentoAdjunto, &entity.TotalFilas,
			&entity.URLBlob, &entity.Extension, &entity.Origen, &entity.Activo,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// GetByDomainType retrieves SCJ cases by domain type
func (r *ScjRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]domain.ScjCase, error) {
	// For SCJ, all cases are of the same domain type, so we just return all
	return r.List(ctx, offset, limit)
}

// GetBySearchParameter retrieves SCJ cases by search parameter
func (r *ScjRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]domain.ScjCase, error) {
	return r.Search(ctx, searchParam, offset, limit)
}

// GetKeywordsByCategory retrieves keywords grouped by category for an SCJ case
func (r *ScjRepository) GetKeywordsByCategory(ctx context.Context, entityID int) (map[domain.KeywordCategory][]string, error) {
	entity, err := r.GetByID(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// SCJ cases don't implement DomainConnector, so we manually extract keywords
	return map[domain.KeywordCategory][]string{
		domain.KeywordCategoryPersonName:  []string{entity.Involucrados},
		domain.KeywordCategoryCompanyName: []string{entity.DescTribunal, entity.DescMateria},
	}, nil
}
