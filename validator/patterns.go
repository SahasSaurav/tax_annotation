package validator

import "regexp"

func ValidateSSN(value string) ValidationResult {
	pattern := `^\d{3}-\d{2}-\d{4}$`
	if matched, _ := regexp.MatchString(pattern, value); !matched {
		return ValidationResult{Valid: false, Message: "invalid SSN format (expected XXX-XX-XXXX)"}
	}
	return ValidationResult{Valid: true}
}

func ValidateZIP(value string) ValidationResult {
	pattern := `^\d{5}(-\d{4})?$`
	if matched, _ := regexp.MatchString(pattern, value); !matched {
		return ValidationResult{Valid: false, Message: "invalid ZIP format (expected XXXXX or XXXXX-XXXX)"}
	}
	return ValidationResult{Valid: true}
}

func ValidateEmail(value string) ValidationResult {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(pattern, value); !matched {
		return ValidationResult{Valid: false, Message: "invalid email format"}
	}
	return ValidationResult{Valid: true}
}
