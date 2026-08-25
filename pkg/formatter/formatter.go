package formatter

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// formatter is the default implementation of the Formatter interface.
// It converts raw Go values into human-readable display strings based
// on the annotation's format configuration.
type formatter struct{}

// New creates a Formatter that satisfies the Formatter interface.
func New() Formatter {
	return &formatter{}
}

// Format converts a raw value into a display string using the specified field type
// and format configuration. If format is nil, it falls back to a default representation
// based on the field type. Returns an error if the value cannot be converted.
func (f *formatter) Format(value interface{}, fieldType annotation.FieldType, format *annotation.Format) (string, error) {
	if value == nil {
		return "", nil
	}
	if format == nil {
		return f.formatDefault(value, fieldType), nil
	}

	switch format.Type {
	case annotation.FormatCurrency:
		return f.formatCurrency(value, format)
	case annotation.FormatNumber:
		return f.formatNumber(value, format)
	case annotation.FormatDate:
		return f.formatDate(value, format)
	case annotation.FormatBoolean:
		return f.formatBoolean(value), nil
	case annotation.FormatSSN:
		return f.formatSSN(value, format)
	case annotation.FormatEIN:
		return f.formatEIN(value, format)
	case annotation.FormatPhone:
		return f.formatPhone(value, format)
	case annotation.FormatZIP:
		return f.formatZIP(value, format)
	case annotation.FormatPercent:
		return f.formatPercent(value, format)
	case annotation.FormatText, "":
		return f.formatText(value, format), nil
	default:
		return f.formatDefault(value, fieldType), nil
	}
}

// formatDefault provides a fallback representation when no explicit format is specified.
// Checkboxes are rendered as "Yes"/"No"; everything else uses fmt.Sprintf.
func (f *formatter) formatDefault(value interface{}, fieldType annotation.FieldType) string {
	switch fieldType {
	case annotation.FieldTypeCheckbox:
		return f.formatBoolean(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}
