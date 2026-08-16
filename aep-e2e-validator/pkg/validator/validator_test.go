package validator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aep-dev/aep-lib-go/pkg/api"
)

func TestCollectionURL(t *testing.T) {
	a := &api.API{ServerURL: "http://localhost:8000"}
	r := &api.Resource{Plural: "books", API: a}

	tests := []struct {
		name   string
		parent string
		want   string
	}{
		{
			name: "without parent",
			want: "http://localhost:8000/books",
		},
		{
			name:   "with parent",
			parent: "shelves/horror",
			want:   "http://localhost:8000/shelves/horror/books",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{parent: tt.parent}
			if got := v.collectionURL(r); got != tt.want {
				t.Errorf("collectionURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtendedClientDo_InjectsHeaders(t *testing.T) {
	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	headers := []Header{
		{Key: "Authorization", Value: "Bearer tok"},
		{Key: "X-Custom", Value: "v1"},
		{Key: "X-Custom", Value: "v2"},
	}
	client := &extendedClient{inner: &http.Client{}, headers: headers}

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := receivedHeaders.Get("Authorization"); got != "Bearer tok" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer tok")
	}
	vals := receivedHeaders.Values("X-Custom")
	if len(vals) != 2 || vals[0] != "v1" || vals[1] != "v2" {
		t.Errorf("X-Custom = %v, want [v1 v2]", vals)
	}
}

func TestExtendedClientDo_NoHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Custom"); got != "" {
			t.Errorf("unexpected header X-Custom = %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &extendedClient{inner: &http.Client{}, headers: nil}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
