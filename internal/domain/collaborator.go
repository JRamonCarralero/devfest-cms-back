package domain

import "github.com/google/uuid"

type Collaborator struct {
	ID      uuid.UUID
	EventID uuid.UUID
	Area    *string
	Person
	Audit
}
