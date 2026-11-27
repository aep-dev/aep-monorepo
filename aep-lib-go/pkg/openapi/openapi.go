package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	OAS2             = "2.0"
	OAS3             = "3.0"
	APPLICATION_JSON = "application/json"
	JSON_MERGE_PATCH = "application/merge-patch+json"
)

type OpenAPI struct {
	// oas 2.0 has swagger in the root.
	Swagger    string                `json:"swagger,omitempty"`
	Info       Info                  `json:"info"`
	OpenAPI    string                `json:"openapi,omitempty"`
	Servers    []Server              `json:"servers,omitempty"`
	Paths      map[string]*PathItem  `json:"paths"`
	Components Components            `json:"components,omitempty"`
	Security   []SecurityRequirement `json:"security,omitempty"`
	// oas 2.0 has definitions in the root.
	Definitions map[string]Schema `json:"definitions,omitempty"`
}

func (o *OpenAPI) OASVersion() string {
	if o.Swagger == "2.0" {
		return OAS2
	} else if o.OpenAPI != "" {
		return OAS3
	}
	return ""
}

func (o *OpenAPI) DereferenceSchema(schema Schema) (*Schema, error) {
	if schema.Ref != "" {
		var refSchema Schema
		if strings.HasPrefix(schema.Ref, "https://") || strings.HasPrefix(schema.Ref, "http://") {
			body, err := readFileOrURL(schema.Ref)
			if err != nil {
				return nil, fmt.Errorf("error fetching external schema %q: %w", schema.Ref, err)
			}
			slog.Debug("Fetched schema body", "body", string(body))

			if err := json.Unmarshal(body, &refSchema); err != nil {
				return nil, fmt.Errorf("error unmarshalingexternal schema %q: %w", schema.Ref, err)
			}
		} else if strings.HasPrefix(schema.Ref, "#") {
			// Handle local schema references
			parts := strings.Split(schema.Ref, "/")
			key := parts[len(parts)-1]
			var ok bool
			switch o.OASVersion() {
			case OAS2:
				refSchema, ok = o.Definitions[key]
				slog.Debug("oasv2.0", "key", key)
				if !ok {
					return nil, fmt.Errorf("schema %q not found", schema.Ref)
				}
			default:
				refSchema, ok = o.Components.Schemas[key]
				if !ok {
					return nil, fmt.Errorf("schema %q not found", schema.Ref)
				}
			}
		} else {
			return nil, fmt.Errorf("unsupported schema reference %q", schema.Ref)
		}
		slog.Debug("ref schema", "schema", refSchema)
		return o.DereferenceSchema(refSchema)
	}
	return &schema, nil
}

func (o *OpenAPI) GetSchemaFromResponse(r Response, contentType string) *Schema {
	switch o.OASVersion() {
	case OAS2:
		return r.Schema
	default:
		ct := r.Content[contentType]
		return ct.Schema
	}
}

