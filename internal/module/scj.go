package module

// Sumprema corte de Justica

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/stuff"
	"io"
	"net/url"
	"strings"
)

var _ domain.DomainConnector[domain.ScjCase] = &Scj{}

type Scj struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

func NewScjDomain() Scj {
	return Scj{
		BaseParh: "https://consultasentenciascj.poderjudicial.gob.do/Home/GetExpedientes",
		Stuff:    *stuff.NewStuff(),
	}
}

func (p *Scj) Search(query string) ([]domain.ScjCase, error) {
	form := url.Values{}
	form.Add("search[value]", query)
	form.Add("Contenido", query)
	form.Add("search[regex]", "false")
	form.Add("start", "0")
	form.Add("length", "10")

	resp, err := p.Stuff.Client.Post(p.BaseParh, form.Encode(), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result domain.ScjSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return result.Data, nil
}

func (dgi *Scj) GetDomainType() domain.DomainType {
	return domain.DomainTypeSCJ
}

// Implement DomainConnector[domain.ScjCase] for Scj
func (p *Scj) ProcessData(data domain.ScjCase) (domain.ScjCase, error) {
	// Process the entity data (e.g., clean, validate, enrich)
	if err := p.ValidateData(data); err != nil {
		return domain.ScjCase{}, err
	}
	return p.TransformData(data), nil
}

func (p *Scj) ValidateData(data domain.ScjCase) error {
	// Validate the entity data
	if data.IDExpediente == 0 {
		return fmt.Errorf("id is required")
	}

	return nil
}

func (p *Scj) TransformData(data domain.ScjCase) domain.ScjCase {
	transformed := data
	transformed.IDExpediente = data.IDExpediente
	transformed.Involucrados = strings.TrimSpace(data.Involucrados)
	transformed.URLBlob = strings.TrimSpace(data.URLBlob)
	transformed.DescMateria = strings.TrimSpace(data.DescMateria)

	return transformed
}

func (p *Scj) GetDataByCategory(data domain.ScjCase, category domain.KeywordCategory) []string {
	result := []string{}

	switch category {
	case domain.KeywordCategoryPersonName:
		result = append(result, data.Involucrados)
	}

	return result
}

func (p *Scj) GetSearchableKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryPersonName,
		domain.KeywordCategoryCompanyName,
	}
}

func (p *Scj) GetFoundKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryPersonName,
	}
}
