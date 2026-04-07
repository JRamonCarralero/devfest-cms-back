package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateTrackDTO struct {
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	EventDate string    `json:"event_date" binding:"required"`
}

type UpdateTrackDTO struct {
	Name      *string `json:"name" binding:"omitempty"`
	EventDate *string `json:"event_date" binding:"omitempty"`
}

type TrackResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Name      string    `json:"name"`
	EventDate string    `json:"event_date"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
