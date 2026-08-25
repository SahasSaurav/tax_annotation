package annotation

// Annotation defines a single field on a tax form page. It specifies the field type,
// the data path to resolve the value from, the position on the page, optional formatting,
// and optional validation rules.
type Annotation struct {
	ID         string      `json:"id"`
	Label      string      `json:"label,omitempty"`
	FieldType  FieldType   `json:"fieldType"`
	Value      ValueRef    `json:"value"`
	Position   Position    `json:"position"`
	Format     *Format     `json:"format,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

// FieldType determines how a value is interpreted and rendered.
type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDate     FieldType = "date"
	FieldTypeCheckbox FieldType = "checkbox"
)

// ValueRef holds a dot-notation path used to look up the field's value
// in the data map (e.g. "employee.ssn" or "wages.box1").
type ValueRef struct {
	Path string `json:"path"`
}

// Position defines where an annotation is placed on the page, using
// top-left origin coordinates in PDF points (1/72 inch).
type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
