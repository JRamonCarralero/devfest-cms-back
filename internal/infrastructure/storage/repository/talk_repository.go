package repository

import (
	"devfest/internal/infrastructure/storage/dbgen"
)

type talkRepository struct {
	queries *dbgen.Queries
}

// NewTalkRepository returns a new TalkRepository
func NewTalkRepository(queries *dbgen.Queries) *talkRepository {
	return &talkRepository{queries: queries}
}

// --- READERS ---

// GetAll
