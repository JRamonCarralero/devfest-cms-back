package dtos

import "github.com/google/uuid"

type CreateSponsorDTO struct {
	EventID    uuid.UUID `json:"event_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	LogoURL    string    `json:"logo_url" binding:"required"`
	WebsiteURL string    `json:"website_url" binding:"required"`
	Tier       int       `json:"tier" binding:"required"`
	CreatedBy  uuid.UUID `json:"-"`
}

type UpdateSponsorDTO struct {
	Name       *string   `json:"name" binding:"omitempty"`
	LogoURL    *string   `json:"logo_url" binding:"omitempty"`
	WebsiteURL *string   `json:"website_url" binding:"omitempty"`
	Tier       *int      `json:"tier" binding:"omitempty"`
	UpdatedBy  uuid.UUID `json:"-"`
}
