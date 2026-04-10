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
	Title       *string  `json:"title" binding:"omitempty"`
	Description *string  `json:"description" binding:"omitempty"`
	Tags        []string `json:"tags" binding:"omitempty"`
}

type CreateTalkSpeakerDTO struct {
	TalkID    uuid.UUID `json:"talk_id" binding:"required"`
	SpeakerID uuid.UUID `json:"speaker_id" binding:"required"`
}

type SpeakerTalkDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Company   string    `json:"company"`
	Bio       string    `json:"bio"`
}

type TalkResponse struct {
	ID          uuid.UUID                   `json:"id"`
	EventID     uuid.UUID                   `json:"event_id"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Tags        []string                    `json:"tags"`
	Speakers    []SpeakerTalkDetailResponse `json:"speakers"`
}
