package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// eventInteractor implements domain.EventUsecase
type eventInteractor struct {
	repo domain.EventRepository
}

// NewEventInteractor is a constructor for eventInteractor
func NewEventInteractor(repo domain.EventRepository) domain.EventUsecase {
	return &eventInteractor{
		repo: repo,
	}
}

// --- Readers ---

// GetEvents returns all events
func (i *eventInteractor) GetEvents(ctx context.Context) ([]domain.Event, error) {
	return i.repo.GetAll(ctx)
}

// GetByID returns an Event by its ID
func (i *eventInteractor) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return i.repo.GetByID(ctx, id)
}

// GetEventBySlug validates slug and returns an Event by its slug
func (i *eventInteractor) GetEventBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	if slug == "" {
		appErr := domain.NewAppError(domain.TypeBadRequest, "event slug is required", nil)
		return nil, appErr
	}

	event, err := i.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if event == nil {
		appErr := domain.NewAppError(domain.TypeNotFound, "event not found", nil)
		return nil, appErr
	}

	return event, nil
}

// GetActiveEvents returns all active Events
func (i *eventInteractor) GetActiveEvents(ctx context.Context) ([]domain.Event, error) {
	return i.repo.GetActiveList(ctx)
}

// GetEventsPaged validates params and returns a page of Events
func (i *eventInteractor) GetEventsPaged(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]domain.Event, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	validOrders := map[string]bool{
		"name_asc":        true,
		"name_desc":       true,
		"created_at_asc":  true,
		"created_at_desc": true,
	}
	if !validOrders[orderBy] {
		orderBy = "created_at_desc"
	}

	return i.repo.ListPaged(ctx, search, page, pageSize, orderBy)
}

// --- Writers ---

// CreateEvent creates a new Event
func (i *eventInteractor) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	createdEvent, err := i.repo.Create(ctx, event)

	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}

// UpdateEvent validates params and updates an Event
func (i *eventInteractor) UpdateEvent(ctx context.Context, id uuid.UUID, evUpdate *domain.UpdateEvent) (*domain.Event, error) {
	event, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if evUpdate.Name != nil {
		event.Name = *evUpdate.Name
	}
	if evUpdate.Slug != nil {
		event.Slug = *evUpdate.Slug
	}
	if evUpdate.IsActive != nil {
		event.IsActive = evUpdate.IsActive
	}

	event.Audit.UpdatedBy = evUpdate.UpdatedBy

	updatedEvent, err := i.repo.Update(ctx, event)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// DeleteEvent deletes an Event by its ID
func (i *eventInteractor) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return i.repo.Delete(ctx, id)
}
