package formatter

import (
	"fmt"
	"math"
	"strings"
)

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		if err != nil {
			return 0, fmt.Errorf("parse %q as number: %w", v, err)
		}
		return f, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

func FormatFloatWithDecimals(f float64, decimals int) string {
	rounded := math.Round(f*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))

	intPart := int64(math.Abs(rounded))
	decPart := math.Abs(rounded) - float64(intPart)

	intStr := FormatIntWithCommas(intPart)
	if decimals <= 0 {
		if rounded < 0 {
			return "-" + intStr
		}
		return intStr
	}

	decStr := fmt.Sprintf("%.*f", decimals, decPart)
	decStr = strings.TrimPrefix(decStr, "0")

	if rounded < 0 {
		return "-" + intStr + decStr
	}
	return intStr + decStr
}

func FormatIntWithCommas(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	parts := []string{}
	for n >= 1000 {
		parts = append([]string{fmt.Sprintf("%03d", n%1000)}, parts...)
		n /= 1000
	}
	parts = append([]string{fmt.Sprintf("%d", n)}, parts...)
	return strings.Join(parts, ",")
}
