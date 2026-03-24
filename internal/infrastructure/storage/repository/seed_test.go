package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func SeedEvent(t *testing.T) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	id := uuid.New()
	uid := id.String()[:8]

	query := `INSERT INTO events (id, name, slug, is_active, created_by, updated_by) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := testPool.Exec(ctx, query,
		id, "Event "+uid, "slug-"+uid, true, uuid.New(), uuid.New())
	require.NoError(t, err)
	return id
}

func SeedPerson(t *testing.T) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	id := uuid.New()
	uid := id.String()[:8]

	query := `INSERT INTO persons (id, first_name, last_name, email, created_by, updated_by) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	email := fmt.Sprintf("test-%s@devfest.com", uid)
	_, err := testPool.Exec(ctx, query,
		id, "Name-"+uid, "Surname-"+uid, email, uuid.New(), uuid.New())
	require.NoError(t, err)
	return id
}
