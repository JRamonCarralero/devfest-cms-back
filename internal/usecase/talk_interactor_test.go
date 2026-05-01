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

func TestTalkUsecase(t *testing.T) {
	ctx := context.Background()

	setup := func() (*mocks.TalkRepository, *mocks.EventRepository, domain.TalkUsecase) {
		tr := new(mocks.TalkRepository)
		er := new(mocks.EventRepository)
		ti := usecase.NewTalkInteractor(tr, er)
		return tr, er, ti
	}

	t.Run("Create - Success", func(t *testing.T) {
		talkRepo, eventRepo, interactor := setup()
		eventID := uuid.New()
		talk := &domain.Talk{EventID: eventID, Title: "Go Routines"}

		eventRepo.On("GetByID", ctx, eventID).Return(&domain.Event{ID: eventID}, nil).Once()
		talkRepo.On("Create", ctx, talk).Return(talk, nil).Once()

		result, err := interactor.Create(ctx, talk)

		assert.NoError(t, err)
		assert.Equal(t, "Go Routines", result.Title)
		eventRepo.AssertExpectations(t)
	})

	t.Run("GetAll - Fail Event Not Found", func(t *testing.T) {
		talkRepo, eventRepo, interactor := setup()
		eventID := uuid.New()

		eventRepo.On("GetByID", ctx, eventID).Return(nil, errors.New("not found")).Once()

		_, err := interactor.GetAll(ctx, eventID)

		assert.Error(t, err)
		talkRepo.AssertNotCalled(t, "GetAll", mock.Anything, mock.Anything)
	})

	t.Run("GetById - Success", func(t *testing.T) {
		talkRepo, _, interactor := setup()
		id := uuid.New()
		expected := &domain.Talk{ID: id, Title: "Concurrency"}

		talkRepo.On("GetById", ctx, id).Return(expected, nil).Once()

		result, err := interactor.GetById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, "Concurrency", result.Title)
	})

	t.Run("Update - Partial Success", func(t *testing.T) {
		talkRepo, _, interactor := setup()
		id := uuid.New()
		newTitle := "New Title"
		existing := &domain.Talk{ID: id, Title: "Old Title", Description: "Keep me"}

		newTags := []string{"Go", "Advanced"}
		updateDTO := &domain.UpdateTalk{
			Title:     &newTitle,
			Tags:      newTags,
			UpdatedBy: uuid.New(),
		}

		talkRepo.On("GetById", ctx, id).Return(existing, nil).Once()
		talkRepo.On("Update", ctx, mock.MatchedBy(func(tk *domain.Talk) bool {
			return tk.Title == newTitle && tk.Description == "Keep me" && len(tk.Tags) == 2
		})).Return(existing, nil).Once()

		result, err := interactor.Update(ctx, id, updateDTO)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		talkRepo, _, interactor := setup()
		id := uuid.New()
		talkRepo.On("Delete", ctx, id).Return(nil).Once()

		err := interactor.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("AddSpeaker - Success", func(t *testing.T) {
		talkRepo, _, interactor := setup()
		speakerRelation := &domain.TalkSpeaker{
			TalkID:    uuid.New(),
			SpeakerID: uuid.New(),
		}

		talkRepo.On("AddSpeaker", ctx, speakerRelation).Return(nil).Once()

		err := interactor.AddSpeaker(ctx, speakerRelation)

		assert.NoError(t, err)
	})
}
