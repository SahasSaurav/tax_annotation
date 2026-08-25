package formatter

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// formatCurrency converts a numeric value to a currency string with a dollar sign prefix
// and two decimal places by default. Decimals, prefix, and suffix can be overridden.
func (f *formatter) formatCurrency(value interface{}, format *annotation.Format) (string, error) {
	num, err := ToFloat64(value)
	if err != nil {
		return "", fmt.Errorf("convert to currency: %w", err)
	}
	decimals := 2
	if format.Decimals != nil {
		decimals = *format.Decimals
	}
	prefix := "$"
	if format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format.Suffix != "" {
		suffix = format.Suffix
	}
	formatted := FormatFloatWithDecimals(num, decimals)
	return prefix + formatted + suffix, nil
}

// formatNumber converts a numeric value to a plain number string with optional
// decimal places, prefix, and suffix. Defaults to zero decimal places.
func (f *formatter) formatNumber(value interface{}, format *annotation.Format) (string, error) {
	num, err := ToFloat64(value)
	if err != nil {
		return "", fmt.Errorf("convert to number: %w", err)
	}
	decimals := 0
	if format.Decimals != nil {
		decimals = *format.Decimals
	}
	prefix := ""
	if format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := ""
	if format.Suffix != "" {
		suffix = format.Suffix
	}
	formatted := FormatFloatWithDecimals(num, decimals)
	return prefix + formatted + suffix, nil
}

// formatPercent converts a decimal fraction (e.g. 0.15) to a percentage string (e.g. "15%").
// The value is multiplied by 100 before formatting. Defaults to no decimals and "%" suffix.
func (f *formatter) formatPercent(value interface{}, format *annotation.Format) (string, error) {
	num, err := ToFloat64(value)
	if err != nil {
		return "", fmt.Errorf("convert to percent: %w", err)
	}
	decimals := 0
	if format != nil && format.Decimals != nil {
		decimals = *format.Decimals
	}
	prefix := ""
	if format != nil && format.Prefix != "" {
		prefix = format.Prefix
	}
	suffix := "%"
	if format != nil && format.Suffix != "" {
		suffix = format.Suffix
	}
	formatted := FormatFloatWithDecimals(num*100, decimals)
	return prefix + formatted + suffix, nil
}
