package openapi

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

type componentRequestBodySchema struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Required   []string       `json:"required,omitempty"`
}

func generateComponents(sch schema.Schema) (Components, error) {
	components := Components{
		Schemas:       map[string]any{},
		RequestBodies: map[string]any{},
		Responses:     map[string]any{},
	}

	for _, typeNode := range sch.GetTypeNodes() {
		desc := "is a domain type defined in UFO RPC with no documentation."
		if typeNode.Doc != nil {
			desc = strings.TrimSpace(strutil.NormalizeIndent(*typeNode.Doc))
		}

		if typeNode.Deprecated != nil {
			desc += "\n\nDeprecated: "
			if *typeNode.Deprecated == "" {
				desc += "This type is deprecated and should not be used in new code."
			} else {
				desc += *typeNode.Deprecated
			}
		}

		properties, requiredFields := generateProperties(typeNode.Fields)

		typeSchema := map[string]any{
			"deprecated":  typeNode.Deprecated != nil,
			"type":        "object",
			"description": desc,
			"properties":  properties,
		}
		if len(requiredFields) > 0 {
			typeSchema["required"] = requiredFields
		}

		components.Schemas[typeNode.Name] = typeSchema
	}

	for _, procNode := range sch.GetProcNodes() {
		name := procNode.Name
		inputName := fmt.Sprintf("%sInput", name)
		outputName := fmt.Sprintf("%sOutput", name)

		inputProperties, inputRequiredFields := generateProperties(procNode.Input)
		components.RequestBodies[inputName] = map[string]any{
			"description": "Request body for the " + name + " procedure",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": componentRequestBodySchema{
						Type:       "object",
						Properties: inputProperties,
						Required:   inputRequiredFields,
					},
				},
			},
		}

		outputProperties, outputRequiredFields := generateOutputProperties(procNode.Output)
		components.Responses[outputName] = map[string]any{
			"description": "Response for the " + name + " procedure both for success and error cases based on the `ok` field.",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": componentRequestBodySchema{
						Type:       "object",
						Properties: outputProperties,
						Required:   outputRequiredFields,
					},
				},
			},
		}
	}

	for _, streamNode := range sch.GetStreamNodes() {
		name := streamNode.Name
		inputName := fmt.Sprintf("%sInput", name)
		outputName := fmt.Sprintf("%sOutput", name)

		inputProperties, inputRequiredFields := generateProperties(streamNode.Input)
		components.RequestBodies[inputName] = map[string]any{
			"description": "Request body for the " + name + " stream",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": componentRequestBodySchema{
						Type:       "object",
						Properties: inputProperties,
						Required:   inputRequiredFields,
					},
				},
			},
		}

		outputProperties, outputRequiredFields := generateOutputProperties(streamNode.Output)
		components.Responses[outputName] = map[string]any{
			"description": "Response for the " + name + " stream, both for success and error cases based on the `ok` field.",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": componentRequestBodySchema{
						Type:       "object",
						Properties: outputProperties,
						Required:   outputRequiredFields,
					},
				},
			},
		}
	}

	return components, nil
}
