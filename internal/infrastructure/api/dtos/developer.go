package dtos

import "github.com/google/uuid"

type CreateDeveloperDTO struct {
	PersonID        uuid.UUID `json:"person_id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription *string   `json:"role_description"`
	CreatedBy       uuid.UUID `json:"-"`
}

type UpdateDeveloperDTO struct {
	RoleDescription *string   `json:"role_description,omitempty"`
	UpdatedBy       uuid.UUID `json:"-"`
}
