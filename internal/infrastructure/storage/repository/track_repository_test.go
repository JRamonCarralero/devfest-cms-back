package repository

import (
	"context"
	"devfest/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewTrackRepository(testQueries)

	eventID := SeedEvent(t)

	t.Run("Create Track", func(t *testing.T) {
		track := &domain.Track{
			Name:      "Main Stage",
			EventID:   eventID,
			EventDate: time.Now().Truncate(time.Hour * 24),
		}

		created, err := repo.Create(ctx, track)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, track.Name, created.Name)
		assert.Equal(t, track.EventID, created.EventID)
	})

	t.Run("GetById", func(t *testing.T) {
		trackName := "Workshop Room"
		created, err := repo.Create(ctx, &domain.Track{
			Name:      trackName,
			EventID:   eventID,
			EventDate: time.Now(),
		})
		require.NoError(t, err)

		found, err := repo.GetById(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, trackName, found.Name)
	})

	t.Run("GetAll", func(t *testing.T) {
		tracks, err := repo.GetAll(ctx, eventID)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tracks), 1)
		for _, tr := range tracks {
			assert.Equal(t, eventID, tr.EventID)
		}
	})

	t.Run("Update Track", func(t *testing.T) {
		created, _ := repo.Create(ctx, &domain.Track{
			Name:      "Old Name",
			EventID:   eventID,
			EventDate: time.Now(),
		})

		created.Name = "Updated Name"
		created.UpdatedBy = uuid.New()

		updated, err := repo.Update(ctx, created)

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
		assert.Equal(t, created.UpdatedBy, updated.UpdatedBy)
	})

	t.Run("Delete Track", func(t *testing.T) {
		created, _ := repo.Create(ctx, &domain.Track{
			Name:      "To Be Deleted",
			EventID:   eventID,
			EventDate: time.Now(),
		})

		err := repo.Delete(ctx, created.ID)
		require.NoError(t, err)

		_, err = repo.GetById(ctx, created.ID)
		assert.Error(t, err)
	})

	t.Run("GetFullEventSchedule", func(t *testing.T) {
		repo.Create(ctx, &domain.Track{
			Name:      "Schedule Track",
			EventID:   eventID,
			EventDate: time.Now(),
		})

		schedule, err := repo.GetFullEventSchedule(ctx, eventID)

		require.NoError(t, err)
		assert.NotNil(t, schedule)
		found := false
		for _, s := range schedule {
			if s.TrackName == "Schedule Track" {
				found = true
				break
			}
		}
		assert.True(t, found, "Debería encontrar el track en el cronograma")
	})
}
