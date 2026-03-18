package usecase_test

import (
	"context"
	"errors"
	"testing"

	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventInteractor(t *testing.T) {
	ctx := context.Background()

	t.Run("GetEvents", func(t *testing.T) {
		t.Run("should return a list of events", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			interactor := usecase.NewEventInteractor(repo)

			mockEvents := []domain.Event{
				{Name: "Event 1", Slug: "event-1"},
				{Name: "Event 2", Slug: "event-2"},
			}

			repo.On("GetAll", ctx).Return(mockEvents, nil)

			res, err := interactor.GetEvents(ctx)

			assert.NoError(t, err)
			assert.Len(t, res, 2)
			assert.Equal(t, mockEvents, res)
		})

		t.Run("should return error when repository fails", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			repo.On("GetAll", ctx).Return(nil, errors.New("db error"))
			interactor := usecase.NewEventInteractor(repo)

			res, err := interactor.GetEvents(ctx)

			assert.Error(t, err)
			assert.Nil(t, res)
			assert.EqualError(t, err, "db error")
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		eventID := uuid.New()

		t.Run("should return an event when ID exists", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			expectedEvent := &domain.Event{
				ID:   eventID,
				Name: "DevFest Special Edition",
			}

			repo.On("GetByID", ctx, eventID).Return(expectedEvent, nil)

			interactor := usecase.NewEventInteractor(repo)
			res, err := interactor.GetByID(ctx, eventID)

			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, expectedEvent.Name, res.Name)
			assert.Equal(t, eventID, res.ID)
		})

		t.Run("should return error when repository fails", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			repo.On("GetByID", ctx, eventID).Return(nil, errors.New("event not found"))

			interactor := usecase.NewEventInteractor(repo)
			res, err := interactor.GetByID(ctx, eventID)

			assert.Error(t, err)
			assert.Nil(t, res)
			assert.EqualError(t, err, "event not found")
		})
	})

	t.Run("GetEventBySlug", func(t *testing.T) {
		t.Run("should return error if slug is empty", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			interactor := usecase.NewEventInteractor(repo)

			res, err := interactor.GetEventBySlug(ctx, "")

			assert.Nil(t, res)
			appErr, ok := err.(*domain.AppError)
			require.True(t, ok, "error should be an AppError")
			assert.Equal(t, domain.TypeBadRequest, appErr.Type)
			assert.Equal(t, "event slug is required", appErr.Message)
		})

		t.Run("should return error if event not found", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			repo.On("GetBySlug", ctx, "missing").Return(nil, nil)
			interactor := usecase.NewEventInteractor(repo)

			res, err := interactor.GetEventBySlug(ctx, "missing")

			assert.Nil(t, res)
			appErr, ok := err.(*domain.AppError)
			require.True(t, ok, "error should be an AppError")
			assert.Equal(t, domain.TypeNotFound, appErr.Type)
			assert.Equal(t, "event not found", appErr.Message)
		})
	})

	t.Run("GetEventsPaged", func(t *testing.T) {
		t.Run("should use default pagination values", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			repo.On("ListPaged", ctx, "", int32(1), int32(10), "created_at_desc").
				Return([]domain.Event{}, int64(0), nil)

			interactor := usecase.NewEventInteractor(repo)
			_, _, err := interactor.GetEventsPaged(ctx, "", 0, 0, "")

			assert.NoError(t, err)
		})
	})

	t.Run("GetActiveEvents", func(t *testing.T) {
		t.Run("should return only active events", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			interactor := usecase.NewEventInteractor(repo)

			activeEvents := []domain.Event{
				{ID: uuid.New(), Name: "Active Event 1"},
				{ID: uuid.New(), Name: "Active Event 2"},
			}

			repo.On("GetActiveList", ctx).Return(activeEvents, nil)

			res, err := interactor.GetActiveEvents(ctx)

			assert.NoError(t, err)
			assert.Len(t, res, 2)
			assert.Equal(t, activeEvents, res)
		})

		t.Run("should return empty list if no active events", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)
			repo.On("GetActiveList", ctx).Return([]domain.Event{}, nil)

			interactor := usecase.NewEventInteractor(repo)
			res, err := interactor.GetActiveEvents(ctx)

			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	})

	t.Run("CreateEvent", func(t *testing.T) {
		userID := uuid.New()
		isActive := true

		dto := dtos.CreateEventDTO{
			Name:      "DevFest 2026",
			Slug:      "devfest-2026",
			IsActive:  &isActive,
			CreatedBy: userID,
		}

		t.Run("should successfully create event with UUID audit", func(t *testing.T) {
			repo := mocks.NewEventRepository(t)

			repo.On("Create", ctx, mock.MatchedBy(func(e *domain.Event) bool {
				return e.Name == dto.Name &&
					e.CreatedBy == userID &&
					e.UpdatedBy == userID &&
					*e.IsActive == true
			})).Return(&domain.Event{ID: uuid.New(), Name: dto.Name}, nil)

			interactor := usecase.NewEventInteractor(repo)
			res, err := interactor.CreateEvent(ctx, dto)

			assert.NoError(t, err)
			assert.NotNil(t, res)
			repo.AssertExpectations(t)
		})
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		eventID := uuid.New()
		editorID := uuid.New()
		repo := mocks.NewEventRepository(t)
		interactor := usecase.NewEventInteractor(repo)

		t.Run("should update name and audit UUID", func(t *testing.T) {
			existing := &domain.Event{
				ID:    eventID,
				Name:  "Old Name",
				Audit: domain.Audit{CreatedBy: uuid.New()},
			}

			newName := "Updated Name"
			dto := dtos.UpdateEventDTO{
				Name:      &newName,
				UpdatedBy: editorID,
			}

			repo.On("GetByID", ctx, eventID).Return(existing, nil).Once()

			repo.On("Update", ctx, mock.MatchedBy(func(e *domain.Event) bool {
				return e.Name == "Updated Name" && e.UpdatedBy == editorID
			})).Return(&domain.Event{ID: eventID, Name: "Updated Name"}, nil).Once()

			res, err := interactor.UpdateEvent(ctx, eventID, dto)

			assert.NoError(t, err)
			assert.Equal(t, "Updated Name", res.Name)
			repo.AssertExpectations(t)
		})
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		id := uuid.New()
		repo := mocks.NewEventRepository(t)
		repo.On("Delete", ctx, id).Return(nil)

		interactor := usecase.NewEventInteractor(repo)
		err := interactor.DeleteEvent(ctx, id)

		assert.NoError(t, err)
	})
}
