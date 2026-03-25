package dtos

import (
	"github.com/google/uuid"
)

type CreatePersonDTO struct {
	FirstName   string  `json:"first_name" binding:"required,min=2"`
	LastName    string  `json:"last_name" binding:"required,min=2"`
	Email       *string `json:"email" binding:"omitempty,email"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
	GithubUser  *string `json:"github_user" binding:"omitempty,url"`
	LinkedinURL *string `json:"linkedin_url" binding:"omitempty,url"`
	TwitterURL  *string `json:"twitter_url" binding:"omitempty,url"`
	WebsiteURL  *string `json:"website_url" binding:"omitempty,url"`
}

type UpdatePersonDTO struct {
	FirstName   *string `json:"first_name" binding:"omitempty,min=2"`
	LastName    *string `json:"last_name" binding:"omitempty,min=2"`
	Email       *string `json:"email" binding:"omitempty,email"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
	GithubUser  *string `json:"github_user" binding:"omitempty,url"`
	LinkedinURL *string `json:"linkedin_url" binding:"omitempty,url"`
	TwitterURL  *string `json:"twitter_url" binding:"omitempty,url"`
	WebsiteURL  *string `json:"website_url" binding:"omitempty,url"`
}

type PersonResponse struct {
	ID uuid.UUID `json:"id"`
	PersonFieldsDTO
	AuditDTO
}
