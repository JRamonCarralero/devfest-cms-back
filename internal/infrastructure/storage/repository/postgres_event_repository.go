package repository

import (
	"context"
	"database/sql"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/storage/dbgen"
)

type PostgresEventRepository struct {
	queries *dbgen.Queries
}

// NewPostgresEventRepository returns a new PostgresEventRepository
func NewPostgresEventRepository(queries *dbgen.Queries) *PostgresEventRepository {
	return &PostgresEventRepository{queries: queries}
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

	event := mapToDomain(row)
	return &event, nil
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
		events[i] = mapToDomain(row)
	}

	return events, total, nil
}

// --- Mappers ---

// mapToDomain maps a dbgen.Event to a domain.Event
func mapToDomain(row dbgen.Event) domain.Event {
	return domain.Event{
		ID:       row.ID.String(),
		Name:     row.Name,
		Slug:     row.Slug,
		IsActive: row.IsActive.Bool,
		Audit: domain.Audit{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			CreatedBy: uuidPtrToString(row.CreatedBy),
			UpdatedBy: uuidPtrToString(row.UpdatedBy),
		},
	}
}
