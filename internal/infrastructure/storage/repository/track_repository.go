package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"
	"encoding/json"

	"github.com/google/uuid"
)

type TrackRepository struct {
	queries *dbgen.Queries
}

// NewTrackRepository returns a new TrackRepository
func NewTrackRepository(queries *dbgen.Queries) *TrackRepository {
	return &TrackRepository{queries: queries}
}

// --- READERS ---

// GetAll
func (r *TrackRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Track, error) {
	rows, err := r.queries.ListTracksByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Track")
	}

	tracks := make([]domain.Track, len(rows))
	for i, row := range rows {
		tracks[i] = *mapToDomainTrack(row)
	}

	return tracks, nil
}

// GetById
func (r *TrackRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Track, error) {
	row, err := r.queries.GetTrackByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Track")
	}

	return mapToDomainTrack(row), nil
}

// GetFullEventSchedule
func (r *TrackRepository) GetFullEventSchedule(ctx context.Context, eventID uuid.UUID) ([]domain.FullTrackSchedule, error) {
	rows, err := r.queries.GetFullEventSchedule(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Track")
	}

	tracks, err := mapToDomainFullTrack(rows)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

// --- WRITERS ---

// Create
func (r *TrackRepository) Create(ctx context.Context, track *domain.Track) (*domain.Track, error) {
	row, err := r.queries.CreateTrack(ctx, dbgen.CreateTrackParams{
		Name:      track.Name,
		EventID:   track.EventID,
		EventDate: utils.ToPgDate(track.EventDate),
	})
	if err != nil {
		return nil, ParseDBError(err, "Track")
	}

	return mapToDomainTrack(row), nil
}

// Update
func (r *TrackRepository) Update(ctx context.Context, track *domain.Track) (*domain.Track, error) {
	row, err := r.queries.UpdateTrack(ctx, dbgen.UpdateTrackParams{
		ID:        track.ID,
		Name:      utils.TextToPgString(track.Name),
		EventDate: utils.ToPgDate(track.EventDate),
		UpdatedBy: track.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Track")
	}

	return mapToDomainTrack(row), nil
}

// Delete
func (r *TrackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteTrack(ctx, id)
	if err != nil {
		return ParseDBError(err, "Track")
	}

	return nil
}

// --- Mappers ---

// mapToDomainTrack
func mapToDomainTrack(row dbgen.Track) *domain.Track {
	return &domain.Track{
		ID:        row.ID,
		Name:      row.Name,
		EventID:   row.EventID,
		EventDate: row.EventDate.Time,
		Audit: domain.Audit{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			CreatedBy: row.CreatedBy,
			UpdatedBy: row.UpdatedBy,
		},
	}
}

// mapToDomainFullTrack
func mapToDomainFullTrack(rows []dbgen.GetFullEventScheduleRow) ([]domain.FullTrackSchedule, error) {
	var schedule []domain.FullTrackSchedule

	for _, row := range rows {
		var domainEntries []domain.ScheduleEntryTrack

		if len(row.Entries) > 0 {
			if err := json.Unmarshal(row.Entries, &domainEntries); err != nil {
				return nil, domain.NewAppError(domain.TypeInternal, "Failed to unmarshal schedule entries", err)
			}
		}

		trackSchedule := domain.FullTrackSchedule{
			TrackID:   row.TrackID,
			TrackName: row.TrackName,
			EventDate: row.EventDate.Time,
			Entries:   domainEntries,
		}

		schedule = append(schedule, trackSchedule)
	}

	return schedule, nil
}
