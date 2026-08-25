package formatter

import (
	"fmt"

	"github.com/sahassauarv/tax-annotation/annotation"
)

func (f *Formatter) formatCurrency(value interface{}, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatNumber(value interface{}, format *annotation.Format) (string, error) {
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

func (f *Formatter) formatPercent(value interface{}, format *annotation.Format) (string, error) {
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
