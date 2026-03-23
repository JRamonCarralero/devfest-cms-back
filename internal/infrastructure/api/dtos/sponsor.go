package dtos

import "github.com/google/uuid"

type CreateSponsorDTO struct {
	EventID    uuid.UUID `json:"event_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	LogoURL    string    `json:"logo_url" binding:"required"`
	WebsiteURL string    `json:"website_url" binding:"required"`
	Tier       int       `json:"tier" binding:"required"`
}

type UpdateSponsorDTO struct {
	Name       *string `json:"name" binding:"omitempty"`
	LogoURL    *string `json:"logo_url" binding:"omitempty"`
	WebsiteURL *string `json:"website_url" binding:"omitempty"`
	Tier       *int    `json:"tier" binding:"omitempty"`
}

type SponsorResponse struct {
	ID            uuid.UUID `json:"id"`
	EventID       uuid.UUID `json:"event_id"`
	Name          string    `json:"name"`
	LogoURL       string    `json:"logo_url"`
	WebsiteURL    string    `json:"website_url"`
	Tier          int       `json:"tier"`
	OrderPriority int       `json:"order_priority"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	CreatedBy     uuid.UUID `json:"created_by"`
	UpdatedBy     uuid.UUID `json:"updated_by"`
}
