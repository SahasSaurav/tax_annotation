package render

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
	"github.com/sahassauarv/tax-annotation/pkg/formatter"
	"github.com/sahassauarv/tax-annotation/pkg/parser"
	"github.com/sahassauarv/tax-annotation/pkg/validator"
)

// formRenderer is the default implementation of the Renderer interface.
// It resolves data paths, validates values, and formats them without
// any output side effects. Use a Writer to present the RenderResult.
type formRenderer struct {
	resolver  parser.PathResolver
	formatter formatter.Formatter
	validator validator.Validator
}

// NewRenderer creates a Renderer with the provided dependencies.
// Accepts interfaces for PathResolver, Formatter, and Validator to allow
// custom implementations and easy testing via mocks.
func NewRenderer(resolver parser.PathResolver, fmtr formatter.Formatter, vld validator.Validator) (Renderer, error) {
	if resolver == nil {
		return nil, fmt.Errorf("resolver cannot be nil")
	}
	if fmtr == nil {
		return nil, fmt.Errorf("formatter cannot be nil")
	}
	if vld == nil {
		return nil, fmt.Errorf("validator cannot be nil")
	}
	return &formRenderer{
		resolver:  resolver,
		formatter: fmtr,
		validator: vld,
	}, nil
}

// RenderedField holds the complete result for a single annotation after rendering.
// It contains the original annotation metadata, the resolved raw value, the formatted
// display string, and the validation outcome so callers can inspect both the output
// and whether it passed all checks.
type RenderedField struct {
	ID               string                    // unique annotation identifier
	Label            string                    // human-readable label for the field
	FieldType        annotation.FieldType      // text, number, date, or checkbox
	Position         annotation.Position       // x, y, width, height on the page
	RawValue         interface{}               // original value from the data map
	FormattedValue   string                    // formatted string ready for display
	ValidationResult validator.ValidationResult // pass/fail with message
	IsValid          bool                      // true if validation passed
	HasValue         bool                      // true if a value was found at the path
}

// RenderResult aggregates the rendered output of an entire form or a single page.
// It provides helper methods to check completeness and filter fields by validity.
type RenderResult struct {
	Fields     map[string]*RenderedField // keyed by annotation ID
	Errors     []string                  // non-fatal errors encountered during render
	FormID     string                    // form identifier (e.g. "W-2")
	FormName   string                    // form display name (e.g. "Wage and Tax Statement")
	RenderTime time.Duration             // total time spent rendering
}

// IsComplete reports whether every field in the result passed validation.
func (r *RenderResult) IsComplete() bool {
	for _, field := range r.Fields {
		if !field.IsValid {
			return false
		}
	}
	return true
}

// GetValidFields returns a subset of fields that passed validation.
func (r *RenderResult) GetValidFields() map[string]*RenderedField {
	result := make(map[string]*RenderedField)
	for id, field := range r.Fields {
		if field.IsValid {
			result[id] = field
		}
	}
	return result
}

// GetInvalidFields returns a subset of fields that failed validation,
// useful for highlighting errors or generating error reports.
func (r *RenderResult) GetInvalidFields() map[string]*RenderedField {
	result := make(map[string]*RenderedField)
	for id, field := range r.Fields {
		if !field.IsValid {
			result[id] = field
		}
	}
	return result
}

// RenderForm processes all pages and annotations in a form, returning a complete
// RenderResult. The context is checked between pages so rendering can be cancelled
// early. Returns an error if the form is nil, has no pages, or if any page fails.
func (ren *formRenderer) RenderForm(ctx context.Context, form *annotation.Form) (*RenderResult, error) {
	if form == nil {
		return nil, fmt.Errorf("form cannot be nil")
	}
	if len(form.Pages) == 0 {
		return nil, fmt.Errorf("form has no pages")
	}

	start := time.Now()
	result := &RenderResult{
		Fields:   make(map[string]*RenderedField),
		FormID:   form.ID,
		FormName: form.Name,
	}

	for _, page := range form.Pages {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("render cancelled: %w", err)
		}

		pageResult, err := ren.RenderPage(ctx, &page)
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %w", page.Number, err)
		}
		for id, field := range pageResult.Fields {
			result.Fields[id] = field
		}
		result.Errors = append(result.Errors, pageResult.Errors...)
	}

	result.RenderTime = time.Since(start)
	return result, nil
}

