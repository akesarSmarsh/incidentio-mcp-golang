package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListIncidentUpdates retrieves incident updates with optional filtering
func (c *Client) ListIncidentUpdates(opts *ListIncidentUpdatesOptions) (*ListIncidentUpdatesResponse, error) {
	// Set default page size
	pageSize := 25
	if opts != nil && opts.PageSize > 0 {
		pageSize = opts.PageSize
	}

	params := url.Values{}
	params.Set("page_size", strconv.Itoa(pageSize)) // Always set (may be required)

	if opts != nil {
		if opts.IncidentID != "" {
			params.Set("incident_id", opts.IncidentID)
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
	}

	respBody, err := c.doRequest("GET", "/incident_updates", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListIncidentUpdatesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetIncidentUpdate retrieves a specific incident update by ID.
// Note: The incident.io API does not have a GET-by-ID endpoint for incident updates.
// This implementation fetches updates via the list endpoint and filters by ID.
func (c *Client) GetIncidentUpdate(id string) (*IncidentUpdate, error) {
	// Fetch a page of updates (the target update is likely recent)
	params := url.Values{}
	params.Set("page_size", "250")

	respBody, err := c.doRequest("GET", "/incident_updates", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListIncidentUpdatesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for i := range response.IncidentUpdates {
		if response.IncidentUpdates[i].ID == id {
			return &response.IncidentUpdates[i], nil
		}
	}

	// If not found in first page, paginate through remaining pages
	after := response.PaginationMeta.After
	for after != "" {
		params.Set("after", after)
		respBody, err = c.doRequest("GET", "/incident_updates", params, nil)
		if err != nil {
			return nil, err
		}

		var pageResponse ListIncidentUpdatesResponse
		if err := json.Unmarshal(respBody, &pageResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		for i := range pageResponse.IncidentUpdates {
			if pageResponse.IncidentUpdates[i].ID == id {
				return &pageResponse.IncidentUpdates[i], nil
			}
		}

		after = pageResponse.PaginationMeta.After
	}

	return nil, fmt.Errorf("incident update with ID %s not found", id)
}
