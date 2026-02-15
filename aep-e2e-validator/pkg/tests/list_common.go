package tests

import (
	"github.com/aep-dev/aep-e2e-validator/pkg/utils"
)

// setupListResources creates multiple resources to ensure pagination can be tested.
func setupListResources(v ValidationActions, ctx *ValidationContext) error {
	// Create 3 resources to ensuring we have enough for 2 pages of size 1 and a 3rd page or just ensuring we have > 1.
	for i := 0; i < 3; i++ {
		resource, err := utils.CreateResource(v, ctx.Resource, ctx.CollectionURL)
		if err != nil {
			return err
		}
		ctx.Resources = append(ctx.Resources, resource)
	}
	return nil
}
