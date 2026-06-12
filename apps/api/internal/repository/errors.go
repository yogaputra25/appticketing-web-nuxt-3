package repository

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn" //nolint:unused // optional dep
	"gorm.io/gorm"
)

// isUniqueViolation detects a Postgres unique constraint error.
// We import pgx via the gorm postgres driver — if unavailable, fall back to string match.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// Direct pgx error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	// Fallback string match (covers gorm wrapped errors)
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "23505")
}

// FindOne helper — re-export gorm.ErrRecordNotFound for handlers.
var ErrRecordNotFound = gorm.ErrRecordNotFound
