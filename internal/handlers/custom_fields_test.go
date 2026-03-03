package handlers

import (
	"testing"
)

func TestListCustomFieldsTool_Schema(t *testing.T) {
	tool := &ListCustomFieldsTool{}

	// Test Name
	if tool.Name() != "list_custom_fields" {
		t.Errorf("Expected name 'list_custom_fields', got %s", tool.Name())
	}

	// Test Description
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}

	// Test InputSchema
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}
}

func TestGetCustomFieldTool_Execute(t *testing.T) {
	tool := &GetCustomFieldTool{}

	// Test missing required id parameter
	t.Run("missing required id", func(t *testing.T) {
		args := map[string]interface{}{}

		_, err := tool.Execute(args)
		if err == nil {
			t.Error("Expected error for missing id parameter")
		}
		if err.Error() != "id parameter is required" {
			t.Errorf("Expected 'id parameter is required' error, got: %v", err)
		}
	})

	// Test id parameter with wrong type
	t.Run("id parameter wrong type", func(t *testing.T) {
		args := map[string]interface{}{
			"id": 123, // Not a string
		}

		_, err := tool.Execute(args)
		if err == nil {
			t.Error("Expected error for wrong type id parameter")
		}
	})

	// Test empty id parameter
	t.Run("empty id parameter", func(t *testing.T) {
		args := map[string]interface{}{
			"id": "",
		}

		_, err := tool.Execute(args)
		if err == nil {
			t.Error("Expected error for empty id parameter")
		}
	})
}

func TestGetCustomFieldTool_Schema(t *testing.T) {
	tool := &GetCustomFieldTool{}

	// Test Name
	if tool.Name() != "get_custom_field" {
		t.Errorf("Expected name 'get_custom_field', got %s", tool.Name())
	}

	// Test Description
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}

	// Test InputSchema
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	properties := schema["properties"].(map[string]interface{})
	if _, ok := properties["id"]; !ok {
		t.Error("Schema should have 'id' property")
	}

	required := schema["required"].([]string)
	if len(required) != 1 || required[0] != "id" {
		t.Error("Schema should require only 'id'")
	}
}

func TestSearchCustomFieldsTool_Execute(t *testing.T) {
	// Note: We can't test full execution without a real client,
	// but SearchCustomFieldsTool doesn't validate required parameters
	// since both query and field_type are optional
}

func TestSearchCustomFieldsTool_Schema(t *testing.T) {
	tool := &SearchCustomFieldsTool{}

	// Test Name
	if tool.Name() != "search_custom_fields" {
		t.Errorf("Expected name 'search_custom_fields', got %s", tool.Name())
	}

	// Test Description
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}

	// Test InputSchema
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	properties := schema["properties"].(map[string]interface{})
	if _, ok := properties["query"]; !ok {
		t.Error("Schema should have 'query' property")
	}
	if _, ok := properties["field_type"]; !ok {
		t.Error("Schema should have 'field_type' property")
	}
}

func TestListCustomFieldOptionsTool_Schema(t *testing.T) {
	tool := &ListCustomFieldOptionsTool{}

	// Test Name
	if tool.Name() != "list_custom_field_options" {
		t.Errorf("Expected name 'list_custom_field_options', got %s", tool.Name())
	}

	// Test Description
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}

	// Test InputSchema
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}
}
