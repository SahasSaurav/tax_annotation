package render

import (
	"context"
	"fmt"
	"strings"

	"github.com/sahassauarv/tax-annotation/pkg/annotation"
)

// TerminalWriter implements Writer and renders a RenderResult to stdout
// with clean line separators, status indicators, and a summary footer.
// This is the default output driver for terminal-based usage.
type TerminalWriter struct{}

// NewTerminalWriter creates a TerminalWriter.
func NewTerminalWriter() *TerminalWriter {
	return &TerminalWriter{}
}

// Write outputs the rendered form to stdout, printing each page's header,
// fields with status indicators, and a footer, followed by a summary.
func (tw *TerminalWriter) Write(ctx context.Context, result *RenderResult, form *annotation.Form) error {
	if result == nil {
		return fmt.Errorf("result cannot be nil")
	}
	if form == nil {
		return fmt.Errorf("form cannot be nil")
	}

	for i, page := range form.Pages {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("write cancelled: %w", err)
		}

		tw.printPageHeader(i, &page)

		for _, ann := range page.Annotations {
			field := result.Fields[ann.ID]
			if field != nil {
				tw.printField(i, field)
			}
		}

		tw.printPageFooter(i, &page)
	}

	tw.printSummary(result)
	return nil
}

// printSummary outputs a formatted summary showing form name, field counts,
// and overall completion state.
func (tw *TerminalWriter) printSummary(result *RenderResult) {
	valid := 0
	invalid := 0
	empty := 0
	for _, field := range result.Fields {
		if !field.IsValid {
			invalid++
		} else if !field.HasValue {
			empty++
		} else {
			valid++
		}
	}

	title := fmt.Sprintf("  %s (%s)", result.FormName, result.FormID)
	stats := fmt.Sprintf("  Fields: %-5d  Valid: %-5d  Invalid: %-5d  Empty: %-5d", len(result.Fields), valid, invalid, empty)
	complete := fmt.Sprintf("  Complete: %-5v  Render Time: %-20v", result.IsComplete(), result.RenderTime)

	w := len(title)
	if len(stats) > w {
		w = len(stats)
	}
	if len(complete) > w {
		w = len(complete)
	}
	w += 4

	line := strings.Repeat("━", w)
	subline := strings.Repeat("─", w)

	fmt.Println()
	fmt.Println(line)
	fmt.Println(title)
	fmt.Println(subline)
	fmt.Println(stats)
	fmt.Println(complete)
	if len(result.Errors) > 0 {
		fmt.Printf("  Errors: %d\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("    - %s\n", err)
		}
	}
	fmt.Println(line)
}

// printDetailed outputs every field with its full metadata including type,
// position, raw value, formatted value, and validation errors.
func (tw *TerminalWriter) printDetailed(result *RenderResult) {
	tw.printSummary(result)
	for id, field := range result.Fields {
		status := "OK"
		if !field.IsValid {
			status = "FAIL"
		} else if !field.HasValue {
			status = "EMPTY"
		}
		fmt.Printf("[%s] %s (%s)\n", status, field.Label, id)
		fmt.Printf("  Type: %s\n", field.FieldType)
		fmt.Printf("  Position: (%.0f, %.0f) %.0fx%.0f\n",
			field.Position.X, field.Position.Y, field.Position.Width, field.Position.Height)
		if field.HasValue {
			fmt.Printf("  Raw Value:   %v\n", field.RawValue)
			fmt.Printf("  Formatted:   %s\n", field.FormattedValue)
		}
		if !field.IsValid {
			fmt.Printf("  Error: %s\n", field.ValidationResult.Message)
		}
		fmt.Println()
	}
}

// printPageHeader prints a styled header for a form page.
func (tw *TerminalWriter) printPageHeader(idx int, page *annotation.Page) {
	title := fmt.Sprintf("  Page %d: %s", page.Number, page.Label)
	w := len(title) + 4
	line := strings.Repeat("━", w)
	subline := strings.Repeat("─", w)

	fmt.Println()
	fmt.Println(line)
	fmt.Println(title)
	fmt.Println(subline)
}

// printPageFooter prints a styled footer for a form page.
func (tw *TerminalWriter) printPageFooter(idx int, page *annotation.Page) {
	summary := fmt.Sprintf("  Page %d complete", page.Number)
	w := len(summary) + 4
	subline := strings.Repeat("─", w)
	line := strings.Repeat("━", w)

	fmt.Println(subline)
	fmt.Println(summary)
	fmt.Println(line)
}

// printField prints a single field with a status indicator.
func (tw *TerminalWriter) printField(idx int, field *RenderedField) {
	status := "✓"
	if !field.IsValid {
		status = "✗"
	} else if !field.HasValue {
		status = "·"
	}

	value := field.FormattedValue
	if value == "" {
		value = "(empty)"
	}

	fmt.Printf("  %s %-38s %s\n", status, field.Label, value)
}
