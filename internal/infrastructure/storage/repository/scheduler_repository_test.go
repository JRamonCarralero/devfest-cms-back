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

func TestSchedulerRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewSchedulerRepository(testQueries)
	trackRepo := NewTrackRepository(testQueries)
	talkRepo := NewTalkRepository(testQueries)

	eventID := SeedEvent(t)

	track, err := trackRepo.Create(ctx, &domain.Track{
		Name:      "Main Hall",
		EventID:   eventID,
		EventDate: time.Now(),
	})
	require.NoError(t, err)

	tags := []string{"Go"}
	talk, err := talkRepo.Create(ctx, &domain.Talk{
		EventID: eventID,
		Title:   "Testing in Go",
		Tags:    tags,
	})
	require.NoError(t, err)

	// --- TESTS ---

	t.Run("Create Schedule Entry", func(t *testing.T) {
		startTime := time.Date(2024, 10, 10, 10, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		sch := &domain.Scheduler{
			StartTime: startTime,
			EndTime:   endTime,
			Room:      "Room A",
			Track:     domain.Track{ID: track.ID},
			Talk:      domain.Talk{ID: talk.ID},
		}

		created, err := repo.Create(ctx, sch)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, "Room A", created.Room)
		assert.Equal(t, track.ID, created.Track.ID)
		assert.Equal(t, talk.ID, created.Talk.ID)
	})

	t.Run("GetByID with Joins", func(t *testing.T) {
		// Creamos una entrada específica
		startTime := time.Date(2024, 10, 10, 11, 0, 0, 0, time.UTC)
		sch, _ := repo.Create(ctx, &domain.Scheduler{
			StartTime: startTime,
			EndTime:   startTime.Add(30 * time.Minute),
			Room:      "Room B",
			Track:     domain.Track{ID: track.ID},
			Talk:      domain.Talk{ID: talk.ID},
		})

		// El GetByID de Scheduler suele hacer JOINs para traer nombres de Track y Talk
		found, err := repo.GetByID(ctx, sch.ID)

		require.NoError(t, err)
		assert.Equal(t, sch.ID, found.ID)
		assert.Equal(t, "Main Hall", found.Track.Name) // Verificamos que el JOIN funcionó
		assert.Equal(t, "Testing in Go", found.Talk.Title)
	})

	t.Run("GetAllByTrack", func(t *testing.T) {
		// Ya tenemos al menos 2 entradas creadas en los tests anteriores para este trackID
		entries, err := repo.GetAllByTrack(ctx, track.ID)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(entries), 2)
		for _, e := range entries {
			assert.Equal(t, track.ID, e.Track.ID)
		}
	})

	t.Run("Update Schedule Entry", func(t *testing.T) {
		sch, _ := repo.Create(ctx, &domain.Scheduler{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
			Room:      "Original Room",
			Track:     domain.Track{ID: track.ID},
			Talk:      domain.Talk{ID: talk.ID},
		})

		sch.Room = "Updated Room"
		updated, err := repo.Update(ctx, sch)

		require.NoError(t, err)
		assert.Equal(t, "Updated Room", updated.Room)
	})

	t.Run("Delete Schedule Entry", func(t *testing.T) {
		sch, _ := repo.Create(ctx, &domain.Scheduler{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
			Track:     domain.Track{ID: track.ID},
			Talk:      domain.Talk{ID: talk.ID},
		})

		err := repo.Delete(ctx, sch.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, sch.ID)
		assert.Error(t, err)
	})
}
