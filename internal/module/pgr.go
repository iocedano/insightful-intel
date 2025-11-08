package module

import (
	"fmt"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/stuff"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

var _ domain.DomainConnector[domain.PGRNews] = &Pgr{}

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

func (*Pgr) GetDomainType() domain.DomainType {
	return domain.DomainTypePGR
}

func (p *Pgr) Search(query string) ([]domain.PGRNews, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("pgr.gob.do"),
	)

	var news []domain.PGRNews

	c.OnHTML("article", func(e *colly.HTMLElement) {
		article := domain.PGRNews{}

		article.URL = e.ChildAttr("h5 a", "href")
		article.Title = e.ChildAttr("h5 a", "title")

		news = append(news, article)

	})
	c.Visit(fmt.Sprintf("%s?s=%s", p.BaseParh, url.QueryEscape(query)))

	return news, nil
}

func (p *Pgr) ProcessData(data domain.PGRNews) (domain.PGRNews, error) {
	if err := p.ValidateData(data); err != nil {
		return domain.PGRNews{}, err
	}
	return p.TransformData(data), nil
}

func (p *Pgr) ValidateData(data domain.PGRNews) error {
	if data.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if data.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

func (p *Pgr) TransformData(data domain.PGRNews) domain.PGRNews {
	transformed := data
	transformed.URL = strings.TrimSpace(data.URL)
	transformed.Title = strings.TrimSpace(data.Title)

	return transformed
}

func (p *Pgr) GetDataByCategory(data domain.PGRNews, category domain.KeywordCategory) []string {
	return []string{}
}

func (p *Pgr) GetSearchableKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryCompanyName,
		domain.KeywordCategoryPersonName,
		domain.KeywordCategoryAddress,
	}
}

func (p *Pgr) GetFoundKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{}
}
