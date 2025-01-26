package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// stringToUUID converts a string ID to pgtype.UUID
func stringToPgUUID(id string) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return pgID, fmt.Errorf("invalid UUID format: %w", err)
	}
	return pgID, nil
}

func stringToGoogleUUID(id string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return parsed, nil
}
