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

	tests := []struct {
		name     string
		resolver parser.PathResolver
		fmtr     formatter.Formatter
		vld      validator.Validator
		wantErr  bool
	}{
		{"valid", resolver, fmtr, vld, false},
		{"nil resolver", nil, fmtr, vld, true},
		{"nil formatter", resolver, nil, vld, true},
		{"nil validator", resolver, fmtr, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRenderer(tt.resolver, tt.fmtr, tt.vld)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r == nil {
				t.Fatal("expected non-nil renderer")
			}
		})
	}
}

func TestRenderForm(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{
		"employee": map[string]interface{}{"name": "John Doe", "ssn": "123-45-6789"},
	}
	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID: "W-2", Name: "Wage and Tax Statement",
		Pages: []annotation.Page{{
			Number: 1, Label: "Page 1",
			Annotations: []annotation.Annotation{
				{ID: "emp_name", Label: "Employee Name", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "employee.name"}, Position: annotation.Position{X: 72, Y: 200, Width: 200, Height: 12}},
				{ID: "emp_ssn", Label: "Employee SSN", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "employee.ssn"}, Position: annotation.Position{X: 72, Y: 220, Width: 200, Height: 12}, Format: &annotation.Format{Type: annotation.FormatSSN}},
			},
		}},
	}

	tests := []struct {
		name    string
		form    *annotation.Form
		ctx     context.Context
		wantErr bool
	}{
		{"valid form", form, ctx, false},
		{"nil form", nil, ctx, true},
		{"empty pages", &annotation.Form{ID: "X", Name: "Y"}, ctx, true},
	}

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	tests = append(tests, struct {
		name    string
		form    *annotation.Form
		ctx     context.Context
		wantErr bool
	}{"cancelled context", form, cancelledCtx, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.RenderForm(tt.ctx, tt.form)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.FormID != "W-2" {
				t.Errorf("FormID: got %s, want W-2", result.FormID)
			}
			if len(result.Fields) != 2 {
				t.Errorf("fields: got %d, want 2", len(result.Fields))
			}
		})
	}
}

func TestRenderAnnotation(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{"amount": 75000.00}
	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	tests := []struct {
		name         string
		ann          *annotation.Annotation
		wantValue    string
		wantHasValue bool
		wantIsValid  bool
		wantErr      bool
	}{
		{
			name:         "with value",
			ann:          &annotation.Annotation{ID: "wages", Label: "Wages", FieldType: annotation.FieldTypeNumber, Value: annotation.ValueRef{Path: "amount"}, Position: annotation.Position{X: 72, Y: 200, Width: 200, Height: 12}, Format: &annotation.Format{Type: annotation.FormatCurrency}},
			wantValue:    "$75,000.00",
			wantHasValue: true,
			wantIsValid:  true,
		},
		{
			name:         "missing value",
			ann:          &annotation.Annotation{ID: "missing", Label: "Missing", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "nonexistent"}, Position: annotation.Position{X: 72, Y: 200, Width: 200, Height: 12}},
			wantHasValue: false,
			wantIsValid:  true,
		},
		{
			name:         "required missing value",
			ann:          &annotation.Annotation{ID: "required", Label: "Required", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "nonexistent"}, Position: annotation.Position{X: 72, Y: 200, Width: 200, Height: 12}, Validation: &annotation.Validation{Required: true}},
			wantHasValue: false,
			wantIsValid:  false,
		},
		{
			name:    "nil annotation",
			ann:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, err := r.RenderAnnotation(ctx, tt.ann)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if field.HasValue != tt.wantHasValue {
				t.Errorf("HasValue: got %v, want %v", field.HasValue, tt.wantHasValue)
			}
			if field.IsValid != tt.wantIsValid {
				t.Errorf("IsValid: got %v, want %v", field.IsValid, tt.wantIsValid)
			}
			if tt.wantValue != "" && field.FormattedValue != tt.wantValue {
				t.Errorf("FormattedValue: got %s, want %s", field.FormattedValue, tt.wantValue)
			}
		})
	}
}

func TestRenderResult(t *testing.T) {
	result := &RenderResult{
		Fields: map[string]*RenderedField{
			"valid":   {IsValid: true, HasValue: true},
			"invalid": {IsValid: false},
			"empty":   {IsValid: true, HasValue: false},
		},
	}

	tests := []struct {
		name      string
		fn        func() bool
		want      bool
	}{
		{"IsComplete returns false", result.IsComplete, false},
		{"IsComplete all valid", (&RenderResult{Fields: map[string]*RenderedField{"a": {IsValid: true, HasValue: true}}}).IsComplete, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

	countTests := []struct {
		name string
		fn   func() map[string]*RenderedField
		want int
	}{
		{"GetValidFields", result.GetValidFields, 2},
		{"GetInvalidFields", result.GetInvalidFields, 1},
	}

	for _, tt := range countTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.fn()); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFormatAllFields(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{"name": "John", "age": 30}
	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID: "TEST", Name: "Test",
		Pages: []annotation.Page{{
			Number: 1,
			Annotations: []annotation.Annotation{
				{ID: "name", Label: "Name", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "name"}},
				{ID: "age", Label: "Age", FieldType: annotation.FieldTypeNumber, Value: annotation.ValueRef{Path: "age"}},
			},
		}},
	}

	tests := []struct {
		name string
		key  string
		want string
	}{
		{"name field", "name", "John"},
		{"age field", "age", "30"},
	}

	formatted, err := r.FormatAllFields(ctx, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if formatted[tt.key] != tt.want {
				t.Errorf("got %s, want %s", formatted[tt.key], tt.want)
			}
		})
	}
}

func TestRenderFieldByID(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{"name": "John"}
	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID: "T", Name: "T",
		Pages: []annotation.Page{{
			Number: 1,
			Annotations: []annotation.Annotation{
				{ID: "name", Label: "Name", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "name"}},
			},
		}},
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"found", "name", false},
		{"not found", "missing", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, err := r.RenderFieldByID(ctx, form, tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if field.FormattedValue != "John" {
				t.Errorf("got %s, want John", field.FormattedValue)
			}
		})
	}
}

func TestGetFieldSummary(t *testing.T) {
	ctx := context.Background()
	data := map[string]interface{}{"name": "John"}
	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()
	r, _ := NewRenderer(resolver, fmtr, vld)

	form := &annotation.Form{
		ID: "T", Name: "Test",
		Pages: []annotation.Page{{
			Number: 1, Label: "Page 1",
			Annotations: []annotation.Annotation{
				{ID: "name", Label: "Name", FieldType: annotation.FieldTypeText, Value: annotation.ValueRef{Path: "name"}},
			},
		}},
	}

	summary, err := r.GetFieldSummary(ctx, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
