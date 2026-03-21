package domain

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	Name      string
	EventDate time.Time
	Audit
}
