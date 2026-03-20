package domain

import "github.com/google/uuid"

type Talk struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description string
	Tags        *[]string
	Speaker
	Audit
}
