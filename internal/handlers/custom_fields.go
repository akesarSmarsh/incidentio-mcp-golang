package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListCustomFieldsTool lists all custom fields
type ListCustomFieldsTool struct {
	apiClient *client.Client
}

func NewListCustomFieldsTool(c *client.Client) *ListCustomFieldsTool {
	return &ListCustomFieldsTool{apiClient: c}
}

func (t *ListCustomFieldsTool) Name() string {
	return "list_custom_fields"
}

func (t *ListCustomFieldsTool) Description() string {
	return "List all custom fields configured in incident.io. Use this to discover what custom fields exist (like \"Affected Team\", \"Priority\", etc.) before filtering incidents by them.\n\n" +
		"WHEN TO USE: If you need to filter incidents by a custom attribute (team, department, priority, etc.) but don't know the exact field name, call this first to see what custom fields are available, then use search_custom_fields or use the ID directly in list_incidents."
}

func (t *ListCustomFieldsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *ListCustomFieldsTool) Execute(args map[string]interface{}) (string, error) {
	resp, err := t.apiClient.ListCustomFields()
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// GetCustomFieldTool retrieves a specific custom field
type GetCustomFieldTool struct {
	apiClient *client.Client
}

func NewGetCustomFieldTool(c *client.Client) *GetCustomFieldTool {
	return &GetCustomFieldTool{apiClient: c}
}

func (t *GetCustomFieldTool) Name() string {
	return "get_custom_field"
}

func (t *GetCustomFieldTool) Description() string {
	return "Get details of a specific custom field by ID, including its options, type, and configuration"
}

func (t *GetCustomFieldTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The custom field ID",
			},
		},
		"required": []string{"id"},
	}
}

func (t *GetCustomFieldTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}

	field, err := t.apiClient.GetCustomField(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// SearchCustomFieldsTool searches for custom fields
type SearchCustomFieldsTool struct {
	apiClient *client.Client
}

func NewSearchCustomFieldsTool(c *client.Client) *SearchCustomFieldsTool {
	return &SearchCustomFieldsTool{apiClient: c}
}

func (t *SearchCustomFieldsTool) Name() string {
	return "search_custom_fields"
}

func (t *SearchCustomFieldsTool) Description() string {
	return "Search for custom fields by name to get their IDs AND options for filtering.\n\n" +
		"IMPORTANT: When you get the custom field, check its 'options' array. If it's a select field:\n" +
		"- The custom_field_value must be the OPTION ID (e.g., '01ABC...'), NOT the option label\n" +
		"- Look for the option with matching 'value' field, then use its 'id'\n\n" +
		"WORKFLOW for team/department filtering:\n" +
		"1. search_custom_fields({\"query\": \"team\"}) → get field and its options\n" +
		"2. Find the option where value=\"Engineering\" → get its id (e.g., '01XYZ...')\n" +
		"3. list_incidents({\"custom_field_id\": \"cf_123\", \"custom_field_value\": \"01XYZ...\"})  ← Use option ID!\n\n" +
		"Example:\n" +
		"User: 'show Engineering team incidents'\n" +
		"→ search_custom_fields({\"query\": \"team\"})\n" +
		"  Returns: {id: 'cf_123', options: [{id: '01ABC', value: 'Engineering'}, ...]}\n" +
		"→ list_incidents({\"custom_field_id\": \"cf_123\", \"custom_field_value\": \"01ABC\"})"
}

func (t *SearchCustomFieldsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Search query to match against custom field names and descriptions",
			},
			"field_type": map[string]interface{}{
				"type":        "string",
				"description": "Filter by field type: single_select, multi_select, text, link, numeric, etc.",
			},
		},
	}
}

func (t *SearchCustomFieldsTool) Execute(args map[string]interface{}) (string, error) {
	query := GetStringArg(args, "query")
	fieldType := GetStringArg(args, "field_type")

	fields, err := t.apiClient.SearchCustomFields(query, fieldType)
	if err != nil {
		return "", err
	}

	response := CreateSimpleResponse(fields, "")
	response["custom_fields"] = fields

	return FormatJSONResponse(response)
}

// ListCustomFieldOptionsTool lists all custom field options
type ListCustomFieldOptionsTool struct {
	apiClient *client.Client
}

func NewListCustomFieldOptionsTool(c *client.Client) *ListCustomFieldOptionsTool {
	return &ListCustomFieldOptionsTool{apiClient: c}
}

func (t *ListCustomFieldOptionsTool) Name() string {
	return "list_custom_field_options"
}

func (t *ListCustomFieldOptionsTool) Description() string {
	return "List all custom field options across all custom fields. Useful for understanding available values for select fields."
}

func (t *ListCustomFieldOptionsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *ListCustomFieldOptionsTool) Execute(args map[string]interface{}) (string, error) {
	options, err := t.apiClient.ListCustomFieldOptions()
	if err != nil {
		return "", err
	}

	response := CreateSimpleResponse(options, "")
	response["custom_field_options"] = options

	return FormatJSONResponse(response)
}
