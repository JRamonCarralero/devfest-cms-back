package domain

import (
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
