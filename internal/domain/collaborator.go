package domain

import (
	"context"

	"github.com/google/uuid"
)

type Collaborator struct {
	ID      uuid.UUID
	EventID uuid.UUID
	Area    *string
	Person
	Audit
}

type UpdateCollaborator struct {
	Area      *string
	UpdatedBy uuid.UUID
}

type CollaboratorUsecase interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Collaborator, error)
	GetById(ctx context.Context, id uuid.UUID) (*Collaborator, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Collaborator, int64, error)
	// Writers
	Create(ctx context.Context, collaborator *Collaborator) (*Collaborator, error)
	Update(ctx context.Context, id uuid.UUID, collaborator *UpdateCollaborator) (*Collaborator, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CollaboratorRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Collaborator, error)
	GetById(ctx context.Context, id uuid.UUID) (*Collaborator, error)
	GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Collaborator, int64, error)
	// Writers
	Create(ctx context.Context, collaborator *Collaborator) (*Collaborator, error)
	Update(ctx context.Context, collaborator *Collaborator) (*Collaborator, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
