package dtos

import (
	"github.com/google/uuid"
)

type CreateOrganizerDTO struct {
	EventID         uuid.UUID `json:"event_id" binding:"required"`
	PersonID        uuid.UUID `json:"person_id" binding:"required"`
	Company         *string   `json:"company"`
	RoleDescription *string   `json:"role_description"`
}

type UpdateOrganizerDTO struct {
	Company         *string `json:"company"`
	RoleDescription *string `json:"role_description"`
}

type OrganizerDetailResponse struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         string    `json:"company"`
	RoleDescription string    `json:"role_description"`
	PersonFieldsDTO
	AuditDTO
}

type OrganizerResponse struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         string    `json:"company"`
	RoleDescription string    `json:"role_description"`
	AuditDTO
}
