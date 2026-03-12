package repository

import (
	"context"
	"devfest/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestEvent(t *testing.T, repo *PostgresEventRepository, name, slug string) *domain.Event {
	active := true
	userID := uuid.New()
	event := &domain.Event{
		Name:     name,
		Slug:     slug,
		IsActive: &active,
		Audit:    domain.Audit{CreatedBy: userID, UpdatedBy: userID},
	}
	created, err := repo.Create(context.Background(), event)
	assert.NoError(t, err)
	return created
}

func TestPostgresEventRepository_Create(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)
	ctx := context.Background()
	createdBy := uuid.New()

	t.Run("Successfully create an event", func(t *testing.T) {
		active := true
		event := &domain.Event{
			Name:     "DevFest Aranjuez 2026",
			Slug:     "devfest-2026",
			IsActive: &active,
			Audit:    domain.Audit{CreatedBy: createdBy},
		}

		createdEvent, err := repo.Create(ctx, event)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, createdEvent.ID)
		assert.Equal(t, event.Name, createdEvent.Name)
		assert.Equal(t, event.Slug, createdEvent.Slug)
	})

	t.Run("Fail to create event with duplicate slug", func(t *testing.T) {
		active := true
		event := &domain.Event{
			Name:     "Duplicate Event",
			Slug:     "devfest-2026",
			IsActive: &active,
			Audit:    domain.Audit{CreatedBy: createdBy},
		}

		_, err := repo.Create(ctx, event)
		assert.Error(t, err)
	})
}

func TestPostgresEventRepository_GetByID(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)
	event := createTestEvent(t, repo, "Get By ID Test", "get-id-slug")

	t.Run("Success", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), event.ID)
		assert.NoError(t, err)
		assert.Equal(t, event.ID, found.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), uuid.New())
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestPostgresEventRepository_ListPaged(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)

	createTestEvent(t, repo, "Apple Event", "apple-slug")
	createTestEvent(t, repo, "Banana Event", "banana-slug")

	t.Run("Search by name", func(t *testing.T) {
		events, total, err := repo.ListPaged(context.Background(), "Apple", 1, 10, "name_asc")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "Apple Event", events[0].Name)
	})

	t.Run("Order by name DESC", func(t *testing.T) {
		events, _, err := repo.ListPaged(context.Background(), "Event", 1, 10, "name_desc")
		assert.NoError(t, err)

		if assert.True(t, len(events) >= 2) {
			assert.Equal(t, "Banana Event", events[0].Name)
		}
	})
}

func TestPostgresEventRepository_Update(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)
	event := createTestEvent(t, repo, "Original Name", "update-slug")

	t.Run("Successfully update fields", func(t *testing.T) {
		newName := "Brand New Name"
		event.Name = newName

		updated, err := repo.Update(context.Background(), event)

		assert.NoError(t, err)
		assert.Equal(t, newName, updated.Name)
		assert.Equal(t, event.ID, updated.ID)
	})
}

func TestPostgresEventRepository_Delete(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)
	event := createTestEvent(t, repo, "Delete Me", "delete-slug")

	t.Run("Successfully delete", func(t *testing.T) {
		err := repo.Delete(context.Background(), event.ID)
		assert.NoError(t, err)

		found, err := repo.GetByID(context.Background(), event.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestPostgresEventRepository_FullLifecycle(t *testing.T) {
	repo := NewPostgresEventRepository(testQueries)
	ctx := context.Background()
	userID := uuid.New()
	active := true

	baseEvent := &domain.Event{
		Name:     "Initial Event",
		Slug:     "initial-slug",
		IsActive: &active,
		Audit:    domain.Audit{CreatedBy: userID, UpdatedBy: userID},
	}
	created, err := repo.Create(ctx, baseEvent)
	assert.NoError(t, err)

	// --- TEST GET BY ID ---
	t.Run("GetByID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, created.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.Name, found.Name)
	})

	// --- TEST GET BY SLUG ---
	t.Run("GetBySlug", func(t *testing.T) {
		found, err := repo.GetBySlug(ctx, "initial-slug")
		assert.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	// --- TEST UPDATE ---
	t.Run("Update name and status", func(t *testing.T) {
		newStatus := false
		created.Name = "Updated Name"
		created.IsActive = &newStatus
		created.UpdatedBy = userID

		updated, err := repo.Update(ctx, created)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
		assert.False(t, *updated.IsActive)
	})

	// --- TEST LIST PAGING ---
	t.Run("ListPaged", func(t *testing.T) {
		_, _ = repo.Create(ctx, &domain.Event{
			Name:     "Second Event",
			Slug:     "second-slug",
			IsActive: &active,
			Audit:    domain.Audit{CreatedBy: userID, UpdatedBy: userID},
		})

		events, total, err := repo.ListPaged(ctx, "", 1, 10, "name_asc")
		assert.NoError(t, err)
		assert.True(t, total >= 2)
		assert.NotEmpty(t, events)
	})

	// --- TEST DELETE ---
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, created.ID)
		assert.NoError(t, err)

		found, err := repo.GetByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}
