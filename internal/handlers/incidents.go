package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListIncidentsTool lists incidents from incident.io
type ListIncidentsTool struct {
	apiClient *client.Client
}

func NewListIncidentsTool(c *client.Client) *ListIncidentsTool {
	return &ListIncidentsTool{apiClient: c}
}

func (t *ListIncidentsTool) Name() string {
	return "list_incidents"
}

func (t *ListIncidentsTool) Description() string {
	return "List incidents with filtering. Returns compact summaries by default to avoid large responses.\n\n" +
		"KEY FEATURES:\n" +
		"- search: Filter by name (e.g., search='speechify' finds all Speechify incidents)\n" +
		"- summary: true (default) returns compact summaries, false returns full details\n" +
		"- page_size: Default 25, increase only if needed\n\n" +
		"RESPONSE FORMAT (summary=true, default):\n" +
		"Returns: reference, name, status, severity, created_at, updated_at, permalink\n" +
		"This is much smaller than full incident objects!\n\n" +
		"EXAMPLES:\n" +
		"Find Speechify incidents: list_incidents({\"search\": \"speechify\"})\n" +
		"Recent incidents: list_incidents({\"created_at_gte\": \"2025-01-28\"})\n" +
		"Full details: list_incidents({\"search\": \"speechify\", \"summary\": false})\n\n" +
		"PAGINATION:\n" +
		"If has_more_results=true, call again with 'after' cursor from pagination_meta.\n\n" +
		"INCIDENT REFERENCE RESOLUTION:\n" +
		"For INC-1691, use get_incident({\"incident_id\": \"1691\"}) for full details.\n\n" +
		"Date format: \"2025-10-15\". Use current year when calculating dates!"
}

func (t *ListIncidentsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"search": map[string]interface{}{
				"type":        "string",
				"description": "Filter incidents by name (case-insensitive substring match). Example: 'speechify' returns all incidents with 'speechify' in the name.",
			},
			"summary": map[string]interface{}{
				"type":        "boolean",
				"description": "Return compact summaries instead of full incident details. Defaults to true. Set to false for full details (warning: large response).",
				"default":     true,
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page. Default is 25. Use 50-100 only if you need more results per page.",
				"default":     25,
				"minimum":     1,
				"maximum":     100,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor from previous response. Get this from pagination_meta.after in the previous response.",
			},
			"status": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by status. Values: triage, active, investigating, monitoring, resolved, closed. Example: ['active', 'triage']",
			},
			"severity_one_of": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by exact severity IDs. Use list_severities to get IDs. Example: ['01ABC123']",
			},
			"severity_gte": map[string]interface{}{
				"type":        "string",
				"description": "Filter by severity rank >= this ID. Returns this severity and all more severe. Example: 'sev_major_id' returns Major, Critical.",
			},
			"severity_lte": map[string]interface{}{
				"type":        "string",
				"description": "Filter by severity rank <= this ID. Returns this severity and all less severe. Example: 'sev_major_id' returns Major, Minor, Low.",
			},
			"created_at_gte": map[string]interface{}{
				"type":        "string",
				"description": "Filter incidents created on or after this date. Format: '2025-10-15' or '2025-10-15T10:30:00Z'. Use current year (2025).",
			},
			"created_at_lte": map[string]interface{}{
				"type":        "string",
				"description": "Filter incidents created on or before this date. Format: '2025-10-15' or '2025-10-15T23:59:59Z'. Use current year (2025).",
			},
			"updated_at_gte": map[string]interface{}{
				"type":        "string",
				"description": "Filter incidents updated on or after this date. Format: '2025-10-15' or '2025-10-15T10:30:00Z'. Use current year (2025).",
			},
			"updated_at_lte": map[string]interface{}{
				"type":        "string",
				"description": "Filter incidents updated on or before this date. Format: '2025-10-15' or '2025-10-15T23:59:59Z'. Use current year (2025).",
			},
			"custom_field_id": map[string]interface{}{
				"type":        "string",
				"description": "Custom field ID to filter by. Must use with custom_field_value. Get ID from search_custom_fields.",
			},
			"custom_field_value": map[string]interface{}{
				"type":        "string",
				"description": "Custom field OPTION ID to match. For select fields, this must be the option's ID (e.g., '01JQ7...'), not the label. Get from the options array of search_custom_fields response.",
			},
		},
	}
}

