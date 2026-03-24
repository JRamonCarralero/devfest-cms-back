package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrganizerRepository struct {
	queries *dbgen.Queries
}

// NewOrganizerRepository returns a new OrganizerRepository
func NewOrganizerRepository(queries *dbgen.Queries) *OrganizerRepository {
	return &OrganizerRepository{queries: queries}
}

// Flat Structs
type organizerFlatRow struct {
	ID              uuid.UUID
	PersonID        uuid.UUID
	EventID         uuid.UUID
	Company         pgtype.Text
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

type organizerFlat struct {
	ID              uuid.UUID
	PersonID        uuid.UUID
	EventID         uuid.UUID
	Company         pgtype.Text
	RoleDescription pgtype.Text
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	CreatedBy       uuid.UUID
	UpdatedBy       uuid.UUID
}

// --- READERS ---

// GetAll
func (r *OrganizerRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Organizer, error) {
	rows, err := r.queries.ListOrganizersByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Organizer")
	}

	organizers := make([]domain.Organizer, len(rows))
	for i, row := range rows {
		organizers[i] = *mapToDomainOrganizer(organizerFlatRow(row))
	}

	return organizers, nil
}

// GetById
func (r *OrganizerRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Organizer, error) {
	row, err := r.queries.GetOrganizerByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Organizer")
	}

	return mapToDomainOrganizer(organizerFlatRow(row)), nil
}

// GetByPersonAndEvent
func (r *OrganizerRepository) GetByPersonAndEvent(ctx context.Context, personID uuid.UUID, eventID uuid.UUID) (uuid.UUID, error) {
	arg := dbgen.GetOrganizerByPersonAndEventParams{
		PersonID: personID,
		EventID:  eventID,
	}
	id, err := r.queries.GetOrganizerByPersonAndEvent(ctx, arg)
	if err != nil {
		return uuid.UUID{}, ParseDBError(err, "Organizer")
	}

	return id, nil
}

// ListPaged
func (r *OrganizerRepository) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Organizer, int64, error) {
	offset := (page - 1) * pageSize

	arg := dbgen.CountOrganizersByEventParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
	}
	total, err := r.queries.CountOrganizersByEvent(ctx, arg)
	if err != nil {
		return nil, 0, ParseDBError(err, "Organizer")
	}

	rows, err := r.queries.ListOrganizersByEventPaged(ctx, dbgen.ListOrganizersByEventPagedParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
		Limit:   pageSize,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Organizer")
	}

	organizers := make([]domain.Organizer, len(rows))
	for i, row := range rows {
		organizers[i] = *mapToDomainOrganizer(organizerFlatRow(row))
	}

	return organizers, total, nil
}

// --- WRITERS ---

// Create
func (r *OrganizerRepository) Create(ctx context.Context, organizer *domain.Organizer) (*domain.Organizer, error) {
	row, err := r.queries.CreateOrganizer(ctx, dbgen.CreateOrganizerParams{
		PersonID:        organizer.Person.ID,
		EventID:         organizer.EventID,
		Company:         utils.PtrToText(organizer.Company),
		RoleDescription: utils.PtrToText(organizer.RoleDescription),
		CreatedBy:       organizer.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Organizer")
	}

	return mapToDomainOrganizerFlat(organizerFlat(row)), nil
}

// Update
func (r *OrganizerRepository) Update(ctx context.Context, organizer *domain.Organizer) (*domain.Organizer, error) {
	row, err := r.queries.UpdateOrganizer(ctx, dbgen.UpdateOrganizerParams{
		ID:              organizer.ID,
		Company:         utils.PtrToText(organizer.Company),
		RoleDescription: utils.PtrToText(organizer.RoleDescription),
		UpdatedBy:       organizer.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Organizer")
	}

	return mapToDomainOrganizerFlat(organizerFlat(row)), nil
}

// Delete
func (r *OrganizerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteOrganizer(ctx, id)
	if err != nil {
		return ParseDBError(err, "Organizer")
	}

	return nil
}

// --- Mappers ---

// mapToDomain maps a organizerFlatRow to a domain.Organizer
func mapToDomainOrganizer(organizerFlatRow organizerFlatRow) *domain.Organizer {
	return &domain.Organizer{
		ID:              organizerFlatRow.ID,
		EventID:         organizerFlatRow.EventID,
		Company:         utils.TextToPtr(organizerFlatRow.Company),
		RoleDescription: utils.TextToPtr(organizerFlatRow.RoleDescription),
		Person: domain.Person{
			ID:          organizerFlatRow.PersonID,
			FirstName:   organizerFlatRow.FirstName,
			LastName:    organizerFlatRow.LastName,
			Email:       utils.TextToPtr(organizerFlatRow.Email),
			AvatarURL:   utils.TextToPtr(organizerFlatRow.AvatarUrl),
			GithubUser:  utils.TextToPtr(organizerFlatRow.GithubUser),
			TwitterURL:  utils.TextToPtr(organizerFlatRow.TwitterUrl),
			LinkedinURL: utils.TextToPtr(organizerFlatRow.LinkedinUrl),
			WebsiteURL:  utils.TextToPtr(organizerFlatRow.WebsiteUrl),
			Audit: domain.Audit{
				CreatedAt: organizerFlatRow.CreatedAt.Time,
				UpdatedAt: organizerFlatRow.UpdatedAt.Time,
				CreatedBy: organizerFlatRow.CreatedBy,
				UpdatedBy: organizerFlatRow.UpdatedBy,
			},
		},
	}
}

// mapToDomain maps a organizerFlat to a domain.Organizer
func mapToDomainOrganizerFlat(organizerFlat organizerFlat) *domain.Organizer {
	return &domain.Organizer{
		ID:              organizerFlat.ID,
		EventID:         organizerFlat.EventID,
		Company:         utils.TextToPtr(organizerFlat.Company),
		RoleDescription: utils.TextToPtr(organizerFlat.RoleDescription),
		Person: domain.Person{
			ID: organizerFlat.PersonID,
		},
		Audit: domain.Audit{
			CreatedAt: organizerFlat.CreatedAt.Time,
			UpdatedAt: organizerFlat.UpdatedAt.Time,
			CreatedBy: organizerFlat.CreatedBy,
			UpdatedBy: organizerFlat.UpdatedBy,
		},
	}
}
