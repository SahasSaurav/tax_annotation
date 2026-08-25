package formatter

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/annotation"
)

type Formatter struct{}

func New() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(value interface{}, fieldType annotation.FieldType, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatDefault(value interface{}, fieldType annotation.FieldType) string {
	switch fieldType {
	case annotation.FieldTypeCheckbox:
		return f.formatBoolean(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}
