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
