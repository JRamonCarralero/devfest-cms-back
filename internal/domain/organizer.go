package domain

import "github.com/google/uuid"

type Organizer struct {
	ID              uuid.UUID
	EventID         uuid.UUID
	Company         *string
	RoleDescription *string
	Person
	Audit
}
