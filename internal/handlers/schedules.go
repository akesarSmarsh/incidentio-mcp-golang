package handlers

import (
	"fmt"
	"time"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListSchedulesTool lists all on-call schedules
type ListSchedulesTool struct {
	*BaseTool
}

func NewListSchedulesTool(c *client.Client) *ListSchedulesTool {
	return &ListSchedulesTool{
		BaseTool: NewBaseTool(c),
	}
}

func (t *ListSchedulesTool) Name() string {
	return "list_schedules"
}

func (t *ListSchedulesTool) Description() string {
	return "List all on-call schedules in incident.io"
}

func (t *ListSchedulesTool) InputSchema() map[string]interface{} {
	properties := t.StandardPaginationProperties()
	return t.StandardInputSchema(properties, []string{})
}

func (t *ListSchedulesTool) Execute(args map[string]interface{}) (string, error) {
	pageSize := t.ValidateOptionalInt(args, "page_size", 25)
	after := t.ValidateOptionalString(args, "after")

	result, err := t.GetClient().ListSchedules(pageSize, after)
	if err != nil {
		return "", fmt.Errorf("failed to list schedules: %w", err)
	}

	response := t.CreatePaginationResponse(result.Schedules, &result.PaginationMeta, len(result.Schedules))
	return t.FormatResponse(response)
}

// GetScheduleTool gets a specific schedule by ID
type GetScheduleTool struct {
	*BaseTool
}

func NewGetScheduleTool(c *client.Client) *GetScheduleTool {
	return &GetScheduleTool{
		BaseTool: NewBaseTool(c),
	}
}

func (t *GetScheduleTool) Name() string {
	return "get_schedule"
}

func (t *GetScheduleTool) Description() string {
	return "Get details of a specific on-call schedule by ID"
}

func (t *GetScheduleTool) InputSchema() map[string]interface{} {
	properties := map[string]interface{}{
		"schedule_id": map[string]interface{}{
			"type":        "string",
			"description": "The schedule ID (e.g., '01JXXESWWYB00TVX9FTK8XNZZJ')",
		},
	}
	return t.StandardInputSchema(properties, []string{"schedule_id"})
}

func (t *GetScheduleTool) Execute(args map[string]interface{}) (string, error) {
	scheduleID, err := t.ValidateRequiredString(args, "schedule_id")
	if err != nil {
		return "", err
	}

	schedule, err := t.GetClient().GetSchedule(scheduleID)
	if err != nil {
		return "", fmt.Errorf("failed to get schedule: %w", err)
	}

	response := t.CreateSimpleResponse(schedule, "")
	return t.FormatResponse(response)
}

// ListScheduleEntriesTool lists on-call entries for a schedule
type ListScheduleEntriesTool struct {
	*BaseTool
}

func NewListScheduleEntriesTool(c *client.Client) *ListScheduleEntriesTool {
	return &ListScheduleEntriesTool{
		BaseTool: NewBaseTool(c),
	}
}

func (t *ListScheduleEntriesTool) Name() string {
	return "list_schedule_entries"
}

func (t *ListScheduleEntriesTool) Description() string {
	return "List on-call entries (who is on-call when) for a specific schedule. Use this to find out who is currently on-call or who will be on-call in a given time range."
}

func (t *ListScheduleEntriesTool) InputSchema() map[string]interface{} {
	properties := t.StandardPaginationProperties()
	properties["schedule_id"] = map[string]interface{}{
		"type":        "string",
		"description": "The schedule ID (e.g., '01JXXESWWYB00TVX9FTK8XNZZJ')",
	}
	properties["from"] = map[string]interface{}{
		"type":        "string",
		"description": "Start time in ISO 8601 format (e.g., '2026-02-20T00:00:00Z'). Defaults to current time if not provided.",
	}
	properties["to"] = map[string]interface{}{
		"type":        "string",
		"description": "End time in ISO 8601 format (e.g., '2026-02-27T00:00:00Z'). Defaults to 7 days from start time if not provided.",
	}
	return t.StandardInputSchema(properties, []string{"schedule_id"})
}

func (t *ListScheduleEntriesTool) Execute(args map[string]interface{}) (string, error) {
	scheduleID, err := t.ValidateRequiredString(args, "schedule_id")
	if err != nil {
		return "", err
	}

	from := t.ValidateOptionalString(args, "from")
	to := t.ValidateOptionalString(args, "to")
	pageSize := t.ValidateOptionalInt(args, "page_size", 25)
	after := t.ValidateOptionalString(args, "after")

	// Default to current time if from is not provided
	if from == "" {
		from = time.Now().Format(time.RFC3339)
	}

	result, err := t.GetClient().ListScheduleEntries(scheduleID, from, to, pageSize, after)
	if err != nil {
		return "", fmt.Errorf("failed to list schedule entries: %w", err)
	}

	// Combine all entries for display (final contains computed on-call entries)
	allEntries := append(result.Final, result.Overrides...)
	allEntries = append(allEntries, result.Scheduled...)

	// Build response with all entry types
	response := map[string]interface{}{
		"schedule_id": scheduleID,
		"final":       result.Final,
		"overrides":   result.Overrides,
		"scheduled":   result.Scheduled,
		"count":       len(allEntries),
	}

	// Find current on-call person from final entries
	now := time.Now()
	for _, entry := range result.Final {
		if entry.StartAt.Before(now) && entry.EndAt.After(now) {
			response["current_on_call"] = map[string]interface{}{
				"user": map[string]interface{}{
					"id":    entry.User.ID,
					"name":  entry.User.Name,
					"email": entry.User.Email,
				},
				"starts_at": entry.StartAt,
				"ends_at":   entry.EndAt,
				"layer_id":  entry.LayerID,
			}
			break
		}
	}

	return t.FormatResponse(response)
}

// GetCurrentOnCallTool gets the current on-call person for a schedule
type GetCurrentOnCallTool struct {
	*BaseTool
}

func NewGetCurrentOnCallTool(c *client.Client) *GetCurrentOnCallTool {
	return &GetCurrentOnCallTool{
		BaseTool: NewBaseTool(c),
	}
}

func (t *GetCurrentOnCallTool) Name() string {
	return "get_current_on_call"
}

func (t *GetCurrentOnCallTool) Description() string {
	return "Get the current on-call person for a specific schedule. This is a convenience tool that returns who is on-call right now."
}

func (t *GetCurrentOnCallTool) InputSchema() map[string]interface{} {
	properties := map[string]interface{}{
		"schedule_id": map[string]interface{}{
			"type":        "string",
			"description": "The schedule ID (e.g., '01JXXESWWYB00TVX9FTK8XNZZJ')",
		},
	}
	return t.StandardInputSchema(properties, []string{"schedule_id"})
}

func (t *GetCurrentOnCallTool) Execute(args map[string]interface{}) (string, error) {
	scheduleID, err := t.ValidateRequiredString(args, "schedule_id")
	if err != nil {
		return "", err
	}

	now := time.Now()
	from := now.Format(time.RFC3339)
	to := now.Add(1 * time.Minute).Format(time.RFC3339)

	result, err := t.GetClient().ListScheduleEntries(scheduleID, from, to, 10, "")
	if err != nil {
		return "", fmt.Errorf("failed to get current on-call: %w", err)
	}

	// Find the entry that covers the current time from final entries
	for _, entry := range result.Final {
		if entry.StartAt.Before(now) && entry.EndAt.After(now) {
			response := map[string]interface{}{
				"schedule_id":   scheduleID,
				"current_time":  now.Format(time.RFC3339),
				"on_call_user": map[string]interface{}{
					"id":    entry.User.ID,
					"name":  entry.User.Name,
					"email": entry.User.Email,
				},
				"shift": map[string]interface{}{
					"starts_at": entry.StartAt,
					"ends_at":   entry.EndAt,
					"layer_id":  entry.LayerID,
				},
			}
			return t.FormatResponse(response)
		}
	}

	return "", fmt.Errorf("no on-call person found for schedule %s at current time", scheduleID)
}
