package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeveloperDTO struct {
	PersonID        uuid.UUID `json:"person_id" binding:"required"`
	EventID         uuid.UUID `json:"event_id" binding:"required"`
	RoleDescription *string   `json:"role_description" binding:"omitempty"`
}

type UpdateDeveloperDTO struct {
	RoleDescription *string `json:"role_description" binding:"omitempty"`
}

type DeveloperDetailResponse struct {
	ID              uuid.UUID `json:"id"`
	PersonID        uuid.UUID `json:"person_id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription string    `json:"role_description"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	AvatarUrl       string    `json:"avatar_url"`
	GithubUser      string    `json:"github_user"`
	TwitterUrl      string    `json:"twitter_url"`
	LinkedinUrl     string    `json:"linkedin_url"`
	WebsiteUrl      string    `json:"website_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uuid.UUID `json:"created_by"`
	UpdatedBy       uuid.UUID `json:"updated_by"`
}

type DeveloperResponse struct {
	ID              uuid.UUID `json:"id"`
	PersonID        uuid.UUID `json:"person_id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription string    `json:"role_description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uuid.UUID `json:"created_by"`
	UpdatedBy       uuid.UUID `json:"updated_by"`
}
