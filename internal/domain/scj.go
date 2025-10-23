package domain

// Sumprema corte de Justica

import (
	"insightful-intel/internal/stuff"
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
