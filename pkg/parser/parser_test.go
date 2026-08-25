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

	tests := []struct {
		name    string
		data    string
		ctx     context.Context
		wantErr bool
		wantID  string
	}{
		{"valid form", validJSON, ctx, false, "W-2"},
		{"invalid JSON", "invalid", ctx, true, ""},
		{"missing ID", `{"name": "Test", "pages": [{"number": 1, "label": "P1"}]}`, ctx, true, ""},
		{"empty pages", `{"id": "X", "name": "Y", "pages": []}`, ctx, true, ""},
	}

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	tests = append(tests, struct {
		name    string
		data    string
		ctx     context.Context
		wantErr bool
		wantID  string
	}{"cancelled context", validJSON, cancelledCtx, true, ""})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, err := p.ParseForm(tt.ctx, []byte(tt.data))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if form.ID != tt.wantID {
				t.Errorf("ID: got %s, want %s", form.ID, tt.wantID)
			}
			if len(form.Pages) != 1 {
				t.Errorf("pages: got %d, want 1", len(form.Pages))
			}
		})
	}
}

func TestLoadData(t *testing.T) {
	ctx := context.Background()
	p := New()

	tests := []struct {
		name    string
		data    string
		ctx     context.Context
		wantErr bool
	}{
		{"valid data", `{"name": "John", "age": 30}`, ctx, false},
		{"invalid JSON", "invalid", ctx, true},
	}

	cancelledCtx, cancel := context.WithCancel(ctx)
	cancel()
	tests = append(tests, struct {
		name    string
		data    string
		ctx     context.Context
		wantErr bool
	}{"cancelled context", `{}`, cancelledCtx, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := p.LoadData(tt.ctx, []byte(tt.data))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if data["name"] != "John" {
				t.Errorf("name: got %v, want John", data["name"])
			}
		})
	}
}

func TestValidateForm(t *testing.T) {
	tests := []struct {
		name    string
		form    *annotation.Form
		wantErr bool
	}{
		{
			name:    "valid form",
			form:    &annotation.Form{ID: "W-2", Name: "Wage and Tax Statement", Pages: []annotation.Page{{Number: 1, Label: "Page 1"}}},
			wantErr: false,
		},
		{
			name:    "missing ID",
			form:    &annotation.Form{Name: "Test", Pages: []annotation.Page{{Number: 1}}},
			wantErr: true,
		},
		{
			name:    "missing name",
			form:    &annotation.Form{ID: "X", Pages: []annotation.Page{{Number: 1}}},
			wantErr: true,
		},
		{
			name:    "no pages",
			form:    &annotation.Form{ID: "X", Name: "Y"},
			wantErr: true,
		},
		{
			name: "duplicate annotation IDs",
			form: &annotation.Form{
				ID: "X", Name: "Y",
				Pages: []annotation.Page{{
					Number:      1,
					Annotations: []annotation.Annotation{{ID: "dup"}, {ID: "dup"}},
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateForm(tt.form)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAnnotation(t *testing.T) {
	form := &annotation.Form{
		Pages: []annotation.Page{{
			Annotations: []annotation.Annotation{{ID: "field1"}, {ID: "field2"}},
		}},
	}

	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"found", "field1", true},
		{"not found", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ann := GetAnnotation(form, tt.id)
			if tt.want && ann == nil {
				t.Fatal("expected to find annotation")
			}
			if !tt.want && ann != nil {
				t.Fatal("expected nil")
			}
		})
	}
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
		"b": map[string]interface{}{"c": "2", "d": "3"},
	}
	overlay := map[string]interface{}{
		"b": map[string]interface{}{"c": "overridden"},
		"e": "new",
	}

	result := OverlayData(base, overlay)

	tests := []struct {
		key      string
		expected interface{}
	}{
		{"a", "1"},
		{"e", "new"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if result[tt.key] != tt.expected {
				t.Errorf("got %v, want %v", result[tt.key], tt.expected)
			}
		})
	}

	bMap, ok := result["b"].(map[string]interface{})
	if !ok {
		t.Fatal("expected b to be a map")
	}

	bTests := []struct {
		key      string
		expected interface{}
	}{
		{"c", "overridden"},
		{"d", "3"},
	}

	for _, tt := range bTests {
		t.Run("b."+tt.key, func(t *testing.T) {
			if bMap[tt.key] != tt.expected {
				t.Errorf("got %v, want %v", bMap[tt.key], tt.expected)
			}
		})
	}
}
