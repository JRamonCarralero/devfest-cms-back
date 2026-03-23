package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrganizerDTO struct {
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         *string   `json:"company"`
	RoleDescription *string   `json:"role_description"`
}

type UpdateOrganizerDTO struct {
	Company         *string `json:"company"`
	RoleDescription *string `json:"role_description"`
}

type OrganizerDetailResponse struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         string    `json:"company"`
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

type OrganizerResponse struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         string    `json:"company"`
	RoleDescription string    `json:"role_description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       uuid.UUID `json:"created_by"`
	UpdatedBy       uuid.UUID `json:"updated_by"`
}
