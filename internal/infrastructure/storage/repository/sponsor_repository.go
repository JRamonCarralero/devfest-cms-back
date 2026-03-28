package repository

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/utils"
	"devfest/internal/infrastructure/storage/dbgen"

	"github.com/google/uuid"
)

type SponsorRepository struct {
	queries *dbgen.Queries
}

// NewSponsorRepository returns a new SponsorRepository
func NewSponsorRepository(queries *dbgen.Queries) *SponsorRepository {
	return &SponsorRepository{queries: queries}
}

// --- READERS ---

// GetAll
func (r *SponsorRepository) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Sponsor, error) {
	rows, err := r.queries.ListSponsorsByEvent(ctx, eventID)
	if err != nil {
		return nil, ParseDBError(err, "Sponsor")
	}

	sponsors := make([]domain.Sponsor, len(rows))
	for i, row := range rows {
		sponsors[i] = *mapToDomainSponsor(row)
	}

	return sponsors, nil
}

// GetById
func (r *SponsorRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Sponsor, error) {
	row, err := r.queries.GetSponsorByID(ctx, id)
	if err != nil {
		return nil, ParseDBError(err, "Sponsor")
	}

	return mapToDomainSponsor(row), nil
}

// ListPaged
func (r *SponsorRepository) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Sponsor, int64, error) {
	offset := (page - 1) * pageSize

	arg := dbgen.CountSponsorsByEventParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
	}
	total, err := r.queries.CountSponsorsByEvent(ctx, arg)
	if err != nil {
		return nil, 0, ParseDBError(err, "Sponsor")
	}

	rows, err := r.queries.ListSponsorsByEventPaged(ctx, dbgen.ListSponsorsByEventPagedParams{
		EventID: eventID,
		Search:  utils.TextToPgString(search),
		Limit:   pageSize,
		Offset:  offset,
	})
	if err != nil {
		return nil, 0, ParseDBError(err, "Sponsor")
	}

	sponsors := make([]domain.Sponsor, len(rows))
	for i, row := range rows {
		sponsors[i] = *mapToDomainSponsor(row)
	}

	return sponsors, total, nil
}

// --- WRITERS ---

// Create
func (r *SponsorRepository) Create(ctx context.Context, sponsor *domain.Sponsor) (*domain.Sponsor, error) {
	row, err := r.queries.CreateSponsor(ctx, dbgen.CreateSponsorParams{
		EventID:       sponsor.EventID,
		Name:          sponsor.Name,
		LogoUrl:       utils.PtrToText(&sponsor.LogoURL),
		WebsiteUrl:    utils.PtrToText(&sponsor.WebsiteURL),
		Tier:          utils.PtrToText(&sponsor.Tier),
		OrderPriority: utils.PtrToInt4(sponsor.OrderPriority),
		CreatedBy:     sponsor.CreatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Sponsor")
	}

	return mapToDomainSponsor(row), nil
}

// Update
func (r *SponsorRepository) Update(ctx context.Context, sponsor *domain.Sponsor) (*domain.Sponsor, error) {
	row, err := r.queries.UpdateSponsor(ctx, dbgen.UpdateSponsorParams{
		ID:            sponsor.ID,
		Name:          utils.PtrToText(&sponsor.Name),
		LogoUrl:       utils.PtrToText(&sponsor.LogoURL),
		WebsiteUrl:    utils.PtrToText(&sponsor.WebsiteURL),
		Tier:          utils.PtrToText(&sponsor.Tier),
		OrderPriority: utils.PtrToInt4(sponsor.OrderPriority),
		UpdatedBy:     sponsor.UpdatedBy,
	})
	if err != nil {
		return nil, ParseDBError(err, "Sponsor")
	}

	return mapToDomainSponsor(row), nil
}

// Delete
func (r *SponsorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteSponsor(ctx, id)
	if err != nil {
		return ParseDBError(err, "Sponsor")
	}

	return nil
}

// --- Mappers ---

// mapToDomain maps a dbgen.Sponsor to a domain.Sponsor
func mapToDomainSponsor(dbSponsor dbgen.Sponsor) *domain.Sponsor {
	return &domain.Sponsor{
		ID:            dbSponsor.ID,
		EventID:       dbSponsor.EventID,
		Name:          dbSponsor.Name,
		LogoURL:       utils.PgStringToText(dbSponsor.LogoUrl),
		WebsiteURL:    utils.PgStringToText(dbSponsor.WebsiteUrl),
		Tier:          utils.PgStringToText(dbSponsor.Tier),
		OrderPriority: utils.Int4ToPtr(dbSponsor.OrderPriority),
		Audit: domain.Audit{
			CreatedAt: dbSponsor.CreatedAt.Time,
			UpdatedAt: dbSponsor.UpdatedAt.Time,
			CreatedBy: dbSponsor.CreatedBy,
			UpdatedBy: dbSponsor.UpdatedBy,
		},
	}
}
