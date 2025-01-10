-- name: GetJobByID :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: GetJobsByStatus :many
SELECT * FROM jobs
WHERE status = $1
ORDER BY created_at DESC;

-- name: CreateJob :one
INSERT INTO jobs (id, job_type, payload)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateJobAndProcess :one
INSERT INTO jobs (id, job_type, payload, status)
VALUES ($1, $2, $3, 'processing')
RETURNING *;

-- name: FinishJob :exec
UPDATE jobs
SET status = $2, result = $3
WHERE id = $1;

-- name: UpdateJobStatus :exec
UPDATE jobs
SET status = $2
WHERE id = $1;
