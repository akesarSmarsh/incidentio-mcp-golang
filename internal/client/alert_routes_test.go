package client

import (
	"net/http"
	"testing"
)

func TestListAlertRoutes(t *testing.T) {
	tests := []struct {
		name           string
		params         *ListAlertRoutesParams
		mockResponse   string
		mockStatusCode int
		wantError      bool
		expectedCount  int
	}{
		{
			name:   "successful list alert routes",
			params: &ListAlertRoutesParams{PageSize: 10},
			mockResponse: `{
				"alert_routes": [
					{
						"id": "ar_123",
						"name": "Critical Alerts Route",
						"enabled": true,
						"conditions": [
							{
								"field": "severity",
								"operation": "equals",
								"value": "critical"
							}
						],
						"escalations": [
							{
								"id": "esc_123",
								"level": 1
							}
						],
						"grouping_keys": ["service"],
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				],
				"pagination_info": {
					"page_size": 10
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
			expectedCount:  1,
		},
		{
			name:           "empty alert routes",
			params:         nil,
			mockResponse:   `{"alert_routes": [], "pagination_info": {"page_size": 25}}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					assertEqual(t, "GET", req.Method)
					assertEqual(t, "Bearer test-api-key", req.Header.Get("Authorization"))

					if tt.params != nil && tt.params.PageSize > 0 {
						assertEqual(t, "10", req.URL.Query().Get("page_size"))
					}

					return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
				},
			}

			client := NewTestClient(mockClient)
			result, err := client.ListAlertRoutes(tt.params)

			if tt.wantError {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			if len(result.AlertRoutes) != tt.expectedCount {
				t.Errorf("expected %d alert routes, got %d", tt.expectedCount, len(result.AlertRoutes))
			}

			if tt.expectedCount > 0 {
				route := result.AlertRoutes[0]
				assertEqual(t, "ar_123", route.ID)
				assertEqual(t, "Critical Alerts Route", route.Name)

				// Verify conditions
				if len(route.Conditions) != 1 {
					t.Errorf("expected 1 condition, got %d", len(route.Conditions))
				} else {
					assertEqual(t, "severity", route.Conditions[0].Field)
					assertEqual(t, "equals", route.Conditions[0].Operation)
					assertEqual(t, "critical", route.Conditions[0].Value)
				}

				// Verify escalations
				if len(route.Escalations) != 1 {
					t.Errorf("expected 1 escalation, got %d", len(route.Escalations))
				} else {
					assertEqual(t, "esc_123", route.Escalations[0].ID)
					if route.Escalations[0].Level != 1 {
						t.Errorf("expected level 1, got %d", route.Escalations[0].Level)
					}
				}
			}
		})
	}
}

func TestGetAlertRoute(t *testing.T) {
	tests := []struct {
		name           string
		routeID        string
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:    "successful get alert route",
			routeID: "ar_123",
			mockResponse: `{
				"alert_route": {
					"id": "ar_123",
					"name": "Critical Alerts Route",
					"enabled": true,
					"conditions": [
						{
							"field": "severity",
							"operation": "equals",
							"value": "critical"
						}
					],
					"escalations": [
						{
							"id": "esc_123",
							"level": 1
						}
					],
					"template": {
						"incident_type": "critical"
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z"
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
		{
			name:           "alert route not found",
			routeID:        "ar_nonexistent",
			mockResponse:   `{"error": "Alert route not found"}`,
			mockStatusCode: http.StatusNotFound,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					assertEqual(t, "GET", req.Method)
					assertEqual(t, "/alert_routes/"+tt.routeID, req.URL.Path)
					return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
				},
			}

			client := NewTestClient(mockClient)
			route, err := client.GetAlertRoute(tt.routeID)

			if tt.wantError {
				assertError(t, err)
				return
			}

			assertNoError(t, err)
			assertEqual(t, tt.routeID, route.ID)
			assertEqual(t, "Critical Alerts Route", route.Name)

			// Verify template
			if route.Template != nil {
				if route.Template["incident_type"] != "critical" {
					t.Error("expected template to contain incident_type:critical")
				}
			}
		})
	}
}
