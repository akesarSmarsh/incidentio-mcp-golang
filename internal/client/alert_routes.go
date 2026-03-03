package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListAlertRoutesParams contains optional parameters for listing alert routes
type ListAlertRoutesParams struct {
	PageSize int
	After    string
}

// ListAlertRoutesResponse represents the response from listing alert routes
type ListAlertRoutesResponse struct {
	AlertRoutes []AlertRoute `json:"alert_routes"`
	Pagination  struct {
		After    string `json:"after,omitempty"`
		PageSize int    `json:"page_size"`
	} `json:"pagination_info"`
}

// ListAlertRoutes returns all alert routes
func (c *Client) ListAlertRoutes(params *ListAlertRoutesParams) (*ListAlertRoutesResponse, error) {
	endpoint := "/alert_routes"

	// Set default page size
	pageSize := 25
	if params != nil && params.PageSize > 0 {
		pageSize = params.PageSize
	}

	v := url.Values{}
	v.Set("page_size", fmt.Sprintf("%d", pageSize)) // Always set (may be required)
	if params != nil && params.After != "" {
		v.Set("after", params.After)
	}

	if len(v) > 0 {
		endpoint = endpoint + "?" + v.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var result ListAlertRoutesResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetAlertRoute returns a specific alert route by ID
func (c *Client) GetAlertRoute(id string) (*AlertRoute, error) {
	endpoint := fmt.Sprintf("/alert_routes/%s", id)

	respBody, err := c.doRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		AlertRoute AlertRoute `json:"alert_route"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result.AlertRoute, nil
}
