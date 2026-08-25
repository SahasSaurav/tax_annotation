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

	t.Run("nil value", func(t *testing.T) {
		result := v.Validate(nil, annotation.FieldTypeText, validation)
		if result.Valid {
			t.Error("expected invalid for nil value")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		result := v.Validate("  ", annotation.FieldTypeText, validation)
		if result.Valid {
			t.Error("expected invalid for empty string")
		}
	})

	t.Run("valid value", func(t *testing.T) {
		result := v.Validate("hello", annotation.FieldTypeText, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})
}

func TestValidateType(t *testing.T) {
	v := New()

	t.Run("string type valid", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeString}
		result := v.Validate("hello", annotation.FieldTypeText, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})

	t.Run("string type invalid", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeString}
		result := v.Validate(123, annotation.FieldTypeText, validation)
		if result.Valid {
			t.Error("expected invalid for number with string type")
		}
	})

	t.Run("number type valid", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeNumber}
		result := v.Validate(123.0, annotation.FieldTypeNumber, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})

	t.Run("number type invalid", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeNumber}
		result := v.Validate("abc", annotation.FieldTypeNumber, validation)
		if result.Valid {
			t.Error("expected invalid for string with number type")
		}
	})
}

func TestValidateNumericBounds(t *testing.T) {
	v := New()

	t.Run("min bound", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeNumber, Min: floatPtr(10)}
		result := v.Validate(5.0, annotation.FieldTypeNumber, validation)
		if result.Valid {
			t.Error("expected invalid for value below min")
		}
	})

	t.Run("max bound", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeNumber, Max: floatPtr(100)}
		result := v.Validate(150.0, annotation.FieldTypeNumber, validation)
		if result.Valid {
			t.Error("expected invalid for value above max")
		}
	})

	t.Run("within bounds", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeNumber, Min: floatPtr(0), Max: floatPtr(100)}
		result := v.Validate(50.0, annotation.FieldTypeNumber, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})
}

func TestValidateStringRules(t *testing.T) {
	v := New()

	t.Run("min length", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeString, MinLength: intPtr(5)}
		result := v.Validate("hi", annotation.FieldTypeText, validation)
		if result.Valid {
			t.Error("expected invalid for string below min length")
		}
	})

	t.Run("max length", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeString, MaxLength: intPtr(3)}
		result := v.Validate("hello", annotation.FieldTypeText, validation)
		if result.Valid {
			t.Error("expected invalid for string above max length")
		}
	})

	t.Run("within length", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeString, MinLength: intPtr(2), MaxLength: intPtr(5)}
		result := v.Validate("hello", annotation.FieldTypeText, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})
}

func TestValidateDateFormat(t *testing.T) {
	v := New()

	t.Run("valid date", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeDate}
		result := v.Validate("2025-01-15", annotation.FieldTypeDate, validation)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		validation := &annotation.Validation{Type: annotation.DataTypeDate}
		result := v.Validate("not-a-date", annotation.FieldTypeDate, validation)
		if result.Valid {
			t.Error("expected invalid for bad date format")
		}
	})
}

func TestValidateAll(t *testing.T) {
	v := New()

	t.Run("with matching pattern", func(t *testing.T) {
		validation := &annotation.Validation{Required: true}
		patterns := map[string]string{
			"email": `^[a-z]+@[a-z]+\.[a-z]+$`,
		}
		result := v.ValidateAll("test@example.com", annotation.FieldTypeText, validation, patterns)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})

	t.Run("with non-matching pattern", func(t *testing.T) {
		validation := &annotation.Validation{Required: true}
		patterns := map[string]string{
			"email": `^[a-z]+@[a-z]+\.[a-z]+$`,
		}
		result := v.ValidateAll("invalid-email", annotation.FieldTypeText, validation, patterns)
		if result.Valid {
			t.Error("expected invalid for non-matching pattern")
		}
	})

	t.Run("nil patterns", func(t *testing.T) {
		validation := &annotation.Validation{Required: true}
		result := v.ValidateAll("test", annotation.FieldTypeText, validation, nil)
		if !result.Valid {
			t.Errorf("expected valid, got invalid: %s", result.Message)
		}
	})
}
