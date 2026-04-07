package domain

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	Name      string
	EventDate time.Time
	Audit
}

type UpdateTrack struct {
	Name      *string `json:"name" binding:"omitempty"`
	EventDate *string `json:"event_date" binding:"omitempty"`
}
