package annotation

// Validation defines the constraints that a field's value must satisfy.
// All specified rules are checked; if any fail, the field is marked invalid.
type Validation struct {
	Required  bool     `json:"required,omitempty"`
	Type      DataType `json:"type,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
}

// DataType specifies the expected Go type for a field value during validation.
type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeDate    DataType = "date"
)
