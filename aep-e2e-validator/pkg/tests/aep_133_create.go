package tests

import (
	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
)

var TestAEP133Create = Test{
	Name:     "aep-133-create",
	Run:      testCreateResource,
	Teardown: testDeleteResource,
}

func testCreateResource(v ValidationActions, ctx *ValidationContext) error {
	resource, err := utils.CreateResource(v, ctx.Resource, ctx.CollectionURL)
	if err != nil {
		return err
	}

	// Store for cleanup/other tests
	ctx.Resources = append(ctx.Resources, resource) // append to slice
	return nil
}
