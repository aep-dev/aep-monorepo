package utils

import (
	"fmt"

	"github.com/aep-dev/aep-lib-go/pkg/api"
)

type Creator interface {
	CreateResource(r *api.Resource, collectionURL string, payload map[string]interface{}) (map[string]interface{}, error)
}

func CreateResource(c Creator, r *api.Resource, collectionURL string) (map[string]interface{}, error) {
	createPayload, err := GenerateCreatePayload(r)
	if err != nil {
		return nil, fmt.Errorf("failed to generate create payload: %w", err)
	}

	resource, err := c.CreateResource(r, collectionURL, createPayload)
	if err != nil {
		return nil, err
	}

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
	return resource, nil
}
