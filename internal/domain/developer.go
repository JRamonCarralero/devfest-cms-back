package domain

import "github.com/google/uuid"

type Developer struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	RoleDescription *string   `json:"role_description"`
	Person
	Audit
}
