package repository

import (
	"database/sql"
	"devfest/internal/domain"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func ParseDBError(err error, entity string) error {
	if err == nil {
		return nil
	}

	// NotFound
	if errors.Is(err, sql.ErrNoRows) {
		return domain.NewAppError(domain.TypeNotFound, entity+" not found", err)
	}

	// Postgers specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // Unique violation
			return domain.NewAppError(domain.TypeAlreadyExists, entity+" already exists", err)
		case "23503": // Foreign key violation
			return domain.NewAppError(domain.TypeBadRequest, "Related entity not found", err)
		}
	}

	// Generic error
	return domain.NewAppError(domain.TypeInternal, "Database error", err)
}
