package tests

import (
	"fmt"
)

var TestAEP135DeleteResource = Test{
	Name:  "aep-135-delete-resource",
	Setup: setupDeleteResource,
	Run:   testDeleteResource,
}

func setupDeleteResource(v ValidationActions, ctx *ValidationContext) error {
	if len(ctx.Resources) == 0 {
		return testCreateResource(v, ctx)
	}
	return nil
}

func testDeleteResource(v ValidationActions, ctx *ValidationContext) error {
	if len(ctx.Resources) == 0 {
		return nil // Nothing to delete
	}

	rName, ok := ctx.Resources[0]["name"].(string)
	if !ok || rName == "" {
		rName, _ = ctx.Resources[0]["path"].(string)
	}
	rURL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, rName)
	if err := v.Delete(rURL); err != nil {
		return err
	}

	// Remove from context
	ctx.Resources = ctx.Resources[1:]

	fmt.Println("   Delete successful.")
	return nil
}
