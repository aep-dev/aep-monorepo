package tests

import (
	"fmt"
)

var TestAEP133Create = Test{
	Name:     "aep-133-create",
	Run:      testCreateResource,
	Teardown: testDeleteResource,
}

func testCreateResource(v ValidationActions, ctx *ValidationContext) error {
	createPayload, err := GenerateCreatePayload(ctx.Resource)
	if err != nil {
		return fmt.Errorf("failed to generate create payload: %w", err)
	}

	resource, err := v.CreateResource(ctx.Resource, ctx.CollectionURL, createPayload)
	if err != nil {
		return err
	}

	// Store for cleanup/other tests
	ctx.Resources = append(ctx.Resources, resource) // append to slice
	rName, ok := resource["name"].(string)
	if !ok || rName == "" {
		rName, ok = resource["path"].(string)
	}

	if rName != "" {
		fmt.Printf("   Created %s\n", rName)
	} else {
		// Fallback logging if really needed, or just normal print
		fmt.Printf("   Created resource (name/path missing)\n")
	}
	return nil
}
