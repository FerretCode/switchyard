-- +goose Up
CREATE TABLE job_receipts (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL,
    retry_count INTEGER NOT NULL DEFAULT 0,
    message VARCHAR(255) NOT NULL DEFAULT '',
    job_name TEXT NOT NULL,
    job_context JSONB NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- +goose Down
DROP TABLE job_receipts;