// IncidentSummary is a lightweight representation for list responses
type IncidentSummary struct {
	Reference string `json:"reference"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Severity  string `json:"severity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Permalink string `json:"permalink"`
}

func (t *ListIncidentsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &client.ListIncidentsOptions{
		PageSize: 25, // Default to 25 for reasonable response sizes
	}

	// Parse search filter (client-side - the incident.io API does not support name filtering)
	searchFilter := ""
	if search, ok := args["search"].(string); ok && search != "" {
		searchFilter = strings.ToLower(search)
	}

	// Parse summary mode (default true)
	summaryMode := true
	if summary, ok := args["summary"].(bool); ok {
		summaryMode = summary
	}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	if after, ok := args["after"].(string); ok && after != "" {
		opts.After = after
	}

	if statuses, ok := args["status"].([]interface{}); ok {
		for _, s := range statuses {
			if str, ok := s.(string); ok {
				opts.Status = append(opts.Status, str)
			}
		}
	}

	if severities, ok := args["severity_one_of"].([]interface{}); ok {
		for _, s := range severities {
			if str, ok := s.(string); ok {
				opts.SeverityOneOf = append(opts.SeverityOneOf, str)
			}
		}
	}

	if severityGte, ok := args["severity_gte"].(string); ok && severityGte != "" {
		opts.SeverityGte = severityGte
	}

	if severityLte, ok := args["severity_lte"].(string); ok && severityLte != "" {
		opts.SeverityLte = severityLte
	}

	if createdAtGte, ok := args["created_at_gte"].(string); ok && createdAtGte != "" {
		opts.CreatedAtGte = createdAtGte
	}

	if createdAtLte, ok := args["created_at_lte"].(string); ok && createdAtLte != "" {
		opts.CreatedAtLte = createdAtLte
	}

	if updatedAtGte, ok := args["updated_at_gte"].(string); ok && updatedAtGte != "" {
		opts.UpdatedAtGte = updatedAtGte
	}

	if updatedAtLte, ok := args["updated_at_lte"].(string); ok && updatedAtLte != "" {
		opts.UpdatedAtLte = updatedAtLte
	}

	// Handle custom field filtering - API format: custom_field[ID][one_of]=option_id
	if customFieldID, ok := args["custom_field_id"].(string); ok && customFieldID != "" {
		if customFieldValue, ok := args["custom_field_value"].(string); ok && customFieldValue != "" {
			if opts.CustomFieldOneOf == nil {
				opts.CustomFieldOneOf = make(map[string]string)
			}
			opts.CustomFieldOneOf[customFieldID] = customFieldValue
		}
	}

	resp, err := t.apiClient.ListIncidents(opts)
	if err != nil {
		return "", err
	}

	// Apply search filter client-side if provided
	filteredIncidents := resp.Incidents
	if searchFilter != "" {
		filteredIncidents = nil
		for _, inc := range resp.Incidents {
			if strings.Contains(strings.ToLower(inc.Name), searchFilter) {
				filteredIncidents = append(filteredIncidents, inc)
			}
		}
	}

	// Build response based on summary mode
	var incidentsData interface{}
	if summaryMode {
		summaries := make([]IncidentSummary, 0, len(filteredIncidents))
		for _, inc := range filteredIncidents {
			summaries = append(summaries, IncidentSummary{
				Reference: inc.Reference,
				Name:      inc.Name,
				Status:    inc.IncidentStatus.Name,
				Severity:  inc.Severity.Name,
				CreatedAt: inc.CreatedAt.Format(time.RFC3339),
				UpdatedAt: inc.UpdatedAt.Format(time.RFC3339),
				Permalink: inc.Permalink,
			})
		}
		incidentsData = summaries
	} else {
		incidentsData = filteredIncidents
	}

	// Create response with prominent pagination info
	response := map[string]interface{}{
		"incidents":       incidentsData,
		"pagination_meta": resp.PaginationMeta,
		"count":           len(filteredIncidents),
	}

	// Add note about search filtering if applied
	if searchFilter != "" {
		response["search_applied"] = searchFilter
		response["search_note"] = fmt.Sprintf("Filtered %d incidents from %d total on this page", len(filteredIncidents), len(resp.Incidents))
	}

	// Add prominent pagination status
	// Use total_record_count to determine if there are more results
	// The "after" cursor is only needed for the next API call, not for determining if more results exist
	recordsFetched := len(resp.Incidents)
	totalRecords := resp.PaginationMeta.TotalRecordCount
	hasMore := recordsFetched < totalRecords

	if hasMore {
		response["has_more_results"] = true
		response["pagination_progress"] = map[string]interface{}{
			"records_fetched":  recordsFetched,
			"total_records":    totalRecords,
			"remaining":        totalRecords - recordsFetched,
			"progress_percent": fmt.Sprintf("%.1f%%", float64(recordsFetched)/float64(totalRecords)*100),
		}
		response["FETCH_NEXT_PAGE"] = map[string]interface{}{
			"action":  "REQUIRED - You must call list_incidents again to get remaining results",
			"after":   resp.PaginationMeta.After,
			"message": fmt.Sprintf("Fetched %d of %d incidents (%.1f%%). Call list_incidents again with after='%s' plus same filters. Repeat until has_more_results=false.", recordsFetched, totalRecords, float64(recordsFetched)/float64(totalRecords)*100, resp.PaginationMeta.After),
		}
	} else {
		response["has_more_results"] = false
		response["pagination_progress"] = map[string]interface{}{
			"records_fetched":  recordsFetched,
			"total_records":    totalRecords,
			"remaining":        0,
			"progress_percent": "100.0%",
		}
		response["pagination_status"] = fmt.Sprintf("COMPLETE - All %d incidents fetched", totalRecords)
	}

	return FormatJSONResponse(response)
}

