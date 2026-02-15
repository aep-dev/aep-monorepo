package tests

import (
	"net/http"

	"github.com/aep-dev/aep-lib-go/pkg/api"
)

type ValidationActions interface {
	CreateResource(r *api.Resource, collectionURL string, payload map[string]interface{}) (map[string]interface{}, error)
	List(url string) (*ListResponse, error)
	Post(url string, body interface{}) (*http.Response, error)
	Patch(url string, body interface{}) (*http.Response, error)
	Get(url string) (map[string]interface{}, error)
	Delete(url string) error
	DeleteReq(url string) (*http.Response, error)
	GenerateID() string
}

type ListResponse struct {
	Resources     []map[string]interface{} `json:"-"`
	NextPageToken string                   `json:"next_page_token"`
}

type ValidationContext struct {
	Resource      *api.Resource
	CollectionURL string
	Resources     []map[string]interface{}
	ListResponse1 *ListResponse
}

type Test struct {
	Name         string
	Precondition func(*ValidationContext) error
	Setup        func(ValidationActions, *ValidationContext) error
	Run          func(ValidationActions, *ValidationContext) error
	Teardown     func(ValidationActions, *ValidationContext) error
}
