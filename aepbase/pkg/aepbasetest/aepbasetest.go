// Package aepbasetest runs in-memory aepbase servers for tests.
//
// It lets clients of AEP APIs (aepcli, aep-lib-go, aep-e2e-validator, ...)
// test against a real AEP-compliant server instead of canned responses:
//
//	srv := aepbasetest.NewServer(t, aepbasetest.Bookstore()...)
//	oas, err := openapi.FetchOpenAPI(srv.OpenAPIURL())
package aepbasetest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/openapi"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/aepbase"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/db"
	"github.com/aep-dev/aep-monorepo/aepbase/pkg/meta"
)

// Server is an aepbase server backed by an in-memory database.
type Server struct {
	*httptest.Server
	// State is the server's aepbase state, for tests that need to register
	// custom methods or inspect the API directly.
	State *aepbase.State
}

// NewServer starts an aepbase server backed by an in-memory database and
// creates each of defs through the resource definition API, in order, so
// parents must come before their children. The server is closed when the test
// ends.
func NewServer(t testing.TB, defs ...meta.ResourceDefinition) *Server {
	t.Helper()
	d, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("aepbasetest: initializing database: %v", err)
	}
	t.Cleanup(func() { d.Close() })

	// The server URL is part of the generated OpenAPI spec, so the listener
	// must exist before the state is created.
	srv := httptest.NewUnstartedServer(nil)
	state := aepbase.NewState(d, "http://"+srv.Listener.Addr().String())
	srv.Config.Handler = state.Handler()
	srv.Start()
	t.Cleanup(srv.Close)

	s := &Server{Server: srv, State: state}
	for _, def := range defs {
		s.CreateResourceDefinition(t, def)
	}
	return s
}

// OpenAPIURL returns the URL of the server's OpenAPI spec.
func (s *Server) OpenAPIURL() string {
	return s.URL + "/openapi.json"
}

// CreateResourceDefinition creates def through the resource definition API,
// using def.Singular as the ID, and fails the test if the request fails.
func (s *Server) CreateResourceDefinition(t testing.TB, def meta.ResourceDefinition) {
	t.Helper()
	body, err := json.Marshal(def)
	if err != nil {
		t.Fatalf("aepbasetest: marshaling resource definition %q: %v", def.Singular, err)
	}
	u := s.URL + "/aep-resource-definitions?id=" + url.QueryEscape(def.Singular)
	resp, err := s.Client().Post(u, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("aepbasetest: creating resource definition %q: %v", def.Singular, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("aepbasetest: creating resource definition %q: %s: %s", def.Singular, resp.Status, b)
	}
}

// Bookstore returns resource definitions for a publisher and its child book,
// matching the resources used by aepbase's bookstore demo. Both allow
// user-specified IDs.
func Bookstore() []meta.ResourceDefinition {
	return []meta.ResourceDefinition{
		{
			Singular:       "publisher",
			Plural:         "publishers",
			UserSettableId: true,
			Schema: openapi.Schema{
				Properties: map[string]openapi.Schema{
					"name":     {Type: "string"},
					"location": {Type: "string"},
				},
			},
		},
		{
			Singular:       "book",
			Plural:         "books",
			Parents:        []string{"publisher"},
			UserSettableId: true,
			Schema: openapi.Schema{
				Properties: map[string]openapi.Schema{
					"title":  {Type: "string"},
					"author": {Type: "string"},
				},
			},
		},
	}
}
