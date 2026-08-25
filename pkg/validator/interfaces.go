package validator

import (
	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// Validator defines the contract for validating raw values against field rules.
type Validator interface {
	// Validate checks a value against the field's type and validation constraints.
	// Returns a ValidationResult indicating pass/fail with an optional message.
	Validate(value interface{}, fieldType annotation.FieldType, validation *annotation.Validation) ValidationResult

	// ValidateAll runs standard validation plus additional regex pattern checks.
	ValidateAll(value interface{}, fieldType annotation.FieldType, validation *annotation.Validation, patterns map[string]string) ValidationResult
}

// ValidationResult holds the outcome of a single validation check.
type ValidationResult struct {
	Valid   bool
	Message string
}
