-- Create enum type for job status
CREATE TYPE job_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'cancelled',
    'retrying'
);

-- Modify existing jobs table or create it with enum
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type VARCHAR(50) NOT NULL,
    status job_status NOT NULL DEFAULT 'pending',
    payload JSONB NOT NULL,
    result JSONB,
    out_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    locked_by VARCHAR(255),
    version INTEGER NOT NULL DEFAULT 0
);

-- Index for faster job type lookups
CREATE INDEX IF NOT EXISTS idx_jobs_type ON jobs(job_type);

-- Index for status-based queries
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);

-- Trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_jobs_updated_at
    BEFORE UPDATE ON jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments to describe the table and important columns
COMMENT ON TABLE jobs IS 'Stores background jobs and their execution status';
COMMENT ON COLUMN jobs.job_type IS 'Type of job (e.g., email, image_processing)';
COMMENT ON COLUMN jobs.status IS 'Current job status (pending, processing, completed, failed)';
COMMENT ON COLUMN jobs.payload IS 'Input data for the job in JSON format';
COMMENT ON COLUMN jobs.result IS 'Output/result data from the job execution';
COMMENT ON COLUMN jobs.version IS 'Optimistic locking version number';
