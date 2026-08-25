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
	ctx := context.Background()
	tw := NewTerminalWriter()

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
		ID:   "W-2",
		Name: "Wage and Tax Statement",
		Pages: []annotation.Page{
			{
				Number: 1,
				Label:  "Page 1",
				Annotations: []annotation.Annotation{
					{ID: "f1", Label: "Name"},
					{ID: "f2", Label: "SSN"},
					{ID: "f3", Label: "Empty"},
				},
			},
		},
	}

	err := tw.Write(ctx, result, form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTerminalWriterWriteNilResult(t *testing.T) {
	tw := NewTerminalWriter()
	err := tw.Write(context.Background(), nil, &annotation.Form{})
	if err == nil {
		t.Fatal("expected error for nil result")
	}
}

func TestTerminalWriterWriteNilForm(t *testing.T) {
	tw := NewTerminalWriter()
	err := tw.Write(context.Background(), &RenderResult{Fields: make(map[string]*RenderedField)}, nil)
	if err == nil {
		t.Fatal("expected error for nil form")
	}
}

func TestTerminalWriterWriteCancelled(t *testing.T) {
	tw := NewTerminalWriter()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := &RenderResult{Fields: make(map[string]*RenderedField)}
	form := &annotation.Form{Pages: []annotation.Page{{Number: 1}}}

	err := tw.Write(ctx, result, form)
	if err == nil {
		t.Fatal("expected error for cancelled context")
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
		FormID:   "W-2",
		FormName: "Test",
	}

	tw.printDetailed(result)
}
