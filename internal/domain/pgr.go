package domain

import (
	"fmt"
	"insightful-intel/internal/stuff"
	"net/url"

	"github.com/gocolly/colly"
)

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
	URL   string
	Title string
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
