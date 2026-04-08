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
	FirstName string
	LastName  string
	AvatarURL string
	Company   string
}

type TalkTrack struct {
	ID          uuid.UUID
	Title       string
	Description string
	Speakers    []SpeakerTrack
}

type ScheduleEntryTrack struct {
	ScheduleID uuid.UUID
	StartTime  time.Time
	EndTime    time.Time
	Room       string
	Talk       TalkTrack
}

type FullTrackSchedule struct {
	TrackID   uuid.UUID
	TrackName string
	EventDate time.Time
	Entries   []ScheduleEntryTrack
}
