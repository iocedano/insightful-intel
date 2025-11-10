package domain

import (
	"time"
)

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
	Certificado          string       `json:"certificado,omitempty"`
	Tipo                 string       `json:"tipo,omitempty"`
	SubTipo              string       `json:"subTipo,omitempty"`
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
	// ID                 int32   `json:"id"`
	SerieExpediente    int32   `json:"serieExpediente"`
	NumeroExpediente   int32   `json:"numeroExpediente"`
	DescripcionColores *string `json:"descripcionColores"`
	Bytes              *string `json:"bytes"`
	CodigoFormato      int32   `json:"codigoFormato"`
	MimeType           string  `json:"mimeType"`
	FileExtension      string  `json:"fileExtension"`
}