// RenderPage processes all annotations on a single page. The context is checked
// before each annotation so long pages can be interrupted.
func (ren *formRenderer) RenderPage(ctx context.Context, page *annotation.Page) (*RenderResult, error) {
	if page == nil {
		return nil, fmt.Errorf("page cannot be nil")
	}

	start := time.Now()
	result := &RenderResult{
		Fields: make(map[string]*RenderedField),
	}

	for _, ann := range page.Annotations {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("render cancelled: %w", err)
		}

		field, err := ren.RenderAnnotation(ctx, &ann)
		if err != nil {
			return nil, fmt.Errorf("failed to render annotation %s: %w", ann.ID, err)
		}
		result.Fields[ann.ID] = field
	}

	result.RenderTime = time.Since(start)
	return result, nil
}

// RenderAnnotation resolves the annotation's value path against the data map,
// validates the raw value against the field's type and constraints, and formats
// it into a display string. If the path has no value and the field is required,
// it marks the field as invalid.
func (ren *formRenderer) RenderAnnotation(ctx context.Context, ann *annotation.Annotation) (*RenderedField, error) {
	if ann == nil {
		return nil, fmt.Errorf("annotation cannot be nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("render cancelled: %w", err)
	}

	field := &RenderedField{
		ID:        ann.ID,
		Label:     ann.Label,
		FieldType: ann.FieldType,
		Position:  ann.Position,
	}

	rawValue, hasValue := ren.resolver.Resolve(ann.Value.Path)
	field.HasValue = hasValue

	if !hasValue {
		field.FormattedValue = ""
		field.IsValid = true
		if ann.Validation != nil && ann.Validation.Required {
			field.IsValid = false
			field.ValidationResult = validator.ValidationResult{Valid: false, Message: "required field missing"}
		}
		return field, nil
	}

	field.RawValue = rawValue

	validationResult := ren.validator.Validate(rawValue, ann.FieldType, ann.Validation)
	field.ValidationResult = validationResult
	field.IsValid = validationResult.Valid

	if field.IsValid {
		formatted, err := ren.formatter.Format(rawValue, ann.FieldType, ann.Format)
		if err != nil {
			field.IsValid = false
			field.ValidationResult = validator.ValidationResult{Valid: false, Message: err.Error()}
		} else {
			if ann.Format != nil && ann.Format.Alignment != "" {
				formatted = formatter.AlignText(formatted, int(ann.Position.Width/2), ann.Format.Alignment)
			}
			field.FormattedValue = formatted
		}
	}

	return field, nil
}

// RenderFieldByID looks up a single annotation by its ID within the form and renders it.
func (ren *formRenderer) RenderFieldByID(ctx context.Context, form *annotation.Form, fieldID string) (*RenderedField, error) {
	ann := parser.GetAnnotation(form, fieldID)
	if ann == nil {
		return nil, fmt.Errorf("annotation %s not found", fieldID)
	}
	return ren.RenderAnnotation(ctx, ann)
}

// FormatAllFields renders every annotation and returns a map of field IDs
// to their formatted display strings. Only includes fields with a valid value.
func (ren *formRenderer) FormatAllFields(ctx context.Context, form *annotation.Form) (map[string]string, error) {
	result := make(map[string]string)
	for _, ann := range parser.AllAnnotations(form) {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("format cancelled: %w", err)
		}

		field, err := ren.RenderAnnotation(ctx, &ann)
		if err != nil {
			return nil, err
		}
		if field.HasValue && field.IsValid {
			result[ann.ID] = field.FormattedValue
		}
	}
	return result, nil
}

// GetFieldSummary renders the full form and builds a plain-text summary with fields
// grouped by page, including status indicators (! for invalid, - for empty, space for valid).
func (ren *formRenderer) GetFieldSummary(ctx context.Context, form *annotation.Form) (string, error) {
	result, err := ren.RenderForm(ctx, form)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "Form: %s (%s)\n", result.FormName, result.FormID)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 50))

	for _, page := range form.Pages {
		fmt.Fprintf(&sb, "\nPage %d: %s\n", page.Number, page.Label)
		fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 40))

		for _, ann := range page.Annotations {
			field := result.Fields[ann.ID]
			status := "  "
			if !field.IsValid {
				status = "! "
			} else if !field.HasValue {
				status = "- "
			} else {
				status = "  "
			}

			value := field.FormattedValue
			if value == "" {
				value = "(empty)"
			}

			sb.WriteString(fmt.Sprintf("%s%-30s %s\n", status, field.Label, value))
		}
	}

	return sb.String(), nil
}
