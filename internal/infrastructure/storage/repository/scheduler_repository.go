package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetScheduleEntryDB struct {
	ID              uuid.UUID          `json:"id"`
	TrackID         uuid.UUID          `json:"track_id"`
	TalkID          pgtype.UUID        `json:"talk_id"`
	StartTime       pgtype.Time        `json:"start_time"`
	EndTime         pgtype.Time        `json:"end_time"`
	Duration        pgtype.Interval    `json:"duration"`
	Room            pgtype.Text        `json:"room"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	TrackName       string             `json:"track_name"`
	EventDate       pgtype.Date        `json:"event_date"`
	TalkTitle       pgtype.Text        `json:"talk_title"`
	TalkDescription pgtype.Text        `json:"talk_description"`
	Speakers        []byte             `json:"speakers"`
}

type SchedulerRepository struct {
	queries *dbgen.Queries
}

// NewSchedulerRepository returns a new SchedulerRepository
func NewSchedulerRepository(queries *dbgen.Queries) *SchedulerRepository {
	return &SchedulerRepository{queries: queries}
}

// --- READERS ---

// GetAllByTrack
func (r *SchedulerRepository) GetAllByTrack(ctx context.Context, trackID uuid.UUID) ([]domain.Scheduler, error) {
	rows, err := r.queries.ListScheduleByTrack(ctx, trackID)
	if err != nil {
		return nil, ParseDBError(err, "Scheduler")
	}

	schedulers := make([]domain.Scheduler, len(rows))
	for i, row := range rows {
		schedulers[i] = *mapToDomainScheduler(row)
	}

	return schedulers, nil
}

// GetTrackById
func (r *SchedulerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Scheduler, error) {
	row, err := r.queries.GetScheduleEntryByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Scheduler")
	}

	return mapToDomainSchedulerRow(row), nil
}

// --- WRITERS ---

// Create
func (r *SchedulerRepository) Create(ctx context.Context, scheduler *domain.Scheduler) (*domain.Scheduler, error) {
	row, err := r.queries.CreateScheduleEntry(ctx, dbgen.CreateScheduleEntryParams{
		StartTime: utils.TimeToPgTime(scheduler.StartTime),
		EndTime:   utils.TimeToPgTime(scheduler.EndTime),
		Room:      utils.TextToPgString(scheduler.Room),
		TrackID:   scheduler.Track.ID,
		TalkID:    utils.UUIDToPgUUID(scheduler.Talk.ID),
	})
	if err != nil {
		return nil, ParseDBError(err, "Scheduler")
	}

	return mapToDomainSchedulerFlat(row), nil
}

// Update
func (r *SchedulerRepository) Update(ctx context.Context, scheduler *domain.Scheduler) (*domain.Scheduler, error) {
	row, err := r.queries.UpdateScheduleEntry(ctx, dbgen.UpdateScheduleEntryParams{
		ID:        scheduler.ID,
		StartTime: utils.TimeToPgTime(scheduler.StartTime),
		EndTime:   utils.TimeToPgTime(scheduler.EndTime),
		Room:      utils.TextToPgString(scheduler.Room),
		TrackID:   utils.UUIDToPgUUID(scheduler.Track.ID),
		TalkID:    utils.UUIDToPgUUID(scheduler.Talk.ID),
	})
	if err != nil {
		return nil, ParseDBError(err, "Scheduler")
	}

	return mapToDomainSchedulerFlat(row), nil
}

// Delete
func (r *SchedulerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteScheduleEntry(ctx, id)
	if err != nil {
		return ParseDBError(err, "Scheduler")
	}
	return nil
}

// --- Mappers ---

func mapToDomainScheduler(row dbgen.ListScheduleByTrackRow) *domain.Scheduler {
	var speakers []domain.SpeakerTalkDetail
	if len(row.Speakers) > 0 {
		if err := json.Unmarshal(row.Speakers, &speakers); err != nil {
			speakers = []domain.SpeakerTalkDetail{}
		}
	}

	return &domain.Scheduler{
		ID:        row.ID,
		StartTime: utils.PgTimeToTime(row.StartTime),
		EndTime:   utils.PgTimeToTime(row.EndTime),
		Duration:  utils.PgIntervalToDuration(row.Duration),
		Room:      utils.PgStringToText(row.Room),
		Track: domain.Track{
			ID: row.TrackID,
		},
		Talk: domain.Talk{
			ID:       utils.PgUUIDToUUID(row.TalkID),
			Speakers: speakers,
		},
		Audit: domain.Audit{
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}

func mapToDomainSchedulerRow(row dbgen.GetScheduleEntryByIDRow) *domain.Scheduler {
	var speakers []domain.SpeakerTalkDetail
	if len(row.Speakers) > 0 {
		if err := json.Unmarshal(row.Speakers, &speakers); err != nil {
			speakers = []domain.SpeakerTalkDetail{}
		}
	}

	return &domain.Scheduler{
		ID:        row.ID,
		StartTime: utils.PgTimeToTime(row.StartTime),
		EndTime:   utils.PgTimeToTime(row.EndTime),
		Duration:  utils.PgIntervalToDuration(row.Duration),
		Room:      utils.PgStringToText(row.Room),
		Track: domain.Track{
			ID:   row.TrackID,
			Name: row.TrackName,
		},
		Talk: domain.Talk{
			ID:          utils.PgUUIDToUUID(row.TalkID),
			Title:       utils.PgStringToText(row.TalkTitle),
			Description: utils.PgStringToText(row.TalkDescription),
			Speakers:    speakers,
		},
		Audit: domain.Audit{
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}

func mapToDomainSchedulerFlat(row dbgen.Scheduler) *domain.Scheduler {
	return &domain.Scheduler{
		ID:        row.ID,
		StartTime: utils.PgTimeToTime(row.StartTime),
		EndTime:   utils.PgTimeToTime(row.EndTime),
		Duration:  utils.PgIntervalToDuration(row.Duration),
		Room:      utils.PgStringToText(row.Room),
		Track: domain.Track{
			ID: row.TrackID,
		},
		Talk: domain.Talk{
			ID: utils.PgUUIDToUUID(row.TalkID),
		},
		Audit: domain.Audit{
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}
