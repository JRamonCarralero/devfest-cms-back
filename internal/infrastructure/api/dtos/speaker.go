package dtos

import (
	"github.com/google/uuid"
)

type CreateSpeakerDTO struct {
	PersonID uuid.UUID `json:"person_id" binding:"required"`
	EventID  uuid.UUID `json:"event_id" binding:"required"`
	Bio      *string   `json:"bio"`
	Company  *string   `json:"company"`
}

type UpdateSpeakerDTO struct {
	Bio     *string `json:"bio,omitempty"`
	Company *string `json:"company,omitempty"`
}

type SpeakerDetailResponse struct {
	ID       uuid.UUID `json:"id"`
	EventID  uuid.UUID `json:"event_id"`
	PersonID uuid.UUID `json:"person_id"`
	Bio      string    `json:"bio"`
	Company  string    `json:"company"`
	PersonFieldsDTO
	AuditDTO
}

type SpeakerResponse struct {
	ID       uuid.UUID `json:"id"`
	EventID  uuid.UUID `json:"event_id"`
	PersonID uuid.UUID `json:"person_id"`
	Bio      string    `json:"bio"`
	Company  string    `json:"company"`
	AuditDTO
}
