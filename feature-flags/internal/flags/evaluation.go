package flags

import (
	"strings"

	"github.com/ferretcode/switchyard/feature-flags/pkg/types"
)

func evaluateFlag(flag types.Flag, context map[string]string) bool {
	if !flag.Enabled {
		return false
	}

	for _, rule := range flag.Rules {
		value, ok := context[rule.Field]
		if !ok {
			continue
		}

		switch rule.Operator {
		case "equals":
			if value == rule.Value {
				return true
			}
		case "contains":
			return strings.Contains(value, rule.Value)
		}
	}

	return false
}
