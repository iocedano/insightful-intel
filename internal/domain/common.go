package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Common struct {
	ID                   uuid.UUID `json:"id"`
	DomainSearchResultID uuid.UUID `json:"domainSearchResultId"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CommonRepository[T any] interface {
	Create(ctx context.Context, entity T) error
	GetByID(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, id string, entity T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]T, error)
	Count(ctx context.Context) (int64, error)
}

type ID = uuid.UUID

func NewID() ID {
	return uuid.New()
}

func NewIDFromString(id string) ID {
	return uuid.MustParse(id)
}
