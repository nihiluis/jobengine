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