package module

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/stuff"
	"io"
	"strings"
	"time"
)

var _ domain.DomainConnector[domain.Entity] = &Onapi{}

type Onapi struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

type OnapiEntityResponse struct {
	ID                int32               `json:"id"`
	SerieExpediente   int32               `json:"serieExpediente"`
	NumeroExpediente  int32               `json:"numeroExpediente"`
	Certificado       string              `json:"certificado,omitempty"`
	Tipo              string              `json:"tipo,omitempty"`
	SubTipo           string              `json:"subTipo,omitempty"`
	Texto             string              `json:"texto"`
	Clases            string              `json:"clases"`
	AplicadoAProteger string              `json:"aplicadoAProteger"`
	Expedicion        string              `json:"expedicion"`
	Vencimiento       string              `json:"vencimiento"`
	EnTramite         bool                `json:"enTramite"`
	Titular           string              `json:"titular"`
	Gestor            string              `json:"gestor"`
	Domicilio         string              `json:"domicilio"`
	Status            string              `json:"status"`
	TipoSigno         string              `json:"tipoSigno"`
	Imagenes          []domain.Image      `json:"imagenes"`
	ListaClases       []domain.ListaClase `json:"listaClases"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
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

// Implement DomainConnector[domain.Entity] for Onapi
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

// Search implements DomainConnector interface by wrapping SearchComercialName
func (o *Onapi) Search(query string) ([]domain.Entity, error) {
	return o.SearchComercialName(query)
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
		"Content-Type":    "application/json",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "gzip, deflate, br, zstd",
		"Referer":         "https://www.onapi.gob.do/busquedas2021/signos/buscar",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"Host":            "www.onapi.gob.do",
	})

	if err != nil {

		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var onapiResponse []OnapiEntityResponse
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	domainEntities := []domain.Entity{}
	for _, onapiEntity := range onapiResponse {
		domainEntity := toDomainEntity(onapiEntity)
		if o.ValidateData(domainEntity) == nil {
			details, err := o.GetDetails(domainEntity.NumeroExpediente, domainEntity.SerieExpediente)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarsha Details JSON response: %w", err)
			}

			domainEntities = append(domainEntities, *details)
		}
	}

	return domainEntities, nil
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
	var onapiResponse OnapiEntityResponse
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	domainEntity := toDomainEntity(onapiResponse)

	return &domainEntity, nil
}

func toDomainEntity(onapiEntity OnapiEntityResponse) domain.Entity {
	return domain.Entity{
		SerieExpediente:   onapiEntity.SerieExpediente,
		NumeroExpediente:  onapiEntity.NumeroExpediente,
		Certificado:       onapiEntity.Certificado,
		Tipo:              onapiEntity.Tipo,
		SubTipo:           onapiEntity.SubTipo,
		Texto:             onapiEntity.Texto,
		Clases:            onapiEntity.Clases,
		AplicadoAProteger: onapiEntity.AplicadoAProteger,
		Expedicion:        onapiEntity.Expedicion,
		Vencimiento:       onapiEntity.Vencimiento,
		EnTramite:         onapiEntity.EnTramite,
		Titular:           onapiEntity.Titular,
		Gestor:            onapiEntity.Gestor,
		Domicilio:         onapiEntity.Domicilio,
		Status:            onapiEntity.Status,
		TipoSigno:         onapiEntity.TipoSigno,
		Imagenes:          onapiEntity.Imagenes,
		ListaClases:       onapiEntity.ListaClases,
		CreatedAt:         onapiEntity.CreatedAt,
		UpdatedAt:         onapiEntity.UpdatedAt,
	}
}
