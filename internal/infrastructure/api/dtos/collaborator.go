package dtos

import (
	"github.com/google/uuid"
)

type CreateCollaboratorDTO struct {
	PersonID  uuid.UUID `json:"person_id" binding:"required"`
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	Area      *string   `json:"area" binding:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateCollaboratorDTO struct {
	Area      *string    `json:"area" binding:"omitempty"`
	PersonID  *uuid.UUID `json:"person_id" binding:"omitempty"`
	UpdatedBy uuid.UUID  `json:"-"`
}
