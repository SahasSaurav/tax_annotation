package render

import (
	"context"
	"testing"
	"time"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
	"github.com/sahassauarv/tax-annotation/pkg/validator"
)

func TestNewTerminalWriter(t *testing.T) {
	tw := NewTerminalWriter()
	if tw == nil {
		t.Fatal("NewTerminalWriter returned nil")
	}
}

func TestTerminalWriterWrite(t *testing.T) {
	tw := NewTerminalWriter()
	ctx := context.Background()

	result := &RenderResult{
		Fields: map[string]*RenderedField{
			"f1": {ID: "f1", Label: "Name", IsValid: true, HasValue: true, FormattedValue: "John"},
			"f2": {ID: "f2", Label: "SSN", IsValid: false, HasValue: true, FormattedValue: "bad", ValidationResult: validator.ValidationResult{Message: "invalid"}},
			"f3": {ID: "f3", Label: "Empty", IsValid: true, HasValue: false},
		},
		FormID:     "W-2",
		FormName:   "Wage and Tax Statement",
		RenderTime: time.Millisecond,
		Errors:     []string{"test error"},
	}

	form := &annotation.Form{
		ID: "W-2", Name: "Wage and Tax Statement",
		Pages: []annotation.Page{{
			Number: 1, Label: "Page 1",
			Annotations: []annotation.Annotation{
				{ID: "f1", Label: "Name"},
				{ID: "f2", Label: "SSN"},
				{ID: "f3", Label: "Empty"},
			},
		}},
	}

	tests := []struct {
		name    string
		result  *RenderResult
		form    *annotation.Form
		ctx     context.Context
		wantErr bool
	}{
		{"valid write", result, form, ctx, false},
		{"nil result", nil, &annotation.Form{}, ctx, true},
		{"nil form", &RenderResult{Fields: make(map[string]*RenderedField)}, nil, ctx, true},
	}

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	tests = append(tests, struct {
		name    string
		result  *RenderResult
		form    *annotation.Form
		ctx     context.Context
		wantErr bool
	}{"cancelled context", &RenderResult{Fields: make(map[string]*RenderedField)}, &annotation.Form{Pages: []annotation.Page{{Number: 1}}}, cancelledCtx, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tw.Write(tt.ctx, tt.result, tt.form)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTerminalWriterPrintDetailed(t *testing.T) {
	tw := NewTerminalWriter()
	result := &RenderResult{
		Fields: map[string]*RenderedField{
			"f1": {ID: "f1", Label: "Name", FieldType: annotation.FieldTypeText,
				Position: annotation.Position{X: 72, Y: 200, Width: 200, Height: 12},
				IsValid: true, HasValue: true, RawValue: "John", FormattedValue: "John"},
			"f2": {ID: "f2", Label: "Bad", FieldType: annotation.FieldTypeText,
				Position: annotation.Position{X: 72, Y: 220, Width: 200, Height: 12},
				IsValid: false, HasValue: true, FormattedValue: "bad",
				ValidationResult: validator.ValidationResult{Message: "invalid format"}},
			"f3": {ID: "f3", Label: "Empty", FieldType: annotation.FieldTypeText,
				Position: annotation.Position{X: 72, Y: 240, Width: 200, Height: 12},
				IsValid: true, HasValue: false},
		},
		FormID: "W-2", FormName: "Test",
	}

	tw.printDetailed(result)
}
