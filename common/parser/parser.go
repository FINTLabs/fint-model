package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FINTprosjektet/fint-model/common/config"
	"github.com/FINTprosjektet/fint-model/common/document"
	"github.com/FINTprosjektet/fint-model/common/types"
	"github.com/FINTprosjektet/fint-model/common/utils"
	"github.com/antchfx/xquery/xml"
)

func GetClasses(owner string, repo string, tag string, filename string, force bool) ([]types.Class, map[string]types.Import, map[string][]types.Class, map[string][]types.Class) {
	doc := document.Get(owner, repo, tag, filename, force)

	var classes []types.Class
	// TODO BUG: packageMap and classMap fail for classes with the same name!
	packageMap := make(map[string]types.Import)
	classMap := make(map[string]types.Class)
	javaPackageClassMap := make(map[string][]types.Class)
	csPackageClassMap := make(map[string][]types.Class)

	classElements := xmlquery.Find(doc, "//element[@type='Class']")
	for _, c := range classElements {

		properties := c.SelectElement("properties")
		var class types.Class

		class.Name = replaceNO(c.SelectAttr("name"))
		class.Abstract = toBool(properties.SelectAttr("isAbstract"))
		class.Extends = getExtends(doc, c)
		class.Attributes = getAttributes(c)
		class.Relations = getRelations(doc, c)
		class.Package = getPackagePath(c, doc)
		class.Namespace = getNamespacePath(c, doc)
		class.Identifiable = identifiable(class.Attributes)
		class.Stereotype = properties.SelectAttr("stereotype")
		class.Documentation = properties.SelectAttr("documentation")
		class.GitTag = tag

		if len(class.Stereotype) == 0 {
			if class.Abstract {
				class.Stereotype = "abstrakt"
			} else {
				class.Stereotype = "datatype"
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

	for i := range classes {
		classes[i].Imports = getImports(classes[i], packageMap)
		classes[i].Using = getUsing(classes[i], packageMap)
		classes[i].Identifiable = identifiableFromExtends(classes[i], classMap)
		javaPackageClassMap[classes[i].Package] = append(javaPackageClassMap[classes[i].Package], classes[i])
		csPackageClassMap[classes[i].Namespace] = append(csPackageClassMap[classes[i].Namespace], classes[i])
	}

	return classes, packageMap, javaPackageClassMap, csPackageClassMap
}

func identifiableFromExtends(class types.Class, classMap map[string]types.Class) bool {
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

func getImports(c types.Class, imports map[string]types.Import) []string {

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

func getUsing(c types.Class, imports map[string]types.Import) []string {

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
	classPkg := getPackage(c)
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
		attrib.List = strings.Compare(a.SelectElement("bounds").SelectAttr("upper"), "*") == 0
		attrib.Optional = !attrib.List && strings.Compare(a.SelectElement("bounds").SelectAttr("lower"), "0") == 0
		attrib.Type = replaceNO(a.SelectElement("properties").SelectAttr("type"))

		attributes = append(attributes, attrib)
	}

	return attributes
}

func getRelations(doc *xmlquery.Node, c *xmlquery.Node) []string {
	var assocs []string
	isAbstract := toBool(c.SelectElement("properties").SelectAttr("isAbstract"))
	if !isAbstract {
		assocs = getAssociations(doc, c)
		assocs = append(assocs, getRecursivelyAssociationsFromExtends(doc, c)...)
	}

	return assocs
}

// TODO: This is actually iterative and not recursive, and works only for linear inheritance.
// TODO: Possible to combine with getAssociationsFromExtends?
func getRecursivelyAssociationsFromExtends(doc *xmlquery.Node, c *xmlquery.Node) []string {
	var assocs []string
	extAssocs, extends := getAssociationsFromExtends(doc, c)
	assocs = append(assocs, extAssocs...)
	for {
		if extends == nil {
			break
		}
		extAssocs, extends = getAssociationsFromExtends(doc, extends)
		assocs = append(assocs, extAssocs...)
	}
	return assocs
}

func getAssociationsFromExtends(doc *xmlquery.Node, c *xmlquery.Node) ([]string, *xmlquery.Node) {
	var assocs []string
	extends := xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target", c.SelectAttr("idref")))
	if len(extends) == 1 {
		assocs = append(assocs, getAssociations(doc, extends[0])...)
		return assocs, extends[0]
	}

	return assocs, nil
}

func getAssociations(doc *xmlquery.Node, c *xmlquery.Node) []string {
	var assocs []string
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source[@idref='%s']/../target/role", c.SelectAttr("idref"))) {
		if len(rr.SelectAttr("name")) > 0 {
			assocs = append(assocs, strings.ToUpper(replaceNO(rr.SelectAttr("name"))))
		}
	}
	for _, rl := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target[@idref='%s']/../source/role", c.SelectAttr("idref"))) {
		if len(rl.SelectAttr("name")) > 0 {
			assocs = append(assocs, strings.ToUpper(replaceNO(rl.SelectAttr("name"))))
		}
	}
	return assocs
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
