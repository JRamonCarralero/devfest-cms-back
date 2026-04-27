package usecase_test

import (
	"context"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/usecase"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTrackUsecase(t *testing.T) {
	ctx := context.Background()

	setup := func() (*mocks.TrackRepository, *mocks.EventRepository, domain.TrackUsecase) {
		tr := new(mocks.TrackRepository)
		er := new(mocks.EventRepository)
		ti := usecase.NewTrackInteractor(tr, er)
		return tr, er, ti
	}

	t.Run("Create - Success", func(t *testing.T) {
		trackRepo, eventRepo, interactor := setup()
		eventID := uuid.New()
		track := &domain.Track{EventID: eventID, Name: "Main Stage"}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		trackRepo.On("Create", ctx, track).Return(track, nil).Once()

		result, err := interactor.Create(ctx, track)

		assert.NoError(t, err)
		assert.Equal(t, "Main Stage", result.Name)
		trackRepo.AssertExpectations(t)
	})

	t.Run("Create - Fail Event Not Found", func(t *testing.T) {
		trackRepo, eventRepo, interactor := setup()
		eventID := uuid.New()
		track := &domain.Track{EventID: eventID}

		eventRepo.On("GetByID", ctx, eventID).Return(nil, errors.New("event not found")).Once()

		result, err := interactor.Create(ctx, track)

		assert.Error(t, err)
		assert.Nil(t, result)
		trackRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		trackRepo, eventRepo, interactor := setup()
		eventID := uuid.New()
		expected := []domain.Track{{ID: uuid.New(), Name: "Track 1"}}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		trackRepo.On("GetAll", ctx, eventID).Return(expected, nil).Once()

		results, err := interactor.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("GetById - Success", func(t *testing.T) {
		trackRepo, _, interactor := setup()
		id := uuid.New()
		trackRepo.On("GetById", ctx, id).Return(&domain.Track{ID: id, Name: "Room 1"}, nil).Once()

		result, err := interactor.GetById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, "Room 1", result.Name)
	})

	t.Run("Update - Success", func(t *testing.T) {
		trackRepo, _, interactor := setup()
		id := uuid.New()
		newName := "New Name"
		existing := &domain.Track{ID: id, Name: "Old"}
		updateDTO := &domain.UpdateTrack{Name: &newName}

		trackRepo.On("GetById", ctx, id).Return(existing, nil).Once()
		trackRepo.On("Update", ctx, mock.MatchedBy(func(tr *domain.Track) bool {
			return tr.Name == newName
		})).Return(existing, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Update - Fail Not Found", func(t *testing.T) {
		trackRepo, _, interactor := setup()
		id := uuid.New()
		trackRepo.On("GetById", ctx, id).Return(nil, errors.New("not found")).Once()

		_, err := interactor.Update(ctx, id, &domain.UpdateTrack{})

		assert.Error(t, err)
		trackRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		trackRepo, _, interactor := setup()
		id := uuid.New()
		trackRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("GetFullEventSchedule - Success", func(t *testing.T) {
		trackRepo, _, interactor := setup()
		eventID := uuid.New()
		trackRepo.On("GetFullEventSchedule", ctx, eventID).Return([]domain.FullTrackSchedule{}, nil).Once()

		_, err := interactor.GetFullEventSchedule(ctx, eventID)

		assert.NoError(t, err)
	})
}
