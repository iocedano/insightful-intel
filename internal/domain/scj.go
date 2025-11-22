package domain

// Sumprema corte de Justica

import (
	"insightful-intel/internal/custom"
	"time"
)

type Scj struct {
	Stuff    custom.Client
	BaseParh string
	PathMap  custom.CustomPathMap
}

func NewScjDomain() Scj {
	return Scj{
		BaseParh: "https://consultasentenciascj.poderjudicial.gob.do/Home/GetExpedientes",
		Stuff:    *custom.NewClient(),
	}
}

type ScjSearchResponse struct {
	Draw            string    `json:"draw"`
	RecordsFiltered int       `json:"recordsFiltered"`
	RecordsTotal    int       `json:"recordsTotal"`
	Data            []ScjCase `json:"data"`
}

type ScjCase struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domain_search_result_id"`
	Linea                int       `json:"linea"`
	AgnoCabecera         int       `json:"agno_cabecera"`
	MesCabecera          int       `json:"mes_cabecera"`
	URLCabecera          string    `json:"url_cabecera"`
	URLCuerpo            string    `json:"url_cuerpo"`
	IDExpediente         int       `json:"id_expediente"`
	NoExpediente         string    `json:"no_expediente"`
	NoSentencia          string    `json:"no_sentencia"`
	NoUnico              string    `json:"no_unico"`
	NoInterno            string    `json:"no_interno"`
	IDTribunal           string    `json:"id_tribunal"`
	DescTribunal         string    `json:"desc_tribunal"`
	IDMateria            string    `json:"id_materia"`
	DescMateria          string    `json:"desc_materia"`
	FechaFallo           string    `json:"fecha_fallo"`
	Involucrados         string    `json:"involucrados"`
	GuidBlob             string    `json:"guid_blob"`
	TipoDocumentoAdjunto string    `json:"tipo_documento_adjunto"`
	TotalFilas           int       `json:"total_filas"`
	URLBlob              string    `json:"url_blob"`
	Extension            string    `json:"extension"`
	Origen               int       `json:"origen"`
	Activo               bool      `json:"activo"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
