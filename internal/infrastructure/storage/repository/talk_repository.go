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

type TalkRepository struct {
	queries *dbgen.Queries
}

// NewTalkRepository returns a new TalkRepository
func NewTalkRepository(queries *dbgen.Queries) *TalkRepository {
	return &TalkRepository{queries: queries}
}

type talkDB struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description pgtype.Text
	Tags        []string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	Speakers    []byte
}

// --- READERS ---

// GetAll
func (r *TalkRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Talk, error) {
	rows, err := r.queries.ListTalksByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Talk")
	}

	talks := make([]domain.Talk, len(rows))
	for i, row := range rows {
		talk, err := mapToDomainTalk(talkDB(row))
		if err != nil {
			return nil, err
		}
		talks[i] = *talk
	}

	return talks, nil
}

// GetById
func (r *TalkRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Talk, error) {
	row, err := r.queries.GetTalkByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Talk")
	}

	return mapToDomainTalk(talkDB(row))
}

// --- WRITERS ---

// Create
func (r *TalkRepository) Create(ctx context.Context, talk *domain.Talk) (*domain.Talk, error) {
	row, err := r.queries.CreateTalk(ctx, dbgen.CreateTalkParams{
		EventID:     talk.EventID,
		Title:       talk.Title,
		Description: utils.TextToPgString(talk.Description),
		Tags:        talk.Tags,
		CreatedBy:   talk.Audit.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Talk")
	}

	return mapToDomainTalkSimple(row), nil
}

// Update
func (r *TalkRepository) Update(ctx context.Context, talk *domain.Talk) (*domain.Talk, error) {
	row, err := r.queries.UpdateTalk(ctx, dbgen.UpdateTalkParams{
		ID:          talk.ID,
		Title:       utils.TextToPgString(talk.Title),
		Description: utils.TextToPgString(talk.Description),
		Tags:        talk.Tags,
		UpdatedBy:   talk.Audit.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Talk")
	}

	return mapToDomainTalkSimple(row), nil
}

// Delete
func (r *TalkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteTalk(ctx, id)
	if err != nil {
		return ParseDBError(err, "Talk")
	}

	return nil
}

// --- TALK-SPEAKERS ---

// AddSpeaker
func (r *TalkRepository) AddSpeaker(ctx context.Context, speaker *domain.TalkSpeaker) error {
	err := r.queries.AddSpeakerToTalk(ctx, dbgen.AddSpeakerToTalkParams{
		TalkID:    speaker.TalkID,
		SpeakerID: speaker.SpeakerID,
		CreatedBy: speaker.CreatedBy,
	})
	if err != nil {
		return ParseDBError(err, "TalkSpeaker")
	}

	return nil
}

// RemoveSpeaker
func (r *TalkRepository) RemoveSpeaker(ctx context.Context, speaker *domain.TalkSpeaker) error {
	err := r.queries.RemoveSpeakerFromTalk(ctx, dbgen.RemoveSpeakerFromTalkParams{
		TalkID:    speaker.TalkID,
		SpeakerID: speaker.SpeakerID,
	})
	if err != nil {
		return ParseDBError(err, "TalkSpeaker")
	}

	return nil
}

// --- Mappers ---

func mapToDomainTalk(row talkDB) (*domain.Talk, error) {
	var speakers []domain.SpeakerTalkDetail
	if len(row.Speakers) > 0 {
		if err := json.Unmarshal(row.Speakers, &speakers); err != nil {
			return nil, domain.NewAppError(domain.TypeInternal, "Failed to unmarshal speakers", err)
		}
	}
	return &domain.Talk{
		ID:          row.ID,
		EventID:     row.EventID,
		Title:       row.Title,
		Description: utils.PgStringToText(row.Description),
		Tags:        row.Tags,
		Speakers:    speakers,
		Audit: domain.Audit{
			CreatedAt: row.CreatedAt.Time,
			CreatedBy: row.CreatedBy,
			UpdatedAt: row.UpdatedAt.Time,
			UpdatedBy: row.UpdatedBy,
		},
	}, nil
}

func mapToDomainTalkSimple(row dbgen.Talk) *domain.Talk {
	return &domain.Talk{
		ID:          row.ID,
		EventID:     row.EventID,
		Title:       row.Title,
		Description: utils.PgStringToText(row.Description),
		Tags:        row.Tags,
		Audit: domain.Audit{
			CreatedAt: row.CreatedAt.Time,
			CreatedBy: row.CreatedBy,
			UpdatedAt: row.UpdatedAt.Time,
			UpdatedBy: row.UpdatedBy,
		},
	}
}
