package formatter

import (
	"fmt"
	"strings"

	"github.com/sahassauarv/tax-annotation/annotation"
)

func (f *Formatter) formatText(value interface{}, format *annotation.Format) string {
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

func (f *Formatter) formatBoolean(value interface{}) string {
	b, ok := value.(bool)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	if b {
		return "Yes"
	}
	return "No"
}

func (f *Formatter) formatSSN(value interface{}, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatEIN(value interface{}, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatPhone(value interface{}, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatZIP(value interface{}, format *annotation.Format) (string, error) {
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

func extractDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

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
