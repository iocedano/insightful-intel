package domain

import (
	"fmt"
	"insightful-intel/internal/stuff"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

var _ GenericConnector[PGRNews] = &Pgr{}

type Pgr struct {
	Stuff    stuff.Stuff
	BaseParh string
	PathMap  stuff.PathMap
}

func NewPgrDomain() Pgr {
	return Pgr{
		BaseParh: "https://pgr.gob.do/",
		Stuff:    *stuff.NewStuff(),
	}
}

type PGRNews struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (*Pgr) GetDomainType() DomainType {
	return DomainTypePGR
}

// doc: https://www.zenrows.com/blog/web-scraping-golang#install-required-libraries

func (p *Pgr) Search(query string) ([]PGRNews, error) {

	c := colly.NewCollector(
		colly.AllowedDomains("pgr.gob.do"),
	)

	var news []PGRNews

	c.OnHTML("article", func(e *colly.HTMLElement) {
		article := PGRNews{}

		article.URL = e.ChildAttr("h5 a", "href")
		article.Title = e.ChildAttr("h5 a", "title")

		news = append(news, article)

	})
	c.Visit(fmt.Sprintf("%s?s=%s", p.BaseParh, url.QueryEscape(query)))

	return news, nil
}

// Implement GenericConnector[PGRNews] for Onapi
func (p *Pgr) ProcessData(data PGRNews) (PGRNews, error) {
	// Process the entity data (e.g., clean, validate, enrich)
	if err := p.ValidateData(data); err != nil {
		return PGRNews{}, err
	}
	return p.TransformData(data), nil
}

func (p *Pgr) ValidateData(data PGRNews) error {
	// Validate the entity data
	if data.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if data.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

func (p *Pgr) TransformData(data PGRNews) PGRNews {
	transformed := data
	transformed.URL = strings.TrimSpace(data.URL)
	transformed.Title = strings.TrimSpace(data.Title)

	return transformed
}

func (p *Pgr) GetDataByCategory(data PGRNews, category DataCategory) []string {

	return []string{}
}

func (p *Pgr) GetListOfSearchableCategory() []DataCategory {
	return []DataCategory{
		DataCategoryCompanyName,
		DataCategoryPersonName,
		DataCategoryAddress,
	}
}

func (p *Pgr) GetListOfRetrievedCategory() []DataCategory {
	return []DataCategory{}
}
