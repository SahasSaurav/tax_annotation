package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type PathResolver struct {
	data map[string]interface{}
}

func NewPathResolver(data map[string]interface{}) *PathResolver {
	return &PathResolver{data: data}
}

func (r *PathResolver) Resolve(path string) (interface{}, bool) {
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

func (r *PathResolver) GetString(path, fallback string) string {
	val, ok := r.Resolve(path)
	if !ok {
		return fallback
	}
	return fmt.Sprintf("%v", val)
}

func (r *PathResolver) GetFloat(path string, fallback float64) float64 {
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

func (r *PathResolver) GetBool(path string, fallback bool) bool {
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

func splitPath(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	return strings.Split(path, ".")
}

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
