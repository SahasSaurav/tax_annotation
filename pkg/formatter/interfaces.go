package formatter

import (
	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// Formatter defines the contract for converting raw Go values into
// formatted display strings based on annotation format configuration.
type Formatter interface {
	// Format converts a raw value into a display string. The fieldType determines
	// the fallback formatting when no explicit format is provided.
	Format(value interface{}, fieldType annotation.FieldType, format *annotation.Format) (string, error)
}
