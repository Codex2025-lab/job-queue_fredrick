-- Worker query
CREATE INDEX idx_jobs_status_created 
ON jobs (status, created_at);

-- JSONB payload
CREATE INDEX idx_jobs_payload 
ON jobs USING GIN (payload);

-- JSONB result
CREATE INDEX idx_jobs_result 
ON jobs USING GIN (result);

-- updated_at
CREATE INDEX idx_jobs_updated_at 
ON jobs (updated_at DESC);