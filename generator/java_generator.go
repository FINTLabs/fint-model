package generator

import (
	"fmt"
	"strings"
)

const JAVA_CLASS_TEMPLATE = "package no.fint.model.%s;\n\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = false)\n" +
	"public class %s  {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_ATTRIBUTE_TEMPLATE = "    private %s %s;\n"

const JAVA_RELATION_TEMPLATE = "    public enum Relasjonsnavn {\n%s\n    }"

const JAVA_EXTENDED_CLASS_TEMPLATE = "package no.fint.model.%s;\n\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public class %s extends %s {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_ABSTRACT_CLASS_TEMPLATE = "package no.fint.model.%s;\n\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
	"@NoArgsConstructor\n" +
	"@EqualsAndHashCode(callSuper = true)\n" +
	"public abstract class %s {\n\n" +
	"%s\n" + // relations
	"%s\n" + // properties
	"}"

const JAVA_ABSTRACT_EXTENDED_CLASS_TEMPLATE = "package no.fint.model.%s;\n\n" +
	"import lombok.AllArgsConstructor;\n" +
	"import lombok.Data;\n" +
	"import lombok.EqualsAndHashCode;\n" +
	"import lombok.NoArgsConstructor;\n\n" +
	"@Data\n" +
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

func GetJavaClass(c Class) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_CLASS_TEMPLATE, c.Package, c.Name, relations, attributes)
}

func GetExtendedJavaClass(c Class) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_EXTENDED_CLASS_TEMPLATE, c.Package, c.Name, c.Extends, relations, attributes)
}

func GetAbstractJavaClass(c Class) string {

	var attributes string
	for _, a := range c.Attributes {
		attributes += fmt.Sprintf(JAVA_ATTRIBUTE_TEMPLATE, getJavaType(a.Type), a.Name)
	}

	relations := ""
	if hasRelations(c.Relations) {

		relations = fmt.Sprintf(JAVA_RELATION_TEMPLATE, fmt.Sprintf("        %s", strings.Join(c.Relations, ",\n        ")))
	}

	return fmt.Sprintf(JAVA_ABSTRACT_CLASS_TEMPLATE, c.Package, c.Name, relations, attributes)
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
