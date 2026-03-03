package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListCatalogTypesTool lists available catalog types
type ListCatalogTypesTool struct {
	apiClient *client.Client
}

func NewListCatalogTypesTool(c *client.Client) *ListCatalogTypesTool {
	return &ListCatalogTypesTool{apiClient: c}
}

func (t *ListCatalogTypesTool) Name() string {
	return "list_catalog_types"
}

func (t *ListCatalogTypesTool) Description() string {
	return "List available catalog types in incident.io (filtered to only include types with TypeName starting with 'Custom' - case insensitive)"
}

func (t *ListCatalogTypesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":                 "object",
		"properties":           map[string]interface{}{},
		"additionalProperties": false,
	}
}

func (t *ListCatalogTypesTool) Execute(args map[string]interface{}) (string, error) {
	result, err := t.apiClient.ListCatalogTypes()
	if err != nil {
		return "", fmt.Errorf("failed to list catalog types: %w", err)
	}

	// Filter catalog types to only include those with TypeName starting with "Custom" (case-insensitive)
	var filteredTypes []client.CatalogType
	for _, catalogType := range result.CatalogTypes {
		if strings.HasPrefix(strings.ToLower(catalogType.TypeName), "custom") {
			filteredTypes = append(filteredTypes, catalogType)
		}
	}

	output := fmt.Sprintf("Found %d catalog types (filtered for Custom* names):\n\n", len(filteredTypes))

	for _, catalogType := range filteredTypes {
		output += fmt.Sprintf("ID: %s\n", catalogType.ID)
		output += fmt.Sprintf("Name: %s\n", catalogType.Name)
		output += fmt.Sprintf("Type Name: %s\n", catalogType.TypeName)
		if catalogType.Description != "" {
			output += fmt.Sprintf("Description: %s\n", catalogType.Description)
		}
		if catalogType.Color != "" {
			output += fmt.Sprintf("Color: %s\n", catalogType.Color)
		}
		if catalogType.Icon != "" {
			output += fmt.Sprintf("Icon: %s\n", catalogType.Icon)
		}
		if len(catalogType.Attributes) > 0 {
			output += fmt.Sprintf("Attributes (%d):\n", len(catalogType.Attributes))
			for _, attr := range catalogType.Attributes {
				output += fmt.Sprintf("  - %s (%s): %s\n", attr.Name, attr.Type, attr.ID)
			}
		}
		output += fmt.Sprintf("Created: %s\n", catalogType.CreatedAt.Format("2006-01-02 15:04:05"))
		output += fmt.Sprintf("Updated: %s\n", catalogType.UpdatedAt.Format("2006-01-02 15:04:05"))
		output += "\n"
	}

	// Also return the raw JSON (only filtered types)
	filteredResult := &client.ListCatalogTypesResponse{
		CatalogTypes: filteredTypes,
		ListResponse: result.ListResponse,
	}
	jsonOutput, err := json.MarshalIndent(filteredResult, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}

// ListCatalogEntriesTool lists catalog entries for a given type
type ListCatalogEntriesTool struct {
	apiClient *client.Client
}

func NewListCatalogEntriesTool(c *client.Client) *ListCatalogEntriesTool {
	return &ListCatalogEntriesTool{apiClient: c}
}

func (t *ListCatalogEntriesTool) Name() string {
	return "list_catalog_entries"
}

func (t *ListCatalogEntriesTool) Description() string {
	return "List catalog entries for a given catalog type. DO NOT use this for finding custom field options - use search_custom_fields or list_custom_fields instead, which include the options array in their response."
}

func (t *ListCatalogEntriesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"catalog_type_id": map[string]interface{}{
				"type":        "string",
				"description": "The catalog type ID to list entries for",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of entries to return per page (default: 25)",
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
			"identifier": map[string]interface{}{
				"type":        "string",
				"description": "Filter by identifier",
			},
		},
		"required":             []interface{}{"catalog_type_id"},
		"additionalProperties": false,
	}
}

func (t *ListCatalogEntriesTool) Execute(args map[string]interface{}) (string, error) {
	catalogTypeID, ok := args["catalog_type_id"].(string)
	if !ok || catalogTypeID == "" {
		return "", fmt.Errorf("catalog_type_id parameter is required")
	}

	opts := client.ListCatalogEntriesOptions{
		CatalogTypeID: catalogTypeID,
	}

	if pageSize, ok := args["page_size"]; ok {
		if ps, ok := pageSize.(float64); ok {
			opts.PageSize = int(ps)
		} else if ps, ok := pageSize.(string); ok {
			if parsed, err := strconv.Atoi(ps); err == nil {
				opts.PageSize = parsed
			}
		}
	}

	if after, ok := args["after"].(string); ok {
		opts.After = after
	}

	if identifier, ok := args["identifier"].(string); ok {
		opts.Identifier = identifier
	}

	result, err := t.apiClient.ListCatalogEntries(opts)
	if err != nil {
		return "", fmt.Errorf("failed to list catalog entries: %w", err)
	}

	output := fmt.Sprintf("Found %d catalog entries for type %s:\n\n", len(result.CatalogEntries), catalogTypeID)

	for _, entry := range result.CatalogEntries {
		output += fmt.Sprintf("ID: %s\n", entry.ID)
		output += fmt.Sprintf("Name: %s\n", entry.Name)
		if len(entry.Aliases) > 0 {
			output += fmt.Sprintf("Aliases: %v\n", entry.Aliases)
		}
		if entry.ExternalID != "" {
			output += fmt.Sprintf("External ID: %s\n", entry.ExternalID)
		}
		output += fmt.Sprintf("Rank: %d\n", entry.Rank)
		if len(entry.AttributeValues) > 0 {
			output += "Attributes:\n"
			for key, value := range entry.AttributeValues {
				if value.Value != nil {
					if value.Value.Literal != "" {
						output += fmt.Sprintf("  %s: %s\n", key, value.Value.Literal)
					} else if value.Value.ID != "" {
						output += fmt.Sprintf("  %s: %s (ID)\n", key, value.Value.ID)
					}
				}
				if len(value.ArrayValue) > 0 {
					output += fmt.Sprintf("  %s: [", key)
					for i, v := range value.ArrayValue {
						if i > 0 {
							output += ", "
						}
						if v.Literal != "" {
							output += v.Literal
						} else if v.ID != "" {
							output += v.ID + " (ID)"
						}
					}
					output += "]\n"
				}
			}
		}
		output += fmt.Sprintf("Created: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
		output += fmt.Sprintf("Updated: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
		output += "\n"
	}

	// Add pagination info
	if result.PaginationMeta.After != "" {
		output += fmt.Sprintf("Pagination: Next page available (after: %s)\n", result.PaginationMeta.After)
	}
	output += fmt.Sprintf("Total entries: %d\n", result.PaginationMeta.TotalRecordCount)

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}
