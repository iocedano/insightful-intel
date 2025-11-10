package domain

import (
	"insightful-intel/internal/stuff"
	"time"
)

var FRAUD_KEYWORDS = []string{
	"fraude",
	"estafa",
	"denuncia",
	"engaño",
	"irregular",
	"contrato",
	"pagos",
	"retrasos",
	"robo",
	"acusa",
	"acusacion",
	"acusado",
}

var SOCIAL_MEDIA_SITES_KEYWORDS = []string{
	"facebook.com",
	"twitter.com",
	"instagram.com",
	"linkedin.com",
	"youtube.com",
	"tiktok.com",
	"linkedin.com/company",
}

var X_IN_URL_KEYWORDS = []string{
	"status", // finds specific posts, threads, or tweet IDs
	"lists",  //  finds public user lists
}

var ADDRESS_KEYWORDS = []string{
	"direccion",
	"domicilio",
	"ubicacion",
	"direccion postal",
}

var ENTITY_NAME_KEYWORDS = []string{
	"empresa",
	"compania",
	"nombre",
	"nombre completo",
	"nombre y apellido",
}

var FILE_TYPE_KEYWORDS = []string{
	"pdf",
	"doc",
	"docx",
	"xls",
	"xlsx",
	"ppt",
	"pptx",
}

// GoogleDocking represents a Google Docking string search connector
type GoogleDocking struct {
	Stuff    stuff.Stuff
	BasePath string
	PathMap  stuff.PathMap
}

// NewGoogleDockingDomain creates a new Google Docking domain instance
func NewGoogleDockingDomain() GoogleDocking {
	return GoogleDocking{
		BasePath: "https://html.duckduckgo.com/html/",
		Stuff:    *stuff.NewStuff(),
	}
}

// GoogleDockingResult represents a search result from Google Docking
type GoogleDockingResult struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domainSearchResultId"`
	SearchParameter      string    `json:"searchParameter"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Relevance            float64   `json:"relevance"`
	Rank                 int       `json:"rank"`
	Keywords             []string  `json:"keywords"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// GoogleDockingSearchParams holds parameters for Google Docking search
type GoogleDockingSearchParams struct {
	Query            string   `json:"query"`
	MaxResults       int      `json:"max_results"`
	MinRelevance     float64  `json:"min_relevance"`
	ExactMatch       bool     `json:"exact_match"`
	CaseSensitive    bool     `json:"case_sensitive"`
	IncludeKeywords  []string `json:"include_keywords"`
	FileTypeKeywords []string `json:"file_type_keywords"`
	InURLKeywords    []string `json:"in_url_keywords"`
	ExcludeKeywords  []string `json:"exclude_keywords"`
	SitesKeywords    []string `json:"sites_keywords"`
}
