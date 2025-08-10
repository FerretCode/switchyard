-- name: GetService :one
SELECT
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count
FROM services
WHERE service_id = $1
LIMIT 1;

-- name: SetServiceEnabled :one
UPDATE services
SET enabled = $1
WHERE service_id = $2
RETURNING
    service_id,
    job_name,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count;

-- name: ListServices :many
SELECT
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count
FROM services
ORDER BY service_id;

-- name: CreateService :one
INSERT INTO services (
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11
)
RETURNING
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count;

-- name: UpdateService :one
UPDATE services
SET
    job_name = COALESCE($2, job_name),
    enabled = COALESCE($3, enabled),
    railway_memory_upscale_threshold = COALESCE($4, railway_memory_upscale_threshold),
    railway_cpu_upscale_threshold = COALESCE($5, railway_cpu_upscale_threshold),
    railway_memory_downscale_threshold = COALESCE($6, railway_memory_downscale_threshold),
    railway_cpu_downscale_threshold = COALESCE($7, railway_cpu_downscale_threshold),
    upscale_cooldown = COALESCE($8, upscale_cooldown),
    downscale_cooldown = COALESCE($9, downscale_cooldown),
    min_replica_count = COALESCE($10, min_replica_count),
    max_replica_count = COALESCE($11, max_replica_count)
WHERE service_id = $1
RETURNING
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count;

-- name: DeleteService :exec
DELETE FROM services
WHERE service_id = $1;

-- name: GetServicesByJobName :many
SELECT
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count
FROM services
WHERE job_name = $1
ORDER BY service_id;

-- name: ListServicesWithJobs :many
SELECT
    service_id,
    job_name,
    enabled,
    railway_memory_upscale_threshold,
    railway_cpu_upscale_threshold,
    railway_memory_downscale_threshold,
    railway_cpu_downscale_threshold,
    upscale_cooldown,
    downscale_cooldown,
    min_replica_count,
    max_replica_count
FROM services
WHERE job_name IS NOT NULL
ORDER BY service_id;
