package service

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/openapi"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/aepbasetest"
)

// newAepbaseServiceCommand starts an aepbase server with the bookstore
// resources and returns a ServiceCommand built from its OpenAPI spec, the same
// way the aepcli binary does.
func newAepbaseServiceCommand(t *testing.T) *ServiceCommand {
	t.Helper()
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)
	oas, err := openapi.FetchOpenAPI(srv.OpenAPIURL())
	if err != nil {
		t.Fatalf("FetchOpenAPI: %v", err)
	}
	a, err := api.GetAPI(oas, "", "")
	if err != nil {
		t.Fatalf("GetAPI: %v", err)
	}
	svc, err := NewServiceCommand(a, nil, false, false, false, "")
	if err != nil {
		t.Fatalf("NewServiceCommand: %v", err)
	}
	return svc
}

// execute runs args and returns the decoded response body, failing the test
// unless the request returned wantStatus. A wantStatus of 0 accepts any 2xx
// status.
func execute(t *testing.T, svc *ServiceCommand, wantStatus int, args ...string) map[string]any {
	t.Helper()
	result, err := svc.Execute(args)
	if err != nil {
		t.Fatalf("Execute(%q): %v", args, err)
	}
	if wantStatus == 0 && result.StatusCode/100 != 2 || wantStatus != 0 && result.StatusCode != wantStatus {
		t.Fatalf("Execute(%q): expected status %d, got %d: %s", args, wantStatus, result.StatusCode, result.Output)
	}
	var body map[string]any
	if result.Output != "" {
		if err := json.Unmarshal([]byte(result.Output), &body); err != nil {
			t.Fatalf("Execute(%q): decoding %q: %v", args, result.Output, err)
		}
	}
	return body
}

func TestAepbase_ResourceLifecycle(t *testing.T) {
	svc := newAepbaseServiceCommand(t)

	// aepbase returns 200 from create while its spec documents 201, so accept
	// any 2xx, as aep-e2e-validator does.
	execute(t, svc, 0, "publisher", "create", "acme", "--name", "Acme Books", "--location", "New York")
	body := execute(t, svc, 0, "book", "--publisher", "acme", "create", "gatsby", "--title", "The Great Gatsby")
	if body["path"] != "publishers/acme/books/gatsby" {
		t.Errorf("create: expected path publishers/acme/books/gatsby, got %v", body["path"])
	}

	execute(t, svc, http.StatusOK, "book", "--publisher", "acme", "update", "gatsby", "--author", "F. Scott Fitzgerald")

	body = execute(t, svc, http.StatusOK, "book", "--publisher", "acme", "get", "gatsby")
	if body["title"] != "The Great Gatsby" || body["author"] != "F. Scott Fitzgerald" {
		t.Errorf("get: unexpected body %v", body)
	}

	body = execute(t, svc, http.StatusOK, "book", "--publisher", "acme", "list")
	if results, _ := body["results"].([]any); len(results) != 1 {
		t.Errorf("list: expected 1 result, got %v", body["results"])
	}

	execute(t, svc, http.StatusNoContent, "book", "--publisher", "acme", "delete", "gatsby")
	execute(t, svc, http.StatusNotFound, "book", "--publisher", "acme", "get", "gatsby")
}

func TestAepbase_ResourceDefinitions(t *testing.T) {
	svc := newAepbaseServiceCommand(t)

	body := execute(t, svc, http.StatusOK, "aep_resource_definition", "list")
	got := map[string]bool{}
	results, _ := body["results"].([]any)
	for _, r := range results {
		got[r.(map[string]any)["singular"].(string)] = true
	}
	if !got["publisher"] || !got["book"] {
		t.Errorf("list: expected publisher and book definitions, got %v", got)
	}
}
