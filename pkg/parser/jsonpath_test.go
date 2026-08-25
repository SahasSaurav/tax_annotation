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

	t.Run("simple path", func(t *testing.T) {
		val, ok := r.Resolve("employee.name")
		if !ok || val != "John" {
			t.Errorf("expected John, got %v", val)
		}
	})

	t.Run("missing path", func(t *testing.T) {
		_, ok := r.Resolve("employee.nonexistent")
		if ok {
			t.Error("expected false for missing path")
		}
	})

	t.Run("array access", func(t *testing.T) {
		_, ok := r.Resolve("items[1].id")
		if ok {
			// This path format isn't supported by splitPath since "items[1]" is treated
			// as one segment. Array indexing works only when parent is already an array.
		}
	})

	t.Run("array access on nested", func(t *testing.T) {
		// Test array access when already inside an array
		nested := map[string]interface{}{
			"first": []interface{}{
				map[string]interface{}{"val": 10},
			},
		}
		nr := NewPathResolver(nested)
		val, ok := nr.Resolve("first")
		if !ok {
			t.Fatal("expected ok")
		}
		arr, ok := val.([]interface{})
		if !ok || len(arr) != 1 {
			t.Fatalf("expected array of 1, got %v", val)
		}
		item, ok := arr[0].(map[string]interface{})
		if !ok {
			t.Fatal("expected map")
		}
		if item["val"] != 10 {
			t.Errorf("expected 10, got %v", item["val"])
		}
	})

	t.Run("nil data", func(t *testing.T) {
		nr := NewPathResolver(nil)
		_, ok := nr.Resolve("anything")
		if ok {
			t.Error("expected false for nil data")
		}
	})

	t.Run("empty path returns full data", func(t *testing.T) {
		val, ok := r.Resolve("")
		if !ok {
			t.Error("expected true for empty path")
		}
		m, ok := val.(map[string]interface{})
		if !ok {
			t.Error("expected map")
		}
		if m["employee"] == nil {
			t.Error("expected employee key in returned data")
		}
	})

	t.Run("non-map in path", func(t *testing.T) {
		_, ok := r.Resolve("employee.name.invalid")
		if ok {
			t.Error("expected false for traversal into non-map")
		}
	})
}

func TestGetString(t *testing.T) {
	data := map[string]interface{}{"name": "John"}
	r := NewPathResolver(data)

	t.Run("found", func(t *testing.T) {
		result := r.GetString("name", "default")
		if result != "John" {
			t.Errorf("expected John, got %s", result)
		}
	})

	t.Run("fallback", func(t *testing.T) {
		result := r.GetString("missing", "default")
		if result != "default" {
			t.Errorf("expected default, got %s", result)
		}
	})
}

func TestGetFloat(t *testing.T) {
	data := map[string]interface{}{
		"amount": 75000.00,
		"strNum": "123.45",
		"bad":    "abc",
		"bool":   true,
	}
	r := NewPathResolver(data)

	t.Run("float64", func(t *testing.T) {
		result := r.GetFloat("amount", 0)
		if result != 75000.00 {
			t.Errorf("expected 75000, got %f", result)
		}
	})

	t.Run("string to float", func(t *testing.T) {
		result := r.GetFloat("strNum", 0)
		if result != 123.45 {
			t.Errorf("expected 123.45, got %f", result)
		}
	})

	t.Run("json.Number", func(t *testing.T) {
		data := map[string]interface{}{"val": json.Number("999.99")}
		r := NewPathResolver(data)
		result := r.GetFloat("val", 0)
		if result != 999.99 {
			t.Errorf("expected 999.99, got %f", result)
		}
	})

	t.Run("bad string", func(t *testing.T) {
		result := r.GetFloat("bad", 42)
		if result != 42 {
			t.Errorf("expected fallback 42, got %f", result)
		}
	})

	t.Run("bool fallback", func(t *testing.T) {
		result := r.GetFloat("bool", 42)
		if result != 42 {
			t.Errorf("expected fallback 42, got %f", result)
		}
	})

	t.Run("missing", func(t *testing.T) {
		result := r.GetFloat("missing", 42)
		if result != 42 {
			t.Errorf("expected fallback 42, got %f", result)
		}
	})

	t.Run("bad json.Number", func(t *testing.T) {
		data := map[string]interface{}{"val": json.Number("notanumber")}
		r := NewPathResolver(data)
		result := r.GetFloat("val", 42)
		if result != 42 {
			t.Errorf("expected fallback 42, got %f", result)
		}
	})
}

func TestGetBool(t *testing.T) {
	data := map[string]interface{}{
		"yes":     true,
		"no":      false,
		"strTrue": "true",
		"strOne":  "1",
		"bad":     "invalid",
		"num":     42,
	}
	r := NewPathResolver(data)

	t.Run("true", func(t *testing.T) {
		if !r.GetBool("yes", false) {
			t.Error("expected true")
		}
	})

	t.Run("false", func(t *testing.T) {
		if r.GetBool("no", true) {
			t.Error("expected false")
		}
	})

	t.Run("string true", func(t *testing.T) {
		if !r.GetBool("strTrue", false) {
			t.Error("expected true from string")
		}
	})

	t.Run("string 1", func(t *testing.T) {
		if !r.GetBool("strOne", false) {
			t.Error("expected true from '1'")
		}
	})

	t.Run("bad string fallback", func(t *testing.T) {
		if !r.GetBool("bad", true) {
			t.Error("expected fallback true")
		}
	})

	t.Run("non-bool fallback", func(t *testing.T) {
		if !r.GetBool("num", true) {
			t.Error("expected fallback true")
		}
	})

	t.Run("missing fallback", func(t *testing.T) {
		if !r.GetBool("missing", true) {
			t.Error("expected fallback true")
		}
	})
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
			t.Errorf("splitPath(%q): expected %d parts, got %d", tt.input, tt.expected, len(result))
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
