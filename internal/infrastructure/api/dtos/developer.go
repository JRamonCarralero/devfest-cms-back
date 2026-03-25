package dtos

import (
	"github.com/google/uuid"
)

type CreateDeveloperDTO struct {
	PersonID        uuid.UUID `json:"person_id" binding:"required"`
	EventID         uuid.UUID `json:"event_id" binding:"required"`
	RoleDescription *string   `json:"role_description" binding:"omitempty"`
}

type UpdateDeveloperDTO struct {
	RoleDescription *string `json:"role_description" binding:"omitempty"`
}

type DeveloperDetailResponse struct {
	ID              uuid.UUID `json:"id"`
	PersonID        uuid.UUID `json:"person_id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription string    `json:"role_description"`
	PersonFieldsDTO
	AuditDTO
}

type DeveloperResponse struct {
	ID              uuid.UUID `json:"id"`
	PersonID        uuid.UUID `json:"person_id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription string    `json:"role_description"`
	AuditDTO
}
