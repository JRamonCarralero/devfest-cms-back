package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// speakerInteractor implements domain.SpeakerUsecase
type speakerInteractor struct {
	spkRepo    domain.SpeakerRepository
	personRepo domain.PersonRepository
	eventRepo  domain.EventRepository
}

// NewSpeakerInteractor creates a new speakerInteractor
func NewSpeakerInteractor(spkRepo domain.SpeakerRepository, personRepo domain.PersonRepository, eventRepo domain.EventRepository) domain.SpeakerUsecase {
	return &speakerInteractor{
		spkRepo:    spkRepo,
		personRepo: personRepo,
		eventRepo:  eventRepo,
	}
}

// --- Readers ---

// Get All
func (s *speakerInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Speaker, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return s.spkRepo.GetAll(ctx, eventID)
}

// GetById
func (s *speakerInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Speaker, error) {
	return s.spkRepo.GetById(ctx, id)
}

// ListPaged
func (s *speakerInteractor) ListPaged(ctx context.Context, eventID uuid.UUID, search string, page, pageSize int32) ([]domain.Speaker, int64, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.spkRepo.ListPaged(ctx, eventID, search, page, pageSize)
}

// --- Writers ---

// Create
func (s *speakerInteractor) Create(ctx context.Context, speaker *domain.Speaker) (*domain.Speaker, error) {
	_, err := s.eventRepo.GetByID(ctx, speaker.EventID)
	if err != nil {
		return nil, err
	}

	_, err = s.personRepo.GetById(ctx, speaker.Person.ID)
	if err != nil {
		return nil, err
	}

	return s.spkRepo.Create(ctx, speaker)
}

// Update
func (s *speakerInteractor) Update(ctx context.Context, id uuid.UUID, upSpeaker *domain.UpdateSpeaker) (*domain.Speaker, error) {
	speaker, err := s.spkRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if upSpeaker.Bio != nil {
		speaker.Bio = upSpeaker.Bio
	}
	if upSpeaker.Company != nil {
		speaker.Company = upSpeaker.Company
	}
	speaker.UpdatedBy = upSpeaker.UpdatedBy

	return s.spkRepo.Update(ctx, speaker)
}

// Delete
func (s *speakerInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return s.spkRepo.Delete(ctx, id)
}
