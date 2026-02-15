package tests

import (
	"fmt"
	"net/http"
)

var TestAEP135DeleteNonExistentResource = Test{
	Name: "aep-135-delete-nonexistent-resource",
	Run:  testDeleteNonExistentResource,
}

func testDeleteNonExistentResource(v ValidationActions, ctx *ValidationContext) error {
	// Generate a random ID that likely does not exist
	randomID := v.GenerateID()
	rURL := fmt.Sprintf("%s/%s", ctx.CollectionURL, randomID)

	respDelete, err := v.DeleteReq(rURL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer respDelete.Body.Close()
	if respDelete.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404, got %d", respDelete.StatusCode)
	}
	fmt.Println("   Got 404 as expected.")
	return nil
}
