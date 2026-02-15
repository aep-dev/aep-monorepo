package tests

import (
	"fmt"

	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
)

var TestAEP132ListResourcesPageToken = Test{
	Name:     "aep-132-list-resources-page-token",
	Setup:    setupListResources,
	Run:      testListResourcesPageToken,
	Teardown: testDeleteResource,
}

func testListResourcesPageToken(v ValidationActions, ctx *ValidationContext) error {
	// Step 1: List with limit 1 to get a page token
	listResp1, err := utils.FetchList(v, ctx.CollectionURL, "", 1)
	if err != nil {
		return err
	}
	if listResp1.NextPageToken == "" {
		return fmt.Errorf("precondition failed: no page token returned with max_page_size=1 (ensure setup created > 1 resource)")
	}

	// Step 2: List using the page token
	listResp2, err := utils.FetchList(v, ctx.CollectionURL, listResp1.NextPageToken, 1)
	if err != nil {
		return err
	}

	// Step 3: Verify we got resources
	if len(listResp2.Resources) < 1 {
		return fmt.Errorf("expected at least 1 resource on second page")
	}

	// Optional: Verify resources in page 2 are different or next in sequence?
	// For now, just verifying we got a response is enough for the basic requirement.

	fmt.Println("   Got second page successfully.")
	return nil
}
