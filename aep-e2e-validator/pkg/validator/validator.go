package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aep-dev/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-lib-go/pkg/openapi"
)

type Validator struct {
	configPath     string
	collection     string
	allCollections bool
	client         *http.Client
}

func NewValidator(configPath, collection string, allCollections bool) *Validator {
	return &Validator{
		configPath:     configPath,
		collection:     collection,
		allCollections: allCollections,
		client:         &http.Client{},
	}
}

func (v *Validator) Run() error {
	doc, err := openapi.FetchOpenAPI(v.configPath)
	if err != nil {
		return fmt.Errorf("failed to fetch OpenAPI spec: %w", err)
	}

	serverURL := ""
	if len(doc.Servers) > 0 {
		serverURL = doc.Servers[0].URL
	}

	aepAPI, err := api.GetAPI(doc, serverURL, "")
	if err != nil {
		return fmt.Errorf("failed to parse API: %w", err)
	}

	if err := api.AddImplicitFieldsAndValidate(aepAPI); err != nil {
		return fmt.Errorf("failed to validate API: %w", err)
	}

	if v.allCollections {
		for _, r := range aepAPI.Resources {
			// Only validate top-level resources
			if len(r.Parents) == 0 {
				if err := v.validateResource(r); err != nil {
					log.Printf("Validation failed for resource %s: %v", r.Singular, err)
					return err
				}
			}
		}
	} else {
		var targetResource *api.Resource
		for _, r := range aepAPI.Resources {
			if r.Plural == v.collection {
				targetResource = r
				break
			}
		}
		if targetResource == nil {
			return fmt.Errorf("collection %s not found in API", v.collection)
		}
		if err := v.validateResource(targetResource); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateResource(r *api.Resource) error {
	fmt.Printf("Starting validation for resource: %s\n", r.Singular)
	ctx := &ValidationContext{
		Resource:      r,
		CollectionURL: fmt.Sprintf("%s/%s", r.API.ServerURL, r.Plural),
	}

	steps := []ValidationStep{
		{Name: "Create first resource", Run: v.stepCreateFirstResource},
		{Name: "Duplicate creation check", Run: v.stepDuplicateCreationCheck},
		{Name: "Create second resource", Run: v.stepCreateSecondResource},
		{Name: "List resources with limit=1", Run: v.stepListResourcesLimit1},
		{Name: "List resources with page token", Run: v.stepListResourcesPageToken},
		{Name: "Update first resource", Run: v.stepUpdateFirstResource},
		{Name: "Get first resource", Run: v.stepGetFirstResource},
		{Name: "Delete first resource", Run: v.stepDeleteFirstResource},
		{Name: "Duplicate delete check", Run: v.stepDuplicateDeleteCheck},
		{Name: "Get deleted resource", Run: v.stepGetDeletedResource},
		{Name: "Delete second resource", Run: v.stepDeleteSecondResource},
		{Name: "List to verify empty", Run: v.stepListVerifyEmpty},
	}

	for i, step := range steps {
		fmt.Printf("%d. %s...\n", i+1, step.Name)
		if err := step.Run(ctx); err != nil {
			return fmt.Errorf("step %d (%s) failed: %w", i+1, step.Name, err)
		}
	}

	return nil
}

func (v *Validator) createResource(r *api.Resource, collectionURL string, payload map[string]interface{}) (map[string]interface{}, error) {
	// If UserSettableID is supported, generate one
	var urlToUse = collectionURL
	if r.Methods.Create != nil && r.Methods.Create.SupportsUserSettableCreate {
		// Wait, multiple calls need different IDs. Random is better.
		// Generator.go doesn't expose random helper directly but we can use simple random here.
		// Let's rely on server assigned if not strictly testing UserSettable.
		// But for Step 2 (duplicate check), we needed it.
		// Let's append ID param.

		// Actually, let's just let server assign if we don't strictly need to control it
		// UNLESS we need to test duplicate creation.
		// For now, appending ID query param if supported to be safe.
		// But I need to generate a unique ID.
		// Using "v.generateID()"
		urlToUse = fmt.Sprintf("%s?%s_id=%s", collectionURL, r.Singular, v.generateID())
	}

	resp, err := v.post(urlToUse, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var createdResource map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createdResource); err != nil {
		return nil, err
	}
	return createdResource, nil
}

func (v *Validator) generateID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("test-id-%d", r.Intn(100000))
}

func (v *Validator) post(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return v.client.Do(req)
}

func (v *Validator) patch(url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return v.client.Do(req)
}

func (v *Validator) getReq(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return v.client.Do(req)
}

func (v *Validator) get(url string) (map[string]interface{}, error) {
	resp, err := v.getReq(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	var resource map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (v *Validator) deleteReq(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return v.client.Do(req)
}

func (v *Validator) delete(url string) error {
	resp, err := v.deleteReq(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

type ListResponse struct {
	Resources     []map[string]interface{} `json:"-"`
	NextPageToken string                   `json:"next_page_token"`
}

func (v *Validator) list(url string) (*ListResponse, error) {
	resp, err := v.getReq(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	// Decode into map first to find the list field (which is plural name)
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	// Attempt to find the list field. It should be the only array field or match plural name?
	// We don't have convenient way to know the field name for list response here easily unless we pass it.
	// Convention: it matches resource plural name (camel case).

	// For now, looking for any array field.
	var resources []map[string]interface{}
	nextToken, _ := raw["next_page_token"].(string)

	for _, val := range raw {
		if list, ok := val.([]interface{}); ok {
			for _, item := range list {
				if r, ok := item.(map[string]interface{}); ok {
					resources = append(resources, r)
				}
			}
			break // Assuming only one list field
		}
	}

	return &ListResponse{
		Resources:     resources,
		NextPageToken: nextToken,
	}, nil
}

func getIDFromResourceName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
