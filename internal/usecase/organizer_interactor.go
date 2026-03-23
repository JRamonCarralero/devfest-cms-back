package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// organizerInteractor implements domain.OrganizerUsecase
type organizerInteractor struct {
	orgRepo    domain.OrganizerRepository
	personRepo domain.PersonRepository
	eventRepo  domain.EventRepository
}

// NewOrganizerInteractor creates a new organizerInteractor
func NewOrganizerInteractor(orgRepo domain.OrganizerRepository, personRepo domain.PersonRepository, eventRepo domain.EventRepository) domain.OrganizerUsecase {
	return &organizerInteractor{
		orgRepo:    orgRepo,
		personRepo: personRepo,
		eventRepo:  eventRepo,
	}
}

// --- Readers ---

// Get All
func (o *organizerInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Organizer, error) {
	_, err := o.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return o.orgRepo.GetAll(ctx, eventID)
}

// GetById
func (o *organizerInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Organizer, error) {
	return o.orgRepo.GetById(ctx, id)
}

// ListPaged
func (o *organizerInteractor) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Organizer, int64, error) {
	_, err := o.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return o.orgRepo.ListPaged(ctx, eventID, search, page, pageSize)
}

// --- Writers ---

// Create
func (o *organizerInteractor) Create(ctx context.Context, organizer *domain.Organizer) (*domain.Organizer, error) {
	_, err := o.eventRepo.GetByID(ctx, organizer.EventID)
	if err != nil {
		return nil, err
	}

	_, err = o.personRepo.GetById(ctx, organizer.Person.ID)
	if err != nil {
		return nil, err
	}

	return o.orgRepo.Create(ctx, organizer)
}

// Update
func (o *organizerInteractor) Update(ctx context.Context, id uuid.UUID, upOrg *domain.UpdateOrganizer) (*domain.Organizer, error) {
	organizer, err := o.orgRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if upOrg.Company != nil {
		organizer.Company = upOrg.Company
	}
	if upOrg.RoleDescription != nil {
		organizer.RoleDescription = upOrg.RoleDescription
	}
	organizer.UpdatedBy = upOrg.UpdatedBy

	return o.orgRepo.Update(ctx, organizer)
}

// Delete
func (o *organizerInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return o.orgRepo.Delete(ctx, id)
}
