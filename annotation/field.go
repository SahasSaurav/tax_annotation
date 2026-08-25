package annotation

type Annotation struct {
	ID         string      `json:"id"`
	Label      string      `json:"label,omitempty"`
	FieldType  FieldType   `json:"fieldType"`
	Value      ValueRef    `json:"value"`
	Position   Position    `json:"position"`
	Format     *Format     `json:"format,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDate     FieldType = "date"
	FieldTypeCheckbox FieldType = "checkbox"
)

type ValueRef struct {
	Path string `json:"path"`
}

type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
