package dtos

import (
	"github.com/google/uuid"
)

type CreateCollaboratorDTO struct {
	PersonID uuid.UUID `json:"person_id" binding:"required"`
	EventID  uuid.UUID `json:"event_id" binding:"required"`
	Area     *string   `json:"area" binding:"required"`
}

type UpdateCollaboratorDTO struct {
	Area *string `json:"area" binding:"omitempty"`
}

type CollaboratorDetailResponse struct {
	ID       uuid.UUID `json:"id"`
	PersonID uuid.UUID `json:"person_id"`
	EventID  uuid.UUID `json:"event_id"`
	Area     string    `json:"area"`
	AuditDTO
	PersonFieldsDTO
}

type CollaboratorResponse struct {
	ID       uuid.UUID `json:"id"`
	PersonID uuid.UUID `json:"person_id"`
	EventID  uuid.UUID `json:"event_id"`
	Area     string    `json:"area"`
	AuditDTO
}
