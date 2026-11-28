package api

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/cases"
	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/constants"
	"github.com/aep-dev/aep-monorepo/aep-lib-go/pkg/openapi"
)

type API struct {
	ServerURL string `json:"server_url"`
	Name      string
	Contact   *Contact
	Schemas   map[string]*openapi.Schema
	// A list of the resources that are exposed by the API.
	//
	// The key "operation" carries a special meaning, and must
	// map to an aep.dev/151 Operation resource.
	Resources map[string]*Resource
}

type Contact struct {
	Name  string
	Email string
	URL   string
}

func GetAPI(api *openapi.OpenAPI, serverURL, pathPrefix string) (*API, error) {
	if api.OASVersion() == "" {
		return nil, fmt.Errorf("unable to detect OAS version. Please add a openapi field or a swagger field")
	}
	slog.Debug("parsing openapi", "pathPrefix", pathPrefix)
	resourceBySingular := make(map[string]*Resource)
	customMethodsByPattern := make(map[string][]*CustomMethod)
	collectionCustomMethodsByPattern := make(map[string][]*CustomMethod)
	// we try to parse the paths to find possible resources, since
	// they may not always be annotated as such.
	for path, pathItem := range api.Paths {
		path = path[len(pathPrefix):]
		slog.Debug("path", "path", path)
		var r Resource
		var sRef *openapi.Schema
		p := getPatternInfo(path)
		if p == nil { // not a resource pattern
			slog.Debug("path is not a resource", "path", path)
			continue
		}
		slog.Debug("parsing path for resource", "path", path)
		if p.CustomMethodName != "" {
			// strip the leading slash and the custom method suffix
			pattern := strings.TrimPrefix(strings.Split(path, ":")[0], "/")
			isCollection := !p.IsResourcePattern
			parsed, err := parseCustomMethods(api, pathItem, p.CustomMethodName, isCollection)
			if err != nil {
				return nil, err
			}
			if isCollection {
				collectionCustomMethodsByPattern[pattern] = append(collectionCustomMethodsByPattern[pattern], parsed...)
			} else {
				customMethodsByPattern[pattern] = append(customMethodsByPattern[pattern], parsed...)
			}
		} else if p.IsResourcePattern {
			// individual resource pattern (get, update, apply, delete)
			if pathItem.Delete != nil {
				opLRO := pathItem.Delete.XAEPLongRunningOperation
				r.Methods.Delete = &DeleteMethod{
					IsLongRunning:   opLRO != nil,
					SupportsForce:   hasParam(pathItem.Delete.Parameters, constants.FIELD_FORCE_NAME, ""),
					SupportsIfMatch: hasHeaderParam(pathItem.Delete.Parameters, constants.HEADER_IF_MATCH_NAME),
				}
				if opLRO != nil && opLRO.Response.Schema != nil {
					sRef = opLRO.Response.Schema
				}
			}
			if pathItem.Get != nil {
				if resp, ok := pathItem.Get.Responses["200"]; ok {
					sRef = getJSONSchemaFromResponse(api, resp)
					r.Methods.Get = &GetMethod{}
				}
			}
			if pathItem.Patch != nil {
				opLRO := pathItem.Patch.XAEPLongRunningOperation
				if resp, ok := pathItem.Patch.Responses["200"]; ok {
					if opLRO != nil && opLRO.Response.Schema != nil {
						sRef = opLRO.Response.Schema
					} else {
						sRef = getJSONSchemaFromResponse(api, resp)
					}
					r.Methods.Update = &UpdateMethod{
						IsLongRunning:   opLRO != nil,
						SupportsIfMatch: hasHeaderParam(pathItem.Patch.Parameters, constants.HEADER_IF_MATCH_NAME),
					}
				}
			}
			if pathItem.Put != nil {
				opLRO := pathItem.Put.XAEPLongRunningOperation
				for _, statusCode := range []string{"200", "201"} {
					if resp, ok := pathItem.Put.Responses[statusCode]; ok {
						if opLRO != nil && opLRO.Response.Schema != nil {
							sRef = opLRO.Response.Schema
						} else {
							sRef = getJSONSchemaFromResponse(api, resp)
						}
						r.Methods.Apply = &ApplyMethod{
							IsLongRunning:   opLRO != nil,
							SupportsIfMatch: hasHeaderParam(pathItem.Put.Parameters, constants.HEADER_IF_MATCH_NAME),
						}
						break
					}
				}
			}
		} else {
			// create method
			if pathItem.Post != nil {
				for _, statusCode := range []string{"200", "201"} {
					// check if there is a query parameter "id"
					if resp, ok := pathItem.Post.Responses[statusCode]; ok {
						opLRO := pathItem.Post.XAEPLongRunningOperation
						if opLRO != nil && opLRO.Response.Schema != nil {
							sRef = opLRO.Response.Schema
						} else {
							sRef = getJSONSchemaFromResponse(api, resp)
						}
						r.Methods.Create = &CreateMethod{
							SupportsUserSettableCreate: hasParam(pathItem.Post.Parameters, constants.FIELD_ID_NAME, ""),
							SupportsRequestID:          hasParam(pathItem.Post.Parameters, constants.FIELD_REQUEST_ID_NAME, ""),
							IsLongRunning:              opLRO != nil,
						}
						break
					}
				}
			}
			// list method
			if pathItem.Get != nil {
				if resp, ok := pathItem.Get.Responses["200"]; ok {
					respSchema := getJSONSchemaFromResponse(api, resp)
					if respSchema == nil {
						slog.Warn(fmt.Sprintf("resource %q has a LIST method with a response schema, but the response schema is nil.", path))
					} else {
						resolvedSchema, err := api.DereferenceSchema(*respSchema)
						if err != nil {
							return nil, fmt.Errorf("error dereferencing schema %q: %v", respSchema.Ref, err)
						}
						itemsField, itemSchema := findListItemsProperty(resolvedSchema)
						if itemSchema != nil {
							sRef = itemSchema
							listMethod := &ListMethod{
								ItemsField: itemsField,
							}
							if _, ok := resolvedSchema.Properties[constants.FIELD_NEXT_PAGE_TOKEN_NAME]; ok {
								listMethod.NextPageTokenField = constants.FIELD_NEXT_PAGE_TOKEN_NAME
							} else if _, ok := resolvedSchema.Properties["nextPageToken"]; ok {
								listMethod.NextPageTokenField = "nextPageToken"
							}
							if _, ok := resolvedSchema.Properties[constants.FIELD_UNREACHABLE_NAME]; ok {
								listMethod.HasUnreachableResources = true
							}
							for _, param := range pathItem.Get.Parameters {
								switch param.Name {
								case constants.FIELD_SKIP_NAME:
									listMethod.SupportsSkip = true
								case constants.FIELD_UNREACHABLE_NAME:
									listMethod.HasUnreachableResources = true
								case constants.FIELD_FILTER_NAME:
									listMethod.SupportsFilter = true
								case constants.FIELD_ORDER_BY_NAME:
									listMethod.SupportsOrderBy = true
								}
							}
							r.Methods.List = listMethod
						} else {
							slog.Warn(fmt.Sprintf("resource %q has a LIST method with a response schema, but the items field is not present or is not an array.", path))
						}
					}
				}
			}
		}
		if sRef != nil && p.CustomMethodName == "" {
			// s should always be a reference to a schema in the components section.
			parts := strings.Split(sRef.Ref, "/")
			key := parts[len(parts)-1]
			dereferencedSchema, err := api.DereferenceSchema(*sRef)
			if err != nil {
				return nil, fmt.Errorf("error dereferencing schema %q: %v", sRef.Ref, err)
			}
			singular := cases.PascalToSnakeCase(key)
			pattern := strings.Split(path, "/")[1:]
			if !p.IsResourcePattern {
				// deduplicate the singular, if applicable
				finalSingular := singular
				parent := ""
				if len(pattern) >= 3 {
					parent = pattern[len(pattern)-3]
					parent = parent[0 : len(parent)-1] // strip curly surrounding
					if strings.HasPrefix(singular, parent) {
						finalSingular = strings.TrimPrefix(singular, parent+"_")
					}
				}
				pattern = append(pattern, fmt.Sprintf("{%s_id}", finalSingular))
			}
			r2, err := getOrPopulateResource(singular, pattern, dereferencedSchema, resourceBySingular, api)
			if err != nil {
				return nil, fmt.Errorf("error populating resource %q: %v", r.Singular, err)
			}
			foldResourceMethods(&r, r2)
		}
	}
	// the custom methods are trickier - because they may not respond with the schema of the resource
	// (which would allow us to map the resource via looking at its reference), we instead map them
	// by pattern across all patterns of each resource.
	for pattern, customMethods := range customMethodsByPattern {
		found := false
		for _, r := range resourceBySingular {
			if resourceHasPattern(r, pattern) {
				r.CustomMethods = append(r.CustomMethods, customMethods...)
				found = true
				break
			}
		}
		if !found {
			slog.Debug(fmt.Sprintf("custom methods with pattern %q have no resource associated with it", pattern))
		}
	}
	for pattern, customMethods := range collectionCustomMethodsByPattern {
		found := false
		for _, r := range resourceBySingular {
			if resourceHasCollectionPattern(r, pattern) {
				r.CustomMethods = append(r.CustomMethods, customMethods...)
				found = true
				break
			}
		}
		if !found {
			slog.Debug(fmt.Sprintf("collection custom methods with pattern %q have no resource associated with it", pattern))
		}
	}
	if serverURL == "" {
		for _, s := range api.Servers {
			serverURL = s.URL + pathPrefix
		}
	}

	if serverURL == "" {
		return nil, fmt.Errorf("no server URL found in openapi, and none was provided")
	}

	// any schemas that are not a resource are added to the API's schemas
	schemas := make(map[string]*openapi.Schema)
	for k, v := range api.Components.Schemas {
		if _, ok := resourceBySingular[k]; ok {
			continue
		}
		schemas[k] = &v
	}

	return &API{
		ServerURL: serverURL,
		Name:      api.Info.Title,
		Contact:   getContact(api.Info.Contact),
		Resources: resourceBySingular,
		Schemas:   schemas,
	}, nil
}

