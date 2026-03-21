package domain

import (
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	ID uuid.UUID
	Track
	Talk
	StarTime time.Time
	EndTime  time.Time
	Duration time.Duration
	Room     *string
	Audit
}
