package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"insightful-intel/internal/database"
	"insightful-intel/internal/domain"
)

// OnapiRepository implements DomainRepository for ONAPI Entity domain type
type OnapiRepository struct {
	db DatabaseAccessor
}

// NewOnapiRepository creates a new ONAPI repository instance
func NewOnapiRepository(db database.Service) *OnapiRepository {
	return &OnapiRepository{
		db: NewDatabaseAdapter(db),
	}
}

// Create inserts a new ONAPI entity
func (r *OnapiRepository) Create(ctx context.Context, entity domain.Entity) error {
	return r.CreateWithDomainSearchResultID(ctx, entity)
}

// CreateWithDomainSearchResultID inserts a new ONAPI entity with a domain search result ID
func (r *OnapiRepository) CreateWithDomainSearchResultID(ctx context.Context, entity domain.Entity) error {
	entity.ID = domain.NewID()

	query := `
		INSERT INTO onapi_entities (
			id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
			texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
			titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	imagenesJSON, _ := json.Marshal(entity.Imagenes)
	listaClasesJSON, _ := json.Marshal(entity.ListaClases)

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.DomainSearchResultID, entity.SerieExpediente, entity.NumeroExpediente, entity.Certificado,
		entity.Tipo, entity.SubTipo, entity.Texto, entity.Clases, entity.AplicadoAProteger,
		entity.Expedicion, entity.Vencimiento, entity.EnTramite, entity.Titular,
		entity.Gestor, entity.Domicilio, entity.Status, entity.TipoSigno,
		imagenesJSON, listaClasesJSON,
	)

	return err
}

// GetByID retrieves an ONAPI entity by its ID
func (r *OnapiRepository) GetByID(ctx context.Context, id string) (domain.Entity, error) {
	query := `
		SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
			   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
			   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
			   created_at, updated_at
		FROM onapi_entities 
		WHERE id = ?
	`

	var entity domain.Entity
	var imagenesJSON, listaClasesJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entity.ID, &entity.DomainSearchResultID, &entity.SerieExpediente, &entity.NumeroExpediente, &entity.Certificado,
		&entity.Tipo, &entity.SubTipo, &entity.Texto, &entity.Clases, &entity.AplicadoAProteger,
		&entity.Expedicion, &entity.Vencimiento, &entity.EnTramite, &entity.Titular,
		&entity.Gestor, &entity.Domicilio, &entity.Status, &entity.TipoSigno,
		&imagenesJSON, &listaClasesJSON,
	)

	if err != nil {
		return domain.Entity{}, err
	}

	// Parse JSON fields
	json.Unmarshal([]byte(imagenesJSON), &entity.Imagenes)
	json.Unmarshal([]byte(listaClasesJSON), &entity.ListaClases)

	return entity, nil
}

// Update modifies an existing ONAPI entity
func (r *OnapiRepository) Update(ctx context.Context, id string, entity domain.Entity) error {
	query := `
		UPDATE onapi_entities SET
			domain_search_result_id = ?, serie_expediente = ?, numero_expediente = ?, certificado = ?, tipo = ?, subtipo = ?,
			texto = ?, clases = ?, aplicado_a_proteger = ?, expedicion = ?, vencimiento = ?, 
			en_tramite = ?, titular = ?, gestor = ?, domicilio = ?, status = ?, tipo_signo = ?,
			imagenes = ?, lista_clases = ?, updated_at = NOW()
		WHERE id = ?
	`

	imagenesJSON, _ := json.Marshal(entity.Imagenes)
	listaClasesJSON, _ := json.Marshal(entity.ListaClases)

	_, err := r.db.ExecContext(ctx, query,
		entity.DomainSearchResultID, entity.SerieExpediente, entity.NumeroExpediente, entity.Certificado,
		entity.Tipo, entity.SubTipo, entity.Texto, entity.Clases, entity.AplicadoAProteger,
		entity.Expedicion, entity.Vencimiento, entity.EnTramite, entity.Titular,
		entity.Gestor, entity.Domicilio, entity.Status, entity.TipoSigno,
		imagenesJSON, listaClasesJSON, id,
	)

	return err
}

// Delete removes an ONAPI entity by its ID
func (r *OnapiRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM onapi_entities WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves multiple ONAPI entities with pagination
func (r *OnapiRepository) List(ctx context.Context, offset, limit int) ([]domain.Entity, error) {
	query := `
		SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
			   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
			   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
			   created_at, updated_at
		FROM onapi_entities 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.Entity
	for rows.Next() {
		var entity domain.Entity
		var imagenesJSON, listaClasesJSON string

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.SerieExpediente, &entity.NumeroExpediente, &entity.Certificado,
			&entity.Tipo, &entity.SubTipo, &entity.Texto, &entity.Clases, &entity.AplicadoAProteger,
			&entity.Expedicion, &entity.Vencimiento, &entity.EnTramite, &entity.Titular,
			&entity.Gestor, &entity.Domicilio, &entity.Status, &entity.TipoSigno,
			&imagenesJSON, &listaClasesJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		json.Unmarshal([]byte(imagenesJSON), &entity.Imagenes)
		json.Unmarshal([]byte(listaClasesJSON), &entity.ListaClases)

		entities = append(entities, entity)
	}

	return entities, nil
}

