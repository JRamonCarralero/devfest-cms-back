package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateSchedulerDTO struct {
	TrackID   uuid.UUID `json:"track_id" binding:"required"`
	TalkID    uuid.UUID `json:"talk_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Room      *string   `json:"room" binding:"required"`
}

type UpdateSchedulerDTO struct {
	StartTime *string `json:"start_time" binding:"omitempty"`
	EndTime   *string `json:"end_time" binding:"omitempty"`
	Room      *string `json:"room" binding:"omitempty"`
}

type SchedulerResponse struct {
	ScheduleID uuid.UUID `json:"schedule_id"`
	TrackID    uuid.UUID `json:"track_id"`
	TalkID     uuid.UUID `json:"talk_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Room       string    `json:"room"`
	AuditDTO
}
