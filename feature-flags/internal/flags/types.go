package flags

import "github.com/ferretcode/switchyard/feature-flags/pkg/types"

const PostgresDuplicateKeyErrorCode = "23505"

type GetFlagResponse struct {
	Name    string       `json:"name"`
	Enabled bool         `json:"enabled"`
	Rules   []types.Rule `json:"rules"`
}

type EvaluateUserContext struct {
	UserContext map[string]string `json:"user_context"`
}
type EvaluateFlagResponse struct {
	EnabledForUser bool `json:"enabled_for_user"`
}
