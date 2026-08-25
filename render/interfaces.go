package render

import (
	"context"

	"github.com/sahassauarv/tax-annotation/annotation"
)

// Renderer defines the contract for computing a RenderResult from a form definition.
// Implementations resolve data paths, validate values, and format them into display
// strings. They produce a RenderResult without any output side effects, allowing
// callers to choose how to present the result (terminal, PDF, HTML, etc.).
type Renderer interface {
	// RenderForm processes all pages and annotations in a form, returning a
	// complete RenderResult with every field resolved, validated, and formatted.
	RenderForm(ctx context.Context, form *annotation.Form) (*RenderResult, error)

	// RenderPage processes all annotations on a single page.
	RenderPage(ctx context.Context, page *annotation.Page) (*RenderResult, error)

	// RenderAnnotation resolves, validates, and formats a single annotation.
	RenderAnnotation(ctx context.Context, ann *annotation.Annotation) (*RenderedField, error)

	// RenderFieldByID looks up an annotation by ID and renders it.
	RenderFieldByID(ctx context.Context, form *annotation.Form, fieldID string) (*RenderedField, error)

	// FormatAllFields renders every annotation and returns a map of ID → formatted string.
	FormatAllFields(ctx context.Context, form *annotation.Form) (map[string]string, error)

	// GetFieldSummary renders the form and builds a plain-text summary grouped by page.
	GetFieldSummary(ctx context.Context, form *annotation.Form) (string, error)
}

// Writer defines the contract for outputting a RenderResult to a destination.
// Implementations decide how to present the rendered data — terminal, PDF, HTML, etc.
// The form is provided alongside the result so writers can access page structure
// and metadata that isn't duplicated in RenderResult.
type Writer interface {
	// Write outputs the rendered result. The form provides page layout context.
	Write(ctx context.Context, result *RenderResult, form *annotation.Form) error
}
