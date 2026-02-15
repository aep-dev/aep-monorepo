package tests

import (
	"fmt"

	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
)

var TestAEP132ListResourcesLimit1 = Test{
	Name:     "aep-132-list-resources-limit-1",
	Setup:    setupListResources,
	Run:      testListResourcesLimit1,
	Teardown: testDeleteResource,
}

func testListResourcesLimit1(v ValidationActions, ctx *ValidationContext) error {
	listResp, err := utils.FetchList(v, ctx.CollectionURL, "", 1)
	if err != nil {
		return err
	}
	if len(listResp.Resources) != 1 {
		return fmt.Errorf("expected 1 resource, got %d", len(listResp.Resources))
	}
	if listResp.NextPageToken == "" {
		return fmt.Errorf("expected next_page_token")
	}
	ctx.ListResponse1 = listResp
	fmt.Println("   Got 1 resource and next_page_token.")
	return nil
}
