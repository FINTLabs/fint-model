package generate

import (
	"strings"
	"testing"

	"github.com/FINTLabs/fint-model/common/types"
)

func TestJavaAttributeValidation(t *testing.T) {
	tests := []struct {
		modelType string
		javaType  string
	}{
		{"date", "Date"},
		{"datetime", "Date"},
		{"dateTime", "Date"},
		{"DateTime", "Date"},
		{"Date", "Date"},
		{"Periode", "@Valid Periode"},
		{"Identifikator", "@Valid Identifikator"},
	}
	generators := []struct {
		name   string
		render func(*types.Class) string
	}{
		{"model", GetJavaClass},
		{"resource", GetJavaResourceClass},
	}

	for _, generator := range generators {
		for _, tt := range tests {
			t.Run(generator.name+"/"+tt.modelType, func(t *testing.T) {
				class := &types.Class{
					Name:       "Example",
					Package:    "no.novari.fint.model.felles",
					Stereotype: "hovedklasse",
					Attributes: []types.Attribute{
						{Name: "value", Type: tt.modelType},
						{Name: "values", Type: tt.modelType, List: true},
					},
				}
				output := generator.render(class)
				for _, want := range []string{
					"@NotNull\n    private " + tt.javaType + " value;",
					"@NotEmpty\n    private List<" + tt.javaType + "> values;",
				} {
					if !strings.Contains(output, want) {
						t.Errorf("generated class does not contain %q", want)
					}
				}
			})
		}
	}
}
