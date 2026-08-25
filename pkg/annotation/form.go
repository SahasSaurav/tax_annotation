package annotation

import "encoding/json"

// Form represents a complete tax form definition loaded from JSON.
// It contains metadata, page dimensions, and the annotations that
// define where values appear on each page.
type Form struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
	Pages       []Page `json:"pages"`
}

// Page represents a single page within a tax form. Each page has its
// own dimensions and a list of annotations that reference data values.
type Page struct {
	Number      int          `json:"number"`
	Label       string       `json:"label,omitempty"`
	Width       float64      `json:"width,omitempty"`
	Height      float64      `json:"height,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

// UnmarshalJSON implements custom JSON deserialization for Form.
// It applies default page dimensions (612x792 points, US Letter)
// when width or height are not specified in the JSON.
func (f *Form) UnmarshalJSON(data []byte) error {
	type Alias Form
	aux := &struct {
		*Alias
	}{
		(*Alias)(f),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	for i := range f.Pages {
		if f.Pages[i].Width == 0 {
			f.Pages[i].Width = 612
		}
		if f.Pages[i].Height == 0 {
			f.Pages[i].Height = 792
		}
	}
	return nil
}
