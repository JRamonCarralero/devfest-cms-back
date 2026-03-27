package dtos

import "github.com/google/uuid"

type CreateSponsorDTO struct {
	EventID       uuid.UUID `json:"event_id" binding:"required"`
	Name          string    `json:"name" binding:"required"`
	LogoURL       string    `json:"logo_url" binding:"required"`
	WebsiteURL    string    `json:"website_url" binding:"required"`
	Tier          string    `json:"tier" binding:"required"`
	OrderPriority *int      `json:"order_priority" binding:"omitempty"`
}

type UpdateSponsorDTO struct {
	Name          *string `json:"name" binding:"omitempty"`
	LogoURL       *string `json:"logo_url" binding:"omitempty"`
	WebsiteURL    *string `json:"website_url" binding:"omitempty"`
	Tier          *string `json:"tier" binding:"omitempty"`
	OrderPriority *int    `json:"order_priority" binding:"omitempty"`
}

type SponsorResponse struct {
	ID            uuid.UUID `json:"id"`
	EventID       uuid.UUID `json:"event_id"`
	Name          string    `json:"name"`
	LogoURL       string    `json:"logo_url"`
	WebsiteURL    string    `json:"website_url"`
	Tier          string    `json:"tier"`
	OrderPriority int       `json:"order_priority"`
	AuditDTO
}
