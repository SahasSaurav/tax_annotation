package render

import (
	"context"
	"testing"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
	"github.com/sahassauarv/tax-annotation/pkg/formatter"
	"github.com/sahassauarv/tax-annotation/pkg/parser"
	"github.com/sahassauarv/tax-annotation/pkg/validator"
)

func TestNewRenderer(t *testing.T) {
	resolver := parser.NewPathResolver(map[string]interface{}{})
	fmtr := formatter.New()
	vld := validator.New()

	t.Run("valid", func(t *testing.T) {
		r, err := NewRenderer(resolver, fmtr, vld)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r == nil {
			t.Fatal("expected non-nil renderer")
		}
	})

	t.Run("nil resolver", func(t *testing.T) {
		_, err := NewRenderer(nil, fmtr, vld)
		if err == nil {
			t.Fatal("expected error for nil resolver")
		}
	})

	t.Run("nil formatter", func(t *testing.T) {
		_, err := NewRenderer(resolver, nil, vld)
		if err == nil {
			t.Fatal("expected error for nil formatter")
		}
	})

	t.Run("nil validator", func(t *testing.T) {
		_, err := NewRenderer(resolver, fmtr, nil)
		if err == nil {
			t.Fatal("expected error for nil validator")
		}
	})
}

func TestRenderForm(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{
		"employee": map[string]interface{}{
			"name": "John Doe",
			"ssn":  "123-45-6789",
		},
	}

	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID:   "W-2",
		Name: "Wage and Tax Statement",
		Pages: []annotation.Page{
			{
				Number: 1,
				Label:  "Page 1",
				Annotations: []annotation.Annotation{
					{
						ID:        "emp_name",
						Label:     "Employee Name",
						FieldType: annotation.FieldTypeText,
						Value:     annotation.ValueRef{Path: "employee.name"},
						Position:  annotation.Position{X: 72, Y: 200, Width: 200, Height: 12},
					},
					{
						ID:        "emp_ssn",
						Label:     "Employee SSN",
						FieldType: annotation.FieldTypeText,
						Value:     annotation.ValueRef{Path: "employee.ssn"},
						Position:  annotation.Position{X: 72, Y: 220, Width: 200, Height: 12},
						Format:    &annotation.Format{Type: annotation.FormatSSN},
					},
				},
			},
		},
	}

	t.Run("valid form", func(t *testing.T) {
		result, err := r.RenderForm(ctx, form)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.FormID != "W-2" {
			t.Errorf("expected form ID W-2, got %s", result.FormID)
		}
		if len(result.Fields) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(result.Fields))
		}
		if result.Fields["emp_name"].FormattedValue != "John Doe" {
			t.Errorf("expected John Doe, got %s", result.Fields["emp_name"].FormattedValue)
		}
		if result.Fields["emp_ssn"].FormattedValue != "123-45-6789" {
			t.Errorf("expected 123-45-6789, got %s", result.Fields["emp_ssn"].FormattedValue)
		}
	})

	t.Run("nil form", func(t *testing.T) {
		_, err := r.RenderForm(ctx, nil)
		if err == nil {
			t.Fatal("expected error for nil form")
		}
	})

	t.Run("empty pages", func(t *testing.T) {
		_, err := r.RenderForm(ctx, &annotation.Form{ID: "X", Name: "Y"})
		if err == nil {
			t.Fatal("expected error for empty pages")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := r.RenderForm(cancelledCtx, form)
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})
}

func TestRenderAnnotation(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{
		"amount": 75000.00,
	}

	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	t.Run("with value", func(t *testing.T) {
		ann := &annotation.Annotation{
			ID:        "wages",
			Label:     "Wages",
			FieldType: annotation.FieldTypeNumber,
			Value:     annotation.ValueRef{Path: "amount"},
			Position:  annotation.Position{X: 72, Y: 200, Width: 200, Height: 12},
			Format:    &annotation.Format{Type: annotation.FormatCurrency},
		}

		field, err := r.RenderAnnotation(ctx, ann)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !field.HasValue {
			t.Error("expected HasValue to be true")
		}
		if field.FormattedValue != "$75,000.00" {
			t.Errorf("expected $75,000.00, got %s", field.FormattedValue)
		}
	})

	t.Run("missing value", func(t *testing.T) {
		ann := &annotation.Annotation{
			ID:        "missing",
			Label:     "Missing",
			FieldType: annotation.FieldTypeText,
			Value:     annotation.ValueRef{Path: "nonexistent"},
			Position:  annotation.Position{X: 72, Y: 200, Width: 200, Height: 12},
		}

		field, err := r.RenderAnnotation(ctx, ann)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if field.HasValue {
			t.Error("expected HasValue to be false")
		}
		if !field.IsValid {
			t.Error("expected IsValid to be true for non-required missing field")
		}
	})

	t.Run("required missing value", func(t *testing.T) {
		ann := &annotation.Annotation{
			ID:        "required",
			Label:     "Required",
			FieldType: annotation.FieldTypeText,
			Value:     annotation.ValueRef{Path: "nonexistent"},
			Position:  annotation.Position{X: 72, Y: 200, Width: 200, Height: 12},
			Validation: &annotation.Validation{Required: true},
		}

		field, err := r.RenderAnnotation(ctx, ann)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if field.IsValid {
			t.Error("expected IsValid to be false for required missing field")
		}
	})

	t.Run("nil annotation", func(t *testing.T) {
		_, err := r.RenderAnnotation(ctx, nil)
		if err == nil {
			t.Fatal("expected error for nil annotation")
		}
	})
}

func TestRenderResult(t *testing.T) {
	result := &RenderResult{
		Fields: map[string]*RenderedField{
			"valid":   {IsValid: true, HasValue: true},
			"invalid": {IsValid: false},
			"empty":   {IsValid: true, HasValue: false},
		},
	}

	t.Run("IsComplete", func(t *testing.T) {
		if result.IsComplete() {
			t.Error("expected not complete due to invalid field")
		}
	})

	t.Run("GetValidFields", func(t *testing.T) {
		valid := result.GetValidFields()
		if len(valid) != 2 {
			t.Errorf("expected 2 valid fields, got %d", len(valid))
		}
	})

	t.Run("GetInvalidFields", func(t *testing.T) {
		invalid := result.GetInvalidFields()
		if len(invalid) != 1 {
			t.Errorf("expected 1 invalid field, got %d", len(invalid))
		}
	})
}

func TestFormatAllFields(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID:   "TEST",
		Name: "Test",
		Pages: []annotation.Page{
			{
				Number: 1,
				Annotations: []annotation.Annotation{
					{ID: "name", Label: "Name", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "name"}},
					{ID: "age", Label: "Age", FieldType: annotation.FieldTypeNumber, Value: annotation.ValueRef{Path: "age"}},
				},
			},
		},
	}

	formatted, err := r.FormatAllFields(ctx, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(formatted) != 2 {
		t.Errorf("expected 2 formatted fields, got %d", len(formatted))
	}
	if formatted["name"] != "John" {
		t.Errorf("expected John, got %s", formatted["name"])
	}
}
