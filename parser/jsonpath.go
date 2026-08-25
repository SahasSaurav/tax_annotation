package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// pathResolver is the default implementation of the PathResolver interface.
// It walks a nested map using dot-notation paths like "employee.ssn".
type pathResolver struct {
	data map[string]interface{}
}

// NewPathResolver creates a PathResolver backed by the given data map.
func NewPathResolver(data map[string]interface{}) PathResolver {
	return &pathResolver{data: data}
}

// Resolve walks the data map using the given dot-notation path and returns
// the value at that location, or false if any part of the path is missing.
// Supports array access via bracket notation (e.g. "items[0].name").
func (r *pathResolver) Resolve(path string) (interface{}, bool) {
	if r.data == nil {
		return nil, false
	}
	parts := splitPath(path)
	var current interface{} = r.data

	for _, part := range parts {
		idx, arrayAccess := parseArrayIndex(part)
		if arrayAccess {
			arr, ok := current.([]interface{})
			if !ok {
				return nil, false
			}
			if idx < 0 || idx >= len(arr) {
				return nil, false
			}
			current = arr[idx]
		} else {
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}
			val, exists := obj[part]
			if !exists {
				return nil, false
			}
			current = val
		}
	}
	return current, true
}

// GetString resolves a path and returns the value formatted as a string.
// Returns fallback if the path cannot be resolved.
func (r *pathResolver) GetString(path, fallback string) string {
	val, ok := r.Resolve(path)
	if !ok {
		return fallback
	}
	return fmt.Sprintf("%v", val)
}

// GetFloat resolves a path and returns the value as a float64.
// Handles json.Number, string, and native numeric types. Returns fallback on failure.
func (r *pathResolver) GetFloat(path string, fallback float64) float64 {
	val, ok := r.Resolve(path)
	if !ok {
		return fallback
	}
	switch v := val.(type) {
	case float64:
		return v
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return fallback
		}
		return f
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fallback
		}
		return f
	default:
		return fallback
	}
}

// GetBool resolves a path and returns the value as a bool.
// Handles string representations ("true", "1", etc.). Returns fallback on failure.
func (r *pathResolver) GetBool(path string, fallback bool) bool {
	val, ok := r.Resolve(path)
	if !ok {
		return fallback
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fallback
		}
		return b
	default:
		return fallback
	}
}

// splitPath breaks a dot-notation path into its individual segments.
func splitPath(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	return strings.Split(path, ".")
}

// parseArrayIndex extracts an integer index from bracket notation (e.g. "items[0]").
// Returns the index and true if brackets were found, or -1 and false otherwise.
func parseArrayIndex(part string) (int, bool) {
	bracketStart := strings.Index(part, "[")
	bracketEnd := strings.Index(part, "]")
	if bracketStart == -1 || bracketEnd == -1 || bracketEnd <= bracketStart+1 {
		return -1, false
	}
	idxStr := part[bracketStart+1 : bracketEnd]
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return -1, false
	}
	return idx, true
}
