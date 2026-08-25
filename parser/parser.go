package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sahassauarv/tax-annotation/annotation"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

//
func (p *Parser) ParseFormFromFile(path string) (*annotation.Form, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read form file %s: %w", path, err)
	}
	return p.ParseForm(data)
}

func (p *Parser) ParseForm(data []byte) (*annotation.Form, error) {
	var form annotation.Form
	if err := json.Unmarshal(data, &form); err != nil {
		return nil, fmt.Errorf("unmarshal form: %w", err)
	}
	if err := ValidateForm(&form); err != nil {
		return nil, fmt.Errorf("validate form: %w", err)
	}
	return &form, nil
}

func (p *Parser) LoadDataFromFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read data file %s: %w", path, err)
	}
	return p.LoadData(data)
}

func (p *Parser) LoadData(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal data: %w", err)
	}
	return result, nil
}

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

func AllAnnotations(f *annotation.Form) []annotation.Annotation {
	var all []annotation.Annotation
	for _, page := range f.Pages {
		all = append(all, page.Annotations...)
	}
	return all
}

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
