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

func TestSponsorUsecase(t *testing.T) {
	sponsorRepo := new(mocks.SponsorRepository)
	personRepo := new(mocks.PersonRepository)
	eventRepo := new(mocks.EventRepository)

	interactor := usecase.NewSponsorInteractor(sponsorRepo, eventRepo)
	ctx := context.Background()

	t.Run("Create - Success", func(t *testing.T) {
		eventID := uuid.New()
		Sponsor := &domain.Sponsor{
			EventID:    eventID,
			Name:       "Google",
			WebsiteURL: "http://www.google.com",
		}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		sponsorRepo.On("Create", ctx, Sponsor).Return(Sponsor, nil).Once()

		result, err := interactor.Create(ctx, Sponsor)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		eventRepo.AssertExpectations(t)
		personRepo.AssertExpectations(t)
		sponsorRepo.AssertExpectations(t)
	})

	t.Run("Create - Fail Event Not Found", func(t *testing.T) {
		sponsorRepo := new(mocks.SponsorRepository)
		eventRepo := new(mocks.EventRepository)
		interactor := usecase.NewSponsorInteractor(sponsorRepo, eventRepo)

		eventID := uuid.New()
		Sponsor := &domain.Sponsor{EventID: eventID}

		eventRepo.On("GetByID", ctx, eventID).Return(nil, errors.New("event not found")).Once()

		result, err := interactor.Create(ctx, Sponsor)

		assert.Error(t, err)
		assert.Nil(t, result)

		personRepo.AssertNotCalled(t, "GetById", mock.Anything, mock.Anything)
		eventRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		sponsorRepo.On("GetAll", ctx, eventID).Return([]domain.Sponsor{{ID: uuid.New()}}, nil).Once()

		results, err := interactor.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("ListPaged - Success with Pagination Fix", func(t *testing.T) {
		eventID := uuid.New()
		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		sponsorRepo.On("ListPaged", ctx, eventID, "search", int32(1), int32(10)).
			Return([]domain.Sponsor{}, int64(0), nil).Once()

		_, _, err := interactor.ListPaged(ctx, eventID, "search", 0, 0)

		assert.NoError(t, err)
	})

	t.Run("Update - Success", func(t *testing.T) {
		id := uuid.New()
		newName := "google Devs"
		existingSponsor := &domain.Sponsor{ID: id, Name: "Google"}
		updateDTO := &domain.UpdateSponsor{
			Name:      &newName,
			UpdatedBy: uuid.New(),
		}

		sponsorRepo.On("GetById", ctx, id).Return(existingSponsor, nil).Once()
		sponsorRepo.On("Update", ctx, mock.MatchedBy(func(c *domain.Sponsor) bool {
			return c.Name == newName
		})).Return(existingSponsor, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.Equal(t, newName, result.Name)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		id := uuid.New()
		sponsorRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})
}
