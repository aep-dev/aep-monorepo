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

	"github.com/aep-dev/aep-e2e-validator/pkg/tests"
	"github.com/aep-dev/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-lib-go/pkg/openapi"
)

type Validator struct {
	configPath     string
	collection     string
	allCollections bool
	testNames      []string
	client         *http.Client
}

func NewValidator(configPath, collection string, allCollections bool, tests []string) *Validator {
	return &Validator{
		configPath:     configPath,
		collection:     collection,
		allCollections: allCollections,
		testNames:      tests,
		client:         &http.Client{},
	}
}

func (v *Validator) Run() int {
	doc, err := openapi.FetchOpenAPI(v.configPath)
	if err != nil {
		log.Printf("failed to fetch OpenAPI spec: %v", err)
		return ExitCodePreconditionFailed // Or some other code for setup failure
	}

	serverURL := ""
	if len(doc.Servers) > 0 {
		serverURL = doc.Servers[0].URL
	}

	aepAPI, err := api.GetAPI(doc, serverURL, "")
	if err != nil {
		log.Printf("failed to parse API: %v", err)
		return ExitCodePreconditionFailed
	}

	if err := api.AddImplicitFieldsAndValidate(aepAPI); err != nil {
		log.Printf("failed to validate API: %v", err)
		return ExitCodePreconditionFailed
	}

	if v.allCollections {
		for _, r := range aepAPI.Resources {
			// Only validate top-level resources
			if len(r.Parents) == 0 {
				if code := v.validateResource(r); code != ExitCodeSuccess {
					log.Printf("Validation failed for resource %s", r.Singular)
					return code
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
			log.Printf("collection %s not found in API", v.collection)
			return ExitCodePreconditionFailed
		}
		if code := v.validateResource(targetResource); code != ExitCodeSuccess {
			return code
		}
	}

	return ExitCodeSuccess
}

func (v *Validator) validateResource(r *api.Resource) int {
	fmt.Printf("Starting validation for resource: %s\n", r.Singular)
	ctx := &tests.ValidationContext{
		Resource:      r,
		CollectionURL: fmt.Sprintf("%s/%s", r.API.ServerURL, r.Plural),
		Resources:     make([]map[string]interface{}, 0),
	}

	availableTests := tests.NewTests()
	var testsToRun []tests.Test

	if len(v.testNames) == 0 {
		testsToRun = availableTests
	} else {
		testMap := make(map[string]tests.Test)
		for _, t := range availableTests {
			testMap[t.Name] = t
		}
		for _, name := range v.testNames {
			if t, ok := testMap[name]; ok {
				testsToRun = append(testsToRun, t)
			} else {
				log.Printf("Test %s not found", name)
				return ExitCodePreconditionFailed // Or warning?
			}
		}
	}

	for i, test := range testsToRun {
		fmt.Printf("%d. %s...\n", i+1, test.Name)

		if test.Precondition != nil {
			if err := test.Precondition(ctx); err != nil {
				fmt.Printf("   Precondition failed: %v\n", err)
				return ExitCodePreconditionFailed
			}
		}

		if test.Setup != nil {
			if err := test.Setup(v, ctx); err != nil {
				fmt.Printf("   Setup failed: %v\n", err)
				// Attempt teardown if setup fail
				if test.Teardown != nil {
					if err := test.Teardown(v, ctx); err != nil {
						fmt.Printf("   Teardown failed: %v\n", err)
						return ExitCodeTeardownFailed
					}
				}
				return ExitCodeSetupFailed
			}
		}

		if err := test.Run(v, ctx); err != nil {
			fmt.Printf("   Failed: %v\n", err)

			// Attempt teardown if test fail
			if test.Teardown != nil {
				if err := test.Teardown(v, ctx); err != nil {
					fmt.Printf("   Teardown failed: %v\n", err)
					return ExitCodeTeardownFailed
				}
			}

			return ExitCodeTestFailed
		}

		if test.Teardown != nil {
			if err := test.Teardown(v, ctx); err != nil {
				fmt.Printf("   Teardown failed: %v\n", err)
				return ExitCodeTeardownFailed
			}
		}
	}

	return ExitCodeSuccess
}

func (v *Validator) CreateResource(r *api.Resource, collectionURL string, payload map[string]interface{}) (map[string]interface{}, error) {
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
		urlToUse = fmt.Sprintf("%s?%s_id=%s", collectionURL, r.Singular, v.GenerateID())
	}

	resp, err := v.Post(urlToUse, payload)
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

func (v *Validator) GenerateID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("test-id-%d", r.Intn(100000))
}

func (v *Validator) Post(url string, body interface{}) (*http.Response, error) {
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

func (v *Validator) Patch(url string, body interface{}) (*http.Response, error) {
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

func (v *Validator) Get(url string) (map[string]interface{}, error) {
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

func (v *Validator) DeleteReq(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return v.client.Do(req)
}

func (v *Validator) Delete(url string) error {
	resp, err := v.DeleteReq(url)
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

func (v *Validator) List(url string) (*tests.ListResponse, error) {
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

	return &tests.ListResponse{
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
