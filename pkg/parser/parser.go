package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// formParser is the default implementation of the Parser interface.
// It reads form definitions and data from disk or byte slices,
// unmarshals them, and runs structural validation.
type formParser struct{}

// New creates a new Parser that satisfies the Parser interface.
func New() Parser {
	return &formParser{}
}

// ParseFormFromFile reads a form definition from the given file path and returns a validated Form.
func (p *formParser) ParseFormFromFile(ctx context.Context, path string) (*annotation.Form, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parse cancelled: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read form file %s: %w", path, err)
	}
	return p.ParseForm(ctx, data)
}

// ParseForm unmarshals JSON data into a Form and validates it.
func (p *formParser) ParseForm(ctx context.Context, data []byte) (*annotation.Form, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parse cancelled: %w", err)
	}

	var form annotation.Form
	if err := json.Unmarshal(data, &form); err != nil {
		return nil, fmt.Errorf("unmarshal form: %w", err)
	}
	if err := ValidateForm(&form); err != nil {
		return nil, fmt.Errorf("validate form: %w", err)
	}
	return &form, nil
}

// LoadDataFromFile reads a JSON data file and returns it as a map.
func (p *formParser) LoadDataFromFile(ctx context.Context, path string) (map[string]interface{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load cancelled: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read data file %s: %w", path, err)
	}
	return p.LoadData(ctx, data)
}

// LoadData unmarshals JSON bytes into a map.
func (p *formParser) LoadData(ctx context.Context, data []byte) (map[string]interface{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load cancelled: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal data: %w", err)
	}
	return result, nil
}

// ValidateForm checks that a form has required fields and no duplicate annotation IDs.
func ValidateForm(f *annotation.Form) error {
	if strings.TrimSpace(f.ID) == "" {
		return fmt.Errorf("form ID is required")
	}
	if strings.TrimSpace(f.Name) == "" {
		return fmt.Errorf("form name is required")
	}
	if len(f.Pages) == 0 {
		return fmt.Errorf("form must have at least one page")
	}
	for i, page := range f.Pages {
		if page.Number <= 0 {
			return fmt.Errorf("page %d: number must be positive", i)
		}
		ids := make(map[string]bool)
		for j, ann := range page.Annotations {
			if strings.TrimSpace(ann.ID) == "" {
				return fmt.Errorf("page %d, annotation %d: ID is required", i, j)
			}
			if ids[ann.ID] {
				return fmt.Errorf("page %d: duplicate annotation ID %q", i, ann.ID)
			}
			ids[ann.ID] = true
		}
	}
	return nil
}

// GetAnnotation finds an annotation by ID within a form. Returns nil if not found.
func GetAnnotation(f *annotation.Form, id string) *annotation.Annotation {
	for i := range f.Pages {
		for j := range f.Pages[i].Annotations {
			if f.Pages[i].Annotations[j].ID == id {
				return &f.Pages[i].Annotations[j]
			}
		}
	}
	return nil
}

// AllAnnotations returns all annotations from all pages in a form.
func AllAnnotations(f *annotation.Form) []annotation.Annotation {
	var all []annotation.Annotation
	for _, page := range f.Pages {
		all = append(all, page.Annotations...)
	}
	return all
}

// OverlayData merges overlay values into base, recursively merging nested maps.
func OverlayData(base, overlay map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		if baseVal, exists := result[k]; exists {
			if baseMap, ok1 := baseVal.(map[string]interface{}); ok1 {
				if overlayMap, ok2 := v.(map[string]interface{}); ok2 {
					result[k] = OverlayData(baseMap, overlayMap)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}
