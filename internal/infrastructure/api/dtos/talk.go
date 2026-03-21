package dtos

import "github.com/google/uuid"

type CreateTalkDTO struct {
	EventID     uuid.UUID `json:"event_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Tags        *[]string `json:"tags" binding:"required"`
	CreatedBy   uuid.UUID `json:"-"`
}

type UpdateTalkDTO struct {
	Title       *string   `json:"title" binding:"omitempty"`
	Description *string   `json:"description" binding:"omitempty"`
	Tags        []string  `json:"tags" binding:"omitempty"`
	UpdatedBy   uuid.UUID `json:"-"`
}

type CreateTalkSpeakerDTO struct {
	TalkID    uuid.UUID `json:"talk_id" binding:"required"`
	PersonID  uuid.UUID `json:"person_id" binding:"required"`
	CreatedBy uuid.UUID `json:"-"`
}
