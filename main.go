package main

import (
	"encoding/json"
	"fmt"
	"log"

	annotation "github.com/sahassauarv/tax-annotation/annotation"
)

func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }

func main() {
	form := annotation.Form{
		ID:      "W-2",
		Name:    "Wage and Tax Statement",
		Version: "2025",
		Pages: []annotation.Page{
			{
				Number: 1,
				Label:  "Employee & Wages",
				Annotations: []annotation.Annotation{
					{
						ID:        "employee_name",
						Label:     "Employee Name",
						FieldType: annotation.FieldTypeText,
						Value:     annotation.ValueRef{Path: "employee.name"},
						Position:  annotation.Position{X: 100, Y: 100, Width: 200, Height: 20},
						Validation: &annotation.Validation{
							Type:      annotation.DataTypeString,
							MaxLength: intPtr(60),
						},
					},
					{
						ID:        "wages",
						Label:     "Wages, Tips, Compensation",
						FieldType: annotation.FieldTypeNumber,
						Value:     annotation.ValueRef{Path: "wages.box1"},
						Position:  annotation.Position{X: 400, Y: 150, Width: 100, Height: 20},
						Format: &annotation.Format{
							Type:      annotation.FormatCurrency,
							Decimals:  intPtr(2),
							Alignment: annotation.AlignRight,
						},
						Validation: &annotation.Validation{
							Type: annotation.DataTypeNumber,
							Min:  floatPtr(0),
						},
					},
				},
			},
			{
				Number: 2,
				Label:  "State Info",
				Annotations: []annotation.Annotation{
					{
						ID:        "state_wages",
						Label:     "State Wages",
						FieldType: annotation.FieldTypeNumber,
						Value:     annotation.ValueRef{Path: "state.wages"},
						Position:  annotation.Position{X: 400, Y: 100, Width: 100, Height: 20},
						Format: &annotation.Format{
							Type:      annotation.FormatCurrency,
							Alignment: annotation.AlignRight,
						},
					},
					{
						ID:        "is_remote_state",
						Label:     "Worked In Remote State",
						FieldType: annotation.FieldTypeCheckbox,
						Value:     annotation.ValueRef{Path: "state.remote"},
						Position:  annotation.Position{X: 100, Y: 140, Width: 12, Height: 12},
						Format:    &annotation.Format{Type: annotation.FormatBoolean},
						Validation: &annotation.Validation{
							Type: annotation.DataTypeBoolean,
						},
					},
					{
						ID:        "filed_on",
						Label:     "Filed On",
						FieldType: annotation.FieldTypeDate,
						Value:     annotation.ValueRef{Path: "state.filedOn"},
						Position:  annotation.Position{X: 400, Y: 180, Width: 100, Height: 20},
						Format: &annotation.Format{
							Type:    annotation.FormatDate,
							Pattern: "01/02/2006",
						},
						Validation: &annotation.Validation{
							Type: annotation.DataTypeDate},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(form, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
