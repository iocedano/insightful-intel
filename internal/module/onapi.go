package module

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/stuff"
	"io"
	"strings"
)

var _ domain.GenericConnector[domain.Entity] = &Onapi{}

type Onapi struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

func (*Onapi) GetDomainType() domain.DomainType {
	return domain.DomainTypeONAPI
}

// onapi endpoint
func NewOnapiDomain() Onapi {
	pm := stuff.PathMap{
		BaseURL: "https://www.onapi.gob.do/busqapi/signos/",
		Paths: map[string]string{
			"firstpage": "",
			"detail":    "byexp",
		},
	}

	return Onapi{
		BaseParh: "https://www.onapi.gob.do/busqapi/signos/",
		Stuff:    *stuff.NewStuff(),
		PathMap:  pm,
	}
}

// Implement GenericConnector[domain.Entity] for Onapi
func (o *Onapi) ProcessData(data domain.Entity) (domain.Entity, error) {
	// Process the entity data (e.g., clean, validate, enrich)
	if err := o.ValidateData(data); err != nil {
		return domain.Entity{}, err
	}
	return o.TransformData(data), nil
}

func (o *Onapi) ValidateData(data domain.Entity) error {
	// Validate the entity data
	if data.NumeroExpediente == 0 {
		return fmt.Errorf("NumeroExpediente is required")
	}
	if data.SerieExpediente == 0 {
		return fmt.Errorf("SerieExpediente is required")
	}
	return nil
}

func (o *Onapi) TransformData(data domain.Entity) domain.Entity {
	transformed := data
	transformed.Texto = strings.TrimSpace(data.Texto)
	transformed.Titular = strings.TrimSpace(data.Titular)
	transformed.Gestor = strings.TrimSpace(data.Gestor)
	transformed.Domicilio = strings.TrimSpace(data.Domicilio)
	return transformed
}

func (o *Onapi) GetDataByCategory(data domain.Entity, category domain.KeywordCategory) []string {
	result := []string{}

	switch category {
	case domain.KeywordCategoryCompanyName:
		result = append(result, data.Texto)
	case domain.KeywordCategoryPersonName:
		result = append(result, data.Titular, data.Gestor)
	case domain.KeywordCategoryAddress:
		result = append(result, data.Domicilio)
	}

	return result
}

func (o *Onapi) GetSearchableKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryCompanyName,
	}
}

func (o *Onapi) GetFoundKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryCompanyName,
		domain.KeywordCategoryPersonName,
		domain.KeywordCategoryAddress,
	}
}

func (o *Onapi) SearchComercialName(query string) ([]domain.Entity, error) {
	response, err := o.Stuff.Client.Get(o.PathMap.GetURLFrom("firstpage"), map[string]string{
		"subtipo":  "",
		"texto":    query,
		"tipo":     "",
		"clases":   "",
		"pageSize": "1000",
		"pageIdx":  "1",
	}, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var onapiResponse []domain.Entity
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	for i, onapiEntity := range onapiResponse {
		if o.ValidateData(onapiEntity) == nil {
			details, err := o.GetDetails(onapiEntity.NumeroExpediente, onapiEntity.SerieExpediente)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
			}
			onapiResponse[i] = *details
		}
	}

	return onapiResponse, nil
}

func (o *Onapi) GetDetails(numero int32, serie int32) (*domain.Entity, error) {
	response, err := o.Stuff.Client.Get(o.PathMap.GetURLFrom("detail"), map[string]string{
		"numero":    fmt.Sprintf("%d", numero),
		"tipoExped": "E",
		"serie":     fmt.Sprintf("%d", serie),
	}, map[string]string{
		"Content-Type": "application/json",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the struct
	var onapiResponse domain.Entity
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &onapiResponse, nil
}
