package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListIncidentRolesTool lists available incident roles
type ListIncidentRolesTool struct {
	apiClient *client.Client
}

func NewListIncidentRolesTool(c *client.Client) *ListIncidentRolesTool {
	return &ListIncidentRolesTool{apiClient: c}
}

func (t *ListIncidentRolesTool) Name() string {
	return "list_available_incident_roles"
}

func (t *ListIncidentRolesTool) Description() string {
	return "List all available incident roles that can be assigned to users"
}

func (t *ListIncidentRolesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
	}
}

func (t *ListIncidentRolesTool) Execute(args map[string]interface{}) (string, error) {
	opts := &client.ListIncidentRolesOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	resp, err := t.apiClient.ListIncidentRoles(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// ListUsersTool lists available users for role assignment
type ListUsersTool struct {
	apiClient *client.Client
}

func NewListUsersTool(c *client.Client) *ListUsersTool {
	return &ListUsersTool{apiClient: c}
}

func (t *ListUsersTool) Name() string {
	return "list_users"
}

func (t *ListUsersTool) Description() string {
	return "List all users available for incident role assignment (automatically paginated to return ALL users)"
}

func (t *ListUsersTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     250,
			},
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Filter users by email address",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListUsersTool) Execute(args map[string]interface{}) (string, error) {
	opts := &client.ListUsersOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	if email, ok := args["email"].(string); ok && email != "" {
		opts.Email = email
	}

	resp, err := t.apiClient.ListUsers(opts)
	if err != nil {
		return "", err
	}

	// Add a helpful message about the results
	var output string
	if opts.Email != "" {
		output = fmt.Sprintf("Users matching email '%s':\n", opts.Email)
	} else {
		output = fmt.Sprintf("Found %d users:\n", len(resp.Users))
	}

	// Format users in a more readable way
	for _, user := range resp.Users {
		output += fmt.Sprintf("\n- Name: %s\n  Email: %s\n  ID: %s\n  Role: %s\n",
			user.Name, user.Email, user.ID, user.Role)
	}

	// Also include the raw JSON
	jsonResult, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return output, nil // Return readable output even if JSON fails
	}

	output += "\n\nRaw JSON response:\n" + string(jsonResult)

	return output, nil
}
