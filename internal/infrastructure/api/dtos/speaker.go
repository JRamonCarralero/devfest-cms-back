package dtos

import "github.com/google/uuid"

type CreateSpeakerDTO struct {
	PersonID  uuid.UUID `json:"person_id" validate:"required"`
	EventID   uuid.UUID `json:"event_id" validate:"required"`
	Bio       *string   `json:"bio" validate:"required"`
	Company   *string   `json:"company" validate:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateSpeakerDTO struct {
	Bio       *string   `json:"bio,omitempty"`
	Company   *string   `json:"company,omitempty"`
	UpdatedBy uuid.UUID `json:"-"`
}
