package domain

import (
	"context"
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
	UpdatedAt time.Time
}

type UpdateScheduler struct {
	StartTime *string
	EndTime   *string
	Room      *string
	UpdatedBy uuid.UUID
}

type SchedulerUsecase interface {
	// Readers
	GetAllByTrack(ctx context.Context, trackID uuid.UUID) ([]Scheduler, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Scheduler, error)
	// Writers
	Create(ctx context.Context, scheduler *Scheduler) (*Scheduler, error)
	Update(ctx context.Context, id uuid.UUID, upScheduler *UpdateScheduler) (*Scheduler, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SchedulerRepository interface {
	// Readers
	GetAllByTrack(ctx context.Context, trackID uuid.UUID) ([]Scheduler, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Scheduler, error)
	// Writers
	Create(ctx context.Context, scheduler *Scheduler) (*Scheduler, error)
	Update(ctx context.Context, scheduler *Scheduler) (*Scheduler, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
