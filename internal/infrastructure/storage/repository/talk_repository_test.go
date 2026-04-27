package repository

import (
	"context"
	"devfest/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTalkRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewTalkRepository(testQueries)

	eventID := SeedEvent(t)
	userID := uuid.New()

	t.Run("Add and Get Speakers in Talk", func(t *testing.T) {
		tags := []string{"Architecture"}
		talk, err := repo.Create(ctx, &domain.Talk{
			EventID:     eventID,
			Title:       "Microservices Patterns",
			Description: "Advanced talk",
			Tags:        tags,
			Audit:       domain.Audit{CreatedBy: userID},
		})
		require.NoError(t, err)

		speakerID := SeedSpeaker(t, eventID)

		ts := &domain.TalkSpeaker{
			TalkID:    talk.ID,
			SpeakerID: speakerID,
			CreatedBy: userID,
		}
		err = repo.AddSpeaker(ctx, ts)
		require.NoError(t, err)

		found, err := repo.GetById(ctx, talk.ID)
		require.NoError(t, err)

		assert.Len(t, found.Speakers, 1)
		assert.Equal(t, speakerID, found.Speakers[0].ID)
	})

	t.Run("Remove Speaker from Talk", func(t *testing.T) {
		talk, _ := repo.Create(ctx, &domain.Talk{EventID: eventID, Title: "Title"})
		speakerID := SeedSpeaker(t, eventID)
		ts := &domain.TalkSpeaker{TalkID: talk.ID, SpeakerID: speakerID}

		repo.AddSpeaker(ctx, ts)

		err := repo.RemoveSpeaker(ctx, ts)
		require.NoError(t, err)

		found, _ := repo.GetById(ctx, talk.ID)
		assert.Empty(t, found.Speakers)
	})
}
