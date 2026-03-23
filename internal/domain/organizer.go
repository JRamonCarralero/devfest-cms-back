package domain

import (
	"context"

	"github.com/google/uuid"
)

type Organizer struct {
	ID              uuid.UUID
	EventID         uuid.UUID
	Company         *string
	RoleDescription *string
	Person
	Audit
}

type UpdateOrganizer struct {
	Company         *string
	RoleDescription *string
	UpdatedBy       uuid.UUID
}

type OrganizerRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Organizer, error)
	GetById(ctx context.Context, id uuid.UUID) (*Organizer, error)
	GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Organizer, int64, error)
	// Writers
	Create(ctx context.Context, organizer *Organizer) (*Organizer, error)
	Update(ctx context.Context, organizer *Organizer) (*Organizer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
