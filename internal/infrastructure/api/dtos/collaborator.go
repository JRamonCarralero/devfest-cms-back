package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateCollaboratorDTO struct {
	PersonID  uuid.UUID `json:"person_id" binding:"required"`
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	Area      *string   `json:"area" binding:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateCollaboratorDTO struct {
	Area     *string    `json:"area" binding:"omitempty"`
	PersonID *uuid.UUID `json:"person_id" binding:"omitempty"`
}

type CollaboratorDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	PersonID    uuid.UUID `json:"person_id"`
	EventID     uuid.UUID `json:"event_id"`
	Area        string    `json:"area"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	AvatarUrl   string    `json:"avatar_url"`
	GithubUser  string    `json:"github_user"`
	TwitterUrl  string    `json:"twitter_url"`
	LinkedinUrl string    `json:"linkedin_url"`
	WebsiteUrl  string    `json:"website_url"`
}

type CollaboratorResponse struct {
	ID        uuid.UUID `json:"id"`
	PersonID  uuid.UUID `json:"person_id"`
	EventID   uuid.UUID `json:"event_id"`
	Area      string    `json:"area"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}
