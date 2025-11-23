package domain

import (
	"time"
)

type PGRNews struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domain_search_result_id"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
