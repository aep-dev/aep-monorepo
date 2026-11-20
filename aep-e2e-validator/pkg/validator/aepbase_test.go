package validator

import (
	"testing"

	"github.com/aep-dev/aep-monorepo/aep-e2e-validator/pkg/tests"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/aepbasetest"
)

// runAgainstAepbase runs v against an aepbase server and fails the test unless
// every test in the suite ran and passed.
func runAgainstAepbase(t *testing.T, v *Validator) {
	t.Helper()
	results, code := v.runTests()
	if results == nil {
		t.Fatalf("validator could not run, exit code %d", code)
	}
	if want := len(tests.NewTests()); len(results) != want {
		t.Errorf("expected %d test results, got %d", want, len(results))
	}
	for _, r := range results {
		if r.Status != StatusPass {
			t.Errorf("%s: %s: %s", r.Name, r.Status, r.Detail)
			for _, l := range r.RequestLogs {
				t.Logf("  %s %s %s -> %d %s", l.Method, l.URL, l.ReqBody, l.RespCode, l.RespBody)
			}
		}
	}
	if code != ExitCodeSuccess {
		t.Errorf("expected exit code %d, got %d", ExitCodeSuccess, code)
	}
}

// TestAepbase_TopLevelCollection verifies that aepbase passes the validation
// suite for a top-level collection.
func TestAepbase_TopLevelCollection(t *testing.T) {
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)
	runAgainstAepbase(t, NewValidator(srv.OpenAPIURL(), "publishers", false, "", nil, nil, true))
}

// TestAepbase_ChildCollection verifies that aepbase passes the validation
// suite for a collection nested under a parent resource.
func TestAepbase_ChildCollection(t *testing.T) {
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)

	parent := NewValidator(srv.OpenAPIURL(), "publishers", false, "", nil, nil, true)
	if _, err := parent.Post(srv.URL+"/publishers?id=acme", map[string]any{"name": "Acme Books"}); err != nil {
		t.Fatalf("creating parent publisher: %v", err)
	}

	runAgainstAepbase(t, NewValidator(srv.OpenAPIURL(), "books", false, "publishers/acme", nil, nil, true))
}
