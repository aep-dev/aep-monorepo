package aepbasetest_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/openapi"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/aepbasetest"
)

func do(t *testing.T, srv *aepbasetest.Server, method, path, contentType, body string) (int, map[string]any) {
	t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, srv.URL+path, r)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s %s response: %v", method, path, err)
	}
	var out map[string]any
	if len(b) > 0 {
		if err := json.Unmarshal(b, &out); err != nil {
			t.Fatalf("%s %s: decoding %q: %v", method, path, b, err)
		}
	}
	return resp.StatusCode, out
}

func TestNewServerRegistersResourceDefinitions(t *testing.T) {
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)

	status, body := do(t, srv, "GET", "/aep-resource-definitions", "", "")
	if status != http.StatusOK {
		t.Fatalf("list resource definitions: expected 200, got %d: %v", status, body)
	}
	got := map[string]bool{}
	for _, r := range body["results"].([]any) {
		got[r.(map[string]any)["singular"].(string)] = true
	}
	for _, want := range []string{"publisher", "book"} {
		if !got[want] {
			t.Errorf("resource definition %q not listed; got %v", want, got)
		}
	}
}

func TestNewServerCRUD(t *testing.T) {
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)

	status, body := do(t, srv, "POST", "/publishers?id=acme", "application/json", `{"name": "Acme", "location": "New York"}`)
	if status/100 != 2 {
		t.Fatalf("create publisher: expected 2xx, got %d: %v", status, body)
	}
	status, body = do(t, srv, "POST", "/publishers/acme/books?id=gatsby", "application/json", `{"title": "The Great Gatsby"}`)
	if status/100 != 2 {
		t.Fatalf("create book: expected 2xx, got %d: %v", status, body)
	}
	if body["path"] != "publishers/acme/books/gatsby" {
		t.Errorf("create book: expected path publishers/acme/books/gatsby, got %v", body["path"])
	}

	status, body = do(t, srv, "PATCH", "/publishers/acme/books/gatsby", "application/merge-patch+json", `{"author": "F. Scott Fitzgerald"}`)
	if status != http.StatusOK {
		t.Fatalf("update book: expected 200, got %d: %v", status, body)
	}

	status, body = do(t, srv, "GET", "/publishers/acme/books/gatsby", "", "")
	if status != http.StatusOK {
		t.Fatalf("get book: expected 200, got %d: %v", status, body)
	}
	if body["title"] != "The Great Gatsby" || body["author"] != "F. Scott Fitzgerald" {
		t.Errorf("get book: unexpected body %v", body)
	}

	status, body = do(t, srv, "GET", "/publishers/acme/books", "", "")
	if status != http.StatusOK {
		t.Fatalf("list books: expected 200, got %d: %v", status, body)
	}
	if n := len(body["results"].([]any)); n != 1 {
		t.Errorf("list books: expected 1 result, got %d", n)
	}

	status, _ = do(t, srv, "DELETE", "/publishers/acme/books/gatsby", "", "")
	if status/100 != 2 {
		t.Fatalf("delete book: expected 2xx, got %d", status)
	}
	status, _ = do(t, srv, "GET", "/publishers/acme/books/gatsby", "", "")
	if status != http.StatusNotFound {
		t.Errorf("get deleted book: expected 404, got %d", status)
	}
}

// TestOpenAPIRoundTrip checks that clients built on aep-lib-go see every
// standard method aepbase serves.
func TestOpenAPIRoundTrip(t *testing.T) {
	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)

	oas, err := openapi.FetchOpenAPI(srv.OpenAPIURL())
	if err != nil {
		t.Fatalf("FetchOpenAPI: %v", err)
	}
	if len(oas.Servers) == 0 || oas.Servers[0].URL != srv.URL {
		t.Errorf("expected OpenAPI server URL %q, got %v", srv.URL, oas.Servers)
	}
	a, err := api.GetAPI(oas, "", "")
	if err != nil {
		t.Fatalf("GetAPI: %v", err)
	}
	for _, singular := range []string{"publisher", "book"} {
		r, ok := a.Resources[singular]
		if !ok {
			t.Errorf("resource %q missing from API", singular)
			continue
		}
		m := r.Methods
		if m.Create == nil || m.Get == nil || m.List == nil || m.Update == nil || m.Delete == nil {
			t.Errorf("resource %q: expected create, get, list, update, and delete methods, got %+v", singular, m)
		}
	}
}

func TestServersAreIsolated(t *testing.T) {
	a := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)
	b := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)

	if status, body := do(t, a, "POST", "/publishers?id=acme", "application/json", `{"name": "Acme"}`); status/100 != 2 {
		t.Fatalf("create publisher: expected 2xx, got %d: %v", status, body)
	}
	if status, _ := do(t, b, "GET", "/publishers/acme", "", ""); status != http.StatusNotFound {
		t.Errorf("expected publisher to exist only on the first server, got %d from the second", status)
	}
}