// GetIncidentTool retrieves a specific incident
type GetIncidentTool struct {
	apiClient *client.Client
}

func NewGetIncidentTool(c *client.Client) *GetIncidentTool {
	return &GetIncidentTool{apiClient: c}
}

func (t *GetIncidentTool) Name() string {
	return "get_incident"
}

func (t *GetIncidentTool) Description() string {
	return "Get details of a specific incident by ID or reference.\n\n" +
		"ACCEPTS BOTH:\n" +
		"- Full incident ID: '01K3VHM0T0ZTMG9JPJ9GESB7XX'\n" +
		"- Short reference: '1691' (just the number from INC-1691)\n\n" +
		"USE THIS TOOL TO:\n" +
		"- Resolve incident references (INC-1691 → use '1691') to get full incident details\n" +
		"- Get the FULL incident ID needed for other API calls (follow-ups, updates, etc.)\n" +
		"- The response contains the full incident ID that other endpoints require\n\n" +
		"EXAMPLES:\n" +
		"- get_incident({\"incident_id\": \"1691\"}) - Get incident by reference number\n" +
		"- get_incident({\"incident_id\": \"01K3VHM0T0ZTMG9JPJ9GESB7XX\"}) - Get incident by full ID"
}

func (t *GetIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID or reference number. Accepts both full IDs (01K3VHM0T0ZTMG9JPJ9GESB7XX) and reference numbers (1691 from INC-1691).",
			},
		},
		"required":             []interface{}{"incident_id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["incident_id"].(string)
	if !ok || id == "" {
		argDetails := make(map[string]interface{})
		for key, value := range args {
			argDetails[key] = value
		}
		return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string. Received parameters: %+v", argDetails)
	}

	incident, err := t.apiClient.GetIncident(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
