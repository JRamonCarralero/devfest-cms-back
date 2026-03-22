package domain

import (
	"context"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       *string
	AvatarURL   *string
	GithubUser  *string
	LinkedinURL *string
	TwitterURL  *string
	WebsiteURL  *string
	Audit
}

type UpdatePerson struct {
	FirstName   *string
	LastName    *string
	Email       *string
	AvatarURL   *string
	GithubUser  *string
	LinkedinURL *string
	TwitterURL  *string
	WebsiteURL  *string
	UpdatedBy   uuid.UUID
}

type PersonUsecase interface {
	// Readers
	GetAll(ctx context.Context) ([]Person, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Person, error)
	GetByEmail(ctx context.Context, email *string) (*Person, error)
	ListPaged(ctx context.Context, search string, page, pageSize int32) ([]Person, int64, error)
	// Writers
	Create(ctx context.Context, person *Person) (*Person, error)
	Update(ctx context.Context, id uuid.UUID, upPerson *UpdatePerson) (*Person, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PersonRepository interface {
	// Readers
	GetAll(ctx context.Context) ([]Person, error)
	GetById(ctx context.Context, id uuid.UUID) (*Person, error)
	GetByEmail(ctx context.Context, email *string) (*Person, error)
	ListPaged(ctx context.Context, search string, page, pageSize int32) ([]Person, int64, error)
	// Writers
	Create(ctx context.Context, person *Person) (*Person, error)
	Update(ctx context.Context, person *Person) (*Person, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
