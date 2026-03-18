package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// uuidPtrToString converts pgtype.UUID to *string
func UuidPtrToString(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}

	id, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return nil
	}

	s := id.String()
	return &s
}

// stringToNullUUID converts *string to pgtype.UUID
func StringToNullUUID(s *string) pgtype.UUID {
	if s == nil || *s == "" {
		return pgtype.UUID{Valid: false}
	}

	id, err := uuid.Parse(*s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}
}

// ToPgBool converts a *bool to pgtype.Bool
func ToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// TextToPtr converts pgtype.Text to *string
func TextToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// PtrToText converts a *string to pgtype.Text
func PtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
