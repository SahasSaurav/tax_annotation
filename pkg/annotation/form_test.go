package annotation

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Run("default dimensions", func(t *testing.T) {
		data := `{
			"id": "W-2",
			"name": "Wage and Tax Statement",
			"version": "2025",
			"pages": [{
				"number": 1,
				"label": "Page 1"
			}]
		}`

		var form Form
		if err := json.Unmarshal([]byte(data), &form); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if form.ID != "W-2" {
			t.Errorf("expected ID W-2, got %s", form.ID)
		}
		if len(form.Pages) != 1 {
			t.Fatalf("expected 1 page, got %d", len(form.Pages))
		}
		if form.Pages[0].Width != 612 {
			t.Errorf("expected default width 612, got %f", form.Pages[0].Width)
		}
		if form.Pages[0].Height != 792 {
			t.Errorf("expected default height 792, got %f", form.Pages[0].Height)
		}
	})

	t.Run("custom dimensions preserved", func(t *testing.T) {
		data := `{
			"id": "X",
			"name": "Test",
			"pages": [{
				"number": 1,
				"width": 500,
				"height": 700
			}]
		}`

		var form Form
		if err := json.Unmarshal([]byte(data), &form); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if form.Pages[0].Width != 500 {
			t.Errorf("expected width 500, got %f", form.Pages[0].Width)
		}
		if form.Pages[0].Height != 700 {
			t.Errorf("expected height 700, got %f", form.Pages[0].Height)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var form Form
		if err := json.Unmarshal([]byte("invalid"), &form); err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}
