package parser

import (
	"encoding/json"
	"testing"
)

func TestNewPathResolver(t *testing.T) {
	data := map[string]interface{}{"a": "1"}
	r := NewPathResolver(data)
	if r == nil {
		t.Fatal("NewPathResolver returned nil")
	}
}

func TestResolve(t *testing.T) {
	data := map[string]interface{}{
		"employee": map[string]interface{}{
			"name": "John",
			"ssn":  "123-45-6789",
			"age":  30,
		},
		"items": []interface{}{
			map[string]interface{}{"id": 1},
			map[string]interface{}{"id": 2},
		},
	}
	r := NewPathResolver(data)

	tests := []struct {
		name     string
		path     string
		wantVal  interface{}
		wantOK   bool
		resolver PathResolver
	}{
		{"simple path", "employee.name", "John", true, r},
		{"nested numeric", "employee.age", 30, true, r},
		{"missing path", "employee.nonexistent", nil, false, r},
		{"non-map in path", "employee.name.invalid", nil, false, r},
		{"nil data", "anything", nil, false, NewPathResolver(nil)},
		{"empty path returns full data", "", data, true, r},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.resolver.Resolve(tt.path)
			if ok != tt.wantOK {
				t.Errorf("ok: got %v, want %v", ok, tt.wantOK)
			}
			if tt.wantOK && val == nil {
				t.Error("expected non-nil value")
			}
		})
	}
}

func TestGetString(t *testing.T) {
	data := map[string]interface{}{"name": "John"}
	r := NewPathResolver(data)

	tests := []struct {
		name     string
		path     string
		fallback string
		expected string
	}{
		{"found", "name", "default", "John"},
		{"fallback", "missing", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.GetString(tt.path, tt.fallback)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		path     string
		fallback float64
		expected float64
	}{
		{"float64", map[string]interface{}{"v": 75000.00}, "v", 0, 75000.00},
		{"string to float", map[string]interface{}{"v": "123.45"}, "v", 0, 123.45},
		{"json.Number", map[string]interface{}{"v": json.Number("999.99")}, "v", 0, 999.99},
		{"bad string", map[string]interface{}{"v": "abc"}, "v", 42, 42},
		{"bool fallback", map[string]interface{}{"v": true}, "v", 42, 42},
		{"missing", map[string]interface{}{}, "v", 42, 42},
		{"bad json.Number", map[string]interface{}{"v": json.Number("notanumber")}, "v", 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewPathResolver(tt.data)
			result := r.GetFloat(tt.path, tt.fallback)
			if result != tt.expected {
				t.Errorf("got %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		path     string
		fallback bool
		expected bool
	}{
		{"true", map[string]interface{}{"v": true}, "v", false, true},
		{"false", map[string]interface{}{"v": false}, "v", true, false},
		{"string true", map[string]interface{}{"v": "true"}, "v", false, true},
		{"string 1", map[string]interface{}{"v": "1"}, "v", false, true},
		{"bad string fallback", map[string]interface{}{"v": "invalid"}, "v", true, true},
		{"non-bool fallback", map[string]interface{}{"v": 42}, "v", true, true},
		{"missing fallback", map[string]interface{}{}, "v", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewPathResolver(tt.data)
			result := r.GetBool(tt.path, tt.fallback)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"a.b.c", 3},
		{"single", 1},
		{"", 0},
		{"  ", 0},
	}
	for _, tt := range tests {
		result := splitPath(tt.input)
		if len(result) != tt.expected {
			t.Errorf("splitPath(%q): got %d parts, want %d", tt.input, len(result), tt.expected)
		}
	}
}

func TestParseArrayIndex(t *testing.T) {
	tests := []struct {
		input   string
		idx     int
		isArray bool
	}{
		{"items[0]", 0, true},
		{"items[5]", 5, true},
		{"items[abc]", -1, false},
		{"items[]", -1, false},
		{"items", -1, false},
		{"items[", -1, false},
	}
	for _, tt := range tests {
		idx, isArray := parseArrayIndex(tt.input)
		if idx != tt.idx || isArray != tt.isArray {
			t.Errorf("parseArrayIndex(%q): got (%d, %v), want (%d, %v)", tt.input, idx, isArray, tt.idx, tt.isArray)
		}
	}
}