func (o *OpenAPI) GetSchemaFromRequestBody(r RequestBody, contentType string) *Schema {
	switch o.OASVersion() {
	case OAS2:
		return r.Schema
	default:
		ct := r.Content[contentType]
		return ct.Schema
	}
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

type Info struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	Contact     Contact `json:"contact,omitempty"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Summary                  string                    `json:"summary,omitempty"`
	Description              string                    `json:"description,omitempty"`
	OperationID              string                    `json:"operationId,omitempty"`
	Tags                     []string                  `json:"tags,omitempty"`
	Deprecated               bool                      `json:"deprecated,omitempty"`
	Parameters               []Parameter               `json:"parameters,omitempty"`
	Responses                map[string]Response       `json:"responses,omitempty"`
	RequestBody              *RequestBody              `json:"requestBody,omitempty"`
	Security                 []SecurityRequirement     `json:"security,omitempty"`
	XAEPLongRunningOperation *XAEPLongRunningOperation `json:"x-aep-long-running-operation,omitempty"`
}

type Parameter struct {
	Ref         string     `json:"$ref,omitempty"`
	Name        string     `json:"name,omitempty"`
	In          string     `json:"in,omitempty"`
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Deprecated  bool       `json:"deprecated,omitempty"`
	Schema      *Schema    `json:"schema,omitempty"`
	Type        string     `json:"type,omitempty"`
	XAEPField   *XAEPField `json:"x-aep-field,omitempty"`
}

type Response struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	// oas 2.0 has the schema in the response.
	Schema *Schema `json:"schema,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required"`
	// oas 2.0 has the schema in the request body.
	Schema *Schema `json:"schema,omitempty"`
}

type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

type Schema struct {
	Type         string        `json:"type,omitempty"`
	Format       string        `json:"format,omitempty"`
	Items        *Schema       `json:"items,omitempty"`
	Properties   Properties    `json:"properties,omitempty"`
	Ref          string        `json:"$ref,omitempty"`
	XAEPResource *XAEPResource `json:"x-aep-resource,omitempty"`
	XAEPField    *XAEPField    `json:"x-aep-field,omitempty"`
	/// Documents the name of the proto message to use for generation.
	/// If unset, proto generation will not create a proto message for this schema.
	XAEPProtoMessageName string   `json:"x-aep-proto-message-name,omitempty"`
	ReadOnly             bool     `json:"readOnly,omitempty"`
	WriteOnly            bool     `json:"writeOnly,omitempty"`
	Required             []string `json:"required,omitempty"`
	Description          string   `json:"description,omitempty"`
	// AdditionalProperties can be a bool and an object - this allows
	// handling for both.
	AdditionalProperties json.RawMessage `json:"additionalProperties,omitempty"`

	// Composition & polymorphism
	OneOf         []Schema       `json:"oneOf,omitempty"`
	AnyOf         []Schema       `json:"anyOf,omitempty"`
	AllOf         []Schema       `json:"allOf,omitempty"`
	Discriminator *Discriminator `json:"discriminator,omitempty"`

	// Enums, defaults, nullability, and lifecycle
	Enum       []any           `json:"enum,omitempty"`
	Default    json.RawMessage `json:"default,omitempty"`
	Nullable   bool            `json:"nullable,omitempty"`
	Deprecated bool            `json:"deprecated,omitempty"`

	// Validation constraints
	MinLength   *int     `json:"minLength,omitempty"`
	MaxLength   *int     `json:"maxLength,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
	Minimum     *float64 `json:"minimum,omitempty"`
	Maximum     *float64 `json:"maximum,omitempty"`
	MinItems    *int     `json:"minItems,omitempty"`
	MaxItems    *int     `json:"maxItems,omitempty"`
	UniqueItems bool     `json:"uniqueItems,omitempty"`
}

// UnmarshalJSON supports both OpenAPI 3.0 string types (with nullable: true)
// and OpenAPI 3.1 type arrays such as ["string", "null"].
func (s *Schema) UnmarshalJSON(data []byte) error {
	type schemaAlias Schema
	aux := &struct {
		RawType json.RawMessage `json:"type,omitempty"`
		*schemaAlias
	}{
		schemaAlias: (*schemaAlias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.RawType) > 0 && string(aux.RawType) != "null" {
		var single string
		if err := json.Unmarshal(aux.RawType, &single); err == nil {
			s.Type = single
		} else {
			var list []string
			if err := json.Unmarshal(aux.RawType, &list); err != nil {
				return fmt.Errorf("invalid schema type %s: %w", string(aux.RawType), err)
			}
			for _, t := range list {
				if t == "null" {
					s.Nullable = true
				} else if s.Type == "" {
					s.Type = t
				}
			}
		}
	}
	return nil
}

type Properties map[string]Schema

type Components struct {
	Schemas         map[string]Schema         `json:"schemas"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// SecurityRequirement maps a security scheme name to the list of required scopes.
type SecurityRequirement map[string][]string

type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type XAEPResource struct {
	Singular string   `json:"singular,omitempty"`
	Plural   string   `json:"plural,omitempty"`
	Patterns []string `json:"patterns,omitempty"`
	Parents  []string `json:"parents,omitempty"`
	Type     string   `json:"type,omitempty"`
}

type XAEPField struct {
	Behavior                   []string `json:"behavior,omitempty"`
	ResourceReference          []string `json:"resource_reference,omitempty"`
	ResourceReferenceChildType []string `json:"resource_reference_child_type,omitempty"`
	FieldNumber                int      `json:"field_number,omitempty"`
}

type XAEPLongRunningOperation struct {
	Response XAEPLongRunningOperationResponse  `json:"response,omitempty"`
	Metadata *XAEPLongRunningOperationMetadata `json:"metadata,omitempty"`
}

type XAEPLongRunningOperationResponse struct {
	Schema *Schema `json:"schema,omitempty"`
}

type XAEPLongRunningOperationMetadata struct {
	Schema *Schema `json:"schema,omitempty"`
}

func FetchOpenAPI(pathOrURL string) (*OpenAPI, error) {
	body, err := readFileOrURL(pathOrURL)
	if err != nil {
		return nil, fmt.Errorf("unable to read file or URL: %w", err)
	}

	var api OpenAPI
	if err := json.Unmarshal(body, &api); err != nil {
		return nil, err
	}

	return &api, nil
}

func readFileOrURL(pathOrURL string) ([]byte, error) {
	if isURL(pathOrURL) {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}

	return os.ReadFile(pathOrURL)
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}
