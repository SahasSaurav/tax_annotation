package annotation

type Format struct {
	Type      FormatType `json:"type"`
	Decimals  *int       `json:"decimals,omitempty"`
	Pattern   string     `json:"pattern,omitempty"`
	Prefix    string     `json:"prefix,omitempty"`
	Suffix    string     `json:"suffix,omitempty"`
	Alignment Alignment  `json:"alignment,omitempty"`
	Multiline bool       `json:"multiline,omitempty"`
}

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

type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)
