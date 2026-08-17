package aep0144

import (
	"testing"

	"github.com/aep-dev/aep-monorepo/api-linter/lint"
)

func TestAddRules(t *testing.T) {
	if err := AddRules(lint.NewRuleRegistry()); err != nil {
		t.Errorf("AddRules got an error: %v", err)
	}
}
