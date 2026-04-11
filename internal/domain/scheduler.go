package domain

import (
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	ID uuid.UUID
	Track
	Talk
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Room      string
	Audit
}

type UpdateScheduler struct {
	StartTime *string
	EndTime   *string
	Room      *string
	UpdatedBy uuid.UUID
}
