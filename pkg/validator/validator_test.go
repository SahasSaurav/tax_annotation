package validator

import (
	"testing"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

func floatPtr(f float64) *float64 { return &f }
func intPtr(i int) *int           { return &i }

func TestNew(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatal("New() returned nil")
	}
}

func TestValidateNilValidation(t *testing.T) {
	v := New()
	result := v.Validate("test", annotation.FieldTypeText, nil)
	if !result.Valid {
		t.Error("expected valid for nil validation")
	}
}

func TestValidateRequired(t *testing.T) {
	v := New()
	validation := &annotation.Validation{Required: true}

	tests := []struct {
		name  string
		value interface{}
		valid bool
	}{
		{"nil value", nil, false},
		{"empty string", "  ", false},
		{"valid value", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.value, annotation.FieldTypeText, validation)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}

func TestValidateType(t *testing.T) {
	v := New()

	tests := []struct {
		name      string
		value     interface{}
		fieldType annotation.FieldType
		dataType  annotation.DataType
		valid     bool
	}{
		{"string type valid", "hello", annotation.FieldTypeText, annotation.DataTypeString, true},
		{"string type invalid", 123, annotation.FieldTypeText, annotation.DataTypeString, false},
		{"number type valid", 123.0, annotation.FieldTypeNumber, annotation.DataTypeNumber, true},
		{"number type invalid", "abc", annotation.FieldTypeNumber, annotation.DataTypeNumber, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := &annotation.Validation{Type: tt.dataType}
			result := v.Validate(tt.value, tt.fieldType, validation)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}

func TestValidateNumericBounds(t *testing.T) {
	v := New()

	tests := []struct {
		name  string
		value interface{}
		min   *float64
		max   *float64
		valid bool
	}{
		{"below min", 5.0, floatPtr(10), nil, false},
		{"above max", 150.0, nil, floatPtr(100), false},
		{"within bounds", 50.0, floatPtr(0), floatPtr(100), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := &annotation.Validation{Type: annotation.DataTypeNumber, Min: tt.min, Max: tt.max}
			result := v.Validate(tt.value, annotation.FieldTypeNumber, validation)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}

func TestValidateStringRules(t *testing.T) {
	v := New()

	tests := []struct {
		name      string
		value     string
		minLength *int
		maxLength *int
		valid     bool
	}{
		{"below min length", "hi", intPtr(5), nil, false},
		{"above max length", "hello", nil, intPtr(3), false},
		{"within length", "hello", intPtr(2), intPtr(5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := &annotation.Validation{Type: annotation.DataTypeString, MinLength: tt.minLength, MaxLength: tt.maxLength}
			result := v.Validate(tt.value, annotation.FieldTypeText, validation)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	v := New()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid date", "2025-01-15", true},
		{"invalid date", "not-a-date", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation := &annotation.Validation{Type: annotation.DataTypeDate}
			result := v.Validate(tt.value, annotation.FieldTypeDate, validation)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}

func TestValidateAll(t *testing.T) {
	v := New()
	validation := &annotation.Validation{Required: true}
	patterns := map[string]string{"email": `^[a-z]+@[a-z]+\.[a-z]+$`}

	tests := []struct {
		name     string
		value    string
		patterns map[string]string
		valid    bool
	}{
		{"matching pattern", "test@example.com", patterns, true},
		{"non-matching pattern", "invalid-email", patterns, false},
		{"nil patterns", "test", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateAll(tt.value, annotation.FieldTypeText, validation, tt.patterns)
			if result.Valid != tt.valid {
				t.Errorf("got valid=%v, want %v (msg: %s)", result.Valid, tt.valid, result.Message)
			}
		})
	}
}
