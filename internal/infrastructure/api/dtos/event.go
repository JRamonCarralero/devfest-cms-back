package dtos

import "github.com/google/uuid"

type CreateEventDTO struct {
	Name      string    `json:"name" binding:"required,min=3,max=100"`
	Slug      string    `json:"slug" binding:"required,lowercase"`
	IsActive  *bool     `json:"is_active" binding:"omitempty"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateEventDTO struct {
	Name      *string   `json:"name" binding:"omitempty,min=3,max=100"`
	Slug      *string   `json:"slug" binding:"omitempty,lowercase"`
	IsActive  *bool     `json:"is_active" binding:"omitempty"`
	UpdatedBy uuid.UUID `json:"-"`
}
