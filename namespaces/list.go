package namespaces

import (
	"github.com/FINTprosjektet/fint-model/common/utils"
	"github.com/FINTprosjektet/fint-model/common/parser"
)



func DistinctNamespaceList(tag string, force bool) []string {
	classes, _ := parser.GetClasses(tag, force)

	var p []string
	for _, c := range classes {
		p = append(p, c.Namespace)
	}

	return utils.Distinct(p)
}
