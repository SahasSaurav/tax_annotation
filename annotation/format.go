package annotation

// Format controls how a raw value is converted to a display string.
// It specifies the output type, decimal precision, prefix/suffix,
// text alignment, and whether the field supports multiline text.
type Format struct {
	Type      FormatType `json:"type"`
	Decimals  *int       `json:"decimals,omitempty"`
	Pattern   string     `json:"pattern,omitempty"`
	Prefix    string     `json:"prefix,omitempty"`
	Suffix    string     `json:"suffix,omitempty"`
	Alignment Alignment  `json:"alignment,omitempty"`
	Multiline bool       `json:"multiline,omitempty"`
}

// FormatType selects which formatter handles the conversion from raw value to string.
type FormatType string

const (
	FormatText     FormatType = "text"
	FormatNumber   FormatType = "number"
	FormatCurrency FormatType = "currency"
	FormatDate     FormatType = "date"
	FormatBoolean  FormatType = "boolean"
	FormatSSN      FormatType = "ssn"
	FormatEIN      FormatType = "ein"
	FormatPhone    FormatType = "phone"
	FormatZIP      FormatType = "zip"
	FormatPercent  FormatType = "percent"
)

// Alignment controls horizontal text alignment within the field's bounding box.
type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)
