-- +goose Up
-- +goose StatementBegin
ALTER TABLE services
    ALTER COLUMN railway_memory_upscale_threshold TYPE DOUBLE PRECISION USING railway_memory_upscale_threshold::DOUBLE PRECISION,
    ALTER COLUMN railway_cpu_upscale_threshold TYPE DOUBLE PRECISION USING railway_cpu_upscale_threshold::DOUBLE PRECISION,
    ALTER COLUMN railway_memory_downscale_threshold TYPE DOUBLE PRECISION USING railway_memory_downscale_threshold::DOUBLE PRECISION,
    ALTER COLUMN railway_cpu_downscale_threshold TYPE DOUBLE PRECISION USING railway_cpu_downscale_threshold::DOUBLE PRECISION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services
    ALTER COLUMN railway_memory_upscale_threshold TYPE NUMERIC(5, 2) USING railway_memory_upscale_threshold::NUMERIC(5, 2),
    ALTER COLUMN railway_cpu_upscale_threshold TYPE NUMERIC(5, 2) USING railway_cpu_upscale_threshold::NUMERIC(5, 2),
    ALTER COLUMN railway_memory_downscale_threshold TYPE NUMERIC(5, 2) USING railway_memory_downscale_threshold::NUMERIC(5, 2),
    ALTER COLUMN railway_cpu_downscale_threshold TYPE NUMERIC(5, 2) USING railway_cpu_downscale_threshold::NUMERIC(5, 2);
-- +goose StatementEnd
