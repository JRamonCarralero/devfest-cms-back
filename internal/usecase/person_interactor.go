package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// personInteractor implements domain.PersonUsecase
type personInteractor struct {
	repo domain.PersonRepository
}

// NewPersonInteractor is a constructor for personInteractor
func NewPersonInteractor(repo domain.PersonRepository) domain.PersonUsecase {
	return &personInteractor{
		repo: repo,
	}
}

// --- Readers ---

// GetAll returns all people
func (i *personInteractor) GetAll(ctx context.Context) ([]domain.Person, error) {
	return i.repo.GetAll(ctx)
}

// GetByID returns a person by its ID
func (i *personInteractor) GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	return i.repo.GetById(ctx, id)
}

// GetByEmail returns a person by its email
func (i *personInteractor) GetByEmail(ctx context.Context, email *string) (*domain.Person, error) {
	if email == nil {
		appErr := domain.NewAppError(domain.TypeBadRequest, "email is required", nil)
		return nil, appErr
	}

	return i.repo.GetByEmail(ctx, email)
}

// ListPaged returns a page of people
func (i *personInteractor) ListPaged(ctx context.Context, search string, page, pageSize int32) ([]domain.Person, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return i.repo.ListPaged(ctx, search, page, pageSize)
}

// --- Writers ---

// Create creates a new person
func (i *personInteractor) Create(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	createdPerson, err := i.repo.Create(ctx, person)

	if err != nil {
		return nil, err
	}

	return createdPerson, nil
}

// Update validates params and updates a person
func (i *personInteractor) Update(ctx context.Context, id uuid.UUID, upPerson *domain.UpdatePerson) (*domain.Person, error) {
	person, err := i.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if upPerson.FirstName != nil {
		person.FirstName = *upPerson.FirstName
	}
	if upPerson.LastName != nil {
		person.LastName = *upPerson.LastName
	}
	if upPerson.Email != nil {
		person.Email = upPerson.Email
	}
	if upPerson.AvatarURL != nil {
		person.AvatarURL = upPerson.AvatarURL
	}
	if upPerson.GithubUser != nil {
		person.GithubUser = upPerson.GithubUser
	}
	if upPerson.LinkedinURL != nil {
		person.LinkedinURL = upPerson.LinkedinURL
	}
	if upPerson.TwitterURL != nil {
		person.TwitterURL = upPerson.TwitterURL
	}
	if upPerson.WebsiteURL != nil {
		person.WebsiteURL = upPerson.WebsiteURL
	}

	person.Audit.UpdatedBy = upPerson.UpdatedBy

	updatedPerson, err := i.repo.Update(ctx, person)
	if err != nil {
		return nil, err
	}

	return updatedPerson, nil
}

// Delete deletes a person by its ID
func (i *personInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return i.repo.Delete(ctx, id)
}
