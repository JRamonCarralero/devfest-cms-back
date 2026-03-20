package domain

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
