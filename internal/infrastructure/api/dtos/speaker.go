package dtos

import "github.com/google/uuid"

type CreateSpeakerDTO struct {
	PersonID uuid.UUID `json:"person_id" validate:"required"`
	EventID  uuid.UUID `json:"event_id" validate:"required"`
	Bio      *string   `json:"bio" validate:"required"`
	Company  *string   `json:"company" validate:"required"`
}

type UpdateSpeakerDTO struct {
	Bio     *string `json:"bio,omitempty"`
	Company *string `json:"company,omitempty"`
}

type SpeakerDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	EventID     uuid.UUID `json:"event_id"`
	PersonID    uuid.UUID `json:"person_id"`
	Bio         string    `json:"bio"`
	Company     string    `json:"company"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	AvatarUrl   string    `json:"avatar_url"`
	GithubUser  string    `json:"github_user"`
	TwitterUrl  string    `json:"twitter_url"`
	LinkedinUrl string    `json:"linkedin_url"`
	WebsiteUrl  string    `json:"website_url"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}

type SpeakerResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	PersonID  uuid.UUID `json:"person_id"`
	Bio       string    `json:"bio"`
	Company   string    `json:"company"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}
