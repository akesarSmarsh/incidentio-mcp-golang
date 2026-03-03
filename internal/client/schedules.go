package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListSchedules retrieves all on-call schedules
func (c *Client) ListSchedules(pageSize int, after string) (*SchedulesResponse, error) {
	params := url.Values{}
	if pageSize > 0 {
		params.Set("page_size", fmt.Sprintf("%d", pageSize))
	}
	if after != "" {
		params.Set("after", after)
	}

	respBody, err := c.doRequest("GET", "/schedules", params, nil)
	if err != nil {
		return nil, err
	}

	var response SchedulesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedules response: %w", err)
	}

	return &response, nil
}

// GetSchedule retrieves a specific schedule by ID
func (c *Client) GetSchedule(scheduleID string) (*Schedule, error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("schedule_id is required")
	}

	path := fmt.Sprintf("/schedules/%s", scheduleID)
	respBody, err := c.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Schedule Schedule `json:"schedule"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	return &response.Schedule, nil
}

// ListScheduleEntries retrieves entries (on-call assignments) for a schedule
// The API endpoint is /v2/schedule_entries with schedule_id as a query parameter
func (c *Client) ListScheduleEntries(scheduleID string, from string, to string, pageSize int, after string) (*ScheduleEntriesResponse, error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("schedule_id is required")
	}

	params := url.Values{}
	params.Set("schedule_id", scheduleID) // schedule_id is a required query param
	if from != "" {
		params.Set("from", from)
	}
	if to != "" {
		params.Set("to", to)
	}
	if pageSize > 0 {
		params.Set("page_size", fmt.Sprintf("%d", pageSize))
	}
	if after != "" {
		params.Set("after", after)
	}

	// API endpoint is /schedule_entries (not /schedules/{id}/entries)
	respBody, err := c.doRequest("GET", "/schedule_entries", params, nil)
	if err != nil {
		return nil, err
	}

	var response ScheduleEntriesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule entries response: %w", err)
	}

	return &response, nil
}
