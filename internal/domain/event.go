package domain

import "context"

type Event struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IsActive bool   `json:"is_active"`
	Audit
}

type EventUsecase interface {
	GetEvents(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]Event, int64, error)
	GetEventBySlug(ctx context.Context, slug string) (*Event, error)
	// CreateEvent(ctx context.Context, event *Event) error
}

type EventRepository interface {
	// Read
	// GetByID(ctx context.Context, id string) (*Event, error)
	GetBySlug(ctx context.Context, slug string) (*Event, error)
	// ListAll(ctx context.Context) ([]Event, error)
	// ListActive(ctx context.Context) ([]Event, error)
	ListPaged(ctx context.Context, search string, limit, offset int32, orderBy string) ([]Event, int64, error)

	// Write
	// Create(ctx context.Context, event *Event) (*Event, error)
	// Update(ctx context.Context, event *Event) (*Event, error)
	// Delete(ctx context.Context, id string) error
}
