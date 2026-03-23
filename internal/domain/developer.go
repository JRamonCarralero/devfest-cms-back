package domain

import (
	"context"

	"github.com/google/uuid"
)

type Developer struct {
	ID              uuid.UUID
	EventID         uuid.UUID
	RoleDescription *string
	Person
	Audit
}

type UpdateDeveloper struct {
	RoleDescription *string
	UpdatedBy       uuid.UUID
}

type DeveloperRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Developer, error)
	GetById(ctx context.Context, id uuid.UUID) (*Developer, error)
	GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Developer, int64, error)
	// Writers
	Create(ctx context.Context, developer *Developer) (*Developer, error)
	Update(ctx context.Context, developer *Developer) (*Developer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
