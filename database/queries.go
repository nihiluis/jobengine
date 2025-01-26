package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nihiluis/jobengine/database/queries"
)

// QueriesImpl wraps the generated queries and adds custom functionality
type QueriesImpl struct {
	db *DB
}

// NewQueries creates a new Queries instance
func NewQueries(db *DB) *QueriesImpl {
	return &QueriesImpl{
		db: db,
	}
}

// GetJobByID wraps the generated GetJobByID query
func (q *QueriesImpl) GetJobByID(ctx context.Context, id string) (*queries.Job, error) {
	uuid, err := stringToGoogleUUID(id)
	if err != nil {
		return nil, err
	}

	job, err := q.db.queries.GetJobByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetJobsByStatus wraps the generated GetJobsByStatus query
func (q *QueriesImpl) GetJobsByStatus(ctx context.Context, status queries.JobStatus) ([]queries.Job, error) {
	return q.db.queries.GetJobsByStatus(ctx, status)
}

// CreateJob creates a new job in the database
func (q *QueriesImpl) CreateJob(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	params := queries.CreateJobParams{
		ID:      uuid.New(),
		JobType: jobType,
		Payload: payloadBytes,
	}

	job, err := q.db.queries.CreateJob(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &job, nil
}

func (q *QueriesImpl) CreateJobAndProcess(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	params := queries.CreateJobAndProcessParams{
		ID:      uuid.New(),
		JobType: jobType,
		Payload: payloadBytes,
	}

	job, err := q.db.queries.CreateJobAndProcess(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &job, nil
}

func (q *QueriesImpl) FinishJob(ctx context.Context, jobIDStr string, status string, message string, result map[string]any) error {
	jobID, err := stringToGoogleUUID(jobIDStr)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	var jobStatus queries.JobStatus
	err = jobStatus.Scan(status)
	if err != nil {
		return fmt.Errorf("invalid job status: %w", err)
	}

	if jobStatus != queries.JobStatusCompleted &&
		jobStatus != queries.JobStatusFailed {
		return fmt.Errorf("invalid job status: %s", status)
	}

	payloadBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	return q.db.queries.FinishJob(ctx, queries.FinishJobParams{
		ID:         jobID,
		Result:     payloadBytes,
		Status:     jobStatus,
		OutMessage: pgtype.Text{String: message, Valid: true},
	})
}
