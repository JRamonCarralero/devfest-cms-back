package repository

import (
	"context"
	"database/sql"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
)

type PostgresEventRepository struct {
	queries *dbgen.Queries
}

// NewPostgresEventRepository returns a new PostgresEventRepository
func NewPostgresEventRepository(queries *dbgen.Queries) *PostgresEventRepository {
	return &PostgresEventRepository{queries: queries}
}

// --- READERS ---

// GetAll returns all Events
func (r *PostgresEventRepository) GetAll(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.queries.ListEvents(ctx)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = *mapToDomain(row)
	}

	return events, nil
}

// GetById returns an Event by its ID
func (r *PostgresEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	row, err := r.queries.GetEventByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, ParseDBError(err, "Event")
	}

	event := *mapToDomain(row)
	return &event, nil
}

// GetBySlug returns an Event by its slug
func (r *PostgresEventRepository) GetBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	row, err := r.queries.GetEventBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, ParseDBError(err, "Event")
	}

	event := *mapToDomain(row)
	return &event, nil
}

// GetActive returns all active Events
func (r *PostgresEventRepository) GetActive(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.queries.ListActiveEvents(ctx)
	if err != nil {
		return nil, ParseDBError(err, "Event")
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = *mapToDomain(row)
	}

	return events, nil
}

// ListPaged returns a page of Events
func (r *PostgresEventRepository) ListPaged(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]domain.Event, int64, error) {
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
		events[i] = *mapToDomain(row)
	}

	return events, total, nil
}

// --- WRITERS ---

// create inserts a new Event
func (r *PostgresEventRepository) Create(ctx context.Context, event *domain.Event) (dbgen.Event, error) {
	params := dbgen.CreateEventParams{
		Name:      event.Name,
		Slug:      event.Slug,
		IsActive:  ToPgBool(event.IsActive),
		CreatedBy: event.CreatedBy,
	}

	return r.queries.CreateEvent(ctx, params)
}

// ToDo Update and Delete

// --- Mappers ---

// mapToDomain maps a dbgen.Event to a domain.Event
func mapToDomain(dbEvent dbgen.Event) *domain.Event {
	return &domain.Event{
		ID:   dbEvent.ID, // Este ya es uuid.UUID (funciona)
		Name: dbEvent.Name,
		Slug: dbEvent.Slug,

		IsActive: &dbEvent.IsActive.Bool,

		Audit: domain.Audit{
			CreatedAt: dbEvent.CreatedAt.Time,
			UpdatedAt: dbEvent.UpdatedAt.Time,
			CreatedBy: dbEvent.CreatedBy,
			UpdatedBy: dbEvent.UpdatedBy,
		},
	}
}
