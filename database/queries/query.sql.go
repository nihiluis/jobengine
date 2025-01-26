// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (id, job_type, payload)
VALUES ($1, $2, $3)
RETURNING id, job_type, status, payload, result, out_message, retry_count, started_at, completed_at, created_at, updated_at, created_by, locked_by, version
`

type CreateJobParams struct {
	ID      uuid.UUID
	JobType string
	Payload []byte
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRow(ctx, createJob, arg.ID, arg.JobType, arg.Payload)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.JobType,
		&i.Status,
		&i.Payload,
		&i.Result,
		&i.OutMessage,
		&i.RetryCount,
		&i.StartedAt,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.LockedBy,
		&i.Version,
	)
	return i, err
}

const createJobAndProcess = `-- name: CreateJobAndProcess :one
INSERT INTO jobs (id, job_type, payload, status)
VALUES ($1, $2, $3, 'processing')
RETURNING id, job_type, status, payload, result, out_message, retry_count, started_at, completed_at, created_at, updated_at, created_by, locked_by, version
`

type CreateJobAndProcessParams struct {
	ID      uuid.UUID
	JobType string
	Payload []byte
}

func (q *Queries) CreateJobAndProcess(ctx context.Context, arg CreateJobAndProcessParams) (Job, error) {
	row := q.db.QueryRow(ctx, createJobAndProcess, arg.ID, arg.JobType, arg.Payload)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.JobType,
		&i.Status,
		&i.Payload,
		&i.Result,
		&i.OutMessage,
		&i.RetryCount,
		&i.StartedAt,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.LockedBy,
		&i.Version,
	)
	return i, err
}

const finishJob = `-- name: FinishJob :exec
UPDATE jobs
SET status = $2, result = $3, out_message = $4
WHERE id = $1
`

type FinishJobParams struct {
	ID         uuid.UUID
	Status     JobStatus
	Result     []byte
	OutMessage pgtype.Text
}

func (q *Queries) FinishJob(ctx context.Context, arg FinishJobParams) error {
	_, err := q.db.Exec(ctx, finishJob,
		arg.ID,
		arg.Status,
		arg.Result,
		arg.OutMessage,
	)
	return err
}

const getJobByID = `-- name: GetJobByID :one
SELECT id, job_type, status, payload, result, out_message, retry_count, started_at, completed_at, created_at, updated_at, created_by, locked_by, version FROM jobs
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetJobByID(ctx context.Context, id uuid.UUID) (Job, error) {
	row := q.db.QueryRow(ctx, getJobByID, id)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.JobType,
		&i.Status,
		&i.Payload,
		&i.Result,
		&i.OutMessage,
		&i.RetryCount,
		&i.StartedAt,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.LockedBy,
		&i.Version,
	)
	return i, err
}

const getJobsByStatus = `-- name: GetJobsByStatus :many
SELECT id, job_type, status, payload, result, out_message, retry_count, started_at, completed_at, created_at, updated_at, created_by, locked_by, version FROM jobs
WHERE status = $1
ORDER BY created_at DESC
`

func (q *Queries) GetJobsByStatus(ctx context.Context, status JobStatus) ([]Job, error) {
	rows, err := q.db.Query(ctx, getJobsByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.JobType,
			&i.Status,
			&i.Payload,
			&i.Result,
			&i.OutMessage,
			&i.RetryCount,
			&i.StartedAt,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.LockedBy,
			&i.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobStatus = `-- name: UpdateJobStatus :exec
UPDATE jobs
SET status = $2
WHERE id = $1
`

type UpdateJobStatusParams struct {
	ID     uuid.UUID
	Status JobStatus
}

func (q *Queries) UpdateJobStatus(ctx context.Context, arg UpdateJobStatusParams) error {
	_, err := q.db.Exec(ctx, updateJobStatus, arg.ID, arg.Status)
	return err
}
