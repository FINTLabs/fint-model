package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FINTLabs/fint-model/common/config"
	"github.com/FINTLabs/fint-model/common/document"
	"github.com/FINTLabs/fint-model/common/types"
	"github.com/FINTLabs/fint-model/common/utils"
	xmlquery "github.com/antchfx/xquery/xml"
)

func GetClasses(owner string, repo string, tag string, filename string, force bool) ([]*types.Class, map[string]types.Import, map[string][]*types.Class, map[string][]*types.Class) {
	doc := document.Get(owner, repo, tag, filename, force)

	var classes []*types.Class
	// TODO BUG: packageMap and classMap fail for classes with the same name!
	packageMap := make(map[string]types.Import)
	classMap := make(map[string]*types.Class)
	javaPackageClassMap := make(map[string][]*types.Class)
	csPackageClassMap := make(map[string][]*types.Class)

	classElements := xmlquery.Find(doc, "//element[@type='Class']")
	for _, classElement := range classElements {

		properties := classElement.SelectElement("properties")
		class := new(types.Class)

		class.Name = replaceNO(classElement.SelectAttr("name"))
		class.Abstract = toBool(properties.SelectAttr("isAbstract"))
		class.Extends = getExtends(doc, classElement)
		class.Attributes = getAttributes(classElement)
		class.Relations = getAssociations(doc, classElement)
		class.ExtendsRelations = getExtendsAssociations(doc, classElement)
		class.Package = getPackagePath(classElement, doc)
		class.Namespace = getNamespacePath(classElement, doc)
		class.Identifiable = identifiable(class.Attributes)
		class.Stereotype = properties.SelectAttr("stereotype")
		class.Documentation = properties.SelectAttr("documentation")
		class.Deprecated = classElement.SelectElement("tags/tag[@name='DEPRECATED']") != nil
		class.GitTag = tag

		if len(class.Stereotype) == 0 {
			if class.Abstract {
				class.Stereotype = "abstrakt"
			}
		}

		imp := types.Import{
			Java:   fmt.Sprintf("%s.%s", class.Package, class.Name),
			CSharp: class.Namespace,
		}
		packageMap[class.Name] = imp

		classes = append(classes, class)
		classMap[class.Name] = class
	}

	packageMap["Date"] = types.Import{
		Java: "java.util.Date",
	}

	for _, class := range classes {
		class.Imports = getImports(class, packageMap)
		class.Using = getUsing(class, packageMap)
		class.Identifiable = identifiableFromExtends(class, classMap)
		class.Resource = isResource(class, classMap)
		javaPackageClassMap[class.Package] = append(javaPackageClassMap[class.Package], class)
		csPackageClassMap[class.Namespace] = append(csPackageClassMap[class.Namespace], class)
		if len(class.Stereotype) == 0 {
			if class.Identifiable {
				class.Stereotype = "hovedklasse"
			} else {
				class.Stereotype = "datatype"
			}
		}
	}

	for _, class := range classes {
		for _, a := range class.Attributes {
			if typ, found := classMap[a.Type]; found {
				if typ.Resource {
					class.Resources = append(class.Resources, a)
				}
			}
		}
	}

	for _, class := range classes {
		if len(class.Extends) > 0 {
			if typ, found := classMap[class.Extends]; found {
				class.ExtendsResource = typ.Resource || len(typ.Resources) > 0
			}
		}
	}

	return classes, packageMap, javaPackageClassMap, csPackageClassMap
}

func isResource(class *types.Class, classMap map[string]*types.Class) bool {
	if len(class.Relations) > 0 {
		return true
	}
	if len(class.Extends) > 0 {
		return isResource(classMap[class.Extends], classMap)
	}
	return false
}

func identifiableFromExtends(class *types.Class, classMap map[string]*types.Class) bool {
	if class.Identifiable {
		return true
	}
	if len(class.Extends) > 0 {
		return identifiableFromExtends(classMap[class.Extends], classMap)
	}
	return false
}

func identifiable(attribs []types.Attribute) bool {

	for _, value := range attribs {
		if value.Type == "Identifikator" {
			return true
		}
	}

	return false

}

func getImports(c *types.Class, imports map[string]types.Import) []string {

	attribs := c.Attributes
	var imps []string
	for _, att := range attribs {
		javaType := types.GetJavaType(att.Type)
		if imports[javaType].Java != c.Package && len(javaType) > 0 {
			imps = append(imps, imports[javaType].Java)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].Java)
	}

	return utils.Distinct(utils.TrimArray(imps))
}

func getUsing(c *types.Class, imports map[string]types.Import) []string {

	attribs := c.Attributes
	var imps []string
	for _, att := range attribs {
		csType := types.GetCSType(att.Type)
		if imports[csType].CSharp != c.Package && len(imports[csType].CSharp) > 0 {
			imps = append(imps, imports[csType].CSharp)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].CSharp)
	}

	return utils.Distinct(utils.TrimArray(imps))

}

func getPackagePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c) // Gets Package ID
	pkgs = append(pkgs, getNameLower(classPkg, doc))

	parentPkg = getParentPackage(classPkg, doc)
	for parentPkg != "" {
		pkgs = append(pkgs, getNameLower(parentPkg, doc))
		parentPkg = getParentPackage(parentPkg, doc)
	}
	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return replaceNO(fmt.Sprintf("%s.%s", config.JAVA_PACKAGE_BASE, strings.Join(pkgs, ".")))

}

func getNamespacePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c)
	pkgs = append(pkgs, getName(classPkg, doc))

	parentPkg = getParentPackage(classPkg, doc)
	for parentPkg != "" {
		pkgs = append(pkgs, getName(parentPkg, doc))
		parentPkg = getParentPackage(parentPkg, doc)
	}
	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return replaceNO(fmt.Sprintf("%s.%s", config.NET_NAMESPACE_BASE, strings.Join(pkgs, ".")))

}

func getName(idref string, doc *xmlquery.Node) string {
	name := ""
	if len(idref) > 0 {
		xpath := fmt.Sprintf("//element[@idref='%s']", idref)
		parent := xmlquery.Find(doc, xpath)

		name = parent[0].SelectAttr("name")
		name = excludeName(name)
	}
	return strings.Replace(name, " ", "", -1)
}

func excludeName(name string) string {
	if name == "FINT" {
		name = strings.Replace(name, "FINT", "", -1)
	}
	if name == "Model" {
		name = strings.Replace(name, "Model", "", -1)
	}
	return name
}

func getNameLower(idref string, doc *xmlquery.Node) string {

	return strings.ToLower(getName(idref, doc))
}

func getParentPackage(idref string, doc *xmlquery.Node) string {
	xpath := fmt.Sprintf("//element[@idref='%s']", idref)

	parent := xmlquery.Find(doc, xpath)

	if len(parent) > 1 {
		fmt.Printf("More than one element with idref %s\n", idref)
		return ""
	}
	if len(parent) < 1 {
		return ""
	}

	model := parent[0].SelectElement("model")
	if model == nil {
		return ""
	}

	return model.SelectAttr("package")
}

func getPackage(c *xmlquery.Node) string {
	return c.SelectElement("model").SelectAttr("package")
}

func getExtends(doc *xmlquery.Node, c *xmlquery.Node) string {

	var extends []string
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target/model[@name]", c.SelectAttr("idref"))) {
		if len(rr.SelectAttr("name")) > 0 {
			extends = append(extends, replaceNO(rr.SelectAttr("name")))
		}
	}

	if len(extends) == 1 {
		return extends[0]
	}

	return ""
}

func getAttributes(c *xmlquery.Node) []types.Attribute {
	var attributes []types.Attribute
	for _, a := range xmlquery.Find(c, "//attributes/attribute") {

		attrib := types.Attribute{}
		attrib.Name = replaceNO(a.SelectAttr("name"))
		attrib.Deprecated = a.SelectElement("tags/tag[@name='DEPRECATED']") != nil
		attrib.List = strings.Compare(a.SelectElement("bounds").SelectAttr("upper"), "*") == 0
		attrib.Optional = strings.Compare(a.SelectElement("bounds").SelectAttr("lower"), "0") == 0
		attrib.Type = replaceNO(a.SelectElement("properties").SelectAttr("type"))

		attributes = append(attributes, attrib)
	}

	return attributes
}

func getAssociations(doc *xmlquery.Node, c *xmlquery.Node) []types.Association {
	var assocs []types.Association

	queries := [...]string{
		fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source[@idref='%s']/../target/role", c.SelectAttr("idref")),
		fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target[@idref='%s']/../source/role", c.SelectAttr("idref")),
	}

	for _, query := range queries {
		for _, relationElement := range xmlquery.Find(doc, query) {
			if len(relationElement.SelectAttr("name")) > 0 {
				classElement := findClassElementByID(doc, relationElement.Parent.SelectAttr("idref"))

				assoc := types.Association{}
				assoc.Name = replaceNO(relationElement.SelectAttr("name"))
				assoc.Target = replaceNO(relationElement.SelectElement("../model").SelectAttr("name"))
				assoc.Multiplicity = relationElement.SelectElement("../type").SelectAttr("multiplicity")
				assoc.Deprecated = relationElement.SelectElement("../../tags/tag[@name='DEPRECATED']") != nil
				assoc.Package = getPackagePath(classElement, doc)

				assocs = append(assocs, assoc)
			}
		}
	}
	return assocs
}

func findClassElementByID(doc *xmlquery.Node, id string) *xmlquery.Node {
	query := fmt.Sprintf("//element[@type='Class'][@idref='%s']", id)
	return xmlquery.FindOne(doc, query)
}

func getExtendsAssociations(doc *xmlquery.Node, c *xmlquery.Node) bool {

	generalizationTargets := xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target[@idref]", c.SelectAttr("idref")))

	if len(generalizationTargets) == 1 && len(generalizationTargets[0].SelectAttr("idref")) > 0 {
		idref := generalizationTargets[0].SelectAttr("idref")
		associationQuery := [...]string{
			fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source[@idref='%s']/../target/role", idref),
			fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target[@idref='%s']/../source/role", idref),
		}

		for _, query := range associationQuery {
			for _, r := range xmlquery.Find(doc, query) {
				if len(r.SelectAttr("name")) > 0 {
					return true
				}
			}
		}

	}

	return false
}

func replaceNO(s string) string {
	r := strings.Replace(s, "æ", "a", -1)
	r = strings.Replace(r, "ø", "o", -1)
	r = strings.Replace(r, "å", "a", -1)
	r = strings.Replace(r, "Æ", "A", -1)
	r = strings.Replace(r, "Ø", "O", -1)
	r = strings.Replace(r, "Å", "A", -1)
	return r
}

func toBool(s string) bool {
	b, err := strconv.ParseBool(s)

	if err != nil {
		return false
	}

	return b
}