func parseCustomMethods(api *openapi.OpenAPI, pathItem *openapi.PathItem, customMethodName string, isCollection bool) ([]*CustomMethod, error) {
	var result []*CustomMethod
	if pathItem.Post != nil {
		if resp, ok := pathItem.Post.Responses["200"]; ok {
			lroDetails := pathItem.Post.XAEPLongRunningOperation
			schema := getJSONSchemaFromResponse(api, resp)
			responseSchema := &openapi.Schema{}
			if lroDetails != nil {
				schema = lroDetails.Response.Schema
			}
			if schema != nil {
				var err error
				responseSchema, err = api.DereferenceSchema(*schema)
				if err != nil {
					return nil, fmt.Errorf("error dereferencing schema %v: %v", schema, err)
				}
			}
			if pathItem.Post.RequestBody == nil {
				return nil, fmt.Errorf("custom method %q has a POST response, but no request body", customMethodName)
			}
			schema = api.GetSchemaFromRequestBody(*pathItem.Post.RequestBody, openapi.APPLICATION_JSON)
			requestSchema := &openapi.Schema{}
			if schema != nil {
				var err error
				requestSchema, err = api.DereferenceSchema(*schema)
				if err != nil {
					return nil, fmt.Errorf("error dereferencing schema %q: %v", schema.Ref, err)
				}
			}
			result = append(result, &CustomMethod{
				Name:               customMethodName,
				Method:             "POST",
				Request:            requestSchema,
				Response:           responseSchema,
				IsLongRunning:      lroDetails != nil,
				IsCollectionMethod: isCollection,
			})
		}
	}
	if pathItem.Get != nil {
		if resp, ok := pathItem.Get.Responses["200"]; ok {
			lroDetails := pathItem.Get.XAEPLongRunningOperation
			schema := getJSONSchemaFromResponse(api, resp)
			responseSchema := &openapi.Schema{}
			if lroDetails != nil {
				schema = lroDetails.Response.Schema
			}
			if schema != nil {
				var err error
				responseSchema, err = api.DereferenceSchema(*schema)
				if err != nil {
					return nil, fmt.Errorf("error dereferencing schema %v: %v", schema.Ref, err)
				}
			}
			result = append(result, &CustomMethod{
				Name:               customMethodName,
				Method:             "GET",
				Response:           responseSchema,
				IsLongRunning:      lroDetails != nil,
				IsCollectionMethod: isCollection,
			})
		}
	}
	return result, nil
}

