// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// isUniqueViolation checks if an error is a PostgreSQL unique constraint violation.
// constraintName is optional - if provided, it also checks the constraint name.
func isUniqueViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if !As(err, &pgErr) {
		return false
	}

	// 23505 is the PostgreSQL error code for unique_violation
	if pgErr.Code != "23505" {
		return false
	}

	if constraintName == "" {
		return true
	}

	return strings.Contains(pgErr.ConstraintName, constraintName)
}

// As is a wrapper for errors.As to work with pgconn.PgError.
func As(err error, target interface{}) bool {
	if err == nil {
		return false
	}

	if t, ok := target.(**pgconn.PgError); ok {
		for err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				*t = pgErr
				return true
			}
			// Try to unwrap
			if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
				err = unwrapper.Unwrap()
			} else {
				return false
			}
		}
	}
	return false
}
