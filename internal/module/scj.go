package module

// Sumprema corte de Justica

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/custom"
	"insightful-intel/internal/domain"
	"io"
	"net/url"
	"regexp"
	"strings"
)

var _ domain.DomainConnector[domain.ScjCase] = &Scj{}

type ScjSearchResponse struct {
	Draw            string            `json:"draw"`
	RecordsFiltered int               `json:"recordsFiltered"`
	RecordsTotal    int               `json:"recordsTotal"`
	Data            []ScjCaseResponse `json:"data"`
}

type ScjCaseResponse struct {
	Linea                int    `json:"linea"`
	AgnoCabecera         int    `json:"agnoCabecera"`
	MesCabecera          int    `json:"mesCabecera"`
	URLCabecera          string `json:"urlCabecera"`
	URLCuerpo            string `json:"urlCuerpo"`
	IDExpediente         int    `json:"idExpediente"`
	NoExpediente         string `json:"noExpediente"`
	NoSentencia          string `json:"noSentencia"`
	NoUnico              string `json:"noUnico"`
	NoInterno            string `json:"noInterno"`
	IDTribunal           string `json:"idTribunal"`
	DescTribunal         string `json:"descTribunal"`
	IDMateria            string `json:"idMateria"`
	DescMateria          string `json:"descMateria"`
	FechaFallo           string `json:"fechaFallo"`
	Involucrados         string `json:"involucrados"`
	GuidBlob             string `json:"guidBlob"`
	TipoDocumentoAdjunto string `json:"tipoDocumentoAdjunto"`
	TotalFilas           int    `json:"totalFilas"`
	UrlBlob              string `json:"urlBlob"`
	Extension            string `json:"extension"`
	Origen               int    `json:"origen"`
	Activo               bool   `json:"activo"`
}

type Scj struct {
	Stuff    custom.Client
	BaseParh string
	PathMap  custom.CustomPathMap
}

var REMOVE_WORDS = map[string]bool{
	"S. A.":    true,
	"S.R.L.":   true,
	"S. R. L.": true,
	"N. V.":    true,
	"N.V.":     true,
	"S.A.":     true,
}

func NewScjDomain() domain.DomainConnector[domain.ScjCase] {
	return &Scj{
		BaseParh: "https://consultasentenciascj.poderjudicial.gob.do/Home/GetExpedientes",
		Stuff:    *custom.NewClient(),
	}
}

func (p *Scj) Search(query string) ([]domain.ScjCase, error) {
	form := url.Values{}
	form.Add("search[value]", query)
	form.Add("Contenido", query)
	form.Add("search[regex]", "false")
	form.Add("start", "0")
	form.Add("length", "10")

	resp, err := p.Stuff.Post(p.BaseParh, form.Encode(), map[string]string{
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

	var result ScjSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	var cases []domain.ScjCase
	for _, item := range result.Data {
		cases = append(cases, p.ToDomain(item))
	}
	return cases, nil
}

func (p *Scj) ToDomain(response ScjCaseResponse) domain.ScjCase {
	return domain.ScjCase{
		Linea:                response.Linea,
		AgnoCabecera:         response.AgnoCabecera,
		MesCabecera:          response.MesCabecera,
		URLCabecera:          response.URLCabecera,
		URLCuerpo:            response.URLCuerpo,
		IDExpediente:         response.IDExpediente,
		NoExpediente:         response.NoExpediente,
		NoSentencia:          response.NoSentencia,
		NoUnico:              response.NoUnico,
		NoInterno:            response.NoInterno,
		IDTribunal:           response.IDTribunal,
		DescTribunal:         response.DescTribunal,
		IDMateria:            response.IDMateria,
		DescMateria:          response.DescMateria,
		FechaFallo:           response.FechaFallo,
		Involucrados:         response.Involucrados,
		GuidBlob:             response.GuidBlob,
		TipoDocumentoAdjunto: response.TipoDocumentoAdjunto,
		TotalFilas:           response.TotalFilas,
		URLBlob:              response.UrlBlob,
		Extension:            response.Extension,
		Origen:               response.Origen,
		Activo:               response.Activo,
	}
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
		// Split by comma or "vs"/"vs." (case-insensitive) with optional spaces
		// Pattern matches: comma with spaces, or "vs"/"vs." with spaces (case-insensitive)
		re := regexp.MustCompile(`(?i)\s*,\s*|\s+vs\.?\s*`)
		// trim the data.Involucrados
		involucrados := strings.TrimSpace(data.Involucrados)
		names := re.Split(involucrados, -1)
		for _, name := range names {
			trimmed := strings.TrimSpace(name)
			if _, ok := REMOVE_WORDS[trimmed]; ok || trimmed == "" {
				continue
			}

			result = append(result, trimmed)
		}
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
