-- +goose Up
-- +goose StatementBegin
ALTER TABLE services
    ADD COLUMN enabled BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services
    DROP COLUMN enabled;
-- +goose StatementEnd
