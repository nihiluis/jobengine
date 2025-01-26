-- +goose Up
-- Add out_message column to jobs table
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS out_message TEXT NOT NULL DEFAULT '';

-- Add comment for the new column
COMMENT ON COLUMN jobs.out_message IS 'Output message or error details from job execution';

-- +goose Down
-- Remove out_message column from jobs table
ALTER TABLE jobs DROP COLUMN IF EXISTS out_message;
