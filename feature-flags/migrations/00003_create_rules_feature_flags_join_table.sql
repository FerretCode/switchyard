-- +goose Up
-- +goose StatementBegin
CREATE TABLE feature_flags_rules (
    feature_flag_id INTEGER NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    rule_id INTEGER NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (feature_flag_id, rule_id)
);

CREATE INDEX idx_feature_flags_rules_feature_flag_id ON feature_flags_rules(feature_flag_id);
CREATE INDEX idx_feature_flags_rules_rule_id ON feature_flags_rules(rule_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feature_flags_rules;
-- +goose StatementEnd
