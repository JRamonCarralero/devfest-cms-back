package dtos

import (
	"github.com/google/uuid"
)

type CreateTrackDTO struct {
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	EventDate string    `json:"event_date" binding:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateTrackDTO struct {
	Name      *string   `json:"name" binding:"omitempty"`
	EventDate *string   `json:"event_date" binding:"omitempty"`
	UpdatedBy uuid.UUID `json:"-"`
}
