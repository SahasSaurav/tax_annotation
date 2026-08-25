package main

import (
	"context"
	"log"

	"github.com/sahassauarv/tax-annotation/formatter"
	"github.com/sahassauarv/tax-annotation/parser"
	"github.com/sahassauarv/tax-annotation/render"
	"github.com/sahassauarv/tax-annotation/validator"
)

func main() {
	ctx := context.Background()

	p := parser.New()
	form, err := p.ParseFormFromFile(ctx, "fixtures/w2.json")
	if err != nil {
		log.Fatal("Failed to parse form:", err)
	}

	data := map[string]interface{}{
		"employee": map[string]interface{}{
			"ssn":     "123-45-6789",
			"name":    "John Doe",
			"address": "123 Main Street",
			"city":    "Springfield",
			"state":   "IL",
			"zip":     "62704",
		},
		"employer": map[string]interface{}{
			"ein":  "12-3456789",
			"name": "Acme Corporation",
		},
		"wages": map[string]interface{}{
			"box1":            75000.00,
			"box2":            12500.00,
			"box3":            75000.00,
			"box4":            4650.00,
			"box5":            75000.00,
			"box6":            1087.50,
			"box12Code":       "D",
			"box12Amount":     19500.00,
			"box13Statutory":  false,
			"box13Retirement": true,
			"box13ThirdParty": false,
		},
		"state": map[string]interface{}{
			"employerState":   "IL",
			"employerStateId": "123-456-789",
			"wages":           75000.00,
			"incomeTax":       3750.00,
		},
		"local": map[string]interface{}{
			"wages":     75000.00,
			"incomeTax": 750.00,
		},
	}

	resolver := parser.NewPathResolver(data)
	fmtr := formatter.New()
	vld := validator.New()

	renderer, err := render.NewRenderer(ctx, resolver, fmtr, vld)
	if err != nil {
		log.Fatal("Failed to create renderer:", err)
	}

	result, err := renderer.RenderForm(ctx, form)
	if err != nil {
		log.Fatal("Failed to render form:", err)
	}

	result.PrintSummary()
}
