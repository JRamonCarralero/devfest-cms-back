package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DeveloperRepository struct {
	queries *dbgen.Queries
}

// NewDeveloperRepository returns a new DeveloperRepository
func NewDeveloperRepository(queries *dbgen.Queries) *DeveloperRepository {
	return &DeveloperRepository{queries: queries}
}

// Flat Structs
type developerFlatRow struct {
	ID              uuid.UUID
	PersonID        uuid.UUID
	EventID         uuid.UUID
	RoleDescription pgtype.Text
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	CreatedBy       uuid.UUID
	UpdatedBy       uuid.UUID
	FirstName       string
	LastName        string
	Email           pgtype.Text
	AvatarUrl       pgtype.Text
	GithubUser      pgtype.Text
	TwitterUrl      pgtype.Text
	LinkedinUrl     pgtype.Text
	WebsiteUrl      pgtype.Text
}

type developerFlat struct {
	ID              uuid.UUID
	PersonID        uuid.UUID
	EventID         uuid.UUID
	RoleDescription pgtype.Text
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	CreatedBy       uuid.UUID
	UpdatedBy       uuid.UUID
}

// --- READERS ---

// GetAll
func (r *DeveloperRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Developer, error) {
	rows, err := r.queries.ListDevelopersByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Developer")
	}

	developers := make([]domain.Developer, len(rows))
	for i, row := range rows {
		developers[i] = *mapToDomainDeveloper(developerFlatRow(row))
	}

	return developers, nil
}

// GetById
func (r *DeveloperRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Developer, error) {
	row, err := r.queries.GetDeveloperByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Developer")
	}

	return mapToDomainDeveloper(developerFlatRow(row)), nil
}

// GetByPersonAndEvent
func (r *DeveloperRepository) GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error) {
	arg := dbgen.GetDeveloperByPersonAndEventParams{
		PersonID: personID,
		EventID:  eventID,
	}
	id, err := r.queries.GetDeveloperByPersonAndEvent(ctx, arg)
	if err != nil {
		return uuid.UUID{}, ParseDBError(err, "Developer")
	}

	return id, nil
}

// ListPaged
func (r *DeveloperRepository) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Developer, int64, error) {
	offset := (page - 1) * pageSize

	arg := dbgen.CountDevelopersByEventParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
	}
	total, err := r.queries.CountDevelopersByEvent(ctx, arg)
	if err != nil {
		return nil, 0, ParseDBError(err, "Developer")
	}

	rows, err := r.queries.ListDevelopersByEventPaged(ctx, dbgen.ListDevelopersByEventPagedParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
		Limit:   pageSize,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Developer")
	}

	developers := make([]domain.Developer, len(rows))
	for i, row := range rows {
		developers[i] = *mapToDomainDeveloper(developerFlatRow(row))
	}

	return developers, total, nil
}

// --- WRITERS ---

// Create
func (r *DeveloperRepository) Create(ctx context.Context, developer *domain.Developer) (*domain.Developer, error) {
	row, err := r.queries.CreateDeveloper(ctx, dbgen.CreateDeveloperParams{
		EventID:         developer.EventID,
		PersonID:        developer.Person.ID,
		RoleDescription: utils.PtrToText(developer.RoleDescription),
		CreatedBy:       developer.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Developer")
	}

	return mapToDomainDeveloperFlat(developerFlat(row)), nil
}

// Update
func (r *DeveloperRepository) Update(ctx context.Context, developer *domain.Developer) (*domain.Developer, error) {
	row, err := r.queries.UpdateDeveloper(ctx, dbgen.UpdateDeveloperParams{
		ID:              developer.ID,
		RoleDescription: utils.PtrToText(developer.RoleDescription),
		UpdatedBy:       developer.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Developer")
	}

	return mapToDomainDeveloperFlat(developerFlat(row)), nil
}

// Delete
func (r *DeveloperRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteDeveloper(ctx, id)
	if err != nil {
		return ParseDBError(err, "Developer")
	}

	return nil
}

// --- MAPPERS ---

// mapToDomain maps a developerFlatRow to a domain.Developer
func mapToDomainDeveloper(dbDeveloper developerFlatRow) *domain.Developer {
	return &domain.Developer{
		ID:              dbDeveloper.ID,
		EventID:         dbDeveloper.EventID,
		RoleDescription: utils.TextToPtr(dbDeveloper.RoleDescription),
		Person: domain.Person{
			ID:          dbDeveloper.PersonID,
			FirstName:   dbDeveloper.FirstName,
			LastName:    dbDeveloper.LastName,
			Email:       utils.TextToPtr(dbDeveloper.Email),
			AvatarURL:   utils.TextToPtr(dbDeveloper.AvatarUrl),
			GithubUser:  utils.TextToPtr(dbDeveloper.GithubUser),
			LinkedinURL: utils.TextToPtr(dbDeveloper.LinkedinUrl),
			TwitterURL:  utils.TextToPtr(dbDeveloper.TwitterUrl),
			WebsiteURL:  utils.TextToPtr(dbDeveloper.WebsiteUrl),
		},
		Audit: domain.Audit{
			CreatedAt: dbDeveloper.CreatedAt.Time,
			UpdatedAt: dbDeveloper.UpdatedAt.Time,
			CreatedBy: dbDeveloper.CreatedBy,
			UpdatedBy: dbDeveloper.UpdatedBy,
		},
	}
}

// mapToDomain maps a developerFlat to a domain.Developer
func mapToDomainDeveloperFlat(dbDeveloper developerFlat) *domain.Developer {
	return &domain.Developer{
		ID:              dbDeveloper.ID,
		EventID:         dbDeveloper.EventID,
		RoleDescription: utils.TextToPtr(dbDeveloper.RoleDescription),
		Person: domain.Person{
			ID: dbDeveloper.PersonID,
		},
		Audit: domain.Audit{
			CreatedAt: dbDeveloper.CreatedAt.Time,
			UpdatedAt: dbDeveloper.UpdatedAt.Time,
			CreatedBy: dbDeveloper.CreatedBy,
			UpdatedBy: dbDeveloper.UpdatedBy,
		},
	}
}
