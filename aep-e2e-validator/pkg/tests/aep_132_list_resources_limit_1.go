package tests

import (
	"fmt"
)

var TestAEP132ListResourcesLimit1 = Test{
	Name: "aep-132-list-resources-limit-1",
	Run:  testListResourcesLimit1,
}

func testListResourcesLimit1(v ValidationActions, ctx *ValidationContext) error {
	listURL := fmt.Sprintf("%s?page_size=1", ctx.CollectionURL)
	listResp, err := v.List(listURL)
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
