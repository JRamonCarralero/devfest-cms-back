package dtos

import (
	"time"

	"github.com/google/uuid"
)

type PersonFieldsDTO struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	GithubUser  string `json:"github_user"`
	LinkedinURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
	WebsiteURL  string `json:"website_url"`
}

type AuditDTO struct {
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
