package domain

import (
	"fmt"
	"insightful-intel/internal/stuff"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var _ DomainConnector[PGRNews] = &Pgr{}

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
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domainSearchResultId"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
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

// Implement DomainConnector[PGRNews] for Pgr
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

func (p *Pgr) GetDataByCategory(data PGRNews, category KeywordCategory) []string {

	return []string{}
}

func (p *Pgr) GetSearchableKeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		KeywordCategoryCompanyName,
		KeywordCategoryPersonName,
		KeywordCategoryAddress,
	}
}

func (p *Pgr) GetFoundKeywordCategories() []KeywordCategory {
	return []KeywordCategory{}
}
