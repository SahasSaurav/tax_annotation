# Tax Annotation

A Go library for parsing, rendering, and validating tax form annotations. It resolves data paths from JSON, formats values into display strings (currency, SSN, EIN, phone, ZIP, dates, etc.), and validates fields against configurable constraints.

## Project Structure

```
.
├── cmd/
│   └── taxrender/
│       └── main.go            # CLI entry point
├── pkg/
│   ├── annotation/            # Form, Page, Annotation, Format, Validation structs
│   ├── parser/                # JSON parsing, path resolution, form validation
│   ├── formatter/             # Value formatting (currency, SSN, EIN, phone, ZIP, date, percent)
│   ├── validator/             # Field validation (required, type, min/max, length)
│   └── render/                # Renderer + Writer interfaces, terminal output
├── fixtures/                  # Sample tax form JSON definitions
├── Makefile                   # Build, run, format, vet targets
├── go.mod
└── README.md
```

## Quick Start

```bash
# Build
make build

# Run with the default W-2 fixture
make run

# Format code
make fmt

# Vet
make vet
```

## Usage

### Basic

```go
ctx := context.Background()

// 1. Parse the form definition
p := parser.New()
form, _ := p.ParseFormFromFile(ctx, "fixtures/w2.json")

// 2. Prepare your data
data := map[string]interface{}{
    "employee": map[string]interface{}{
        "name": "John Doe",
        "ssn":  "123-45-6789",
    },
    "wages": map[string]interface{}{
        "box1": 75000.00,
    },
}

// 3. Create dependencies via interfaces
resolver := parser.NewPathResolver(data)
fmtr := formatter.New()
vld := validator.New()

// 4. Render
renderer, _ := render.NewRenderer(resolver, fmtr, vld)
result, _ := renderer.RenderForm(ctx, form)

// 5. Write output
writer := render.NewTerminalWriter()
writer.Write(ctx, result, form)
```

### Dependency Injection

All core components implement interfaces, so you can swap implementations:

```go
// Custom formatter for testing
type mockFormatter struct{}
func (m *mockFormatter) Format(v interface{}, ft annotation.FieldType, f *annotation.Format) (string, error) {
    return "MOCK", nil
}

renderer, _ := render.NewRenderer(resolver, &mockFormatter{}, vld)
```

### Available Interfaces

| Interface       | Package          | Methods                                           |
|-----------------|------------------|---------------------------------------------------|
| `Parser`        | `pkg/parser`     | `ParseFormFromFile`, `ParseForm`, `LoadDataFromFile`, `LoadData` |
| `PathResolver`  | `pkg/parser`     | `Resolve`, `GetString`, `GetFloat`, `GetBool`     |
| `Formatter`     | `pkg/formatter`  | `Format`                                          |
| `Validator`     | `pkg/validator`  | `Validate`, `ValidateAll`                         |
| `Renderer`      | `pkg/render`     | `RenderForm`, `RenderPage`, `RenderAnnotation`    |
| `Writer`        | `pkg/render`     | `Write`                                           |

### Adding a Custom Writer

Implement the `Writer` interface to add new output formats:

```go
type PDFWriter struct{ ... }

func (pw *PDFWriter) Write(ctx context.Context, result *render.RenderResult, form *annotation.Form) error {
    // Generate PDF from result.Fields
}

// Use it:
writer := &PDFWriter{...}
writer.Write(ctx, result, form)
```

### Format Types

| Type       | Example Output         | Description                    |
|------------|------------------------|--------------------------------|
| `currency` | `$75,000.00`           | Dollar sign + 2 decimals       |
| `number`   | `75,000`               | Plain number with commas       |
| `percent`  | `15%`                  | Decimal × 100 with % suffix    |
| `ssn`      | `123-45-6789`          | XXX-XX-XXXX format             |
| `ein`      | `12-3456789`           | XX-XXXXXXX format              |
| `phone`    | `(555) 123-4567`       | (XXX) XXX-XXXX format          |
| `zip`      | `62704`                | 5-digit or XXXXX-XXXX          |
| `date`     | `01/15/2025`           | Configurable pattern           |
| `boolean`  | `Yes` / `No`           | Boolean display                |
| `text`     | `raw value`            | Plain text with prefix/suffix  |

### Validation Rules

```json
{
    "validation": {
        "required": true,
        "type": "number",
        "min": 0,
        "max": 999999,
        "minLength": 1,
        "maxLength": 100
    }
}
```

### Fixture JSON Format

Each fixture defines a form with pages and annotations:

```json
{
    "id": "W-2",
    "name": "Wage and Tax Statement",
    "version": "2024",
    "pages": [
        {
            "number": 1,
            "label": "Employee Information",
            "annotations": [
                {
                    "id": "employee_ssn",
                    "label": "Employee SSN",
                    "fieldType": "text",
                    "value": { "path": "employee.ssn" },
                    "position": { "x": 72, "y": 200, "width": 200, "height": 12 },
                    "format": { "type": "ssn" },
                    "validation": { "required": true }
                }
            ]
        }
    ]
}
```

## Supported Forms

| Form    | Fixture File            | Pages |
|---------|-------------------------|-------|
| W-2     | `fixtures/w2.json`      | 3     |
| 1099-INT| `fixtures/1099-int.json`| 3     |
| Schedule A | `fixtures/schedule-a.json` | 6 |
| W-4     | `fixtures/w4.json`      | 5     |

## Adding a New Form

1. Create a JSON fixture in `fixtures/` following the format above
2. Add sample data in `cmd/taxrender/main.go`
3. Update the fixture path in `cmd/taxrender/main.go`
4. Run `make run`
