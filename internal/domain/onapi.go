package domain

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/stuff"
	"io"
)

type Onapi struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

// onapi endpoint
func NewOnapiDomain() Onapi {
	pm := stuff.PathMap{
		BaseURL: "https://www.onapi.gov.do/busquedas/api/signos/",
		Paths: map[string]string{
			"firstpage": "firstpage",
			"detail":    "detalle",
		},
	}

	return Onapi{
		BaseParh: "https://www.onapi.gov.do/busquedas/api/signos/",
		Stuff:    *stuff.NewStuff(),
		PathMap:  pm,
	}
}

type SearchComercialNameBodyResponse struct {
	Data []Entity `json:"data"`
}

type DetailsBodyResponse struct {
	Data Entity `json:"data"`
}

type Entity struct {
	Certificado    string `json:"Certificado"`
	Tipo           string `json:"Tipo"`
	Texto          string `json:"Texto"`
	Clases         string `json:"Clases"`
	Actividad      string `json:"Actividad"`
	Expedicion     string `json:"Expedicion"`
	Id             int32  `json:"Id"`
	Estado         string `json:"Estado"`
	EnTramite      bool   `json:"EnTramite"`
	Domicilio      string `json:"Domicilio"`
	Titular        string `json:"Titular"`
	Gestor         string `json:"Gestor"`
	Vencimiento    string `json:"Vencimiento"`
	Serie          int32  `json:"Serie"`
	Numero         int32  `json:"Numero"`
	TipoExpediente string `json:"TipoExpediente"`
	Secuencia      string `json:"Secuencia"`
	Expediente     string `json:"Expediente"`
	ListaClases    string `json:"ListaClases"`
}

func (o *Onapi) SearchComercialName(query string) (*SearchComercialNameBodyResponse, error) {
	response, err := o.Stuff.Client.Get(o.PathMap.GetURLFrom("firstpage"), map[string]string{
		"tipoBusqueda": "0",
		"texto":        query,
		"tipo":         "NO/NO",
		"clase":        "0",
		"pgSize":       "1000",
		"pgIndex":      "1",
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
	var onapiResponse SearchComercialNameBodyResponse
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &onapiResponse, nil
}

func (o *Onapi) GetDetails(id int32) (*DetailsBodyResponse, error) {
	response, err := o.Stuff.Client.Get(o.PathMap.GetURLFrom("detail"), map[string]string{
		"id": fmt.Sprintf("%d", id),
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
	var onapiResponse DetailsBodyResponse
	if err := json.Unmarshal(body, &onapiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &onapiResponse, nil
}
