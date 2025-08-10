-- name: GetService :one
SELECT * FROM services
WHERE service_id = $1 LIMIT 1;

-- name: ListServices :many
SELECT * FROM services
ORDER BY service_id;

-- name: CreateService :one
INSERT INTO services (
    service_id, job_name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET job_name = $2
WHERE service_id = $1
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services
WHERE service_id = $1;

-- name: GetServicesByJobName :many
SELECT * FROM services
WHERE job_name = $1
ORDER BY service_id;

-- name: ListServicesWithJobs :many
SELECT * FROM services
WHERE job_name IS NOT NULL
ORDER BY service_id;

-- name: CreateJobReceipt :one
INSERT INTO job_receipts (
    job_id,
    status,
    retry_count,
    message,
    job_name,
    job_context,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetJobReceiptByID :one
SELECT * FROM job_receipts
WHERE id = $1;

-- name: GetJobReceiptByJobID :one
SELECT * FROM job_receipts
WHERE job_id = $1;

-- name: ListJobReceipts :many
SELECT * FROM job_receipts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateJobReceiptByID :one
UPDATE job_receipts
SET
    status = $2,
    retry_count = $3,
    message = $4,
    job_name = $5,
    job_context = $6,
    updated_at = $7
WHERE id = $1
RETURNING *;

-- name: UpdateJobReceiptByJobID :one
UPDATE job_receipts
SET
    status = $2,
    retry_count = $3,
    message = $4,
    job_name = $5,
    job_context = $6,
    updated_at = $7
WHERE job_id = $1
RETURNING *;

-- name: DeleteJobReceiptByID :exec
DELETE FROM job_receipts
WHERE id = $1;

-- name: DeleteJobReceiptByJobID :exec
DELETE FROM job_receipts
WHERE job_id = $1;

-- name: AggregateJobReceiptsByJobID :one
SELECT
    job_name,
    COUNT(*) AS total_receipts,
    COUNT(*) FILTER (WHERE status = 'ok') AS ok_count,
    COUNT(*) FILTER (WHERE status = 'error') AS error_count,
    MIN(retry_count) AS min_retry_count,
    MAX(retry_count) AS max_retry_count,
    AVG(retry_count) AS avg_retry_count,
    MIN(created_at) AS earliest_created_at,
    MAX(updated_at) AS latest_updated_at
FROM job_receipts
WHERE job_name = $1
GROUP BY job_name
LIMIT 1;

-- name: SetServiceJobName :one
UPDATE services
SET job_name = $1
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