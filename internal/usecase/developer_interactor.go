package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// developerInteractor implements domain.DeveloperUsecase
type developerInteractor struct {
	devRepo    domain.DeveloperRepository
	personRepo domain.PersonRepository
	eventRepo  domain.EventRepository
}

// NewDeveloperInteractor creates a new developerInteractor
func NewDeveloperInteractor(devRepo domain.DeveloperRepository, personRepo domain.PersonRepository, eventRepo domain.EventRepository) domain.DeveloperUsecase {
	return &developerInteractor{
		devRepo:    devRepo,
		personRepo: personRepo,
		eventRepo:  eventRepo,
	}
}

// --- Readers ---

// Get All
func (d *developerInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Developer, error) {
	_, err := d.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return d.devRepo.GetAll(ctx, eventID)
}

// GetById
func (d *developerInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Developer, error) {
	return d.devRepo.GetById(ctx, id)
}

// ListPaged
func (d *developerInteractor) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Developer, int64, error) {
	_, err := d.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return d.devRepo.ListPaged(ctx, eventID, search, page, pageSize)
}

// --- Writers ---

// Create
func (d *developerInteractor) Create(ctx context.Context, developer *domain.Developer) (*domain.Developer, error) {
	_, err := d.eventRepo.GetByID(ctx, developer.EventID)
	if err != nil {
		return nil, err
	}

	_, err = d.personRepo.GetById(ctx, developer.Person.ID)
	if err != nil {
		return nil, err
	}

	return d.devRepo.Create(ctx, developer)
}

// Update
func (d *developerInteractor) Update(ctx context.Context, id uuid.UUID, upDev *domain.UpdateDeveloper) (*domain.Developer, error) {
	developer, err := d.devRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if upDev.RoleDescription != nil {
		developer.RoleDescription = upDev.RoleDescription
	}
	developer.UpdatedBy = upDev.UpdatedBy

	return d.devRepo.Update(ctx, developer)
}

// Delete
func (d *developerInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return d.devRepo.Delete(ctx, id)
}
