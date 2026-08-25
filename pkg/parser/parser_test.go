package parser

import (
	"context"
	"testing"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestParseForm(t *testing.T) {
	ctx := context.Background()
	p := New()

	validJSON := `{
		"id": "W-2",
		"name": "Wage and Tax Statement",
		"version": "2024",
		"pages": [{
			"number": 1,
			"label": "Page 1",
			"annotations": [{
				"id": "test_field",
				"label": "Test Field",
				"fieldType": "text",
				"value": { "path": "test.path" },
				"position": { "x": 72, "y": 200, "width": 200, "height": 12 }
			}]
		}]
	}`

	t.Run("valid form", func(t *testing.T) {
		form, err := p.ParseForm(ctx, []byte(validJSON))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if form.ID != "W-2" {
			t.Errorf("expected ID W-2, got %s", form.ID)
		}
		if form.Name != "Wage and Tax Statement" {
			t.Errorf("expected name Wage and Tax Statement, got %s", form.Name)
		}
		if len(form.Pages) != 1 {
			t.Fatalf("expected 1 page, got %d", len(form.Pages))
		}
		if len(form.Pages[0].Annotations) != 1 {
			t.Fatalf("expected 1 annotation, got %d", len(form.Pages[0].Annotations))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := p.ParseForm(ctx, []byte("invalid"))
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		invalidJSON := `{"name": "Test", "pages": [{"number": 1, "label": "P1"}]}`
		_, err := p.ParseForm(ctx, []byte(invalidJSON))
		if err == nil {
			t.Fatal("expected error for missing ID")
		}
	})

	t.Run("empty pages", func(t *testing.T) {
		invalidJSON := `{"id": "X", "name": "Y", "pages": []}`
		_, err := p.ParseForm(ctx, []byte(invalidJSON))
		if err == nil {
			t.Fatal("expected error for empty pages")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := p.ParseForm(cancelledCtx, []byte(validJSON))
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})
}

func TestLoadData(t *testing.T) {
	ctx := context.Background()
	p := New()

	t.Run("valid data", func(t *testing.T) {
		json := `{"name": "John", "age": 30}`
		data, err := p.LoadData(ctx, []byte(json))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if data["name"] != "John" {
			t.Errorf("expected name John, got %v", data["name"])
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := p.LoadData(ctx, []byte("invalid"))
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()
		_, err := p.LoadData(cancelledCtx, []byte(`{}`))
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})
}

func TestValidateForm(t *testing.T) {
	validForm := &annotation.Form{
		ID:   "W-2",
		Name: "Wage and Tax Statement",
		Pages: []annotation.Page{
			{Number: 1, Label: "Page 1"},
		},
	}

	t.Run("valid form", func(t *testing.T) {
		if err := ValidateForm(validForm); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing ID", func(t *testing.T) {
		f := &annotation.Form{Name: "Test", Pages: []annotation.Page{{Number: 1}}}
		if err := ValidateForm(f); err == nil {
			t.Fatal("expected error for missing ID")
		}
	})

	t.Run("missing name", func(t *testing.T) {
		f := &annotation.Form{ID: "X", Pages: []annotation.Page{{Number: 1}}}
		if err := ValidateForm(f); err == nil {
			t.Fatal("expected error for missing name")
		}
	})

	t.Run("no pages", func(t *testing.T) {
		f := &annotation.Form{ID: "X", Name: "Y"}
		if err := ValidateForm(f); err == nil {
			t.Fatal("expected error for no pages")
		}
	})

	t.Run("duplicate annotation IDs", func(t *testing.T) {
		f := &annotation.Form{
			ID:   "X",
			Name: "Y",
			Pages: []annotation.Page{
				{
					Number: 1,
					Annotations: []annotation.Annotation{
						{ID: "dup"},
						{ID: "dup"},
					},
				},
			},
		}
		if err := ValidateForm(f); err == nil {
			t.Fatal("expected error for duplicate IDs")
		}
	})
}

func TestGetAnnotation(t *testing.T) {
	form := &annotation.Form{
		Pages: []annotation.Page{
			{
				Annotations: []annotation.Annotation{
					{ID: "field1"},
					{ID: "field2"},
				},
			},
		},
	}

	t.Run("found", func(t *testing.T) {
		ann := GetAnnotation(form, "field1")
		if ann == nil {
			t.Fatal("expected to find annotation")
		}
		if ann.ID != "field1" {
			t.Errorf("expected ID field1, got %s", ann.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		ann := GetAnnotation(form, "nonexistent")
		if ann != nil {
			t.Fatal("expected nil for nonexistent annotation")
		}
	})
}

func TestAllAnnotations(t *testing.T) {
	form := &annotation.Form{
		Pages: []annotation.Page{
			{Annotations: []annotation.Annotation{{ID: "a"}, {ID: "b"}}},
			{Annotations: []annotation.Annotation{{ID: "c"}}},
		},
	}

	all := AllAnnotations(form)
	if len(all) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(all))
	}
}

func TestOverlayData(t *testing.T) {
	base := map[string]interface{}{
		"a": "1",
		"b": map[string]interface{}{
			"c": "2",
			"d": "3",
		},
	}

	overlay := map[string]interface{}{
		"b": map[string]interface{}{
			"c": "overridden",
		},
		"e": "new",
	}

	result := OverlayData(base, overlay)

	if result["a"] != "1" {
		t.Errorf("expected a=1, got %v", result["a"])
	}
	if result["e"] != "new" {
		t.Errorf("expected e=new, got %v", result["e"])
	}

	bMap, ok := result["b"].(map[string]interface{})
	if !ok {
		t.Fatal("expected b to be a map")
	}
	if bMap["c"] != "overridden" {
		t.Errorf("expected b.c=overridden, got %v", bMap["c"])
	}
	if bMap["d"] != "3" {
		t.Errorf("expected b.d=3, got %v", bMap["d"])
	}
}
