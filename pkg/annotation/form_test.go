package annotation

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantID      string
		wantWidth   float64
		wantHeight  float64
		wantPages   int
		wantErr     bool
	}{
		{
			name:       "default dimensions",
			data:       `{"id":"W-2","name":"Wage and Tax Statement","version":"2025","pages":[{"number":1,"label":"Page 1"}]}`,
			wantID:     "W-2",
			wantWidth:  612,
			wantHeight: 792,
			wantPages:  1,
		},
		{
			name:       "custom dimensions preserved",
			data:       `{"id":"X","name":"Test","pages":[{"number":1,"width":500,"height":700}]}`,
			wantID:     "X",
			wantWidth:  500,
			wantHeight: 700,
			wantPages:  1,
		},
		{
			name:    "invalid JSON",
			data:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var form Form
			err := json.Unmarshal([]byte(tt.data), &form)
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
			if len(form.Pages) != tt.wantPages {
				t.Fatalf("pages: got %d, want %d", len(form.Pages), tt.wantPages)
			}
			if form.Pages[0].Width != tt.wantWidth {
				t.Errorf("Width: got %f, want %f", form.Pages[0].Width, tt.wantWidth)
			}
			if form.Pages[0].Height != tt.wantHeight {
				t.Errorf("Height: got %f, want %f", form.Pages[0].Height, tt.wantHeight)
			}
		})
	}
}
