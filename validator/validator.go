package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// formValidator is the default implementation of the Validator interface.
// It checks values against type, range, length, and format constraints.
type formValidator struct{}

// New creates a Validator that satisfies the Validator interface.
func New() Validator {
	return &formValidator{}
}

// Validate checks a value against the field's type and validation constraints.
// It validates in order: required check, type check, numeric bounds, string length,
// and date format. All specified rules must pass for the result to be valid.
func (v *formValidator) Validate(value interface{}, fieldType annotation.FieldType, validation *annotation.Validation) ValidationResult {
	if validation == nil {
		return ValidationResult{Valid: true}
	}

	if validation.Required {
		if value == nil {
			return ValidationResult{Valid: false, Message: "field is required"}
		}
		if s, ok := value.(string); ok && strings.TrimSpace(s) == "" {
			return ValidationResult{Valid: false, Message: "field is required"}
		}
	}

	if value == nil {
		return ValidationResult{Valid: true}
	}

	if validation.Type != "" {
		if result := v.validateType(value, validation.Type); !result.Valid {
			return result
		}
	}

	if validation.Type == annotation.DataTypeNumber || fieldType == annotation.FieldTypeNumber {
		if result := v.validateNumericBounds(value, validation); !result.Valid {
			return result
		}
	}

	if validation.Type == annotation.DataTypeString || fieldType == annotation.FieldTypeText {
		if result := v.validateStringRules(value, validation); !result.Valid {
			return result
		}
	}

	if validation.Type == annotation.DataTypeDate || fieldType == annotation.FieldTypeDate {
		if result := v.validateDateFormat(value); !result.Valid {
			return result
		}
	}

	return ValidationResult{Valid: true}
}

// ValidateAll runs standard validation plus additional regex pattern checks.
// If any pattern does not match, the result is invalid.
func (v *formValidator) ValidateAll(value interface{}, fieldType annotation.FieldType, validation *annotation.Validation, patterns map[string]string) ValidationResult {
	result := v.Validate(value, fieldType, validation)
	if !result.Valid {
		return result
	}

	if patterns != nil && value != nil {
		if s, ok := value.(string); ok {
			for name, pattern := range patterns {
				if matched, _ := regexp.MatchString(pattern, s); !matched {
					return ValidationResult{Valid: false, Message: fmt.Sprintf("does not match %s format", name)}
				}
			}
		}
	}

	return ValidationResult{Valid: true}
}
