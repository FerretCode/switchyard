-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_reports (
    id SERIAL PRIMARY KEY,
    service_id TEXT NOT NULL,
    deployment_id TEXT NOT NULL,
    environment_id TEXT NOT NULL,
    message TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);

CREATE INDEX idx_incident_reports_service_id ON incident_reports(service_id);
CREATE INDEX idx_incident_reports_environment_id ON incident_reports(environment_id);
CREATE INDEX idx_incident_reports_timestamp ON incident_reports(timestamp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incident_reports;
-- +goose StatementEnd
