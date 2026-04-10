package domain

import (
	"context"
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

type TalkUsecase interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Talk, error)
	GetById(ctx context.Context, id uuid.UUID) (*Talk, error)
	// Writers
	Create(ctx context.Context, talk *Talk) (*Talk, error)
	Update(ctx context.Context, id uuid.UUID, upTalk *UpdateTalk) (*Talk, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Talk-Speaker
	AddSpeaker(ctx context.Context, speaker *TalkSpeaker) error
	RemoveSpeaker(ctx context.Context, speaker *TalkSpeaker) error
}

type TalkRepository interface {
	// Readers
	GetAll(ctx context.Context, eventID uuid.UUID) ([]Talk, error)
	GetById(ctx context.Context, id uuid.UUID) (*Talk, error)
	// Writers
	Create(ctx context.Context, talk *Talk) (*Talk, error)
	Update(ctx context.Context, talk *Talk) (*Talk, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Talk-Speaker
	AddSpeaker(ctx context.Context, speaker *TalkSpeaker) error
	RemoveSpeaker(ctx context.Context, speaker *TalkSpeaker) error
}
