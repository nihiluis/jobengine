// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
	JobStatusRetrying   JobStatus = "retrying"
)

func (e *JobStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobStatus(s)
	case string:
		*e = JobStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for JobStatus: %T", src)
	}
	return nil
}

type NullJobStatus struct {
	JobStatus JobStatus
	Valid     bool // Valid is true if JobStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullJobStatus) Scan(value interface{}) error {
	if value == nil {
		ns.JobStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.JobStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullJobStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.JobStatus), nil
}

// Stores background jobs and their execution status
type Job struct {
	ID uuid.UUID
	// Type of job (e.g., email, image_processing)
	JobType string
	// Current job status (pending, processing, completed, failed)
	Status JobStatus
	// Input data for the job in JSON format
	Payload []byte
	// Output/result data from the job execution
	Result      []byte
	OutMessage  pgtype.Text
	RetryCount  int32
	StartedAt   pgtype.Timestamptz
	CompletedAt pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	CreatedBy   pgtype.Text
	LockedBy    pgtype.Text
	// Optimistic locking version number
	Version int32
}
