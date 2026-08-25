package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/sahassauarv/tax-annotation/annotation"
	"github.com/sahassauarv/tax-annotation/formatter"
	"github.com/sahassauarv/tax-annotation/parser"
	"github.com/sahassauarv/tax-annotation/validator"
)

type Renderer struct {
	resolver   *parser.PathResolver
	formatter  *formatter.Formatter
	validator  *validator.Validator
}

func NewRenderer(data map[string]interface{}) *Renderer {
	return &Renderer{
		resolver:   parser.NewPathResolver(data),
		formatter:  formatter.New(),
		validator:  validator.New(),
	}
}

type RenderedField struct {
	ID               string
	Label            string
	FieldType        annotation.FieldType
	Position         annotation.Position
	RawValue         interface{}
	FormattedValue   string
	ValidationResult validator.ValidationResult
	IsValid          bool
	HasValue         bool
}

type RenderResult struct {
	Fields     map[string]*RenderedField
	Errors     []string
	FormID     string
	FormName   string
	RenderTime time.Duration
}

func (r *RenderResult) IsComplete() bool {
	for _, field := range r.Fields {
		if !field.IsValid {
			return false
		}
	}
	return true
}

func (r *RenderResult) GetValidFields() map[string]*RenderedField {
	result := make(map[string]*RenderedField)
	for id, field := range r.Fields {
		if field.IsValid {
			result[id] = field
		}
	}
	return result
}

func (r *RenderResult) GetInvalidFields() map[string]*RenderedField {
	result := make(map[string]*RenderedField)
	for id, field := range r.Fields {
		if !field.IsValid {
			result[id] = field
		}
	}
	return result
}

func (r *RenderResult) PrintSummary() {
	valid := 0
	invalid := 0
	empty := 0
	for _, field := range r.Fields {
		if !field.IsValid {
			invalid++
		} else if !field.HasValue {
			empty++
		} else {
			valid++
		}
	}

	fmt.Printf("=== Render Summary: %s (%s) ===\n", r.FormName, r.FormID)
	fmt.Printf("Total Fields: %d\n", len(r.Fields))
	fmt.Printf("  Valid:   %d\n", valid)
	fmt.Printf("  Invalid: %d\n", invalid)
	fmt.Printf("  Empty:   %d\n", empty)
	fmt.Printf("Complete: %v\n", r.IsComplete())
	fmt.Printf("Render Time: %v\n", r.RenderTime)
	if len(r.Errors) > 0 {
		fmt.Printf("Errors: %d\n", len(r.Errors))
		for _, err := range r.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
	fmt.Println()
}

func (r *RenderResult) PrintDetailed() {
	r.PrintSummary()
	for id, field := range r.Fields {
		status := "OK"
		if !field.IsValid {
			status = "FAIL"
		} else if !field.HasValue {
			status = "EMPTY"
		}
		fmt.Printf("[%s] %s (%s)\n", status, field.Label, id)
		fmt.Printf("  Type: %s\n", field.FieldType)
		fmt.Printf("  Position: (%.0f, %.0f) %.0fx%.0f\n",
			field.Position.X, field.Position.Y, field.Position.Width, field.Position.Height)
		if field.HasValue {
			fmt.Printf("  Raw Value:   %v\n", field.RawValue)
			fmt.Printf("  Formatted:   %s\n", field.FormattedValue)
		}
		if !field.IsValid {
			fmt.Printf("  Error: %s\n", field.ValidationResult.Message)
		}
		fmt.Println()
	}
}

func (ren *Renderer) RenderForm(form *annotation.Form) *RenderResult {
	start := time.Now()
	result := &RenderResult{
		Fields:   make(map[string]*RenderedField),
		FormID:   form.ID,
		FormName: form.Name,
	}

	for _, page := range form.Pages {
		pageResult := ren.RenderPage(&page)
		for id, field := range pageResult.Fields {
			result.Fields[id] = field
		}
		result.Errors = append(result.Errors, pageResult.Errors...)
	}

	result.RenderTime = time.Since(start)
	return result
}

func (ren *Renderer) RenderPage(page *annotation.Page) *RenderResult {
	start := time.Now()
	result := &RenderResult{
		Fields: make(map[string]*RenderedField),
	}

	for _, ann := range page.Annotations {
		field := ren.RenderAnnotation(&ann)
		result.Fields[ann.ID] = field
	}

	result.RenderTime = time.Since(start)
	return result
}

func (ren *Renderer) RenderAnnotation(ann *annotation.Annotation) *RenderedField {
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
		return field
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

	return field
}

func (ren *Renderer) RenderFieldByID(form *annotation.Form, fieldID string) *RenderedField {
	ann := parser.GetAnnotation(form, fieldID)
	if ann == nil {
		return nil
	}
	return ren.RenderAnnotation(ann)
}

func (ren *Renderer) FormatAllFields(form *annotation.Form) map[string]string {
	result := make(map[string]string)
	for _, ann := range parser.AllAnnotations(form) {
		field := ren.RenderAnnotation(&ann)
		if field.HasValue && field.IsValid {
			result[ann.ID] = field.FormattedValue
		}
	}
	return result
}

func (ren *Renderer) GetFieldSummary(form *annotation.Form) string {
	result := ren.RenderForm(form)
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

	return sb.String()
}
