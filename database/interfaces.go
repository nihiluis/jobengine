package database

import (
	"context"

	"github.com/nihiluis/jobengine/database/queries"
)

type Queries interface {
	CreateJob(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error)
	CreateJobAndProcess(ctx context.Context, jobType string, payload map[string]any) (*queries.Job, error)
	GetJobByID(ctx context.Context, id string) (*queries.Job, error)
	GetJobsByStatus(ctx context.Context, status queries.JobStatus) ([]queries.Job, error)
	FinishJob(ctx context.Context, jobIDStr string, status string, message string, result map[string]any) error
}
