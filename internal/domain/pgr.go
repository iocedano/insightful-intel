package domain

import (
	"time"
)

type PGRNews struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domainSearchResultId"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
