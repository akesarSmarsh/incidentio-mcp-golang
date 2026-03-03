package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListAlertRoutesTool lists alert routes from incident.io
type ListAlertRoutesTool struct {
	apiClient *client.Client
}

func NewListAlertRoutesTool(c *client.Client) *ListAlertRoutesTool {
	return &ListAlertRoutesTool{apiClient: c}
}

func (t *ListAlertRoutesTool) Name() string {
	return "list_alert_routes"
}

func (t *ListAlertRoutesTool) Description() string {
	return "List alert routes from incident.io with optional pagination"
}

func (t *ListAlertRoutesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page",
				"minimum":     1,
				"maximum":     250,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListAlertRoutesTool) Execute(args map[string]interface{}) (string, error) {
	params := &client.ListAlertRoutesParams{}

	if pageSize, ok := args["page_size"].(float64); ok {
		params.PageSize = int(pageSize)
	}
	if after, ok := args["after"].(string); ok {
		params.After = after
	}

	result, err := t.apiClient.ListAlertRoutes(params)
	if err != nil {
		return "", fmt.Errorf("failed to list alert routes: %w", err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}

// GetAlertRouteTool gets details of a specific alert route
type GetAlertRouteTool struct {
	apiClient *client.Client
}

func NewGetAlertRouteTool(c *client.Client) *GetAlertRouteTool {
	return &GetAlertRouteTool{apiClient: c}
}

func (t *GetAlertRouteTool) Name() string {
	return "get_alert_route"
}

func (t *GetAlertRouteTool) Description() string {
	return "Get details of a specific alert route by ID"
}

func (t *GetAlertRouteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The alert route ID",
				"minLength":   1,
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *GetAlertRouteTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("alert route ID is required")
	}

	alertRoute, err := t.apiClient.GetAlertRoute(id)
	if err != nil {
		return "", fmt.Errorf("failed to get alert route: %w", err)
	}

	output, err := json.MarshalIndent(alertRoute, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}
