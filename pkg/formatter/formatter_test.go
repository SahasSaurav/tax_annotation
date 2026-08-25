package formatter

import (
	"testing"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Fatal("New() returned nil")
	}
}

func TestFormatNil(t *testing.T) {
	f := New()
	result, err := f.Format(nil, annotation.FieldTypeText, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for nil, got %s", result)
	}
}

func TestFormatCurrency(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatCurrency}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"integer", 75000, "$75,000.00"},
		{"float", 75000.00, "$75,000.00"},
		{"negative", -1250.50, "$-1,250.50"},
		{"zero", 0, "$0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.value, annotation.FieldTypeNumber, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatNumber}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"integer", 75000, "75,000"},
		{"float", 75000.00, "75,000"},
		{"small", 123, "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.value, annotation.FieldTypeNumber, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatPercent(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatPercent}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"decimal", 0.15, "15%"},
		{"whole", 1.0, "100%"},
		{"zero", 0.0, "0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.value, annotation.FieldTypeNumber, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatSSN(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatSSN}

	result, err := f.Format("123456789", annotation.FieldTypeText, format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "123-45-6789" {
		t.Errorf("expected 123-45-6789, got %s", result)
	}
}

func TestFormatEIN(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatEIN}

	result, err := f.Format("123456789", annotation.FieldTypeText, format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "12-3456789" {
		t.Errorf("expected 12-3456789, got %s", result)
	}
}

func TestFormatPhone(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatPhone}

	result, err := f.Format("5551234567", annotation.FieldTypeText, format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "(555) 123-4567" {
		t.Errorf("expected (555) 123-4567, got %s", result)
	}
}

func TestFormatZIP(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatZIP}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"5 digit", "62704", "62704"},
		{"9 digit", "627041234", "62704-1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.value, annotation.FieldTypeText, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatBoolean(t *testing.T) {
	f := New()
	format := &annotation.Format{Type: annotation.FormatBoolean}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"true", true, "Yes"},
		{"false", false, "No"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.value, annotation.FieldTypeCheckbox, format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatText(t *testing.T) {
	f := New()

	t.Run("with prefix and suffix", func(t *testing.T) {
		format := &annotation.Format{
			Type:   annotation.FormatText,
			Prefix: "Box ",
			Suffix: ":",
		}
		result, err := f.Format("12", annotation.FieldTypeText, format)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Box 12:" {
			t.Errorf("expected 'Box 12:', got %s", result)
		}
	})

	t.Run("without prefix suffix", func(t *testing.T) {
		format := &annotation.Format{Type: annotation.FormatText}
		result, err := f.Format("hello", annotation.FieldTypeText, format)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("expected hello, got %s", result)
		}
	})
}

func TestFormatDefault(t *testing.T) {
	f := New()

	t.Run("checkbox default", func(t *testing.T) {
		result, err := f.Format(true, annotation.FieldTypeCheckbox, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Yes" {
			t.Errorf("expected Yes, got %s", result)
		}
	})

	t.Run("text default", func(t *testing.T) {
		result, err := f.Format("hello", annotation.FieldTypeText, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("expected hello, got %s", result)
		}
	})
}

func TestFormatDate(t *testing.T) {
	f := New()
	format := &annotation.Format{
		Type:    annotation.FormatDate,
		Pattern: "01/02/2006",
	}

	result, err := f.Format("2025-01-15", annotation.FieldTypeDate, format)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "01/15/2025" {
		t.Errorf("expected 01/15/2025, got %s", result)
	}
}

func TestAlignText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		width     int
		alignment annotation.Alignment
		expected  string
	}{
		{"left", "hi", 10, annotation.AlignLeft, "hi        "},
		{"center", "hi", 10, annotation.AlignCenter, "    hi    "},
		{"right", "hi", 10, annotation.AlignRight, "        hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignText(tt.text, tt.width, tt.alignment)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
