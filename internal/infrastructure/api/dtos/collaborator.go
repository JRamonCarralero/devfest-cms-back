package dtos

import (
	"github.com/google/uuid"
)

type CreateCollaboratorDTO struct {
	PersonID  uuid.UUID `json:"person_id" validate:"required"`
	EventID   uuid.UUID `json:"event_id" validate:"required"`
	Area      *string   `json:"area" validate:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateCollaboratorDTO struct {
	Area      *string    `json:"area,omitempty"`
	PersonID  *uuid.UUID `json:"person_id,omitempty"`
	UpdatedBy uuid.UUID  `json:"-"`
}
