package formatter

import (
	"fmt"
	"strings"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// formatText returns the value as a string with optional prefix and suffix.
func (f *formatter) formatText(value interface{}, format *annotation.Format) string {
	s := fmt.Sprintf("%v", value)
	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	return prefix + s + suffix
}

// formatBoolean converts a bool to "Yes" or "No". Non-bool values are
// rendered using fmt.Sprintf.
func (f *formatter) formatBoolean(value interface{}) string {
	b, ok := value.(bool)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	if b {
		return "Yes"
	}
	return "No"
}

// formatSSN formats a 9-digit value as XXX-XX-XXXX. Strips non-digit characters
// before formatting. Returns an error if the input does not contain exactly 9 digits.
func (f *formatter) formatSSN(value interface{}, format *annotation.Format) (string, error) {
	s := strings.TrimSpace(fmt.Sprintf("%v", value))
	digits := extractDigits(s)

	if len(digits) != 9 {
		return "", fmt.Errorf("SSN must contain exactly 9 digits, got %d", len(digits))
	}

	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	return prefix + digits[:3] + "-" + digits[3:5] + "-" + digits[5:] + suffix, nil
}

// formatEIN formats a 9-digit value as XX-XXXXXXX. Strips non-digit characters
// before formatting. Returns an error if the input does not contain exactly 9 digits.
func (f *formatter) formatEIN(value interface{}, format *annotation.Format) (string, error) {
	s := strings.TrimSpace(fmt.Sprintf("%v", value))
	digits := extractDigits(s)

	if len(digits) != 9 {
		return "", fmt.Errorf("EIN must contain exactly 9 digits, got %d", len(digits))
	}

	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	return prefix + digits[:2] + "-" + digits[2:] + suffix, nil
}

// formatPhone formats a 10-digit value as (XXX) XXX-XXXX. A leading "1" is
// stripped if present (11-digit US numbers). Returns an error for other lengths.
func (f *formatter) formatPhone(value interface{}, format *annotation.Format) (string, error) {
	s := strings.TrimSpace(fmt.Sprintf("%v", value))
	digits := extractDigits(s)

	if len(digits) == 11 && digits[0] == '1' {
		digits = digits[1:]
	}
	if len(digits) != 10 {
		return "", fmt.Errorf("phone must contain 10 digits, got %d", len(digits))
	}

	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	return prefix + "(" + digits[:3] + ") " + digits[3:6] + "-" + digits[6:] + suffix, nil
}

// formatZIP formats a 5 or 9 digit ZIP code. 9-digit codes are formatted as
// XXXXX-XXXX. Returns an error for other lengths.
func (f *formatter) formatZIP(value interface{}, format *annotation.Format) (string, error) {
	s := strings.TrimSpace(fmt.Sprintf("%v", value))
	digits := extractDigits(s)

	if len(digits) != 5 && len(digits) != 9 {
		return "", fmt.Errorf("ZIP must be 5 or 9 digits, got %d", len(digits))
	}

	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	if len(digits) == 9 {
		return prefix + digits[:5] + "-" + digits[5:] + suffix, nil
	}
	return prefix + digits + suffix, nil
}

// extractDigits strips all non-digit characters from a string.
func extractDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

// AlignText pads text to the given width using the specified alignment.
// If the text is already wider than the target, it is returned as-is.
func AlignText(text string, width int, alignment annotation.Alignment) string {
	runeLen := len([]rune(text))
	if runeLen >= width {
		return text
	}
	padding := width - runeLen
	switch alignment {
	case annotation.AlignRight:
		return strings.Repeat(" ", padding) + text
	case annotation.AlignCenter:
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default:
		return text + strings.Repeat(" ", padding)
	}
}
