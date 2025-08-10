-- +goose Up
-- +goose StatementBegin
ALTER TABLE services
    ADD COLUMN railway_memory_upscale_threshold NUMERIC(5, 2) NOT NULL DEFAULT 0.80,
    ADD COLUMN railway_cpu_upscale_threshold NUMERIC(5, 2) NOT NULL DEFAULT 0.80,
    ADD COLUMN railway_memory_downscale_threshold NUMERIC(5, 2) NOT NULL DEFAULT 0.20,
    ADD COLUMN railway_cpu_downscale_threshold NUMERIC(5, 2) NOT NULL DEFAULT 0.20,
    ADD COLUMN upscale_cooldown VARCHAR(255) NOT NULL DEFAULT '1m',
    ADD COLUMN downscale_cooldown VARCHAR(255) NOT NULL DEFAULT '2m',
    ADD COLUMN min_replica_count INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN max_replica_count INTEGER NOT NULL DEFAULT 10;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services
    DROP COLUMN railway_memory_upscale_threshold,
    DROP COLUMN railway_cpu_upscale_threshold,
    DROP COLUMN railway_memory_downscale_threshold,
    DROP COLUMN railway_cpu_downscale_threshold,
    DROP COLUMN upscale_cooldown,
    DROP COLUMN downscale_cooldown,
    DROP COLUMN min_replica_count,
    DROP COLUMN max_replica_count;
-- +goose StatementEnd