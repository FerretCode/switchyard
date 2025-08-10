CREATE TABLE services (
    service_id VARCHAR(255) NOT NULL PRIMARY KEY,
    job_name VARCHAR(255),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,

    railway_memory_upscale_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.80,
    railway_cpu_upscale_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.80,
    railway_memory_downscale_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.20,
    railway_cpu_downscale_threshold DOUBLE PRECISION NOT NULL DEFAULT 0.20,
    upscale_cooldown VARCHAR(255) NOT NULL DEFAULT '1m',
    downscale_cooldown VARCHAR(255) NOT NULL DEFAULT '2m',
    min_replica_count INTEGER NOT NULL DEFAULT 1,
    max_replica_count INTEGER NOT NULL DEFAULT 10
);


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