package packages

import (
	"github.com/FINTprosjektet/fint-model/generator"
	"github.com/FINTprosjektet/fint-model/common/utils"
)

func DistinctList(tag string) []string {
	classes := generator.GetClasses(tag)

	var p []string
	for _, c := range classes {
		p = append(p, c.Package)
	}

	return utils.Distinct(p)
}
