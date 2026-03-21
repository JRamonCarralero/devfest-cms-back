package domain

import (
	"time"

	"github.com/google/uuid"
)

type Talk struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description string
	Tags        *[]string
	Speakers    []Speaker
	Audit
}

type TalkSpeaker struct {
	ID        uuid.UUID
	TalkID    uuid.UUID
	SpeakerID uuid.UUID
	CreatedBy uuid.UUID
	CreatedAt time.Time
}
