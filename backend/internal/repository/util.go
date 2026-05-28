package repository

import (
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// itoa is a short alias used to build positional SQL placeholders.
func itoa(i int) string { return strconv.Itoa(i) }

// nowYear returns the current calendar year, used for order numbering.
func nowYear() int { return time.Now().Year() }

// IsUniqueViolation reports whether err is a PostgreSQL unique-constraint error.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
