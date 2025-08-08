-- +goose Up
-- +goose StatementBegin
CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    feature_flag_id INTEGER,
    field VARCHAR(255) NOT NULL,
    operator VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rules
-- +goose StatementEnd
