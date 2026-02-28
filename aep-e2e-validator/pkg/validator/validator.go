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
	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
	"github.com/aep-dev/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-lib-go/pkg/openapi"
)

type Validator struct {
	configPath     string
	collection     string
	allCollections bool
	testNames      []string
	client         *extendedClient
}

func NewValidator(configPath, collection string, allCollections bool, tests []string, headers []Header) *Validator {
	return &Validator{
		configPath:     configPath,
		collection:     collection,
		allCollections: allCollections,
		testNames:      tests,
		client:         &extendedClient{inner: &http.Client{}, headers: headers},
	}
}

func (v *Validator) Run() int {
	start := time.Now()

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

	var allResults []TestResult

	if v.allCollections {
		for _, r := range aepAPI.Resources {
			// Only validate top-level resources
			if len(r.Parents) == 0 {
				results := v.validateResource(r)
				allResults = append(allResults, results...)
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
		results := v.validateResource(targetResource)
		allResults = append(allResults, results...)
	}

	printSummary(allResults, time.Since(start))
	return worstExitCode(allResults)
}

func (v *Validator) validateResource(r *api.Resource) []TestResult {
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
			}
		}
	}

	// Global Setup: clean up collection
	fmt.Println("Running Global Setup...")
	if err := v.cleanupCollection(r); err != nil {
		fmt.Printf("   Global Setup failed: %v\n", err)
		results := make([]TestResult, len(testsToRun))
		for i, t := range testsToRun {
			results[i] = TestResult{Name: t.Name, Status: StatusError, Detail: fmt.Sprintf("global setup failed: %v", err)}
		}
		return results
	}

	var results []TestResult
	for i, test := range testsToRun {
		fmt.Printf("%d. %s...\n", i+1, test.Name)
		testStart := time.Now()

		v.client.clearLogs()

		if test.Precondition != nil {
			if err := test.Precondition(ctx); err != nil {
				fmt.Printf("   Skipped: %v\n", err)
				results = append(results, TestResult{Name: test.Name, Status: StatusSkip, Detail: err.Error(), Duration: time.Since(testStart)})
				continue
			}
		}

		if test.Setup != nil {
			if err := test.Setup(v, ctx); err != nil {
				fmt.Printf("   Setup failed: %v\n", err)
				v.client.printLogs()
				if test.Teardown != nil {
					_ = test.Teardown(v, ctx)
				}
				results = append(results, TestResult{Name: test.Name, Status: StatusError, Detail: fmt.Sprintf("setup: %v", err), Duration: time.Since(testStart)})
				continue
			}
		}

		if err := test.Run(v, ctx); err != nil {
			fmt.Printf("   Failed: %v\n", err)
			v.client.printLogs()
			if test.Teardown != nil {
				_ = test.Teardown(v, ctx)
			}
			results = append(results, TestResult{Name: test.Name, Status: StatusFail, Detail: err.Error(), Duration: time.Since(testStart)})
			continue
		}

		if test.Teardown != nil {
			if err := test.Teardown(v, ctx); err != nil {
				fmt.Printf("   Teardown failed: %v\n", err)
				v.client.printLogs()
				results = append(results, TestResult{Name: test.Name, Status: StatusError, Detail: fmt.Sprintf("teardown: %v", err), Duration: time.Since(testStart)})
				continue
			}
		}

		results = append(results, TestResult{Name: test.Name, Status: StatusPass, Duration: time.Since(testStart)})
	}

	// Global Teardown: clean up collection
	fmt.Println("Running Global Teardown...")
	if err := v.cleanupCollection(r); err != nil {
		fmt.Printf("   Global Teardown failed: %v\n", err)
	}

	return results
}

func (v *Validator) cleanupCollection(r *api.Resource) error {
	collectionURL := fmt.Sprintf("%s/%s", r.API.ServerURL, r.Plural)
	pageToken := ""

	for {
		listResp, err := utils.FetchList(v, collectionURL, pageToken, 0)
		if err != nil {
			// If 404, the collection likely doesn't exist, which is fine (empty).
			// But List should return 200 with 0 items usually.
			// Let's Log strict error for now.
			return fmt.Errorf("failed to list resources: %w", err)
		}

		for _, resource := range listResp.Resources {
			rName, ok := resource["name"].(string)
			if !ok || rName == "" {
				rName, _ = resource["path"].(string)
			}
			if rName == "" {
				// Skipping resource without name/path
				continue
			}

			delURL := fmt.Sprintf("%s/%s", r.API.ServerURL, rName)
			// We try to delete. If it fails, we log but maybe don't abort entire cleanup?
			// Ideally we want to clean everything.
			if err := v.Delete(delURL); err != nil {
				if strings.Contains(err.Error(), "status 404") {
					continue
				}
				fmt.Printf("   Warning: failed to delete during cleanup %s: %v\n", rName, err)
			}
		}

		pageToken = listResp.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return nil
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

func (v *Validator) List(url string) (*utils.ListResponse, error) {
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

	var resources []map[string]interface{}
	nextToken, _ := raw["next_page_token"].(string)

	if list, ok := raw["results"].([]interface{}); ok {
		for _, item := range list {
			if r, ok := item.(map[string]interface{}); ok {
				resources = append(resources, r)
			}
		}
	}

	return &utils.ListResponse{
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
