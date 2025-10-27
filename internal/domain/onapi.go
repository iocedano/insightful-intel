package domain

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/stuff"
	"io"
	"strings"
	"time"
)

var _ GenericConnector[Entity] = &Onapi{}

type Onapi struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

func (*Onapi) GetDomainType() DomainType {
	return DomainTypeONAPI
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

// Implement GenericConnector[Entity] for Onapi
func (o *Onapi) ProcessData(data Entity) (Entity, error) {
	// Process the entity data (e.g., clean, validate, enrich)
	if err := o.ValidateData(data); err != nil {
		return Entity{}, err
	}
	return o.TransformData(data), nil
}

func (o *Onapi) ValidateData(data Entity) error {
	// Validate the entity data
	if data.NumeroExpediente == 0 {
		return fmt.Errorf("NumeroExpediente is required")
	}
	if data.SerieExpediente == 0 {
		return fmt.Errorf("SerieExpediente is required")
	}
	return nil
}

func (o *Onapi) TransformData(data Entity) Entity {
	transformed := data
	transformed.Texto = strings.TrimSpace(data.Texto)
	transformed.Titular = strings.TrimSpace(data.Titular)
	transformed.Gestor = strings.TrimSpace(data.Gestor)
	transformed.Domicilio = strings.TrimSpace(data.Domicilio)
	return transformed
}

func (o *Onapi) GetDataByCategory(data Entity, category KeywordCategory) []string {
	result := []string{}

	switch category {
	case KeywordCategoryCompanyName:
		result = append(result, data.Texto)
	case KeywordCategoryPersonName:
		result = append(result, data.Titular, data.Gestor)
	case KeywordCategoryAddress:
		result = append(result, data.Domicilio)
	}

	return result
}

func (o *Onapi) GetSearchableKeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		KeywordCategoryCompanyName,
	}
}

func (o *Onapi) GetFoundKeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		KeywordCategoryCompanyName,
		KeywordCategoryPersonName,
		KeywordCategoryAddress,
	}
}

type SearchComercialNameBodyResponse struct {
	Data []Entity `json:"data"`
}

type DetailsBodyResponse struct {
	Data Entity `json:"data"`
}

type Entity struct {
	ID                   ID           `json:"id"`
	DomainSearchResultID ID           `json:"domainSearchResultId"`
	SerieExpediente      int32        `json:"serieExpediente"`
	NumeroExpediente     int32        `json:"numeroExpediente"`
	Certificado          string       `json:"certificado"`
	Tipo                 string       `json:"tipo"`
	SubTipo              string       `json:"subTipo"`
	Texto                string       `json:"texto"`
	Clases               string       `json:"clases"`
	AplicadoAProteger    string       `json:"aplicadoAProteger"`
	Expedicion           string       `json:"expedicion"`
	Vencimiento          string       `json:"vencimiento"`
	EnTramite            bool         `json:"enTramite"`
	Titular              string       `json:"titular"`
	Gestor               string       `json:"gestor"`
	Domicilio            string       `json:"domicilio"`
	Status               string       `json:"status"`
	TipoSigno            string       `json:"tipoSigno"`
	Imagenes             []Image      `json:"imagenes"`
	ListaClases          []ListaClase `json:"listaClases"`
	CreatedAt            time.Time    `json:"createdAt"`
	UpdatedAt            time.Time    `json:"updatedAt"`
}

type ListaClase struct {
	Numero    int32  `json:"numero"`
	Productos string `json:"productos"`
}

type Image struct {
	ID                 int32   `json:"id"`
	SerieExpediente    int32   `json:"serieExpediente"`
	NumeroExpediente   int32   `json:"numeroExpediente"`
	DescripcionColores *string `json:"descripcionColores"`
	Bytes              *string `json:"bytes"`
	CodigoFormato      int32   `json:"codigoFormato"`
	MimeType           string  `json:"mimeType"`
	FileExtension      string  `json:"fileExtension"`
}

func (o *Onapi) SearchComercialName(query string) ([]Entity, error) {
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

	var onapiResponse []Entity
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

func (o *Onapi) GetDetails(numero int32, serie int32) (*Entity, error) {
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
	var onapiResponse Entity
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &onapiResponse, nil
}
