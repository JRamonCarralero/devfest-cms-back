package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateTrackDTO struct {
	EventID   uuid.UUID `json:"event_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	EventDate string    `json:"event_date" binding:"required"`
}

type UpdateTrackDTO struct {
	Name      *string `json:"name" binding:"omitempty"`
	EventDate *string `json:"event_date" binding:"omitempty"`
}

type TrackResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	Name      string    `json:"name"`
	EventDate string    `json:"event_date"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Detailed Track Response

type SpeakerTrackResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
}

type TalkTrackResponse struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Speakers    []SpeakerTrackResponse `json:"speakers"`
}

type ScheduleEntryTrackResponse struct {
	ScheduleID uuid.UUID         `json:"schedule_id"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Room       string            `json:"room"`
	Talk       TalkTrackResponse `json:"talk"`
}

type FullTrackScheduleResponse struct {
	TrackID   uuid.UUID                    `json:"track_id"`
	TrackName string                       `json:"track_name"`
	EventDate time.Time                    `json:"event_date"`
	Entries   []ScheduleEntryTrackResponse `json:"entries"`
}
