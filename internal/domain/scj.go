package domain

// Sumprema corte de Justica

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/stuff"
	"io"
	"net/url"
	"strings"
)

var _ GenericConnector[ScjCase] = &Scj{}

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

type ScjSearchResponse struct {
	Draw            string    `json:"draw"`
	RecordsFiltered int       `json:"recordsFiltered"`
	RecordsTotal    int       `json:"recordsTotal"`
	Data            []ScjCase `json:"data"`
}

type ScjCase struct {
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
	URLBlob              string `json:"urlBlob"`
	Extension            string `json:"extension"`
	Origen               int    `json:"origen"`
	Activo               bool   `json:"activo"`
}

func (p *Scj) Search(query string) ([]ScjCase, error) {
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

	var result ScjSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return result.Data, nil
}

func (dgi *Scj) GetName() string {
	return "SCJ"
}

// Implement GenericConnector[ScjCase] for Onapi
func (p *Scj) ProcessData(data ScjCase) (ScjCase, error) {
	// Process the entity data (e.g., clean, validate, enrich)
	if err := p.ValidateData(data); err != nil {
		return ScjCase{}, err
	}
	return p.TransformData(data), nil
}

func (p *Scj) ValidateData(data ScjCase) error {
	// Validate the entity data
	if data.IDExpediente == 0 {
		return fmt.Errorf("id is required")
	}

	return nil
}

func (p *Scj) TransformData(data ScjCase) ScjCase {
	transformed := data
	transformed.IDExpediente = data.IDExpediente
	transformed.Involucrados = strings.TrimSpace(data.Involucrados)
	transformed.URLBlob = strings.TrimSpace(data.URLBlob)
	transformed.DescMateria = strings.TrimSpace(data.DescMateria)

	return transformed
}

func (p *Scj) GetDataByCategory(data ScjCase, category DataCategory) []string {
	result := []string{}

	switch category {
	case DataCategoryPersonName:
		result = append(result, data.Involucrados)
	}

	return result
}

func (p *Scj) GetListOfSearchableCategory() []DataCategory {
	return []DataCategory{
		DataCategoryPersonName,
		DataCategoryCompanyName,
	}
}

func (p *Scj) GetListOfRetrievedCategory() []DataCategory {
	return []DataCategory{
		DataCategoryPersonName,
	}
}
