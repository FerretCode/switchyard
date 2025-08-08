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
