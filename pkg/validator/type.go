package validator

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// validateType checks that a value matches the expected data type.
// Strings must be Go strings, numbers can be numeric types or numeric strings,
// booleans can be Go bools or string representations, and dates must match common patterns.
func (v *formValidator) validateType(value interface{}, dataType annotation.DataType) ValidationResult {
	switch dataType {
	case annotation.DataTypeString:
		if _, ok := value.(string); !ok {
			return ValidationResult{Valid: false, Message: "expected string"}
		}
	case annotation.DataTypeNumber:
		switch value.(type) {
		case float64, float32, int, int64, int32, uint, uint64:
		case string:
			_, err := fmt.Sscanf(value.(string), "%f", new(float64))
			if err != nil {
				return ValidationResult{Valid: false, Message: "expected number"}
			}
		default:
			return ValidationResult{Valid: false, Message: "expected number"}
		}
	case annotation.DataTypeBoolean:
		if _, ok := value.(bool); !ok {
			if s, ok := value.(string); ok {
				if s != "true" && s != "false" && s != "True" && s != "False" && s != "1" && s != "0" {
					return ValidationResult{Valid: false, Message: "expected boolean"}
				}
			} else {
				return ValidationResult{Valid: false, Message: "expected boolean"}
			}
		}
	case annotation.DataTypeDate:
		if result := v.validateDateFormat(value); !result.Valid {
			return result
		}
	}
	return ValidationResult{Valid: true}
}
