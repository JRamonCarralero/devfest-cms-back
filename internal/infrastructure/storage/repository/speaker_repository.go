package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SpeakerRepository struct {
	queries *dbgen.Queries
}

// NewSpeakerRepository returns a new SpeakerRepository
func NewSpeakerRepository(queries *dbgen.Queries) *SpeakerRepository {
	return &SpeakerRepository{queries: queries}
}

// Flat Structs

type speakerFlatRow struct {
	ID          uuid.UUID
	PersonID    uuid.UUID
	EventID     uuid.UUID
	Bio         pgtype.Text
	Company     pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	FirstName   string
	LastName    string
	Email       pgtype.Text
	AvatarUrl   pgtype.Text
	GithubUser  pgtype.Text
	TwitterUrl  pgtype.Text
	LinkedinUrl pgtype.Text
	WebsiteUrl  pgtype.Text
}

type speakerFlat struct {
	ID        uuid.UUID
	PersonID  uuid.UUID
	EventID   uuid.UUID
	Bio       pgtype.Text
	Company   pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

// --- READERS ---

// GetAll
func (r *SpeakerRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Speaker, error) {
	rows, err := r.queries.ListSpeakersByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Speaker")
	}

	speakers := make([]domain.Speaker, len(rows))
	for i, row := range rows {
		speakers[i] = *mapToDomainSpeaker(speakerFlatRow(row))
	}

	return speakers, nil
}

// GetById
func (r *SpeakerRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Speaker, error) {
	row, err := r.queries.GetSpeakerByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Speaker")
	}

	return mapToDomainSpeaker(speakerFlatRow(row)), nil
}

// GetByPersonAndEvent
func (r *SpeakerRepository) GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error) {
	arg := dbgen.GetSpeakerByPersonAndEventParams{
		PersonID: personID,
		EventID:  eventID,
	}
	id, err := r.queries.GetSpeakerByPersonAndEvent(ctx, arg)
	if err != nil {
		return uuid.UUID{}, ParseDBError(err, "Speaker")
	}

	return id, nil
}

// ListPaged
func (r *SpeakerRepository) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Speaker, int64, error) {
	offset := (page - 1) * pageSize

	arg := dbgen.CountSpeakersByEventParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
	}
	total, err := r.queries.CountSpeakersByEvent(ctx, arg)
	if err != nil {
		return nil, 0, ParseDBError(err, "Speaker")
	}

	rows, err := r.queries.ListSpeakersByEventPaged(ctx, dbgen.ListSpeakersByEventPagedParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
		Limit:   pageSize,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Speaker")
	}

	speakers := make([]domain.Speaker, len(rows))
	for i, row := range rows {
		speakers[i] = *mapToDomainSpeaker(speakerFlatRow(row))
	}

	return speakers, total, nil
}

// --- WRITERS ---

// Create
func (r *SpeakerRepository) Create(ctx context.Context, speaker *domain.Speaker) (*domain.Speaker, error) {
	row, err := r.queries.CreateSpeaker(ctx, dbgen.CreateSpeakerParams{
		PersonID:  speaker.Person.ID,
		EventID:   speaker.EventID,
		Bio:       utils.PtrToText(speaker.Bio),
		Company:   utils.PtrToText(speaker.Company),
		CreatedBy: speaker.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Speaker")
	}

	return mapToDomainSpeakerFlat(speakerFlat(row)), nil
}

// Update
func (r *SpeakerRepository) Update(ctx context.Context, speaker *domain.Speaker) (*domain.Speaker, error) {
	row, err := r.queries.UpdateSpeaker(ctx, dbgen.UpdateSpeakerParams{
		ID:        speaker.ID,
		Bio:       utils.PtrToText(speaker.Bio),
		Company:   utils.PtrToText(speaker.Company),
		UpdatedBy: speaker.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Speaker")
	}

	return mapToDomainSpeakerFlat(speakerFlat(row)), nil
}

// Delete
func (r *SpeakerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteSpeaker(ctx, id)
	if err != nil {
		return ParseDBError(err, "Speaker")
	}

	return nil
}

// --- Mappers ---

// mapToDomain maps a speakerFlatRow to a domain.Speaker
func mapToDomainSpeaker(speakerFlatRow speakerFlatRow) *domain.Speaker {
	return &domain.Speaker{
		ID:      speakerFlatRow.ID,
		EventID: speakerFlatRow.EventID,
		Bio:     utils.TextToPtr(speakerFlatRow.Bio),
		Company: utils.TextToPtr(speakerFlatRow.Company),
		Person: domain.Person{
			ID:          speakerFlatRow.PersonID,
			FirstName:   speakerFlatRow.FirstName,
			LastName:    speakerFlatRow.LastName,
			Email:       utils.TextToPtr(speakerFlatRow.Email),
			AvatarURL:   utils.TextToPtr(speakerFlatRow.AvatarUrl),
			GithubUser:  utils.TextToPtr(speakerFlatRow.GithubUser),
			TwitterURL:  utils.TextToPtr(speakerFlatRow.TwitterUrl),
			LinkedinURL: utils.TextToPtr(speakerFlatRow.LinkedinUrl),
			WebsiteURL:  utils.TextToPtr(speakerFlatRow.WebsiteUrl),
		},
		Audit: domain.Audit{
			CreatedAt: speakerFlatRow.CreatedAt.Time,
			UpdatedAt: speakerFlatRow.UpdatedAt.Time,
			CreatedBy: speakerFlatRow.CreatedBy,
			UpdatedBy: speakerFlatRow.UpdatedBy,
		},
	}
}

// mapToDomain maps a speakerFlat to a domain.Speaker
func mapToDomainSpeakerFlat(speakerFlat speakerFlat) *domain.Speaker {
	return &domain.Speaker{
		ID:      speakerFlat.ID,
		EventID: speakerFlat.EventID,
		Bio:     utils.TextToPtr(speakerFlat.Bio),
		Company: utils.TextToPtr(speakerFlat.Company),
		Person: domain.Person{
			ID: speakerFlat.PersonID,
		},
		Audit: domain.Audit{
			CreatedAt: speakerFlat.CreatedAt.Time,
			UpdatedAt: speakerFlat.UpdatedAt.Time,
			CreatedBy: speakerFlat.CreatedBy,
			UpdatedBy: speakerFlat.UpdatedBy,
		},
	}
}
