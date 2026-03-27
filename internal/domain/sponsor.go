package domain

import (
	"context"

	"github.com/google/uuid"
)

type Sponsor struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	Name          string
	LogoURL       string
	WebsiteURL    string
	Tier          string
	OrderPriority *int
	Audit
}

type UpdateSponsor struct {
	Name          *string
	LogoURL       *string
	WebsiteURL    *string
	Tier          *string
	OrderPriority *int
	UpdatedBy     uuid.UUID
}

type SponsorUsecase interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Sponsor, error)
	GetById(ctx context.Context, id uuid.UUID) (*Sponsor, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Sponsor, int64, error)
	// Writers
	Create(ctx context.Context, sponsor *Sponsor) (*Sponsor, error)
	Update(ctx context.Context, id uuid.UUID, upSponsor *UpdateSponsor) (*Sponsor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SponsorRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Sponsor, error)
	GetById(ctx context.Context, id uuid.UUID) (*Sponsor, error)
	ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]Sponsor, int64, error)
	// Writers
	Create(ctx context.Context, sponsor *Sponsor) (*Sponsor, error)
	Update(ctx context.Context, sponsor *Sponsor) (*Sponsor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
