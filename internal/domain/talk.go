package domain

import (
	"time"

	"github.com/google/uuid"
)

type SpeakerTalkDetail struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	AvatarURL string
	Company   string
	Bio       string
}

type Talk struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description string
	Tags        *[]string
	Speakers    []SpeakerTalkDetail
	Audit
}

type UpdateTalk struct {
	Title       *string
	Description *string
	Tags        []string
	UpdatedBy   uuid.UUID
}

type TalkSpeaker struct {
	ID        uuid.UUID
	TalkID    uuid.UUID
	SpeakerID uuid.UUID
	CreatedBy uuid.UUID
	CreatedAt time.Time
}
