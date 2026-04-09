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

type UpdateTrack struct {
	Name      *string
	EventDate *string
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
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Room       string    `json:"room"`
	Talk       TalkTrack `json:"talk"`
}

type FullTrackSchedule struct {
	TrackID   uuid.UUID            `json:"track_id"`
	TrackName string               `json:"track_name"`
	EventDate time.Time            `json:"event_date"`
	Entries   []ScheduleEntryTrack `json:"entries"`
}
