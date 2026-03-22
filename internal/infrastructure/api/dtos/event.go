package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventDTO struct {
	Name      string    `json:"name" binding:"required,min=3,max=100"`
	Slug      string    `json:"slug" binding:"required,lowercase"`
	IsActive  *bool     `json:"is_active" binding:"omitempty"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateEventDTO struct {
	Name     *string `json:"name" binding:"omitempty,min=3,max=100"`
	Slug     *string `json:"slug" binding:"omitempty,lowercase"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}

type EventResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	IsActive  *bool     `json:"is_active"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
