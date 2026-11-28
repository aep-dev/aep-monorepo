package openapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaUnmarshalTypeAndNullability(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantType     string
		wantNullable bool
		wantErr      bool
	}{
		{
			name:         "plain string type",
			input:        `{"type": "string"}`,
			wantType:     "string",
			wantNullable: false,
		},
		{
			name:         "OAS 3.0 nullable true",
			input:        `{"type": "string", "nullable": true}`,
			wantType:     "string",
			wantNullable: true,
		},
		{
			name:         "OAS 3.1 type array string and null",
			input:        `{"type": ["string", "null"]}`,
			wantType:     "string",
			wantNullable: true,
		},
		{
			name:         "OAS 3.1 type array null first",
			input:        `{"type": ["null", "integer"]}`,
			wantType:     "integer",
			wantNullable: true,
		},
		{
			name:         "OAS 3.1 type array without null",
			input:        `{"type": ["string", "integer"]}`,
			wantType:     "string",
			wantNullable: false,
		},
		{
			name:         "omitted type",
			input:        `{"description": "any value"}`,
			wantType:     "",
			wantNullable: false,
		},
		{
			name:         "null type value",
			input:        `{"type": null}`,
			wantType:     "",
			wantNullable: false,
		},
		{
			name:    "invalid numeric type",
			input:   `{"type": 123}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Schema
			err := json.Unmarshal([]byte(tt.input), &s)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantType, s.Type)
			assert.Equal(t, tt.wantNullable, s.Nullable)
		})
	}
}

func TestSchemaExpandedFields(t *testing.T) {
	raw := `{
		"type": "object",
		"deprecated": true,
		"properties": {
			"status": {
				"type": ["string", "null"],
				"enum": ["ACTIVE", "ARCHIVED"],
				"default": "ACTIVE",
				"minLength": 1,
				"maxLength": 32,
				"pattern": "^[A-Z]+$"
			},
			"score": {
				"type": "number",
				"minimum": 0,
				"maximum": 100.5,
				"writeOnly": true
			},
			"tags": {
				"type": "array",
				"items": {"type": "string"},
				"minItems": 0,
				"maxItems": 10,
				"uniqueItems": true
			},
			"payload": {
				"oneOf": [
					{"$ref": "#/components/schemas/TextPayload"},
					{"$ref": "#/components/schemas/BinaryPayload"}
				],
				"discriminator": {
					"propertyName": "kind",
					"mapping": {
						"text": "#/components/schemas/TextPayload",
						"binary": "#/components/schemas/BinaryPayload"
					}
				}
			},
			"combined": {
				"allOf": [{"$ref": "#/components/schemas/Base"}],
				"anyOf": [{"type": "string"}, {"type": "integer"}]
			}
		}
	}`

	var s Schema
	require.NoError(t, json.Unmarshal([]byte(raw), &s))

	assert.Equal(t, "object", s.Type)
	assert.True(t, s.Deprecated)

	status := s.Properties["status"]
	assert.Equal(t, "string", status.Type)
	assert.True(t, status.Nullable)
	assert.Equal(t, []any{"ACTIVE", "ARCHIVED"}, status.Enum)
	assert.Equal(t, `"ACTIVE"`, string(status.Default))
	require.NotNil(t, status.MinLength)
	assert.Equal(t, 1, *status.MinLength)
	require.NotNil(t, status.MaxLength)
	assert.Equal(t, 32, *status.MaxLength)
	assert.Equal(t, "^[A-Z]+$", status.Pattern)

	score := s.Properties["score"]
	require.NotNil(t, score.Minimum)
	assert.Equal(t, 0.0, *score.Minimum)
	require.NotNil(t, score.Maximum)
	assert.Equal(t, 100.5, *score.Maximum)
	assert.True(t, score.WriteOnly)

	tags := s.Properties["tags"]
	require.NotNil(t, tags.MinItems)
	assert.Equal(t, 0, *tags.MinItems)
	require.NotNil(t, tags.MaxItems)
	assert.Equal(t, 10, *tags.MaxItems)
	assert.True(t, tags.UniqueItems)

	payload := s.Properties["payload"]
	require.Len(t, payload.OneOf, 2)
	assert.Equal(t, "#/components/schemas/TextPayload", payload.OneOf[0].Ref)
	require.NotNil(t, payload.Discriminator)
	assert.Equal(t, "kind", payload.Discriminator.PropertyName)
	assert.Equal(t, "#/components/schemas/BinaryPayload", payload.Discriminator.Mapping["binary"])

	combined := s.Properties["combined"]
	require.Len(t, combined.AllOf, 1)
	require.Len(t, combined.AnyOf, 2)
}

func TestOpenAPISecurityAndLROMetadata(t *testing.T) {
	raw := `{
		"openapi": "3.1.0",
		"info": {"title": "Bookstore", "version": "1.0.0"},
		"security": [{"bearerAuth": []}],
		"components": {
			"schemas": {
				"Book": {"type": "object"},
				"ArchiveMetadata": {
					"type": "object",
					"properties": {"progress": {"type": "integer"}}
				}
			},
			"securitySchemes": {
				"bearerAuth": {
					"type": "http",
					"scheme": "bearer",
					"bearerFormat": "JWT"
				},
				"apiKey": {
					"type": "apiKey",
					"name": "X-API-Key",
					"in": "header"
				},
				"oauth": {
					"type": "oauth2",
					"flows": {
						"clientCredentials": {
							"tokenUrl": "https://auth.example.com/token",
							"scopes": {"books:write": "Write books"}
						}
					}
				},
				"oidc": {
					"type": "openIdConnect",
					"openIdConnectUrl": "https://auth.example.com/.well-known/openid-configuration"
				}
			}
		},
		"paths": {
			"/books/{book_id}:archive": {
				"post": {
					"operationId": "ArchiveBook",
					"tags": ["books"],
					"deprecated": true,
					"security": [{"oauth": ["books:write"]}],
					"parameters": [
						{"$ref": "#/components/parameters/BookId", "deprecated": true}
					],
					"x-aep-long-running-operation": {
						"response": {
							"schema": {"$ref": "#/components/schemas/Book"}
						},
						"metadata": {
							"schema": {"$ref": "#/components/schemas/ArchiveMetadata"}
						}
					},
					"responses": {
						"200": {"description": "OK"}
					}
				}
			}
		}
	}`

	var doc OpenAPI
	require.NoError(t, json.Unmarshal([]byte(raw), &doc))

	require.Len(t, doc.Security, 1)
	assert.Contains(t, doc.Security[0], "bearerAuth")

	require.Len(t, doc.Components.SecuritySchemes, 4)
	assert.Equal(t, "http", doc.Components.SecuritySchemes["bearerAuth"].Type)
	assert.Equal(t, "bearer", doc.Components.SecuritySchemes["bearerAuth"].Scheme)
	assert.Equal(t, "JWT", doc.Components.SecuritySchemes["bearerAuth"].BearerFormat)
	assert.Equal(t, "X-API-Key", doc.Components.SecuritySchemes["apiKey"].Name)
	assert.Equal(t, "header", doc.Components.SecuritySchemes["apiKey"].In)
	require.NotNil(t, doc.Components.SecuritySchemes["oauth"].Flows)
	require.NotNil(t, doc.Components.SecuritySchemes["oauth"].Flows.ClientCredentials)
	assert.Equal(t, "https://auth.example.com/token", doc.Components.SecuritySchemes["oauth"].Flows.ClientCredentials.TokenURL)
	assert.Equal(t, "https://auth.example.com/.well-known/openid-configuration", doc.Components.SecuritySchemes["oidc"].OpenIDConnectURL)

	op := doc.Paths["/books/{book_id}:archive"].Post
	require.NotNil(t, op)
	assert.Equal(t, []string{"books"}, op.Tags)
	assert.True(t, op.Deprecated)
	require.Len(t, op.Security, 1)
	assert.Equal(t, []string{"books:write"}, op.Security[0]["oauth"])
	require.Len(t, op.Parameters, 1)
	assert.Equal(t, "#/components/parameters/BookId", op.Parameters[0].Ref)
	assert.True(t, op.Parameters[0].Deprecated)

	require.NotNil(t, op.XAEPLongRunningOperation)
	require.NotNil(t, op.XAEPLongRunningOperation.Response.Schema)
	assert.Equal(t, "#/components/schemas/Book", op.XAEPLongRunningOperation.Response.Schema.Ref)
	require.NotNil(t, op.XAEPLongRunningOperation.Metadata)
	require.NotNil(t, op.XAEPLongRunningOperation.Metadata.Schema)
	assert.Equal(t, "#/components/schemas/ArchiveMetadata", op.XAEPLongRunningOperation.Metadata.Schema.Ref)

	// Verify metadata is omitted when nil during JSON marshaling.
	lroNoMeta := XAEPLongRunningOperation{
		Response: XAEPLongRunningOperationResponse{Schema: &Schema{Type: "object"}},
	}
	marshaled, err := json.Marshal(lroNoMeta)
	require.NoError(t, err)
	assert.NotContains(t, string(marshaled), "metadata")
}
