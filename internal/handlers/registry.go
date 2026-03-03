package handlers

import (
	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ToolRegistry manages tool registration and provides utilities for common patterns
type ToolRegistry struct {
	tools map[string]Handler
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Handler),
	}
}

// RegisterTool registers a tool with the registry
func (r *ToolRegistry) RegisterTool(name string, tool Handler) {
	r.tools[name] = tool
}

// GetTools returns all registered tools
func (r *ToolRegistry) GetTools() map[string]Handler {
	return r.tools
}

// RegisterIncidentTools registers all incident-related tools
func (r *ToolRegistry) RegisterIncidentTools(c *client.Client) {
	r.RegisterTool("list_incidents", NewListIncidentsTool(c))
	r.RegisterTool("get_incident", NewGetIncidentTool(c))
	r.RegisterTool("list_incident_statuses", NewListIncidentStatusesTool(c))
	r.RegisterTool("list_incident_types", NewListIncidentTypesTool(c))
}

// RegisterIncidentUpdateTools registers all incident update tools
func (r *ToolRegistry) RegisterIncidentUpdateTools(c *client.Client) {
	r.RegisterTool("list_incident_updates", NewListIncidentUpdatesTool(c))
	r.RegisterTool("get_incident_update", NewGetIncidentUpdateTool(c))
}

// RegisterFollowUpTools registers all follow-up tools
func (r *ToolRegistry) RegisterFollowUpTools(c *client.Client) {
	r.RegisterTool("list_follow_ups", NewListFollowUpsTool(c))
	r.RegisterTool("get_follow_up", NewGetFollowUpTool(c))
}

// RegisterAlertTools registers all alert tools
func (r *ToolRegistry) RegisterAlertTools(c *client.Client) {
	r.RegisterTool("list_alerts", NewListAlertsTool(c))
	r.RegisterTool("get_alert", NewGetAlertTool(c))
	r.RegisterTool("list_incident_alerts", NewListIncidentAlertsTool(c))
}

// RegisterActionTools registers all action tools
func (r *ToolRegistry) RegisterActionTools(c *client.Client) {
	r.RegisterTool("list_actions", NewListActionsTool(c))
	r.RegisterTool("get_action", NewGetActionTool(c))
}

// RegisterRoleTools registers all role-related tools
func (r *ToolRegistry) RegisterRoleTools(c *client.Client) {
	r.RegisterTool("list_available_incident_roles", NewListIncidentRolesTool(c))
	r.RegisterTool("list_users", NewListUsersTool(c))
}

// RegisterWorkflowTools registers all workflow tools
func (r *ToolRegistry) RegisterWorkflowTools(c *client.Client) {
	r.RegisterTool("list_workflows", NewListWorkflowsTool(c))
	r.RegisterTool("get_workflow", NewGetWorkflowTool(c))
}

// RegisterAlertRouteTools registers all alert route tools
func (r *ToolRegistry) RegisterAlertRouteTools(c *client.Client) {
	r.RegisterTool("list_alert_routes", NewListAlertRoutesTool(c))
	r.RegisterTool("get_alert_route", NewGetAlertRouteTool(c))
}

// RegisterAlertSourceTools registers all alert source tools
func (r *ToolRegistry) RegisterAlertSourceTools(c *client.Client) {
	r.RegisterTool("list_alert_sources", NewListAlertSourcesTool(c))
}

// RegisterCatalogTools registers all catalog tools
func (r *ToolRegistry) RegisterCatalogTools(c *client.Client) {
	r.RegisterTool("list_catalog_types", NewListCatalogTypesTool(c))
	r.RegisterTool("list_catalog_entries", NewListCatalogEntriesTool(c))
}

// RegisterCustomFieldTools registers all custom field tools
func (r *ToolRegistry) RegisterCustomFieldTools(c *client.Client) {
	r.RegisterTool("list_custom_fields", NewListCustomFieldsTool(c))
	r.RegisterTool("get_custom_field", NewGetCustomFieldTool(c))
	r.RegisterTool("search_custom_fields", NewSearchCustomFieldsTool(c))
	r.RegisterTool("list_custom_field_options", NewListCustomFieldOptionsTool(c))
}

// RegisterSeverityTools registers all severity tools
func (r *ToolRegistry) RegisterSeverityTools(c *client.Client) {
	r.RegisterTool("list_severities", NewListSeveritiesTool(c))
	r.RegisterTool("get_severity", NewGetSeverityTool(c))
}

// RegisterScheduleTools registers all schedule tools
func (r *ToolRegistry) RegisterScheduleTools(c *client.Client) {
	r.RegisterTool("list_schedules", NewListSchedulesTool(c))
	r.RegisterTool("get_schedule", NewGetScheduleTool(c))
	r.RegisterTool("list_schedule_entries", NewListScheduleEntriesTool(c))
	r.RegisterTool("get_current_on_call", NewGetCurrentOnCallTool(c))
}

// RegisterAllTools registers all available tools
func (r *ToolRegistry) RegisterAllTools(c *client.Client) {
	r.RegisterIncidentTools(c)
	r.RegisterIncidentUpdateTools(c)
	r.RegisterFollowUpTools(c)
	r.RegisterAlertTools(c)
	r.RegisterActionTools(c)
	r.RegisterRoleTools(c)
	r.RegisterWorkflowTools(c)
	r.RegisterAlertRouteTools(c)
	r.RegisterAlertSourceTools(c)
	r.RegisterCatalogTools(c)
	r.RegisterCustomFieldTools(c)
	r.RegisterSeverityTools(c)
	r.RegisterScheduleTools(c)
}
