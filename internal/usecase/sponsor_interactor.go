package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// sponsorInteractor implements domain.SponsorUsecase
type sponsorInteractor struct {
	sponsorRepository domain.SponsorRepository
	eventRepo         domain.EventRepository
}

// NewSponsorInteractor creates a new sponsorInteractor
func NewSponsorInteractor(sponsorRepository domain.SponsorRepository, eventRepo domain.EventRepository) domain.SponsorUsecase {
	return &sponsorInteractor{
		sponsorRepository: sponsorRepository,
		eventRepo:         eventRepo,
	}
}

// --- Readers ---

// Get All
func (s *sponsorInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Sponsor, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return s.sponsorRepository.GetAll(ctx, eventID)
}

// GetById
func (s *sponsorInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Sponsor, error) {
	return s.sponsorRepository.GetById(ctx, id)
}

// ListPaged
func (s *sponsorInteractor) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Sponsor, int64, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.sponsorRepository.ListPaged(ctx, eventID, search, page, pageSize)
}

// --- Writers ---

// Create
func (s *sponsorInteractor) Create(ctx context.Context, sponsor *domain.Sponsor) (*domain.Sponsor, error) {
	_, err := s.eventRepo.GetByID(ctx, sponsor.EventID)
	if err != nil {
		return nil, err
	}

	return s.sponsorRepository.Create(ctx, sponsor)
}

// Update
func (s *sponsorInteractor) Update(ctx context.Context, id uuid.UUID, upSponsor *domain.UpdateSponsor) (*domain.Sponsor, error) {
	sponsor, err := s.sponsorRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if upSponsor.Name != nil {
		sponsor.Name = *upSponsor.Name
	}
	if upSponsor.LogoURL != nil {
		sponsor.LogoURL = *upSponsor.LogoURL
	}
	if upSponsor.WebsiteURL != nil {
		sponsor.WebsiteURL = *upSponsor.WebsiteURL
	}
	if upSponsor.Tier != nil {
		sponsor.Tier = *upSponsor.Tier
	}
	if upSponsor.OrderPriority != nil {
		sponsor.OrderPriority = upSponsor.OrderPriority
	}
	sponsor.Audit.UpdatedBy = upSponsor.UpdatedBy

	return s.sponsorRepository.Update(ctx, sponsor)
}

// Delete
func (s *sponsorInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return s.sponsorRepository.Delete(ctx, id)
}
