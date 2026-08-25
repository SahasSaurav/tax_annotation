package validator

import "regexp"

// validateDateFormat checks that a string value matches one of the common
// date formats: YYYY-MM-DD, MM/DD/YYYY, or MM-DD-YYYY.
func (v *formValidator) validateDateFormat(value interface{}) ValidationResult {
	s, ok := value.(string)
	if !ok {
		return ValidationResult{Valid: true}
	}

	datePatterns := []string{
		`^\d{4}-\d{2}-\d{2}$`,
		`^\d{2}/\d{2}/\d{4}$`,
		`^\d{2}-\d{2}-\d{4}$`,
	}

	for _, pattern := range datePatterns {
		if matched, _ := regexp.MatchString(pattern, s); matched {
			return ValidationResult{Valid: true}
		}
	}

	return ValidationResult{Valid: false, Message: "invalid date format"}
}
