package usecase

import (
	"context"
	"devfest/internal/domain"

	"github.com/google/uuid"
)

// trackInteractor
type trackInteractor struct {
	trackRepository domain.TrackRepository
	eventRepo       domain.EventRepository
}

// NewTrackInteractor creates a new trackInteractor
func NewTrackInteractor(trackRepository domain.TrackRepository, eventRepo domain.EventRepository) domain.TrackUsecase {
	return &trackInteractor{
		trackRepository: trackRepository,
		eventRepo:       eventRepo,
	}
}

// --- Readers ---

// Get All
func (t *trackInteractor) GetAll(ctx context.Context, eventID uuid.UUID) ([]domain.Track, error) {
	_, err := t.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return t.trackRepository.GetAll(ctx, eventID)
}

// GetById
func (t *trackInteractor) GetById(ctx context.Context, id uuid.UUID) (*domain.Track, error) {
	return t.trackRepository.GetById(ctx, id)
}

// GetFullEventSchedule
func (t *trackInteractor) GetFullEventSchedule(ctx context.Context, eventID uuid.UUID) ([]domain.FullTrackSchedule, error) {
	return t.trackRepository.GetFullEventSchedule(ctx, eventID)
}

// --- Writers ---

// Create
func (t *trackInteractor) Create(ctx context.Context, track *domain.Track) (*domain.Track, error) {
	_, err := t.eventRepo.GetByID(ctx, track.EventID)
	if err != nil {
		return nil, err
	}

	return t.trackRepository.Create(ctx, track)
}

// Update
func (t *trackInteractor) Update(ctx context.Context, id uuid.UUID, update *domain.UpdateTrack) (*domain.Track, error) {
	track, err := t.trackRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if update.Name != nil && *update.Name != "" {
		track.Name = *update.Name
	}
	if update.EventDate != nil {
		track.EventDate = *update.EventDate
	}
	track.UpdatedBy = update.UpdatedBy

	return t.trackRepository.Update(ctx, track)
}

// Delete
func (t *trackInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return t.trackRepository.Delete(ctx, id)
}