func getJSONSchemaFromResponse(api *openapi.OpenAPI, resp openapi.Response) *openapi.Schema {
	if s := api.GetSchemaFromResponse(resp, openapi.APPLICATION_JSON); s != nil {
		return s
	}
	return api.GetSchemaFromResponse(resp, openapi.JSON_MERGE_PATCH)
}

func hasParam(params []openapi.Parameter, name, in string) bool {
	for _, param := range params {
		if param.Name != name {
			continue
		}
		if in != "" && param.In != in {
			continue
		}
		return true
	}
	return false
}

func hasHeaderParam(params []openapi.Parameter, name string) bool {
	for _, param := range params {
		if strings.EqualFold(param.Name, name) && (param.In == "" || strings.EqualFold(param.In, "header")) {
			return true
		}
	}
	return false
}

// findListItemsProperty deterministically selects the resource array property
// from a List response schema, preferring "results" and ignoring "unreachable".
func findListItemsProperty(resolvedSchema *openapi.Schema) (string, *openapi.Schema) {
	if resolvedSchema == nil || len(resolvedSchema.Properties) == 0 {
		return "", nil
	}
	if prop, ok := resolvedSchema.Properties[constants.FIELD_RESULTS_NAME]; ok && prop.Type == "array" && prop.Items != nil {
		return constants.FIELD_RESULTS_NAME, prop.Items
	}
	keys := make([]string, 0, len(resolvedSchema.Properties))
	for k := range resolvedSchema.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k == constants.FIELD_UNREACHABLE_NAME {
			continue
		}
		prop := resolvedSchema.Properties[k]
		if prop.Type == "array" && prop.Items != nil {
			return k, prop.Items
		}
	}
	return "", nil
}

