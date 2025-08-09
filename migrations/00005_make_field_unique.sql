-- +goose Up
-- +goose StatementBegin
ALTER TABLE rules
ADD CONSTRAINT unique_field_per_feature_flag UNIQUE (feature_flag_id, field);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rules
DROP CONSTRAINT unique_field_per_feature_flag;
-- +goose StatementEnd
