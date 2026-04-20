package usecase

import (
	"context"
	"devfest/internal/domain"
	"time"

	"github.com/google/uuid"
)

// schedulerInteractor implements domain.SchedulerUsecase
type schedulerInteractor struct {
	schedulerRepository domain.SchedulerRepository
	eventRepo           domain.EventRepository
}

// NewSchedulerInteractor creates a new schedulerInteractor
func NewSchedulerInteractor(schedulerRepository domain.SchedulerRepository, eventRepo domain.EventRepository) domain.SchedulerUsecase {
	return &schedulerInteractor{
		schedulerRepository: schedulerRepository,
		eventRepo:           eventRepo,
	}
}

// --- Readers ---

// GetAllByTrack
func (s *schedulerInteractor) GetAllByTrack(ctx context.Context, trackID uuid.UUID) ([]domain.Scheduler, error) {
	return s.schedulerRepository.GetAllByTrack(ctx, trackID)
}

// GetByID
func (s *schedulerInteractor) GetByID(ctx context.Context, id uuid.UUID) (*domain.Scheduler, error) {
	return s.schedulerRepository.GetByID(ctx, id)
}

// --- Writers ---

// Create
func (s *schedulerInteractor) Create(ctx context.Context, scheduler *domain.Scheduler) (*domain.Scheduler, error) {
	_, err := s.eventRepo.GetByID(ctx, scheduler.Track.EventID)
	if err != nil {
		return nil, err
	}

	return s.schedulerRepository.Create(ctx, scheduler)
}

// Update
func (s *schedulerInteractor) Update(ctx context.Context, id uuid.UUID, upScheduler *domain.UpdateScheduler) (*domain.Scheduler, error) {
	scheduler, err := s.schedulerRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if upScheduler.Room != nil {
		scheduler.Room = *upScheduler.Room
	}
	if upScheduler.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *upScheduler.StartTime)
		if err != nil {
			return nil, domain.NewAppError(domain.TypeInternal, "Error parsing startDate", err)
		}
		scheduler.StartTime = startTime
	}
	if upScheduler.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *upScheduler.EndTime)
		if err != nil {
			return nil, domain.NewAppError(domain.TypeInternal, "Error parsing endDate", err)
		}
		scheduler.EndTime = endTime
	}
	scheduler.UpdatedBy = upScheduler.UpdatedBy

	return s.schedulerRepository.Update(ctx, scheduler)
}

// Delete
func (s *schedulerInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	return s.schedulerRepository.Delete(ctx, id)
}
