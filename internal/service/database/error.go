package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsErrConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}
