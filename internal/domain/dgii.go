package domain

import (
	"bytes"
	"insightful-intel/internal/stuff"
	"io"
	"net/url"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
)

type Dgii struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

type Register struct {
	RNCRNC                string
	RazonSocial           string
	NombreComercial       string
	Categoria             string
	RegimenPagos          string
	FacturadorElectronico string
	LicenciaComercial     string
	Estado                string
}

func NewDgiiDomain() Dgii {
	return Dgii{
		BaseParh: "https://dgii.gov.do/app/WebApps/ConsultasWeb2/ConsultasWeb/consultas/rnc.aspx",
		Stuff:    *stuff.NewStuff(),
	}
}

func (dgi *Dgii) GetRegister(query string) ([]Register, error) {
	response, err := dgi.Stuff.Client.Get(dgi.BaseParh, nil, map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
	})

	defer response.Body.Close()

	data := make(map[string]string)
	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	doc, _ := html.Parse(bytes.NewReader(body))
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					data[attr.Val] = ""
					for _, attrValue := range n.Attr {
						if attrValue.Key == "value" {
							data[attr.Val] = attrValue.Val
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	data["ctl00$MainContent$txtName"] = "bank"
	data["ctl00$MainContent$cbIncludeCeased"] = "on"
	data["ctl00$smMain"] = "ctl00$cphMain$upBusqueda|ctl00$cphMain$btnBuscarPorRazonSocial"
	data["ctl00$cphMain$txtRNCCedula"] = ""
	data["ctl00$cphMain$txtRazonSocial"] = query
	data["ctl00$cphMain$btnBuscarPorRazonSocial"] = "Buscar"
	data["ctl00$cphMain$hidActiveTab"] = "razonsocial"
	data["ctl00$cphMain$hidActiveTab"] = "ctl00$cphMain$upBusqueda|ctl00$cphMain$btnBuscarPorRazonSocial"
	data["__EVENTARGUMENT"] = "Page$2"
	data["__EVENTTARGET"] = "ctl00$cphMain$gvBuscRazonSocial"
	data["__ASYNCPOST"] = "true"

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}

	resp, err := dgi.Stuff.Client.Post(dgi.BaseParh, formData.Encode(), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
	})

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		// return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, _ = html.Parse(bytes.NewReader(body))
	var rows []*html.Node
	f = func(n *html.Node) {
		if n.Data == "tr" && n.Parent.Data == "tbody" && (len(n.Attr) == 1 && n.Attr[0].Val == "TbRow") {
			rows = append(rows, n)

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	result := []Register{}

	for _, row := range rows {
		if row.Type == html.ElementNode && row.Data == "tr" {
			var cells []string
			for cell := row.FirstChild; cell != nil; cell = cell.NextSibling {
				if cell.Type == html.ElementNode && (cell.Data == "th" || cell.Data == "td") {
					cells = append(cells, strings.TrimSpace(getInnerText(cell)))
				}
			}

			spew.Dump(cells)
			result = append(result, Register{
				RNCRNC:                cells[0],
				RazonSocial:           cells[1],
				NombreComercial:       cells[2],
				Categoria:             cells[3],
				RegimenPagos:          cells[4],
				Estado:                cells[5],
				FacturadorElectronico: cells[6],
				LicenciaComercial:     cells[7],
			})
		}
	}

	spew.Dump(result)

	return result, nil
}

func getInnerText(n *html.Node) string {
	var buf strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return buf.String()
}
