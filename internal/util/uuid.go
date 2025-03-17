package util

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUuid(idString string) (*pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(idString)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid")
	}

	return &id, nil
}
