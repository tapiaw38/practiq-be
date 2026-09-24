-- Submit jobs table for persisting async submit job state
CREATE TABLE IF NOT EXISTS submit_jobs (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('practice_sheet', 'notebook')),
    status TEXT NOT NULL,
    error_code TEXT,
    message TEXT,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on created_at for potential cleanup queries
CREATE INDEX IF NOT EXISTS idx_submit_jobs_created_at ON submit_jobs(created_at);
