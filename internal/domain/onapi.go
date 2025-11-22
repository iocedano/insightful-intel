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
	DomainSearchResultID ID           `json:"domain_search_result_id"`
	SerieExpediente      int32        `json:"serie_expediente"`
	NumeroExpediente     int32        `json:"numero_expediente"`
	Certificado          string       `json:"certificado,omitempty"`
	Tipo                 string       `json:"tipo,omitempty"`
	SubTipo              string       `json:"sub_tipo,omitempty"`
	Texto                string       `json:"texto"`
	Clases               string       `json:"clases"`
	AplicadoAProteger    string       `json:"aplicado_a_proteger"`
	Expedicion           string       `json:"expedicion"`
	Vencimiento          string       `json:"vencimiento"`
	EnTramite            bool         `json:"en_tramite"`
	Titular              string       `json:"titular"`
	Gestor               string       `json:"gestor"`
	Domicilio            string       `json:"domicilio"`
	Status               string       `json:"status"`
	TipoSigno            string       `json:"tipo_signo"`
	Imagenes             []Image      `json:"imagenes"`
	ListaClases          []ListaClase `json:"lista_clases"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

type ListaClase struct {
	Numero    int32  `json:"numero"`
	Productos string `json:"productos"`
}

type Image struct {
	// ID                 int32   `json:"id"`
	SerieExpediente    int32   `json:"serie_expediente"`
	NumeroExpediente   int32   `json:"numero_expediente"`
	DescripcionColores *string `json:"descripcion_colores"`
	Bytes              *string `json:"bytes"`
	CodigoFormato      int32   `json:"codigo_formato"`
	MimeType           string  `json:"mime_type"`
	FileExtension      string  `json:"file_extension"`
}
