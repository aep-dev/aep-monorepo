package tests

import (
	"fmt"
	"io"
	"net/http"
)

var TestAEP134UpdateResource = Test{
	Name:     "aep-134-update-resource",
	Setup:    setupUpdateResource,
	Run:      testUpdateResource,
	Teardown: testDeleteResource,
}

func setupUpdateResource(v ValidationActions, ctx *ValidationContext) error {
	if len(ctx.Resources) == 0 {
		return testCreateResource(v, ctx)
	}
	return nil
}

func testUpdateResource(v ValidationActions, ctx *ValidationContext) error {
	updatePayload, err := GenerateCreatePayload(ctx.Resource)
	if err != nil {
		return fmt.Errorf("failed to generate update payload: %w", err)
	}

	rName, ok := ctx.Resources[0]["name"].(string)
	if !ok || rName == "" {
		rName, _ = ctx.Resources[0]["path"].(string)
	}
	rURL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, rName)

	respUpdate, err := v.Patch(rURL, updatePayload)
	if err != nil {
		return err
	}
	defer respUpdate.Body.Close()
	if respUpdate.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respUpdate.Body)
		return fmt.Errorf("update returned %d: %s", respUpdate.StatusCode, string(body))
	}
	fmt.Println("   Update successful.")
	return nil
}
