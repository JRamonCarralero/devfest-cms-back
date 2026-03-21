package dtos

import "github.com/google/uuid"

type CreateOrganizerDTO struct {
	EventID         uuid.UUID `json:"event_id"`
	PersonID        uuid.UUID `json:"person_id"`
	Company         *string   `json:"company"`
	RoleDescription *string   `json:"role_description"`
	CreatedBy       uuid.UUID `json:"-"`
}

type UpdateOrganizerDTO struct {
	Company         *string   `json:"company"`
	RoleDescription *string   `json:"role_description"`
	UpdatedBy       uuid.UUID `json:"-"`
}
