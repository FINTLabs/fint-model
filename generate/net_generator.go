package generate

import (
	"fmt"
	"strings"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/common/utils"
)

const CSHARP_CLASS_TEMPLATE = 	"%s\n" +
	"namespace %s\n" +
	"{\n" +
	"    public class %s\n" +
	"    {\n" +
	"%s\n" + // relations
	"%s\n\n" + // attributes
	"    }\n" +
	"}"

const CSHARP_EXTENDED_CLASS_TEMPLATE = "%s\n" +
	"namespace %s\n" +
	"{\n" +
	"    public class %s : %s\n" +
	"    {\n" +
	"%s\n" + // relations
	"%s\n\n" + // attributes
	"    }\n" +
	"}"

const CSHARP_ABSTRACT_CLASS_TEMPLATE = "%s\n" +
	"namespace %s\n" +
	"{\n" +
	"    public abstract class %s\n" +
	"    {\n" +
	"%s\n" + // relations
	"%s\n\n" + // attributes
	"    }\n" +
	"}"

const CSHARP_EXTENDED_ABSTRACT_CLASS_TEMPLATE = "%s\n" +
	"namespace %s\n" +
	"{\n" +
	"    public abstract class %s : %s\n" +
	"    {\n" +
	"%s\n" + // relations
	"%s\n\n" + // attributes
	"    }\n" +
	"}"

const CSHARP_ATTRIBUTE_TEMPLATE = "        public %s %s { get; set; }\n"

const CSHARP_RELATION_TEMPLATE = "        public enum Relasjonsnavn\n        {\n%s\n        }"

var CSHARP_TYPE_MAP = map[string]string{
	"string":   "string",
	"boolean":  "bool",
	"date":     "DateTime",
	"dateTime": "DateTime",
	"double":   "double",
}

func GetCSharpClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(CSHARP_ATTRIBUTE_TEMPLATE, getCSharpType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(CSHARP_RELATION_TEMPLATE, fmt.Sprintf("            %s", strings.Join(c.Relations, ",\n            ")))
	}

	return fmt.Sprintf(CSHARP_CLASS_TEMPLATE, getCSharpImports(c, impMap), c.Namespace, c.Name, relations, attributes)
}

func GetAbstractCSharpClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(CSHARP_ATTRIBUTE_TEMPLATE, getCSharpType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(CSHARP_RELATION_TEMPLATE, fmt.Sprintf("            %s", strings.Join(c.Relations, ",\n            ")))
	}

	return fmt.Sprintf(CSHARP_ABSTRACT_CLASS_TEMPLATE, getCSharpImports(c, impMap), c.Namespace, c.Name, relations, attributes)
}

func GetExtendedCSharpClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(CSHARP_ATTRIBUTE_TEMPLATE, getCSharpType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(CSHARP_RELATION_TEMPLATE, fmt.Sprintf("            %s", strings.Join(c.Relations, ",\n            ")))
	}

	return fmt.Sprintf(CSHARP_EXTENDED_CLASS_TEMPLATE, getCSharpImports(c, impMap), c.Namespace, c.Name, c.Extends, relations, attributes)
}

func GetExtendedAbstractCSharpClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(CSHARP_ATTRIBUTE_TEMPLATE, getCSharpType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(CSHARP_RELATION_TEMPLATE, fmt.Sprintf("            %s", strings.Join(c.Relations, ",\n            ")))
	}

	return fmt.Sprintf(CSHARP_EXTENDED_ABSTRACT_CLASS_TEMPLATE, getCSharpImports(c, impMap), c.Namespace, c.Name, c.Extends, relations, attributes)
}

func getCSharpType(t string) string {

	value, ok := JAVA_TYPE_MAP[t]
	if ok {
		return value
	} else {
		return t
	}
}

func getCSharpImports(c parser.Class, imports map[string]parser.Import) string {

	attribs := c.Attributes
	var imps []string
	for _, value := range attribs {
		if imports[value.Type].CSharp != c.Package && len(imports[value.Type].CSharp) > 0 {
			imp := fmt.Sprintf("using %s;", imports[value.Type].CSharp)
			imps = append(imps, imp)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, fmt.Sprintf("using %s;", imports[c.Extends].CSharp))
	}

	return strings.Join(utils.Distinct(imps), "\n")
}
