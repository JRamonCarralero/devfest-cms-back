package dtos

import "github.com/google/uuid"

type CreateDeveloperDTO struct {
	PersonID        uuid.UUID `json:"person_id" binding:"required"`
	EventID         uuid.UUID `json:"event_id" binding:"required"`
	RoleDescription *string   `json:"role_description" binding:"omitempty"`
	CreatedBy       uuid.UUID `json:"-"`
}

type UpdateDeveloperDTO struct {
	RoleDescription *string   `json:"role_description" binding:"omitempty"`
	UpdatedBy       uuid.UUID `json:"-"`
}
