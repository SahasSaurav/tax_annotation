package parser

import (
	"context"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// Parser defines the contract for loading and validating tax form definitions.
type Parser interface {
	// ParseFormFromFile reads a form definition from disk and returns a validated Form.
	ParseFormFromFile(ctx context.Context, path string) (*annotation.Form, error)

	// ParseForm unmarshals JSON bytes into a Form and validates it.
	ParseForm(ctx context.Context, data []byte) (*annotation.Form, error)

	// LoadDataFromFile reads a JSON data file and returns a map.
	LoadDataFromFile(ctx context.Context, path string) (map[string]interface{}, error)

	// LoadData unmarshals JSON bytes into a map.
	LoadData(ctx context.Context, data []byte) (map[string]interface{}, error)
}

// PathResolver defines the contract for resolving dot-notation paths against a data map.
type PathResolver interface {
	// Resolve walks the data map using the given path and returns the value, if found.
	Resolve(path string) (interface{}, bool)

	// GetString resolves a path and returns the value as a string, or fallback on failure.
	GetString(path, fallback string) string

	// GetFloat resolves a path and returns the value as a float64, or fallback on failure.
	GetFloat(path string, fallback float64) float64

	// GetBool resolves a path and returns the value as a bool, or fallback on failure.
	GetBool(path string, fallback bool) bool
}
