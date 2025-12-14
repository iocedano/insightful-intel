package domain

import (
	"insightful-intel/internal/custom"
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
	"fraude inmobiliaria",
	"fraude inmobiliaria en republica dominicana",
	"terrenos irregulares",
	"incumplimiento de contrato",
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

// GoogleDorking represents a Google Docking string search connector
type GoogleDorking struct {
	Stuff    custom.Client
	BasePath string
	PathMap  custom.CustomPathMap
}

// NewGoogleDorkingDomain creates a new Google Docking domain instance
func NewGoogleDorkingDomain() GoogleDorking {
	return GoogleDorking{
		BasePath: "https://html.duckduckgo.com/html/",
		Stuff:    *custom.NewClient(),
	}
}

// GoogleDorkingResult represents a search result from Google Docking
type GoogleDorkingResult struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domainSearchResultId"`
	SearchParameter      string    `json:"searchParameter" omitempty`
	URL                  string    `json:"link"`
	Title                string    `json:"title"`
	Description          string    `json:"snippet"`
	Relevance            float64   `json:"relevance" omitempty`
	Rank                 int       `json:"rank" omitempty`
	Keywords             []string  `json:"keywords" omitempty`
	CreatedAt            time.Time `json:"createdAt" omitempty`
	UpdatedAt            time.Time `json:"updatedAt" omitempty`
}

// GoogleDorkingSearchParams holds parameters for Google Docking search
type GoogleDorkingSearchParams struct {
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
