package formatter

import (
	"fmt"
	"time"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// formatDate converts a date value to a string using the pattern specified in the
// format configuration. Defaults to "01/02/2006" (US short date). Accepts string
// dates in multiple common formats and time.Time values.
func (f *formatter) formatDate(value interface{}, format *annotation.Format) (string, error) {
	pattern := "01/02/2006"
	if format.Pattern != "" {
		pattern = format.Pattern
	}

	switch v := value.(type) {
	case string:
		t, err := parseTimeFlexible(v)
		if err != nil {
			return v, nil
		}
		return t.Format(pattern), nil
	case time.Time:
		return v.Format(pattern), nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

// parseTimeFlexible tries multiple common date formats to parse a string into time.Time.
// Returns an error if none of the formats match.
func parseTimeFlexible(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"01-02-2006",
		"2006/01/02",
		"January 2, 2006",
		"Jan 2, 2006",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range formats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q", s)
}
