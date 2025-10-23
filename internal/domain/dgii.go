package domain

type Register struct {
	RNC                   string `json:"rnc"`
	RazonSocial           string `json:"razonSocial"`
	NombreComercial       string `json:"nombreComercial"`
	Categoria             string `json:"categoria"`
	RegimenPagos          string `json:"regimenPagos"`
	FacturadorElectronico string `json:"facturadorElectronico"`
	LicenciaComercial     string `json:"licenciaComercial"`
	Estado                string `json:"estado"`
}
