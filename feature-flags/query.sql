-- name: CreateFeatureFlag :one
INSERT INTO feature_flags (name, enabled)
VALUES ($1, $2)
RETURNING *;

-- name: GetFeatureFlagByName :one
SELECT * FROM feature_flags
WHERE name = $1;

-- name: DeleteFeatureFlag :exec
DELETE FROM feature_flags
WHERE id = $1;


-- name: GetFeatureFlagByNameWithRules :many
SELECT 
    ff.id as feature_flag_id,
    ff.name as feature_flag_name,
    ff.enabled as feature_flag_enabled,
    ff.created_at as feature_flag_created_at,
    ff.updated_at as feature_flag_updated_at,
    r.id as rule_id,
    r.field as rule_field,
    r.operator as rule_operator,
    r.value as rule_value,
    r.created_at as rule_created_at,
    r.updated_at as rule_updated_at
FROM feature_flags ff
LEFT JOIN feature_flags_rules ffr ON ff.id = ffr.feature_flag_id
LEFT JOIN rules r ON ffr.rule_id = r.id
WHERE ff.name = $1
ORDER BY r.id;

-- name: BulkAssociateFeatureFlagWithRules :exec
INSERT INTO feature_flags_rules (feature_flag_id, rule_id)
SELECT $1, unnest($2::int[])
ON CONFLICT DO NOTHING;

-- name: BulkCreateRulesForFeatureFlag :many
INSERT INTO rules (feature_flag_id, field, operator, value)
SELECT 
    $1::int as feature_flag_id,
    unnest($2::text[]) as field,
    unnest($3::text[]) as operator,
    unnest($4::text[]) as value
RETURNING *;

-- name: UpsertFeatureFlagByNameWithRules :one
WITH upserted_flag AS (
    INSERT INTO feature_flags (name, enabled)
    VALUES ($1, $2)
    ON CONFLICT (name) DO UPDATE
        SET enabled = EXCLUDED.enabled,
            updated_at = CURRENT_TIMESTAMP
    RETURNING id, name, enabled, created_at, updated_at
),
deleted_associations AS (
    DELETE FROM feature_flags_rules
    WHERE feature_flag_id = (SELECT id FROM upserted_flag)
),
deleted_rules AS (
    DELETE FROM rules
    WHERE feature_flag_id = (SELECT id FROM upserted_flag)
),
inserted_rules AS (
    INSERT INTO rules (feature_flag_id, field, operator, value)
    SELECT 
        (SELECT id FROM upserted_flag),
        unnest($3::text[]),
        unnest($4::text[]),
        unnest($5::text[])
    WHERE cardinality($3::text[]) > 0
    RETURNING id, feature_flag_id
),
inserted_associations AS (
    INSERT INTO feature_flags_rules (feature_flag_id, rule_id)
    SELECT ir.feature_flag_id, ir.id FROM inserted_rules ir
    RETURNING feature_flag_id, rule_id
)
SELECT * FROM upserted_flag;
