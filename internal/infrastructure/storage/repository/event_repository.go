package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
)

type EventRepository struct {
	queries *dbgen.Queries
}

// NewEventRepository returns a new EventRepository
func NewEventRepository(queries *dbgen.Queries) *EventRepository {
	return &EventRepository{queries: queries}
}

// --- READERS ---

// GetAll returns all Events
func (r *EventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.queries.ListEvents(ctx)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = *mapToDomainEvent(row)
	}

	return events, nil
}

// GetById returns an Event by its ID
func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	row, err := r.queries.GetEventByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	return mapToDomainEvent(row), nil
}

// GetBySlug returns an Event by its slug
func (r *EventRepository) GetBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	row, err := r.queries.GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	return mapToDomainEvent(row), nil
}

// GetActive returns all active Events
func (r *EventRepository) GetActiveList(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.queries.ListActiveEvents(ctx)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = *mapToDomainEvent(row)
	}

	return events, nil
}

// ListPaged returns a page of Events
func (r *EventRepository) ListPaged(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]domain.Event, int64, error) {
	offset := (page - 1) * pageSize

	total, err := r.queries.CountEvents(ctx, search)
	if err != nil {
		return nil, 0, ParseDBError(err, "Event")
	}

	rows, err := r.queries.ListEventsPaged(ctx, dbgen.ListEventsPagedParams{
		Column1: search,
		Limit:   pageSize,
		Offset:  offset,
		Column4: orderBy,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Event")
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = *mapToDomainEvent(row)
	}

	return events, total, nil
}

// --- WRITERS ---

// create inserts a new Event
func (r *EventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	params := dbgen.CreateEventParams{
		Name:      event.Name,
		Slug:      event.Slug,
		IsActive:  ToPgBool(event.IsActive),
		CreatedBy: event.CreatedBy,
	}

	row, err := r.queries.CreateEvent(ctx, params)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}
	return mapToDomainEvent(row), nil
}

// Update updates an Event
func (r *EventRepository) Update(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	params := dbgen.UpdateEventParams{
		ID:        event.ID,
		Name:      event.Name,
		Slug:      event.Slug,
		IsActive:  ToPgBool(event.IsActive),
		UpdatedBy: event.Audit.UpdatedBy,
	}

	row, err := r.queries.UpdateEvent(ctx, params)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}
	return mapToDomainEvent(row), nil
}

// Delete deletes an Event
func (r *EventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEvent(ctx, id)
	if err != nil {
		return ParseDBError(err, "Event")
	}
	return nil
}

// --- Mappers ---

// mapToDomain maps a dbgen.Event to a domain.Event
func mapToDomainEvent(dbEvent dbgen.Event) *domain.Event {
	isActive := dbEvent.IsActive.Bool
	return &domain.Event{
		ID:   dbEvent.ID,
		Name: dbEvent.Name,
		Slug: dbEvent.Slug,

		IsActive: &isActive,

		Audit: domain.Audit{
			CreatedAt: dbEvent.CreatedAt.Time,
			UpdatedAt: dbEvent.UpdatedAt.Time,
			CreatedBy: dbEvent.CreatedBy,
			UpdatedBy: dbEvent.UpdatedBy,
		},
	}
}
