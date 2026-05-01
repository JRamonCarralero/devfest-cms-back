package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// talkInteractor implements domain.TalkUsecase
type talkInteractor struct {
	talkRepository domain.TalkRepository
	eventRepo      domain.EventRepository
}

// NewTalkInteractor returns a new TalkInteractor
func NewTalkInteractor(talkRepository domain.TalkRepository, eventRepo domain.EventRepository) domain.TalkUsecase {
	return &talkInteractor{
		talkRepository: talkRepository,
		eventRepo:      eventRepo,
	}
}

// --- Readers ---

// Get All
func (t *talkInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Talk, error) {
	_, err := t.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return t.talkRepository.GetAll(ctx, eventID)
}

// Get By ID
func (t *talkInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Talk, error) {
	return t.talkRepository.GetById(ctx, id)
}

// --- Writers ---

// Create
func (t *talkInteractor) Create(ctx context.Context, talk *domain.Talk) (*domain.Talk, error) {
	_, err := t.eventRepo.GetByID(ctx, talk.EventID)
	if err != nil {
		return nil, err
	}

	return t.talkRepository.Create(ctx, talk)
}

// Update
func (t *talkInteractor) Update(ctx context.Context, id uuid.UUID, upTalk *domain.UpdateTalk) (*domain.Talk, error) {
	talk, err := t.talkRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if upTalk.Title != nil {
		talk.Title = *upTalk.Title
	}
	if upTalk.Description != nil {
		talk.Description = *upTalk.Description
	}
	if len(upTalk.Tags) > 0 {
		talk.Tags = upTalk.Tags
	}
	talk.UpdatedBy = upTalk.UpdatedBy

	return t.talkRepository.Update(ctx, talk)
}

// Delete
func (t *talkInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return t.talkRepository.Delete(ctx, id)
}

// --- Talk-Speaker ---

// AddSpeaker
func (t *talkInteractor) AddSpeaker(ctx context.Context, speaker *domain.TalkSpeaker) error {
	return t.talkRepository.AddSpeaker(ctx, speaker)
}

// RemoveSpeaker
func (t *talkInteractor) RemoveSpeaker(ctx context.Context, speaker *domain.TalkSpeaker) error {
	return t.talkRepository.RemoveSpeaker(ctx, speaker)
}
