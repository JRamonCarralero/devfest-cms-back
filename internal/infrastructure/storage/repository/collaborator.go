package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CollaboratorRepository struct {
	queries *dbgen.Queries
}

// NewCollaboratorRepository returns a new CollaboratorRepository
func NewCollaboratorRepository(queries *dbgen.Queries) *CollaboratorRepository {
	return &CollaboratorRepository{queries: queries}
}

// Flat Structs
type collaboratorFlatRow struct {
	ID          uuid.UUID
	PersonID    uuid.UUID
	EventID     uuid.UUID
	Area        pgtype.Text
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

type collaboratorFlat struct {
	ID        uuid.UUID
	PersonID  uuid.UUID
	EventID   uuid.UUID
	Area      pgtype.Text
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

// --- READERS ---

// GetAll returns all collaborators
func (r *CollaboratorRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Collaborator, error) {
	rows, err := r.queries.ListCollaboratorsByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Collaborator")
	}

	collaborators := make([]domain.Collaborator, len(rows))
	for i, row := range rows {
		collaborators[i] = *mapToDomainCollaborator(collaboratorFlatRow(row))
	}

	return collaborators, nil
}

// GetById
func (r *CollaboratorRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Collaborator, error) {
	row, err := r.queries.GetCollaboratorByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Collaborator")
	}

	return mapToDomainCollaborator(collaboratorFlatRow(row)), nil
}

// GetByPersonAndEvent
func (r *CollaboratorRepository) GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error) {
	arg := dbgen.GetCollaboratorByPersonAndEventParams{
		PersonID: personID,
		EventID:  eventID,
	}
	id, err := r.queries.GetCollaboratorByPersonAndEvent(ctx, arg)
	if err != nil {
		return uuid.UUID{}, ParseDBError(err, "Collaborator")
	}

	return id, nil
}

// ListPaged
func (r *CollaboratorRepository) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Collaborator, int64, error) {
	offset := (page - 1) * pageSize

	arg := dbgen.CountCollaboratorsByEventParams{
		EventID: eventID,
		Search:  TextToPgString(search),
	}
	total, err := r.queries.CountCollaboratorsByEvent(ctx, arg)
	if err != nil {
		return nil, 0, ParseDBError(err, "Collaborator")
	}

	rows, err := r.queries.ListCollaboratorsByEventPaged(ctx, dbgen.ListCollaboratorsByEventPagedParams{
		EventID: eventID,
		Search:  TextToPgString(search),
		Limit:   pageSize,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Collaborator")
	}

	collaborators := make([]domain.Collaborator, len(rows))
	for i, row := range rows {
		collaborators[i] = *mapToDomainCollaborator(collaboratorFlatRow(row))
	}

	return collaborators, total, nil
}

// --- WRITERS ---

// Create
func (r *CollaboratorRepository) Create(ctx context.Context, collaborator *domain.Collaborator) (*domain.Collaborator, error) {
	row, err := r.queries.CreateCollaborator(ctx, dbgen.CreateCollaboratorParams{
		EventID:   collaborator.EventID,
		Area:      PtrToText(collaborator.Area),
		PersonID:  collaborator.Person.ID,
		CreatedBy: collaborator.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Collaborator")
	}

	return mapToDomainCollaboratorFlat(collaboratorFlat(row)), nil
}

// Update
func (r *CollaboratorRepository) Update(ctx context.Context, collaborator *domain.Collaborator) (*domain.Collaborator, error) {
	row, err := r.queries.UpdateCollaborator(ctx, dbgen.UpdateCollaboratorParams{
		ID:        collaborator.ID,
		Area:      PtrToText(collaborator.Area),
		UpdatedBy: collaborator.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Collaborator")
	}

	return mapToDomainCollaboratorFlat(collaboratorFlat(row)), nil
}

// Delete
func (r *CollaboratorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteCollaborator(ctx, id)
	if err != nil {
		return ParseDBError(err, "Collaborator")
	}

	return nil
}

// --- Mappers ---

// mapToDomain maps a collaboratorFlatRow to a domain.Collaborator
func mapToDomainCollaborator(dbCollaborator collaboratorFlatRow) *domain.Collaborator {
	return &domain.Collaborator{
		ID:      dbCollaborator.ID,
		EventID: dbCollaborator.EventID,
		Area:    TextToPtr(dbCollaborator.Area),
		Person: domain.Person{
			ID:          dbCollaborator.PersonID,
			FirstName:   dbCollaborator.FirstName,
			LastName:    dbCollaborator.LastName,
			Email:       TextToPtr(dbCollaborator.Email),
			AvatarURL:   TextToPtr(dbCollaborator.AvatarUrl),
			GithubUser:  TextToPtr(dbCollaborator.GithubUser),
			LinkedinURL: TextToPtr(dbCollaborator.LinkedinUrl),
			TwitterURL:  TextToPtr(dbCollaborator.TwitterUrl),
			WebsiteURL:  TextToPtr(dbCollaborator.WebsiteUrl),
		},
		Audit: domain.Audit{
			CreatedAt: dbCollaborator.CreatedAt.Time,
			UpdatedAt: dbCollaborator.UpdatedAt.Time,
			CreatedBy: dbCollaborator.CreatedBy,
			UpdatedBy: dbCollaborator.UpdatedBy,
		},
	}
}

// mapToDomain maps a collaboratorFlat to a domain.Collaborator
func mapToDomainCollaboratorFlat(dbCollaborator collaboratorFlat) *domain.Collaborator {
	return &domain.Collaborator{
		ID:      dbCollaborator.ID,
		EventID: dbCollaborator.EventID,
		Area:    TextToPtr(dbCollaborator.Area),
		Person: domain.Person{
			ID: dbCollaborator.PersonID,
		},
		Audit: domain.Audit{
			CreatedAt: dbCollaborator.CreatedAt.Time,
			UpdatedAt: dbCollaborator.UpdatedAt.Time,
			CreatedBy: dbCollaborator.CreatedBy,
			UpdatedBy: dbCollaborator.UpdatedBy,
		},
	}
}
