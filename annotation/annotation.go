package annotation

type Form struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"Description"`
	Version     string `json:"version"`
	Pages       []Page `json:"pages"`
}

type Page struct {
	Number      int          `json:"number"`
	Label       string       `json:"label,omitempty"` 
	Annotations []Annotation `json:"annotations,omitempty"`
}

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

type Format struct {
	Type      FormatType `json:"type"`
	Decimals  *int       `json:"decimals,omitempty"`
	Pattern   string     `json:"pattern,omitempty"`
	Prefix    string     `json:"prefix,omitempty"`
	Suffix    string     `json:"suffix,omitempty"`
	Alignment Alignment  `json:"alignment,omitempty"`
}

type FormatType string

const (
	FormatText     FormatType = "text"
	FormatNumber   FormatType = "number"
	FormatCurrency FormatType = "currency"
	FormatDate     FormatType = "date"
	FormatBoolean  FormatType = "boolean"
)

type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

type Validation struct {
	Required  bool     `json:"required,omitempty"`
	Type      DataType `json:"type,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
}

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeDate    DataType = "date"
)
