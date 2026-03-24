package repository

import (
	"context"
	"devfest/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpeakerRepository(t *testing.T) {
	repo := NewSpeakerRepository(testQueries)
	ctx := context.Background()

	t.Run("Create and GetById", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)

		bio := "AI Expert"
		company := "Google"
		userID := uuid.New()

		newColl := &domain.Speaker{
			EventID: eventID,
			Bio:     &bio,
			Company: &company,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: userID},
		}

		created, err := repo.Create(ctx, newColl)
		assert.NoError(t, err)

		found, err := repo.GetById(ctx, created.ID)
		assert.NoError(t, err)

		assert.Contains(t, found.Person.FirstName, "Name-")
		assert.Equal(t, bio, *found.Bio)
		assert.Equal(t, company, *found.Company)
	})

	t.Run("Update", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)
		bioInitial := "AI Expert"
		bioUpdated := "Multiagent Engineer"
		userID := uuid.New()

		coll, _ := repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Bio:     &bioInitial,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: userID},
		})

		coll.Bio = &bioUpdated
		coll.UpdatedBy = userID
		updated, err := repo.Update(ctx, coll)

		assert.NoError(t, err)
		assert.Equal(t, bioUpdated, *updated.Bio)
	})

	t.Run("ListPaged", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)
		bio := "AI Expert"

		_, err := repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Bio:     &bio,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: uuid.New()},
		})
		require.NoError(t, err)

		results, total, err := repo.ListPaged(ctx, eventID, "", 1, 10)
		assert.NoError(t, err)
		assert.True(t, total >= 1)
		assert.NotEmpty(t, results)
		assert.Equal(t, eventID, results[0].EventID)
	})

	t.Run("GetByPersonAndEvent", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)

		created, _ := repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: uuid.New()},
		})

		id, err := repo.GetByPersonAndEvent(ctx, personID, eventID)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, id)
	})

	t.Run("Delete", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)

		created, _ := repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: uuid.New()},
		})

		err := repo.Delete(ctx, created.ID)
		assert.NoError(t, err)

		_, err = repo.GetById(ctx, created.ID)
		assert.Error(t, err)
	})

	t.Run("GetAll - Should return all Speakers for an event", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID1 := SeedPerson(t)
		personID2 := SeedPerson(t)

		repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID1},
			Audit:   domain.Audit{CreatedBy: uuid.New()},
		})
		repo.Create(ctx, &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID2},
			Audit:   domain.Audit{CreatedBy: uuid.New()},
		})

		results, err := repo.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		for _, c := range results {
			assert.Equal(t, eventID, c.EventID)
			assert.NotEmpty(t, c.Person.FirstName)
		}
	})

	t.Run("Constraint - Should fail if Speaker already exists", func(t *testing.T) {
		eventID := SeedEvent(t)
		personID := SeedPerson(t)
		userID := uuid.New()

		coll := &domain.Speaker{
			EventID: eventID,
			Person:  domain.Person{ID: personID},
			Audit:   domain.Audit{CreatedBy: userID},
		}

		_, err := repo.Create(ctx, coll)
		assert.NoError(t, err)

		_, err = repo.Create(ctx, coll)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists", "Debería lanzar un error de duplicado")
	})
}
