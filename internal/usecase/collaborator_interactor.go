package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// collaboratorInteractor implements domain.CollaboratorUsecase
type collaboratorInteractor struct {
	colRepo    domain.CollaboratorRepository
	personRepo domain.PersonRepository
	eventRepo  domain.EventRepository
}

// NewCollaboratorInteractor creates a new collaboratorInteractor
func NewCollaboratorInteractor(colRepo domain.CollaboratorRepository, personRepo domain.PersonRepository, eventRepo domain.EventRepository) domain.CollaboratorUsecase {
	return &collaboratorInteractor{
		colRepo:    colRepo,
		personRepo: personRepo,
		eventRepo:  eventRepo,
	}
}

// --- Readers ---

// Get All
func (c *collaboratorInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Collaborator, error) {
	_, err := c.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return c.colRepo.GetAll(ctx, eventID)
}

// GetById
func (c *collaboratorInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Collaborator, error) {
	return c.colRepo.GetById(ctx, id)
}

// ListPaged
func (c *collaboratorInteractor) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Collaborator, int64, error) {
	_, err := c.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return c.colRepo.ListPaged(ctx, eventID, search, page, pageSize)
}

// --- Writers ---

// Create
func (c *collaboratorInteractor) Create(ctx context.Context, collaborator *domain.Collaborator) (*domain.Collaborator, error) {
	_, err := c.eventRepo.GetByID(ctx, collaborator.EventID)
	if err != nil {
		return nil, err
	}

	_, err = c.personRepo.GetById(ctx, collaborator.Person.ID)
	if err != nil {
		return nil, err
	}

	return c.colRepo.Create(ctx, collaborator)
}

// Update
func (c *collaboratorInteractor) Update(ctx context.Context, id uuid.UUID, updCol *domain.UpdateCollaborator) (*domain.Collaborator, error) {
	collaborator, err := c.colRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if updCol.Area != nil {
		collaborator.Area = updCol.Area
	}
	if updCol.PersonID != nil {
		collaborator.Person.ID = *updCol.PersonID
	}
	collaborator.UpdatedBy = updCol.UpdatedBy

	return c.colRepo.Update(ctx, collaborator)
}

// Delete
func (c *collaboratorInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return c.colRepo.Delete(ctx, id)
}
