package handlers

import (
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListIncidentUpdatesTool lists incident updates
type ListIncidentUpdatesTool struct {
	apiClient *client.Client
}

func NewListIncidentUpdatesTool(c *client.Client) *ListIncidentUpdatesTool {
	return &ListIncidentUpdatesTool{apiClient: c}
}

func (t *ListIncidentUpdatesTool) Name() string {
	return "list_incident_updates"
}

func (t *ListIncidentUpdatesTool) Description() string {
	return "List incident updates (status messages posted during an incident).\n\n" +
		"CRITICAL: If user mentions an incident reference like 'INC-1691', you MUST first resolve it to the full incident ID!\n\n" +
		"WORKFLOW FOR INCIDENT REFERENCES:\n" +
		"1. User says: 'updates for INC-1691'\n" +
		"2. FIRST: get_incident({\"incident_id\": \"1691\"}) to get the full incident details\n" +
		"3. Extract the full incident ID from the response (e.g., '01K3VHM0T0ZTMG9JPJ9GESB7XX')\n" +
		"4. THEN: list_incident_updates({\"incident_id\": \"01K3VHM0T0ZTMG9JPJ9GESB7XX\"})\n\n" +
		"IMPORTANT: This endpoint requires the FULL incident ID (01K3VHM0T0ZTMG9JPJ9GESB7XX), NOT the short reference (1691)!\n\n" +
		"EXAMPLES:\n" +
		"- User: 'updates for INC-1691' → get_incident({\"incident_id\": \"1691\"}) → list_incident_updates({\"incident_id\": \"01K3VHM0T0ZTMG9JPJ9GESB7XX\"})"
}

func (t *ListIncidentUpdatesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Filter updates by incident ID. IMPORTANT: This must be the FULL incident ID (e.g., '01K3VHM0T0ZTMG9JPJ9GESB7XX'), NOT the short reference (e.g., '1691' from INC-1691). Use get_incident first to resolve references.",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListIncidentUpdatesTool) Execute(args map[string]interface{}) (string, error) {
	opts := &client.ListIncidentUpdatesOptions{
		IncidentID: GetStringArg(args, "incident_id"),
		PageSize:   GetIntArg(args, "page_size", 25),
	}

	resp, err := t.apiClient.ListIncidentUpdates(opts)
	if err != nil {
		return "", err
	}

	return FormatJSONResponse(resp)
}

// GetIncidentUpdateTool gets a specific incident update
type GetIncidentUpdateTool struct {
	apiClient *client.Client
}

func NewGetIncidentUpdateTool(c *client.Client) *GetIncidentUpdateTool {
	return &GetIncidentUpdateTool{apiClient: c}
}

func (t *GetIncidentUpdateTool) Name() string {
	return "get_incident_update"
}

func (t *GetIncidentUpdateTool) Description() string {
	return "Get details of a specific incident update by ID"
}

func (t *GetIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident update ID",
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	id := GetStringArg(args, "id")
	if id == "" {
		return "", fmt.Errorf("id parameter is required")
	}

	update, err := t.apiClient.GetIncidentUpdate(id)
	if err != nil {
		return "", err
	}

	return FormatJSONResponse(update)
}
