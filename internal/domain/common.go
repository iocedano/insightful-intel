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

type contextKey string

const executionIDKey contextKey = "executionID"

// GetExecutionID retrieves the execution ID from the context
func GetExecutionID(ctx context.Context) (string, bool) {
	executionID, ok := ctx.Value(executionIDKey).(string)
	return executionID, ok
}

func SetExecutionID(ctx context.Context, executionID string) context.Context {
	return context.WithValue(ctx, executionIDKey, executionID)
}
