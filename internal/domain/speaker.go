package domain

import (
	"context"

	"github.com/google/uuid"
)

type Speaker struct {
	ID      uuid.UUID
	EventID uuid.UUID
	Person
	Bio     *string
	Company *string
	Audit
}

type UpdateSpeaker struct {
	Bio       *string
	Company   *string
	UpdatedBy uuid.UUID
}

type SpeakerRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Speaker, error)
	GetById(ctx context.Context, id uuid.UUID) (*Speaker, error)
	GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Speaker, int64, error)
	// Writers
	Create(ctx context.Context, speaker *Speaker) (*Speaker, error)
	Update(ctx context.Context, speaker *Speaker) (*Speaker, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
