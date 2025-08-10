# Switchyard API reference

## Autoscaler

### POST /autoscale/upsert-service

Upserts (inserts or updates) an autoscaling service configuration.

#### Request Body:

```json
{
  "service_id": "string",
  "job_name": "string",
  "enabled": true,
  "railway_memory_upscale_threshold": 0.80,
  "railway_cpu_upscale_threshold": 0.80,
  "railway_memory_downscale_threshold": 0.20,
  "railway_cpu_downscale_threshold": 0.20,
  "upscale_cooldown": "1m",
  "downscale_cooldown": "2m",
  "min_replica_count": 1,
  "max_replica_count": 10
}
```

### POST /register-service

Registers a new autoscaling service. Same schema as `/upsert-service`

### DELETE /unregister-service/{id}

Unregisters a service by its ID. If the service is enabled, this route will disable it instead of deleting it to preserve the service configuration. If you want to delete the service, disable the service first, then call this route.

Path parameters:
- `id` (string) - The Railway service ID.

### GET /list-services

Returns a list of registered services and their scaling configurations.

Response schema:

```json
{
  "services": [
    {
      "service_id": "string",
      "project_id": "string",
      "job_name": "string",
      "service_name": "string",
      "environment_name": "string",
      "environment_id": "string",
      "replicas": 3,
      "min_replicas": 1,
      "max_replicas": 5,
      "cpu_upscale_threshold": 0.80,
      "memory_upscale_threshold": 0.80,
      "cpu_downscale_threshold": 0.20,
      "memory_downscale_threshold": 0.20,
      "upscale_cooldown": "1m",
      "downscale_cooldown": "2m",
      "last_scaled_at": "2025-08-10T15:04:05Z",
      "enabled": true
    }
  ]
}
```

### PATCH /set-service-enabled/{id}

Enables or disables a service.

Path parameters:

- `id` (string) - The Railway service ID

Request body:

```json
{ "enabled": true }
```

## Configurator

### POST /configure/{service}

Updates the environment variable configuration and redeploys a given Switchyard service.

Path Parameters:
- `service` (string) - Name of the service

Request body:

`application/json` object that maps environment variable names to values.

Currently supported services:

- scheduler
- autoscale
- feature-flags
- incident
- locomotive

## Feature Flags

### POST /flags/create

Creates a new feature flag

Request body:

```json
{
  "name": "string",
  "enabled": true,
  "rules": [
    { "field": "string", "operator": "string", "value": "string" }
  ]
}
```

Currently supported operators:

- equals
- contains

### GET /flags/get/{name}

Fetches a feature flag by name.

Path parameters:
- `name` (string) - The feature flag name

Response:

```json
{
  "name": "string",
  "enabled": true,
  "rules": [ { "field": "string", "operator": "string", "value": "string" } ]
}
```

### GET /flags/list

Returns a list of feature flags

```json
[
    {
        "name": "string",
        "enabled": true,
        "rules": [ { "field": "string", "operator": "string", "value": "string" } ]
    }
]
```

### PATCH /flags/toggle-feature-flag/{name}?enabled={true/false}

Toggles the enabled status of a feature flag

Path parameters:

- `name` (string) - The feature flag name

Query parameters:

- `enabled` (boolean as string) - "true" or "false"; what to set the feature flag enabled field to

### PATCH /flags/update/{name}

Update a feature flag

Request body:

```json
{
    "name": "string",
    "enabled": true,
    "rules": [
        {
            field: "test-field",
            "operator": "equals",
            "value": "test"
        }
    ]
}
```

### PATCH /flags/upsert-service/{name}

Updates or inserts rules for a feature flag

Request body:

```json
{
  "rules": [ { "field": "string", "operator": "string", "value": "string" } ]
}
```

### DELETE /flags/delete/{name}

Deletes a feature flag.

Path parameters:

- `name` (string) - The feature flag name

### POST /evaluate/{name}

Evaluate whether a feature flag should be enabled for a given user context

Request body:

```json
{ "user_context": { "key": "value" } }
```

Response:

```json
{ "enabled_for_user": true }
```

## Incident

### GET /incident/list-incident-reports

Returns a list of incident reports

Response:

```json
{
    "deployment_reports": [
        {
            "service_id": "string",
            "deployment_id": "string",
            "environment_id": "string",
            "message": "string",
            "timestamp": "unix timestamp (s)"
        }
    ],
    "general_reports": [
        {
            "message": "string",
            "timestamp": "unix timestamp (s)"
        }
    ]
}
```

## Scheduler

For registering services, use the autoscale upsert/insert routes, and specify the job name.

### POST /schedule-job

Request body:

```json
{
    "job_name": "string",
    "job_context": { "key": "value" }
}
```

### GET /get-job-statistics/{name}

Path parameters:

- `name` (string) - Job name