// Count returns the total number of ONAPI entities
func (r *OnapiRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM onapi_entities`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// Search performs a search query on ONAPI entities
func (r *OnapiRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.Entity, error) {
	searchQuery := `
		SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
			   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
			   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
			   created_at, updated_at
		FROM onapi_entities 
		WHERE texto LIKE ? OR titular LIKE ? OR gestor LIKE ? OR domicilio LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery,
		searchPattern, searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []domain.Entity
	for rows.Next() {
		var entity domain.Entity
		var imagenesJSON, listaClasesJSON string

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.SerieExpediente, &entity.NumeroExpediente, &entity.Certificado,
			&entity.Tipo, &entity.SubTipo, &entity.Texto, &entity.Clases, &entity.AplicadoAProteger,
			&entity.Expedicion, &entity.Vencimiento, &entity.EnTramite, &entity.Titular,
			&entity.Gestor, &entity.Domicilio, &entity.Status, &entity.TipoSigno,
			&imagenesJSON, &listaClasesJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		json.Unmarshal([]byte(imagenesJSON), &entity.Imagenes)
		json.Unmarshal([]byte(listaClasesJSON), &entity.ListaClases)

		entities = append(entities, entity)
	}

	return entities, nil
}

// SearchByCategory performs a search within a specific keyword category
func (r *OnapiRepository) SearchByCategory(ctx context.Context, category domain.KeywordCategory, query string, offset, limit int) ([]domain.Entity, error) {
	var searchQuery string
	searchPattern := "%" + query + "%"

	switch category {
	case domain.KeywordCategoryCompanyName:
		searchQuery = `
			SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
				   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
				   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
				   created_at, updated_at
			FROM onapi_entities 
			WHERE texto LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryPersonName:
		searchQuery = `
			SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
				   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
				   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
				   created_at, updated_at
			FROM onapi_entities 
			WHERE titular LIKE ? OR gestor LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	case domain.KeywordCategoryAddress:
		searchQuery = `
			SELECT id, domain_search_result_id, serie_expediente, numero_expediente, certificado, tipo, subtipo,
				   texto, clases, aplicado_a_proteger, expedicion, vencimiento, en_tramite,
				   titular, gestor, domicilio, status, tipo_signo, imagenes, lista_clases,
				   created_at, updated_at
			FROM onapi_entities 
			WHERE domicilio LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
	default:
		return []domain.Entity{}, fmt.Errorf("unsupported category: %s", category)
	}

	var rows *sql.Rows
	var err error

	if category == domain.KeywordCategoryPersonName {
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

	var entities []domain.Entity
	for rows.Next() {
		var entity domain.Entity
		var imagenesJSON, listaClasesJSON string

		err := rows.Scan(
			&entity.ID, &entity.DomainSearchResultID, &entity.SerieExpediente, &entity.NumeroExpediente, &entity.Certificado,
			&entity.Tipo, &entity.SubTipo, &entity.Texto, &entity.Clases, &entity.AplicadoAProteger,
			&entity.Expedicion, &entity.Vencimiento, &entity.EnTramite, &entity.Titular,
			&entity.Gestor, &entity.Domicilio, &entity.Status, &entity.TipoSigno,
			&imagenesJSON, &listaClasesJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		json.Unmarshal([]byte(imagenesJSON), &entity.Imagenes)
		json.Unmarshal([]byte(listaClasesJSON), &entity.ListaClases)

		entities = append(entities, entity)
	}

	return entities, nil
}

// GetByDomainType retrieves ONAPI entities by domain type
func (r *OnapiRepository) GetByDomainType(ctx context.Context, domainType domain.DomainType, offset, limit int) ([]domain.Entity, error) {
	// For ONAPI, all entities are of the same domain type, so we just return all
	return r.List(ctx, offset, limit)
}

// GetBySearchParameter retrieves ONAPI entities by search parameter
func (r *OnapiRepository) GetBySearchParameter(ctx context.Context, searchParam string, offset, limit int) ([]domain.Entity, error) {
	return r.Search(ctx, searchParam, offset, limit)
}

// GetKeywordsByCategory retrieves keywords grouped by category for an ONAPI entity
func (r *OnapiRepository) GetKeywordsByCategory(ctx context.Context, entityID string) (map[domain.KeywordCategory][]string, error) {
	entity, err := r.GetByID(ctx, entityID)
	if err != nil {
		return nil, err
	}

	onapi := domain.NewOnapiDomain()
	return map[domain.KeywordCategory][]string{
		domain.KeywordCategoryCompanyName: onapi.GetDataByCategory(entity, domain.KeywordCategoryCompanyName),
		domain.KeywordCategoryPersonName:  onapi.GetDataByCategory(entity, domain.KeywordCategoryPersonName),
		domain.KeywordCategoryAddress:     onapi.GetDataByCategory(entity, domain.KeywordCategoryAddress),
	}, nil
}
