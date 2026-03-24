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

func TestDeveloperUsecase(t *testing.T) {
	colRepo := new(mocks.DeveloperRepository)
	personRepo := new(mocks.PersonRepository)
	eventRepo := new(mocks.EventRepository)

	interactor := usecase.NewDeveloperInteractor(colRepo, personRepo, eventRepo)
	ctx := context.Background()

	t.Run("Create - Success", func(t *testing.T) {
		eventID := uuid.New()
		personID := uuid.New()
		Developer := &domain.Developer{
			EventID: eventID,
			Person:  domain.Person{ID: personID},
		}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		personRepo.On("GetById", ctx, personID).Return(&domain.Person{ID: personID}, nil).Once()
		colRepo.On("Create", ctx, Developer).Return(Developer, nil).Once()

		result, err := interactor.Create(ctx, Developer)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		eventRepo.AssertExpectations(t)
		personRepo.AssertExpectations(t)
		colRepo.AssertExpectations(t)
	})

	t.Run("Create - Fail Event Not Found", func(t *testing.T) {
		colRepo := new(mocks.DeveloperRepository)
		personRepo := new(mocks.PersonRepository)
		eventRepo := new(mocks.EventRepository)
		interactor := usecase.NewDeveloperInteractor(colRepo, personRepo, eventRepo)

		eventID := uuid.New()
		Developer := &domain.Developer{EventID: eventID}

		eventRepo.On("GetByID", ctx, eventID).Return(nil, errors.New("event not found")).Once()

		result, err := interactor.Create(ctx, Developer)

		assert.Error(t, err)
		assert.Nil(t, result)

		personRepo.AssertNotCalled(t, "GetById", mock.Anything, mock.Anything)
		eventRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		colRepo.On("GetAll", ctx, eventID).Return([]domain.Developer{{ID: uuid.New()}}, nil).Once()

		results, err := interactor.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("ListPaged - Success with Pagination Fix", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		colRepo.On("ListPaged", ctx, eventID, "search", int32(1), int32(10)).
			Return([]domain.Developer{}, int64(0), nil).Once()

		_, _, err := interactor.ListPaged(ctx, eventID, "search", 0, 0)

		assert.NoError(t, err)
	})

	t.Run("Update - Success", func(t *testing.T) {
		id := uuid.New()
		newRoleDescription := "Backend"
		existingColl := &domain.Developer{ID: id, RoleDescription: nil}
		updateDTO := &domain.UpdateDeveloper{
			RoleDescription: &newRoleDescription,
			UpdatedBy:       uuid.New(),
		}

		colRepo.On("GetById", ctx, id).Return(existingColl, nil).Once()
		colRepo.On("Update", ctx, mock.MatchedBy(func(c *domain.Developer) bool {
			return *c.RoleDescription == newRoleDescription
		})).Return(existingColl, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.Equal(t, newRoleDescription, *result.RoleDescription)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		id := uuid.New()
		colRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})
}
