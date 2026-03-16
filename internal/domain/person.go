package domain

import (
	"context"
	"devfest/internal/infrastructure/api/dtos"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       *string   `json:"email,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	GithubUser  *string   `json:"github_user,omitempty"`
	LinkedinURL *string   `json:"linkedin_url,omitempty"`
	TwitterURL  *string   `json:"twitter_url,omitempty"`
	WebsiteURL  *string   `json:"website_url,omitempty"`
	Audit
}

type PersonUsecase interface {
	// Readers
	GetAll(ctx context.Context) ([]Person, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Person, error)
	GetByEmail(ctx context.Context, email *string) (*Person, error)
	ListPaged(ctx context.Context, search string, page, pageSize int32) ([]Person, int64, error)
	// Writers
	Create(ctx context.Context, dto dtos.CreatePersonDTO) (*Person, error)
	Update(ctx context.Context, dto dtos.UpdatePersonDTO) (*Person, error)
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
