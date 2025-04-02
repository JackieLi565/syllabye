package database

import (
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParsePgUuid(uuid string) (pgtype.UUID, error) {
	var entityUuid pgtype.UUID
	if err := entityUuid.Scan(uuid); err != nil {
		return pgtype.UUID{}, util.ErrMalformed
	}

	return entityUuid, nil
}
