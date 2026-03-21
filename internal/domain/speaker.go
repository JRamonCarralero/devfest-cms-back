package domain

import "github.com/google/uuid"

type Speaker struct {
	ID      uuid.UUID
	EventID uuid.UUID
	Person
	Bio     *string
	Company *string
	Audit
}
