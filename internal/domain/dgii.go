package domain

import "time"

type Register struct {
	ID                    ID        `json:"id"`
	DomainSearchResultID  ID        `json:"domain_search_result_id"`
	RNC                   string    `json:"rnc"`
	RazonSocial           string    `json:"razon_social"`
	NombreComercial       string    `json:"nombre_comercial"`
	Categoria             string    `json:"categoria"`
	RegimenPagos          string    `json:"regimen_pagos"`
	FacturadorElectronico string    `json:"facturador_electronico"`
	LicenciaComercial     string    `json:"licencia_comercial"`
	Estado                string    `json:"estado"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
