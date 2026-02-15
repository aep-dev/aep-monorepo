package tests

import (
	"fmt"
)

var TestAEP132ListResourcesPageToken = Test{
	Name: "aep-132-list-resources-page-token",
	Run:  testListResourcesPageToken,
}

func testListResourcesPageToken(v ValidationActions, ctx *ValidationContext) error {
	if ctx.ListResponse1 == nil || ctx.ListResponse1.NextPageToken == "" {
		return fmt.Errorf("precondition failed: no page token available")
	}

	listURL2 := fmt.Sprintf("%s?page_size=1&page_token=%s", ctx.CollectionURL, ctx.ListResponse1.NextPageToken)
	listResp2, err := v.List(listURL2)
	if err != nil {
		return err
	}
	if len(listResp2.Resources) < 1 {
		return fmt.Errorf("expected at least 1 resource (the second one)")
	}
	fmt.Println("   Got second page.")
	return nil
}
