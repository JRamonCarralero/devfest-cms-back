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

func TestSpeakerUsecase(t *testing.T) {
	colRepo := new(mocks.SpeakerRepository)
	personRepo := new(mocks.PersonRepository)
	eventRepo := new(mocks.EventRepository)

	interactor := usecase.NewSpeakerInteractor(colRepo, personRepo, eventRepo)
	ctx := context.Background()

	t.Run("Create - Success", func(t *testing.T) {
		eventID := uuid.New()
		personID := uuid.New()
		Speaker := &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID},
		}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		personRepo.On("GetById", ctx, personID).Return(&domain.Person{ID: personID}, nil).Once()
		colRepo.On("Create", ctx, Speaker).Return(Speaker, nil).Once()

		result, err := interactor.Create(ctx, Speaker)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		eventRepo.AssertExpectations(t)
		personRepo.AssertExpectations(t)
		colRepo.AssertExpectations(t)
	})

	t.Run("Create - Fail Event Not Found", func(t *testing.T) {
		colRepo := new(mocks.SpeakerRepository)
		personRepo := new(mocks.PersonRepository)
		eventRepo := new(mocks.EventRepository)
		interactor := usecase.NewSpeakerInteractor(colRepo, personRepo, eventRepo)

		eventID := uuid.New()
		Speaker := &domain.Speaker{EventID: eventID}

		eventRepo.On("GetByID", ctx, eventID).Return(nil, errors.New("event not found")).Once()

		result, err := interactor.Create(ctx, Speaker)

		assert.Error(t, err)
		assert.Nil(t, result)

		personRepo.AssertNotCalled(t, "GetById", mock.Anything, mock.Anything)
		eventRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		colRepo.On("GetAll", ctx, eventID).Return([]domain.Speaker{{ID: uuid.New()}}, nil).Once()

		results, err := interactor.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("ListPaged - Success with Pagination Fix", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		colRepo.On("ListPaged", ctx, eventID, "search", int32(1), int32(10)).
			Return([]domain.Speaker{}, int64(0), nil).Once()

		_, _, err := interactor.ListPaged(ctx, eventID, "search", 0, 0)

		assert.NoError(t, err)
	})

	t.Run("Update - Success", func(t *testing.T) {
		id := uuid.New()
		newBio := "New Bio"
		existingColl := &domain.Speaker{ID: id, Bio: nil}
		updateDTO := &domain.UpdateSpeaker{
			Bio:       &newBio,
			UpdatedBy: uuid.New(),
		}

		colRepo.On("GetById", ctx, id).Return(existingColl, nil).Once()
		colRepo.On("Update", ctx, mock.MatchedBy(func(c *domain.Speaker) bool {
			return *c.Bio == newBio
		})).Return(existingColl, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.Equal(t, newBio, *result.Bio)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		id := uuid.New()
		colRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})
}
