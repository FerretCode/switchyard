-- name: CreateFeatureFlag :one
INSERT INTO feature_flags (name, enabled)
VALUES ($1, $2)
RETURNING *;

-- name: GetFeatureFlag :one
SELECT * FROM feature_flags
WHERE id = $1;

-- name: GetFeatureFlagByName :one
SELECT * FROM feature_flags
WHERE name = $1;

-- name: ListFeatureFlags :many
SELECT * FROM feature_flags
ORDER BY name;

-- name: ListEnabledFeatureFlags :many
SELECT * FROM feature_flags
WHERE enabled = true
ORDER BY name;

-- name: UpdateFeatureFlag :one
UPDATE feature_flags
SET name = $2, enabled = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateFeatureFlagEnabled :one
UPDATE feature_flags
SET enabled = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteFeatureFlag :exec
DELETE FROM feature_flags
WHERE id = $1;

-- Rules Operations

-- name: CreateRule :one
INSERT INTO rules (feature_flag_id, field, operator, value)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRule :one
SELECT * FROM rules
WHERE id = $1;

-- name: ListRules :many
SELECT * FROM rules
ORDER BY id;

-- name: ListRulesByFeatureFlag :many
SELECT * FROM rules
WHERE feature_flag_id = $1
ORDER BY id;

-- name: UpdateRule :one
UPDATE rules
SET feature_flag_id = $2, field = $3, operator = $4, value = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules
WHERE id = $1;

-- name: DeleteRulesByFeatureFlag :exec
DELETE FROM rules
WHERE feature_flag_id = $1::int4;

-- Feature Flag Rules Junction Table Operations

-- name: AssociateFeatureFlagWithRule :exec
INSERT INTO feature_flags_rules (feature_flag_id, rule_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DisassociateFeatureFlagFromRule :exec
DELETE FROM feature_flags_rules
WHERE feature_flag_id = $1 AND rule_id = $2;

-- name: DisassociateAllRulesFromFeatureFlag :exec
DELETE FROM feature_flags_rules
WHERE feature_flag_id = $1;

-- name: DisassociateRuleFromAllFeatureFlags :exec
DELETE FROM feature_flags_rules
WHERE rule_id = $1;

-- Complex Queries with Joins

-- name: GetFeatureFlagWithRules :many
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
WHERE ff.id = $1
ORDER BY r.id;

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

-- name: ListAllFeatureFlagsWithRules :many
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
ORDER BY ff.name, r.id;

-- name: ListEnabledFeatureFlagsWithRules :many
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
WHERE ff.enabled = true
ORDER BY ff.name, r.id;

-- name: GetRuleWithFeatureFlags :many
SELECT 
    r.id as rule_id,
    r.field as rule_field,
    r.operator as rule_operator,
    r.value as rule_value,
    r.created_at as rule_created_at,
    r.updated_at as rule_updated_at,
    ff.id as feature_flag_id,
    ff.name as feature_flag_name,
    ff.enabled as feature_flag_enabled,
    ff.created_at as feature_flag_created_at,
    ff.updated_at as feature_flag_updated_at
FROM rules r
LEFT JOIN feature_flags_rules ffr ON r.id = ffr.rule_id
LEFT JOIN feature_flags ff ON ffr.feature_flag_id = ff.id
WHERE r.id = $1
ORDER BY ff.name;

-- Utility Queries

-- name: CountFeatureFlags :one
SELECT COUNT(*) FROM feature_flags;

-- name: CountEnabledFeatureFlags :one
SELECT COUNT(*) FROM feature_flags WHERE enabled = true;

-- name: CountRules :one
SELECT COUNT(*) FROM rules;

-- name: CountRulesForFeatureFlag :one
SELECT COUNT(*) FROM feature_flags_rules WHERE feature_flag_id = $1;

-- name: CountFeatureFlagsForRule :one
SELECT COUNT(*) FROM feature_flags_rules WHERE rule_id = $1;

-- name: FeatureFlagExists :one
SELECT EXISTS(SELECT 1 FROM feature_flags WHERE id = $1);

-- name: FeatureFlagExistsByName :one
SELECT EXISTS(SELECT 1 FROM feature_flags WHERE name = $1);

-- name: RuleExists :one
SELECT EXISTS(SELECT 1 FROM rules WHERE id = $1);

-- name: AssociationExists :one
SELECT EXISTS(SELECT 1 FROM feature_flags_rules WHERE feature_flag_id = $1 AND rule_id = $2);

-- Search Queries

-- name: SearchFeatureFlagsByName :many
SELECT * FROM feature_flags
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name;

-- name: SearchRulesByField :many
SELECT * FROM rules
WHERE field ILIKE '%' || $1 || '%'
ORDER BY field, id;

-- name: SearchRulesByOperator :many
SELECT * FROM rules
WHERE operator = $1
ORDER BY field, id;

-- name: SearchRulesByValue :many
SELECT * FROM rules
WHERE value ILIKE '%' || $1 || '%'
ORDER BY field, id;

-- Bulk Operations

-- name: BulkEnableFeatureFlags :exec
UPDATE feature_flags
SET enabled = true, updated_at = CURRENT_TIMESTAMP
WHERE id = ANY($1::int[]);

-- name: BulkDisableFeatureFlags :exec
UPDATE feature_flags
SET enabled = false, updated_at = CURRENT_TIMESTAMP
WHERE id = ANY($1::int[]);

-- name: BulkDeleteFeatureFlags :exec
DELETE FROM feature_flags
WHERE id = ANY($1::int[]);

-- name: BulkDeleteRules :exec
DELETE FROM rules
WHERE id = ANY($1::int[]);

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
