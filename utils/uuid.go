package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// StringToPgUUID converts a string ID to pgtype.UUID
func StringToPgUUID(id string) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return pgID, fmt.Errorf("invalid UUID format: %w", err)
	}
	return pgID, nil
}

// StringToGoogleUUID converts a string ID to uuid.UUID
func StringToGoogleUUID(id string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return parsed, nil
}
