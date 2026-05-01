package domain

import (
	"context"
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

type UpdateTrack struct {
	Name      *string
	EventDate *time.Time
	UpdatedBy uuid.UUID
}

// Detailed Track

type SpeakerTrack struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
}

type TalkTrack struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Speakers    []SpeakerTrack `json:"speakers"`
}

type ScheduleEntryTrack struct {
	ScheduleID uuid.UUID `json:"schedule_id"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Room       string    `json:"room"`
	Talk       TalkTrack `json:"talk"`
}

type FullTrackSchedule struct {
	TrackID   uuid.UUID            `json:"track_id"`
	TrackName string               `json:"track_name"`
	EventDate time.Time            `json:"event_date"`
	Entries   []ScheduleEntryTrack `json:"entries"`
}

type TrackUsecase interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Track, error)
	GetById(ctx context.Context, id uuid.UUID) (*Track, error)
	GetFullEventSchedule(ctx context.Context, eventID uuid.UUID) ([]FullTrackSchedule, error)
	// Writers
	Create(ctx context.Context, track *Track) (*Track, error)
	Update(ctx context.Context, id uuid.UUID, update *UpdateTrack) (*Track, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TrackRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Track, error)
	GetById(ctx context.Context, trackID uuid.UUID) (*Track, error)
	GetFullEventSchedule(ctx context.Context, eventID uuid.UUID) ([]FullTrackSchedule, error)
	// Writers
	Create(ctx context.Context, track *Track) (*Track, error)
	Update(ctx context.Context, track *Track) (*Track, error)
	Delete(ctx context.Context, trackID uuid.UUID) error
}
