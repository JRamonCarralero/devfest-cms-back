package dtos

import "github.com/google/uuid"

type CreateSchedulerDTO struct {
	TrackID   uuid.UUID `json:"track_id" binding:"required"`
	TalkID    uuid.UUID `json:"talk_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Room      *string   `json:"room" binding:"required"`
	CreatedBy uuid.UUID `json:"-"`
}

type UpdateSchedulerDTO struct {
	StartTime *string   `json:"start_time" binding:"omitempty"`
	EndTime   *string   `json:"end_time" binding:"omitempty"`
	Room      *string   `json:"room" binding:"omitempty"`
	UpdatedBy uuid.UUID `json:"-"`
}
