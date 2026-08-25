package validator

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// validateStringRules checks that a string value's length falls within the
// min/max bounds specified in the validation rules.
func (v *formValidator) validateStringRules(value interface{}, validation *annotation.Validation) ValidationResult {
	s, ok := value.(string)
	if !ok {
		return ValidationResult{Valid: true}
	}

	if validation.MinLength != nil && len(s) < *validation.MinLength {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("length %d is below minimum %d", len(s), *validation.MinLength)}
	}

	if validation.MaxLength != nil && len(s) > *validation.MaxLength {
		return ValidationResult{Valid: false, Message: fmt.Sprintf("length %d exceeds maximum %d", len(s), *validation.MaxLength)}
	}

	return ValidationResult{Valid: true}
}