func resourceHasPattern(r *Resource, pattern string) bool {
	for _, p := range r.GetPatterns() {
		if p == pattern {
			return true
		}
	}
	return false
}

func resourceHasCollectionPattern(r *Resource, collectionPattern string) bool {
	for _, elems := range r.AllPatternElems() {
		if len(elems) >= 2 && strings.Join(elems[:len(elems)-1], "/") == collectionPattern {
			return true
		}
	}
	return false
}

func (s *API) GetResource(resource string) (*Resource, error) {
	r, ok := (*s).Resources[resource]
	if !ok {
		return nil, fmt.Errorf("Resource %q not found", resource)
	}
	return r, nil
}

type PatternInfo struct {
	// if true, the pattern represents an individual resource,
	// otherwise it represents a path to a collection of resources
	IsResourcePattern bool
	CustomMethodName  string
}

// getPatternInfo returns true if the path is an alternating pairing of collection and id,
// and returns the collection names if so.
func getPatternInfo(path string) *PatternInfo {
	customMethodName := ""
	if strings.Contains(path, ":") {
		parts := strings.Split(path, ":")
		path = parts[0]
		customMethodName = parts[1]
	}
	// we ignore the first segment, which is empty.
	pattern := strings.Split(path, "/")[1:]
	for i, segment := range pattern {
		// check if segment is wrapped in curly brackets
		wrapped := strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
		wantWrapped := i%2 == 1
		if wrapped != wantWrapped {
			return nil
		}
	}
	return &PatternInfo{
		IsResourcePattern: len(pattern)%2 == 0,
		CustomMethodName:  customMethodName,
	}
}

