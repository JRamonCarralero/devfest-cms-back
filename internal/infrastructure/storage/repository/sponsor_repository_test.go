package repository

import (
	"context"
	"devfest/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSponsorRepository(t *testing.T) {
	repo := NewSponsorRepository(testQueries)
	ctx := context.Background()

	t.Run("Create and GetById", func(t *testing.T) {
		eventID := SeedEvent(t)

		url := "http://www.google.com"
		name := "Google"
		tier := "Platinum"
		userID := uuid.New()

		newSponsor := &domain.Sponsor{
			EventID:    eventID,
			Name:       name,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: userID},
		}

		created, err := repo.Create(ctx, newSponsor)
		assert.NoError(t, err)

		found, err := repo.GetById(ctx, created.ID)
		assert.NoError(t, err)

		assert.Contains(t, found.Name, "Google")
		assert.Equal(t, tier, found.Tier)
		assert.Equal(t, url, found.LogoURL)
	})

	t.Run("Update", func(t *testing.T) {
		eventID := SeedEvent(t)
		nameUpdated := "GoogleDevs"
		userID := uuid.New()

		url := "http://www.google.com"
		name := "Google"
		tier := "Platinum"

		newSponsor, _ := repo.Create(ctx, &domain.Sponsor{
			EventID:    eventID,
			Name:       name,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: userID},
		})

		newSponsor.Name = nameUpdated
		newSponsor.UpdatedBy = userID
		updated, err := repo.Update(ctx, newSponsor)

		assert.NoError(t, err)
		assert.Equal(t, nameUpdated, updated.Name)
	})

	t.Run("ListPaged", func(t *testing.T) {
		eventID := SeedEvent(t)

		url := "http://www.google.com"
		name := "Google"
		tier := "Platinum"

		_, err := repo.Create(ctx, &domain.Sponsor{
			EventID:    eventID,
			Name:       name,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: uuid.New()},
		})
		require.NoError(t, err)

		results, total, err := repo.ListPaged(ctx, eventID, "", 1, 10)
		assert.NoError(t, err)
		assert.True(t, total >= 1)
		assert.NotEmpty(t, results)
		assert.Equal(t, eventID, results[0].EventID)
	})

	t.Run("Delete", func(t *testing.T) {
		eventID := SeedEvent(t)

		url := "http://www.google.com"
		name := "Google"
		tier := "Platinum"

		created, _ := repo.Create(ctx, &domain.Sponsor{
			EventID:    eventID,
			Name:       name,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: uuid.New()},
		})

		err := repo.Delete(ctx, created.ID)
		assert.NoError(t, err)

		_, err = repo.GetById(ctx, created.ID)
		assert.Error(t, err)
	})

	t.Run("GetAll - Should return all Sponsors for an event", func(t *testing.T) {
		eventID := SeedEvent(t)

		url := "http://www.google.com"
		name := "Google"
		name2 := "GoogleDevs"
		tier := "Platinum"

		repo.Create(ctx, &domain.Sponsor{
			EventID:    eventID,
			Name:       name,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: uuid.New()},
		})
		repo.Create(ctx, &domain.Sponsor{
			EventID:    eventID,
			Name:       name2,
			LogoURL:    url,
			WebsiteURL: url,
			Tier:       tier,
			Audit:      domain.Audit{CreatedBy: uuid.New()},
		})

		results, err := repo.GetAll(ctx, eventID)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		for _, c := range results {
			assert.Equal(t, eventID, c.EventID)
			assert.NotEmpty(t, c.Name)
		}
	})
}
