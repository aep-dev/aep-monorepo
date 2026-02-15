package tests

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
)

var TestAEP133DuplicateCreationCheck = Test{
	Name:     "aep-133-duplicate-creation-check",
	Setup:    setupDuplicateCreationCheck,
	Run:      testDuplicateCreationCheck,
	Teardown: testDeleteResource, // Clean up the one created in Setup
}

func setupDuplicateCreationCheck(v ValidationActions, ctx *ValidationContext) error {
	if len(ctx.Resources) == 0 {
		return testCreateResource(v, ctx)
	}
	return nil
}

func testDuplicateCreationCheck(v ValidationActions, ctx *ValidationContext) error {
	r := ctx.Resource
	if r.Methods.Create != nil && r.Methods.Create.SupportsUserSettableCreate {
		fmt.Println("   Attempting duplicate creation...")
		r1Name, ok := ctx.Resources[0]["name"].(string)
		if !ok || r1Name == "" {
			r1Name, _ = ctx.Resources[0]["path"].(string)
		}
		r1ID := getIDFromResourceName(r1Name)
		createPayload, _ := utils.GenerateCreatePayload(r)

		urlWithID := fmt.Sprintf("%s?id=%s", ctx.CollectionURL, r1ID)
		resp, err := v.Post(urlWithID, createPayload)
		if err != nil {
			return fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusConflict && resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("expected 409/400 for duplicate creation, got %d", resp.StatusCode)
		}
		fmt.Println("   Duplicate creation rejected as expected.")
	} else {
		fmt.Println("   Skipping duplicate check (client-assigned ID not supported or uncheckable).")
	}
	return nil
}

func getIDFromResourceName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
