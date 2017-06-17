package packages

import (
	"github.com/FINTprosjektet/fint-model/generator"
	"github.com/FINTprosjektet/fint-model/common/utils"
)

func DistinctPackageList(tag string) []string {
	classes := generator.GetClasses(tag)

	var p []string
	for _, c := range classes {
		p = append(p, c.Package)
	}

	return utils.Distinct(p)
}

func DistinctNamespaceList(tag string) []string {
	classes := generator.GetClasses(tag)

	var p []string
	for _, c := range classes {
		p = append(p, c.Namespace)
	}

	return utils.Distinct(p)
}