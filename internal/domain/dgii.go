package domain

import "time"

type Register struct {
	ID                    ID        `json:"id"`
	DomainSearchResultID  ID        `json:"domainSearchResultId"`
	RNC                   string    `json:"rnc"`
	RazonSocial           string    `json:"razonSocial"`
	NombreComercial       string    `json:"nombreComercial"`
	Categoria             string    `json:"categoria"`
	RegimenPagos          string    `json:"regimenPagos"`
	FacturadorElectronico string    `json:"facturadorElectronico"`
	LicenciaComercial     string    `json:"licenciaComercial"`
	Estado                string    `json:"estado"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
