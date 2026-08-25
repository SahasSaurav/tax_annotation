package annotation

import "encoding/json"

type Form struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
	Pages       []Page `json:"pages"`
}

type Page struct {
	Number      int          `json:"number"`
	Label       string       `json:"label,omitempty"`
	Width       float64      `json:"width,omitempty"`
	Height      float64      `json:"height,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

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
