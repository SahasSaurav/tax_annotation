package validator

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/annotation"
	"github.com/sahassauarv/tax-annotation/formatter"
)

// validateNumericBounds checks that a numeric value falls within the min/max range
// specified in the validation rules. Converts the value to float64 for comparison.
func (v *formValidator) validateNumericBounds(value interface{}, validation *annotation.Validation) ValidationResult {
	num, err := formatter.ToFloat64(value)
	if err != nil {
		return ValidationResult{Valid: false, Message: "cannot convert to number"}
	}

	if validation.Min != nil && num < *validation.Min {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("value %v is below minimum %v", num, *validation.Min)}
	}
	if validation.Max != nil && num > *validation.Max {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("value %v exceeds maximum %v", num, *validation.Max)}
	}

	return ValidationResult{Valid: true}
}
