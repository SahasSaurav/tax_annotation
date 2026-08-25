package validator

import "regexp"

func (v *Validator) validateDateFormat(value interface{}) ValidationResult {
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
