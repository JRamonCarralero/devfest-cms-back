package domain

import (
	"context"
	"devfest/internal/infrastructure/api/dtos"

	"github.com/google/uuid"
)

type Event struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	IsActive *bool     `json:"is_active"`
	Audit
}

type EventUsecase interface {
	GetEvents(ctx context.Context) ([]Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetEventBySlug(ctx context.Context, slug string) (*Event, error)
	GetActiveEvents(ctx context.Context) ([]Event, error)
	GetEventsPaged(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]Event, int64, error)
	CreateEvent(ctx context.Context, dto dtos.CreateEventDTO) (*Event, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, dto dtos.UpdateEventDTO) (*Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type EventRepository interface {
	// Readers
	GetAll(ctx context.Context) ([]Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetBySlug(ctx context.Context, slug string) (*Event, error)
	GetActiveList(ctx context.Context) ([]Event, error)
	ListPaged(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]Event, int64, error)
	// Writers
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
