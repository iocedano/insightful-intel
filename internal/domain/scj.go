package domain

// Sumprema corte de Justica

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/stuff"
	"io"
	"net/url"
)

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
