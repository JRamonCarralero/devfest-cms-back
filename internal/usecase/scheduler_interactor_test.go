package usecase_test

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/usecase"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSchedulerUsecase(t *testing.T) {
	ctx := context.Background()

	setup := func() (*mocks.SchedulerRepository, *mocks.TrackRepository, *mocks.TalkRepository, domain.SchedulerUsecase) {
		sr := new(mocks.SchedulerRepository)
		trr := new(mocks.TrackRepository)
		tlr := new(mocks.TalkRepository)
		si := usecase.NewSchedulerInteractor(sr, trr, tlr)
		return sr, trr, tlr, si
	}

	t.Run("Create - Success", func(t *testing.T) {
		schRepo, trackRepo, talkRepo, interactor := setup()
		trackID := uuid.New()
		talkID := uuid.New()
		eventID := uuid.New()
		scheduler := &domain.Scheduler{
			Track: domain.Track{ID: trackID, EventID: eventID},
			Talk:  domain.Talk{ID: talkID, EventID: eventID},
			Room:  "Hall 1",
		}

		trackRepo.On("GetById", ctx, trackID).Return(&domain.Track{ID: trackID}, nil).Once()
		talkRepo.On("GetById", ctx, talkID).Return(&domain.Talk{ID: talkID}, nil).Once()
		schRepo.On("Create", ctx, scheduler).Return(scheduler, nil).Once()

		result, err := interactor.Create(ctx, scheduler)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetAllByTrack - Success", func(t *testing.T) {
		schRepo, _, _, interactor := setup()
		trackID := uuid.New()
		expected := []domain.Scheduler{{ID: uuid.New(), Room: "Lab 2"}}

		schRepo.On("GetAllByTrack", ctx, trackID).Return(expected, nil).Once()

		results, err := interactor.GetAllByTrack(ctx, trackID)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("Update - Success with Date Parsing", func(t *testing.T) {
		schRepo, _, _, interactor := setup()
		id := uuid.New()

		newStartStr := "2026-04-27T10:00:00Z"
		newEndStr := "2026-04-27T11:00:00Z"
		newRoom := "Room Red"

		existing := &domain.Scheduler{ID: id, Room: "Old Room"}
		updateDTO := &domain.UpdateScheduler{
			Room:      &newRoom,
			StartTime: &newStartStr,
			EndTime:   &newEndStr,
			UpdatedBy: uuid.New(),
		}

		schRepo.On("GetByID", ctx, id).Return(existing, nil).Once()
		schRepo.On("Update", ctx, mock.MatchedBy(func(s *domain.Scheduler) bool {
			return s.Room == newRoom &&
				s.StartTime.Format(time.RFC3339) == newStartStr &&
				s.EndTime.Format(time.RFC3339) == newEndStr
		})).Return(existing, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Update - Fail Invalid Date Format", func(t *testing.T) {
		schRepo, _, _, interactor := setup()
		id := uuid.New()
		badDate := "27-04-2026 10:00"

		existing := &domain.Scheduler{ID: id}
		updateDTO := &domain.UpdateScheduler{StartTime: &badDate}

		schRepo.On("GetByID", ctx, id).Return(existing, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.Error(t, err)
		assert.Nil(t, result)
		schRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Update - Fail Not Found", func(t *testing.T) {
		schRepo, _, _, interactor := setup()
		id := uuid.New()
		schRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found")).Once()

		result, err := interactor.Update(ctx, id, &domain.UpdateScheduler{})

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		schRepo, _, _, interactor := setup()
		id := uuid.New()
		schRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})
}
