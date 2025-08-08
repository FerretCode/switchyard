-- +goose Up
-- +goose StatementBegin
CREATE TABLE services (
    service_id VARCHAR(255) NOT NULL PRIMARY KEY,
    job_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE services;
-- +goose StatementEnd
