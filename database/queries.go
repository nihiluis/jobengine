package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nihiluis/jobengine/database/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

// Queries wraps the generated queries and adds custom functionality
type Queries struct {
	db *DB
}

// NewQueries creates a new Queries instance
func NewQueries(db *DB) *Queries {
	return &Queries{
		db: db,
	}
}

// GetJobByID wraps the generated GetJobByID query
func (q *Queries) GetJobByID(ctx context.Context, id string) (queries.Job, error) {
	// Convert string ID to pgtype.UUID
	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return queries.Job{}, err
	}

	return q.db.queries.GetJobByID(ctx, pgID)
}

// GetJobsByStatus wraps the generated GetJobsByStatus query
func (q *Queries) GetJobsByStatus(ctx context.Context, status queries.JobStatus) ([]queries.Job, error) {
	return q.db.queries.GetJobsByStatus(ctx, status)
}

// CreateJob creates a new job in the database
func (q *Queries) CreateJob(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	params := queries.CreateJobParams{
		JobType: jobType,
		Payload: payloadBytes,
	}

	job, err := q.db.queries.CreateJob(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return &job, nil
}
