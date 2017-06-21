package generate

import (
	"fmt"
	"strings"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/common/utils"
)

const JAVA_CLASS_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = false)\n" +
	"public class %s  {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_CLASS_IDENTIFIABLE_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import no.fint.model.relation.Identifiable;\n" +
	"import com.fasterxml.jackson.annotation.JsonIgnore;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = false)\n" +
	"public class %s implements Identifiable  {\n\n" +
	"%s\n" + // relations
	"%s\n\n" + // properties
	"    @JsonIgnore\n" +
	"    @Override\n" +
	"    public String getId() {\n" +
	"        return this.get{fixme}().getIdentifikatorverdi();\n" +
	"    }\n" +
	"}"

const JAVA_ATTRIBUTE_TEMPLATE = "    private %s %s;\n"

const JAVA_RELATION_TEMPLATE = "    public enum Relasjonsnavn {\n%s\n    }"

const JAVA_EXTENDED_CLASS_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public class %s extends %s {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_EXTENDED_CLASS_IDENTIFIABLE_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import no.fint.model.relation.Identifiable;\n" +
	"import com.fasterxml.jackson.annotation.JsonIgnore;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public class %s extends %s implements Identifiable {\n\n" +
	"%s\n" + // relations
	"%s\n\n" + // properties
	"    @JsonIgnore\n" +
	"    @Override\n" +
	"    public String getId() {\n" +
	"        return this.get{fixme}().getIdentifikatorverdi();\n" +
	"    }\n" +
	"}"

const JAVA_ABSTRACT_CLASS_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public abstract class %s {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_ABSTRACT_EXTENDED_CLASS_TEMPLATE = "packages %s;\n\n" +
	"%s\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@AllArgsConstructor\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public abstract class %s extends %s {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

var JAVA_TYPE_MAP = map[string]string{
	"string":   "String",
	"boolean":  "boolean",
	"date":     "Date",
	"dateTime": "Date",
	"double":   "double",
}

func GetJavaClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_CLASS_TEMPLATE, c.Package, getImports(c, impMap), c.Name, relations, attributes)
}

func getImports(c parser.Class, imports map[string]parser.Import) string {

	attribs := c.Attributes
	var imps []string
	for _, value := range attribs {
		if imports[value.Type].Java != c.Package && len(imports[value.Type].Java) > 0 {
			imp := fmt.Sprintf("import %s;", imports[value.Type].Java)
			imps = append(imps, imp)
		}
	}


	return strings.Join(utils.Distinct(imps), "\n")
}

func GetJavaClassIdentifiable(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_CLASS_IDENTIFIABLE_TEMPLATE, c.Package, getImports(c, impMap), c.Name, relations, attributes)
}

func GetExtendedJavaClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_EXTENDED_CLASS_TEMPLATE, c.Package, getImports(c, impMap), c.Name, c.Extends, relations, attributes)
}

func GetExtendedJavaClassIdentifiable(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_EXTENDED_CLASS_IDENTIFIABLE_TEMPLATE, c.Package, getImports(c, impMap), c.Name, c.Extends, relations, attributes)
}

func GetAbstractJavaClass(c parser.Class, impMap map[string]parser.Import) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_ABSTRACT_CLASS_TEMPLATE, c.Package, getImports(c, impMap), c.Name, relations, attributes)
}

func getJavaType(t string) string {

	value, ok := JAVA_TYPE_MAP[t]
	if ok {
		return value
	} else {
		return t
	}
}

func hasRelations(s []string) bool {

	if len(s) > 0 {
		return true
	}
	return false
}
