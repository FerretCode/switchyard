-- name: CreateIncidentReport :one
INSERT INTO incident_reports (
    service_id,
    deployment_id,
    environment_id,
    message,
    timestamp
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetIncidentReport :one
SELECT *
FROM incident_reports
WHERE id = $1
LIMIT 1;

-- name: ListIncidentReports :many
SELECT *
FROM incident_reports
ORDER BY timestamp DESC
LIMIT $1 OFFSET $2;

-- name: ListIncidentReportsByService :many
SELECT *
FROM incident_reports
WHERE service_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: DeleteIncidentReport :exec
DELETE FROM incident_reports
WHERE id = $1;

-- name: ListIncidentReportsWithServiceID :many
SELECT
    service_id,
    deployment_id,
    environment_id,
    message,
    timestamp
FROM incident_reports
WHERE service_id IS NOT NULL AND service_id <> ''
ORDER BY timestamp DESC
LIMIT $1;

-- name: ListIncidentReportsWithoutServiceID :many
SELECT
    service_id,
    deployment_id,
    environment_id,
    message,
    timestamp
FROM incident_reports
WHERE service_id IS NULL OR service_id = ''
ORDER BY timestamp DESC
LIMIT $1;
