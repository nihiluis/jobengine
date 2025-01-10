package database

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// stringToUUID converts a string ID to pgtype.UUID
func stringToUUID(id string) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return pgID, fmt.Errorf("invalid UUID format: %w", err)
	}
	return pgID, nil
}
