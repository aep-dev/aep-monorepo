package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aep-dev/aep-lib-go/pkg/api"
	"github.com/aep-dev/aep-lib-go/pkg/openapi"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func GenerateCreatePayload(r *api.Resource) (map[string]interface{}, error) {
    payload := make(map[string]interface{})

    if r.Schema == nil {
        return nil, fmt.Errorf("resource schema is nil")
    }

    // Only populate required fields for now, or fields that seem important
    for propName, propSchema := range r.Schema.Properties {
        if isSystemField(propName) {
            continue
        }

        // TODO: check if required, or just generate everything that isn't readOnly
        if !propSchema.ReadOnly {
             payload[propName] = generateValue(propSchema)
        }
    }

    return payload, nil
}

func isSystemField(name string) bool {
    switch name {
    case "name", "createTime", "updateTime", "deleteTime", "uid", "etag":
        return true
    }
    return false
}

func generateValue(schema openapi.Schema) interface{} {
    switch schema.Type {
    case "string":
        return fmt.Sprintf("test-%s-%d", "string", rand.Intn(10000))
    case "integer":
        return rand.Intn(100)
    case "boolean":
        return true
    case "object":
         // recursive generation not implemented for simplicity in this pass
         return map[string]interface{}{}
    case "array":
        return []interface{}{}
    }
    return nil
}
