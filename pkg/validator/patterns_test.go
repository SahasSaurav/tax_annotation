package validator

import "testing"

func TestValidateSSN(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"123-45-6789", true},
		{"000-00-0000", true},
		{"123456789", false},
		{"123-45-678", false},
		{"12-345-6789", false},
		{"", false},
		{"abc-de-fghi", false},
	}

	for _, tt := range tests {
		result := ValidateSSN(tt.input)
		if result.Valid != tt.valid {
			t.Errorf("ValidateSSN(%q): got valid=%v, want %v (msg: %s)", tt.input, result.Valid, tt.valid, result.Message)
		}
	}
}

func TestValidateZIP(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"62704", true},
		{"62704-1234", true},
		{"00000", true},
		{"1234", false},
		{"123456", false},
		{"12345-123", false},
		{"", false},
		{"abcde", false},
	}

	for _, tt := range tests {
		result := ValidateZIP(tt.input)
		if result.Valid != tt.valid {
			t.Errorf("ValidateZIP(%q): got valid=%v, want %v (msg: %s)", tt.input, result.Valid, tt.valid, result.Message)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"user@example.com", true},
		{"test.email@domain.co", true},
		{"name+tag@gmail.com", true},
		{"user@example", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
		{"plaintext", false},
	}

	for _, tt := range tests {
		result := ValidateEmail(tt.input)
		if result.Valid != tt.valid {
			t.Errorf("ValidateEmail(%q): got valid=%v, want %v (msg: %s)", tt.input, result.Valid, tt.valid, result.Message)
		}
	}
}
