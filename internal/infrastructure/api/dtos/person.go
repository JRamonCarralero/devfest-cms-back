package dtos

import "github.com/google/uuid"

type CreatePersonRequest struct {
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Email       *string   `json:"email" binding:"omitempty,email"`
	AvatarURL   *string   `json:"avatar_url" binding:"omitempty,url"`
	GithubUser  *string   `json:"github_user"`
	LinkedinURL *string   `json:"linkedin_url" binding:"omitempty,url"`
	TwitterURL  *string   `json:"twitter_url" binding:"omitempty,url"`
	WebsiteURL  *string   `json:"website_url" binding:"omitempty,url"`
	CreatedBy   uuid.UUID `json:"-"`
}

type UpdatePersonRequest struct {
	FirstName   *string   `json:"first_name" binding:"required"`
	LastName    *string   `json:"last_name" binding:"required"`
	Email       *string   `json:"email" binding:"omitempty,email"`
	AvatarURL   *string   `json:"avatar_url"`
	GithubUser  *string   `json:"github_user"`
	LinkedinURL *string   `json:"linkedin_url"`
	TwitterURL  *string   `json:"twitter_url"`
	WebsiteURL  *string   `json:"website_url"`
	UpdatedBy   uuid.UUID `json:"-"`
}
