package parser

import (
	"github.com/antchfx/xquery/xml"
	"fmt"
	"strings"
	"strconv"
	"github.com/FINTprosjektet/fint-model/common/document"
	"github.com/FINTprosjektet/fint-model/common/utils"
	"github.com/FINTprosjektet/fint-model/common/config"
	"github.com/FINTprosjektet/fint-model/common/types"
)

func GetClasses(tag string, force bool) ([]types.Class, map[string]types.Import) {
	doc := document.Get(tag, force)

	var classes []types.Class
	packageMap := make(map[string]types.Import)

	for _, c := range xmlquery.Find(doc, "//element[@type='Class']") {

		var class types.Class

		class.Name = replaceNO(c.SelectAttr("name"))
		class.Abstract = toBool(c.SelectElement("properties").SelectAttr("isAbstract"))
		class.Extends = getExtends(doc, c)
		class.Attributes = getAttributes(c)
		class.Relations = getAssociations(doc, c)
		class.Package = getPackagePath(c, doc)
		class.Namespace = getNamespacePath(c, doc)
		class.Identifiable = identifiable(class.Attributes)

		imp := types.Import{
			Java:   fmt.Sprintf("%s.%s", class.Package, class.Name),
			CSharp: class.Namespace,
		}
		packageMap[class.Name] = imp

		class.Imports = getImports(class, packageMap)
		class.Using = getUsing(class, packageMap)

		classes = append(classes, class)
	}
	return classes, packageMap
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
	for _, value := range attribs {
		if imports[value.Type].Java != c.Package && len(imports[value.Type].Java) > 0 {
			//imp := fmt.Sprintf("import %s;", imports[value.Type].Java)
			imps = append(imps, imports[value.Type].Java)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].Java)
	}

	//return strings.Join(utils.Distinct(imps), "\n")
	return utils.Distinct(imps)
}

func getUsing(c types.Class, imports map[string]types.Import) []string {

	attribs := c.Attributes
	var imps []string
	for _, value := range attribs {
		if imports[value.Type].CSharp != c.Package && len(imports[value.Type].CSharp) > 0 {
			//imp := fmt.Sprintf("using %s;", imports[value.Type].CSharp)
			imps = append(imps, imports[value.Type].CSharp)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].CSharp)
	}

	//return strings.Join(utils.Distinct(imps), "\n")
	return utils.Distinct(imps)

}

func getPackagePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c)
	pkgs = append(pkgs, getNameLower(classPkg, doc))

	// TODO: This needs to be done much better
	parentPkg = getParentPackage(classPkg, doc)
	parentPkg2 := getParentPackage(parentPkg, doc)
	parentPkg3 := getParentPackage(parentPkg2, doc)

	pkgs = append(pkgs, getNameLower(parentPkg, doc), getNameLower(parentPkg2, doc), getNameLower(parentPkg3, doc))

	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return fmt.Sprintf("%s.%s", config.JAVA_PACKAGE_BASE, strings.Join(pkgs, "."))

}

func getNamespacePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c)
	pkgs = append(pkgs, getName(classPkg, doc))

	// TODO: This needs to be done much better
	parentPkg = getParentPackage(classPkg, doc)
	parentPkg2 := getParentPackage(parentPkg, doc)
	parentPkg3 := getParentPackage(parentPkg2, doc)

	pkgs = append(pkgs, getName(parentPkg, doc), getName(parentPkg2, doc), getName(parentPkg3, doc))

	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return fmt.Sprintf("%s.%s", config.NET_NAMESPACE_BASE, strings.Join(pkgs, "."))

}

func getName(idref string, doc *xmlquery.Node) string {
	name := ""
	if len(idref) > 0 {
		xpath := fmt.Sprintf("//element[@idref='%s']", idref)
		parent := xmlquery.Find(doc, xpath)

		name = parent[0].SelectAttr("name")
		name = strings.Replace(name, "FINT", "", -1)
		name = strings.Replace(name, "Model", "", -1)
	}
	return strings.Replace(name, " ", "", -1)
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
		fmt.Printf("Could not find any elements with idref %s\n", idref)
		return ""
	}

	return parent[0].SelectElement("model").SelectAttr("package")
}

func getPackage(c *xmlquery.Node) string {
	return c.SelectElement("model").SelectAttr("package")
}

func getExtends(doc *xmlquery.Node, c *xmlquery.Node) string {

	var extends []string
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source/model[@name='%s']/../../target/model[@name]", c.SelectAttr("name"))) {
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
		attrib.Type = a.SelectElement("properties").SelectAttr("type")

		attributes = append(attributes, attrib)
	}

	return attributes
}

func getAssociations(doc *xmlquery.Node, c *xmlquery.Node) []string {
	var assocs []string
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source/model[@name='%s']/../../target/role", c.SelectAttr("name"))) {
		if len(rr.SelectAttr("name")) > 0 {
			assocs = append(assocs, strings.ToUpper(replaceNO(rr.SelectAttr("name"))))
		}
	}
	for _, rl := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target/model[@name='%s']/../../source/role", c.SelectAttr("name"))) {
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

	return r
}

func toBool(s string) bool {
	b, err := strconv.ParseBool(s)

	if err != nil {
		return false
	}

	return b
}
