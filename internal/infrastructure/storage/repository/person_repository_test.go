package repository

import (
	"context"
	"devfest/internal/domain"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonRepository_Create(t *testing.T) {
	ctx := context.Background()

	cleanTable(t, "persons")

	repo := NewPersonRepository(testQueries)
	createdBy := uuid.New()

	t.Run("Should create a person with all fields", func(t *testing.T) {
		email := "test@devfest.com"
		person := &domain.Person{
			FirstName: "John",
			LastName:  "Doe",
			Email:     &email,
			Audit:     domain.Audit{CreatedBy: createdBy},
		}

		created, err := repo.Create(ctx, person)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, *person.Email, *created.Email)
	})

	t.Run("Should handle optional fields as NULL", func(t *testing.T) {
		person := &domain.Person{
			FirstName: "Minimal",
			LastName:  "User",
			Email:     nil,
			Audit:     domain.Audit{CreatedBy: createdBy},
		}

		created, err := repo.Create(ctx, person)

		require.NoError(t, err)
		assert.Nil(t, created.Email)
	})
}

func TestPersonRepository_Read(t *testing.T) {
	ctx := context.Background()
	repo := NewPersonRepository(testQueries)
	cleanTable(t, "persons")

	email := "jane@devfest.com"
	person, err := repo.Create(ctx, &domain.Person{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     &email,
		Audit:     domain.Audit{CreatedBy: uuid.New()},
	})
	require.NoError(t, err)

	t.Run("GetById - Should return person if exists", func(t *testing.T) {
		found, err := repo.GetById(ctx, person.ID)
		require.NoError(t, err)
		assert.Equal(t, person.ID, found.ID)
		assert.Equal(t, "Jane", found.FirstName)
	})

	t.Run("GetByEmail - Should return person if exists", func(t *testing.T) {
		found, err := repo.GetByEmail(ctx, &email)
		require.NoError(t, err)
		assert.Equal(t, person.ID, found.ID)
	})

	t.Run("GetAll - Should return all persons", func(t *testing.T) {
		list, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, list)
		assert.Equal(t, 1, len(list))
	})
}

func TestPersonRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := NewPersonRepository(testQueries)
	cleanTable(t, "persons")

	email := "original@test.com"
	person, _ := repo.Create(ctx, &domain.Person{
		FirstName: "Original",
		LastName:  "Name",
		Email:     &email,
		Audit:     domain.Audit{CreatedBy: uuid.New()},
	})

	t.Run("Should update fields and audit", func(t *testing.T) {
		newName := "Updated"
		newEditor := uuid.New()
		person.FirstName = newName
		person.Email = nil
		person.Audit.UpdatedBy = newEditor

		updated, err := repo.Update(ctx, person)

		require.NoError(t, err)
		assert.Equal(t, newName, updated.FirstName)
		assert.Nil(t, updated.Email)
		assert.Equal(t, newEditor, updated.Audit.UpdatedBy)
	})
}

func TestPersonRepository_ListPaged(t *testing.T) {
	ctx := context.Background()
	repo := NewPersonRepository(testQueries)
	cleanTable(t, "persons")

	for i := 0; i < 5; i++ {
		email := fmt.Sprintf("user%d@test.com", i)
		_, _ = repo.Create(ctx, &domain.Person{
			FirstName: fmt.Sprintf("User%d", i),
			LastName:  "Test",
			Email:     &email,
			Audit:     domain.Audit{CreatedBy: uuid.New()},
		})
	}

	t.Run("Should return paginated results", func(t *testing.T) {
		list, total, err := repo.ListPaged(ctx, "", 1, 2)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Equal(t, 2, len(list))
	})

	t.Run("Should filter by search term", func(t *testing.T) {
		list, total, err := repo.ListPaged(ctx, "User0", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "User0", list[0].FirstName)
	})
}

func TestPersonRepository_Delete(t *testing.T) {
	ctx := context.Background()
	repo := NewPersonRepository(testQueries)
	cleanTable(t, "persons")

	person, _ := repo.Create(ctx, &domain.Person{
		FirstName: "To Delete",
		LastName:  "User",
		Audit:     domain.Audit{CreatedBy: uuid.New()},
	})

	t.Run("Should delete person", func(t *testing.T) {
		err := repo.Delete(ctx, person.ID)
		assert.NoError(t, err)

		found, err := repo.GetById(ctx, person.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func cleanTable(t *testing.T, tableName string) {
	_, err := testPool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	if err != nil {
		t.Fatalf("failed to truncate table %s: %v", tableName, err)
	}
}