// getOrPopulateResource populates the resource via a variety of means:
// - if the resource already exists in the map, it returns it
// - if the schema has the x-aep-resource annotation, it parses the resource
// - otherwise, it attempts to infer the resource from the schema and name.
func getOrPopulateResource(singular string, pattern []string, s *openapi.Schema, resourceBySingular map[string]*Resource, api *openapi.OpenAPI) (*Resource, error) {
	if r, ok := resourceBySingular[singular]; ok {
		r.addPattern(pattern)
		return r, nil
	}
	var r *Resource
	hasETag := false
	if s != nil && s.Properties != nil {
		_, hasETag = s.Properties[constants.FIELD_ETAG_NAME]
	}
	// use the X-AEP-Resource annotation to populate the resource,
	// if it exists.
	if s.XAEPResource != nil {
		parents := []*Resource{}
		for _, parentSingular := range s.XAEPResource.Parents {
			parentSchema, ok := api.Components.Schemas[parentSingular]
			if !ok {
				return nil, fmt.Errorf("resource %q parent %q not found", singular, parentSingular)
			}
			parentResource, err := getOrPopulateResource(parentSingular, []string{}, &parentSchema, resourceBySingular, api)
			if err != nil {
				return nil, fmt.Errorf("error parsing resource %q parent %q: %v", singular, parentSingular, err)
			}
			parents = append(parents, parentResource)
		}
		r = &Resource{
			Singular:        s.XAEPResource.Singular,
			Plural:          s.XAEPResource.Plural,
			Parents:         s.XAEPResource.Parents,
			parentResources: parents,
			Children:        []*Resource{},
			HasETag:         hasETag,
			Schema:          s,
		}
		for _, rawPattern := range s.XAEPResource.Patterns {
			trimmed := strings.TrimPrefix(rawPattern, "/")
			if trimmed != "" {
				r.addPattern(strings.Split(trimmed, "/"))
			}
		}
		if len(r.patterns) == 0 {
			r.addPattern(pattern)
		}
		for _, parentResource := range parents {
			parentResource.Children = append(parentResource.Children, r)
		}
	} else {
		// best effort otherwise
		r = &Resource{
			Schema:          s,
			Singular:        singular,
			Parents:         []string{},
			parentResources: []*Resource{},
			Children:        []*Resource{},
			Plural:          plural(singular),
			HasETag:         hasETag,
		}
		r.addPattern(pattern)
	}
	// update the resource map
	resourceBySingular[singular] = r
	return r, nil
}

func foldResourceMethods(from, into *Resource) {
	if from.Methods.Get != nil {
		into.Methods.Get = from.Methods.Get
	}
	if from.Methods.List != nil {
		into.Methods.List = from.Methods.List
	}
	if from.Methods.Apply != nil {
		into.Methods.Apply = from.Methods.Apply
	}
	if from.Methods.Create != nil {
		into.Methods.Create = from.Methods.Create
	}
	if from.Methods.Update != nil {
		into.Methods.Update = from.Methods.Update
	}
	if from.Methods.Delete != nil {
		into.Methods.Delete = from.Methods.Delete
	}
}

func getContact(contact openapi.Contact) *Contact {
	if contact.Name != "" || contact.Email != "" || contact.URL != "" {
		return &Contact{
			Name:  contact.Name,
			Email: contact.Email,
			URL:   contact.URL,
		}
	}
	return nil
}

// plural returns the plural form of a singular resource name
// This is a simple implementation that just adds 's' to the end
// of the singular form, which works for most cases.
func plural(singular string) string {
	return singular + "s"
}
