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

type importCandidate struct {
	Package string
	Import  types.Import
}

func qualifiedKey(pkg string, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func makeImportForClass(c *types.Class) types.Import {
	return types.Import{
		Java:   qualifiedKey(c.Package, c.Name),
		CSharp: c.Namespace,
	}
}

func buildImportNameMap(imports map[string]types.Import) map[string][]importCandidate {
	result := make(map[string][]importCandidate)
	for key, imp := range imports {
		lastDot := strings.LastIndex(key, ".")
		if lastDot < 0 {
			continue
		}
		name := key[lastDot+1:]
		pkg := key[:lastDot]
		result[name] = append(result[name], importCandidate{Package: pkg, Import: imp})
	}
	return result
}

func pickPackage(packageContext string, candidatePackage string, currentBest string) (string, bool) {
	if candidatePackage != packageContext && !strings.HasPrefix(candidatePackage, packageContext+".") {
		return currentBest, false
	}
	if currentBest == "" || len(candidatePackage) < len(currentBest) {
		return candidatePackage, false
	}
	if len(candidatePackage) == len(currentBest) {
		return currentBest, true
	}
	return currentBest, false
}

func resolveClassCandidate(packageContext string, candidates []*types.Class) (*types.Class, bool) {
	var best *types.Class
	bestPackage := ""
	tie := false
	for _, c := range candidates {
		var t bool
		bestPackage, t = pickPackage(packageContext, c.Package, bestPackage)
		tie = tie || t
		if bestPackage == c.Package {
			best = c
		}
	}
	return best, best != nil && !tie
}

func resolveImportCandidate(packageContext string, candidates []importCandidate) (types.Import, bool) {
	var best types.Import
	bestPackage := ""
	tie := false
	for _, c := range candidates {
		var t bool
		bestPackage, t = pickPackage(packageContext, c.Package, bestPackage)
		tie = tie || t
		if bestPackage == c.Package {
			best = c.Import
		}
	}
	return best, len(bestPackage) > 0 && !tie
}

func GetClasses(owner string, repo string, tag string, filename string, force bool) ([]*types.Class, map[string]types.Import, map[string][]*types.Class, map[string][]*types.Class) {
	doc := document.Get(owner, repo, tag, filename, force)

	var classes []*types.Class
	packageMap := make(map[string]types.Import)
	classMap := make(map[string]*types.Class)
	classNameMap := make(map[string][]*types.Class)
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

		imp := makeImportForClass(class)
		key := qualifiedKey(class.Package, class.Name)
		packageMap[key] = imp

		classes = append(classes, class)
		classMap[key] = class
		classNameMap[class.Name] = append(classNameMap[class.Name], class)
	}

	for name, list := range classNameMap {
		if len(list) == 1 {
			imp := makeImportForClass(list[0])
			packageMap[name] = imp
			classMap[name] = list[0]
		}
	}

	packageMap["Date"] = types.Import{
		Java: "java.util.Date",
	}

	importNameMap := buildImportNameMap(packageMap)

	for _, class := range classes {
		class.Imports = getImports(class, packageMap, importNameMap)
		class.Using = getUsing(class, packageMap, importNameMap)
		class.Identifiable = identifiableFromExtends(class, classMap, classNameMap)
		class.ExtendsIdentifiable = extendsIdentifiable(class, classMap, classNameMap)
		class.Writable = isWritable(class.Attributes)
		class.Resource = isResource(class, classMap, classNameMap)
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
			if typ, found := findClass(a.Type, class.Package, classMap, classNameMap); found {
				if typ.Resource {
					class.Resources = append(class.Resources, a)
				}
			}
		}
	}

	for _, class := range classes {
		if len(class.Extends) > 0 {
			if typ, found := findClass(class.Extends, class.Package, classMap, classNameMap); found {
				class.ExtendsResource = typ.Resource || len(typ.Resources) > 0
			}
		}
	}

	return classes, packageMap, javaPackageClassMap, csPackageClassMap
}

func isWritable(attribs []types.Attribute) bool {

	for _, value := range attribs {
		if value.Writable {
			return true
		}
	}
	return false
}

func findClass(className string, packageContext string, classMap map[string]*types.Class, classNameMap map[string][]*types.Class) (*types.Class, bool) {
	if class, found := classMap[className]; found {
		return class, true
	}
	if class, found := classMap[qualifiedKey(packageContext, className)]; found {
		return class, true
	}
	candidates := classNameMap[className]
	if len(candidates) == 1 {
		return candidates[0], true
	}
	return resolveClassCandidate(packageContext, candidates)
}

func isResource(class *types.Class, classMap map[string]*types.Class, classNameMap map[string][]*types.Class) bool {
	if len(class.Relations) > 0 {
		return true
	}
	if len(class.Extends) > 0 {
		if extendedClass, found := findClass(class.Extends, class.Package, classMap, classNameMap); found {
			return isResource(extendedClass, classMap, classNameMap)
		}
	}
	return false
}

func identifiableFromExtends(class *types.Class, classMap map[string]*types.Class, classNameMap map[string][]*types.Class) bool {
	if class.Identifiable {
		return true
	}
	if len(class.Extends) > 0 {
		if extendedClass, found := findClass(class.Extends, class.Package, classMap, classNameMap); found {
			return identifiableFromExtends(extendedClass, classMap, classNameMap)
		}
	}
	return false
}

func extendsIdentifiable(class *types.Class, classMap map[string]*types.Class, classNameMap map[string][]*types.Class) bool {
	if len(class.Extends) > 0 {
		if extendedClass, found := findClass(class.Extends, class.Package, classMap, classNameMap); found {
			return identifiableFromExtends(extendedClass, classMap, classNameMap)
		}
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

func findImport(typeName string, packageContext string, imports map[string]types.Import, importNameMap map[string][]importCandidate) (types.Import, bool) {
	if imp, found := imports[typeName]; found {
		return imp, true
	}
	if imp, found := imports[qualifiedKey(packageContext, typeName)]; found {
		return imp, true
	}
	return resolveImportCandidate(packageContext, importNameMap[typeName])
}

func getImports(c *types.Class, imports map[string]types.Import, importNameMap map[string][]importCandidate) []string {

	attribs := c.Attributes
	self := qualifiedKey(c.Package, c.Name)
	var imps []string
	for _, att := range attribs {
		javaType := types.GetJavaType(att.Type)
		if len(javaType) > 0 {
			imp, found := findImport(javaType, c.Package, imports, importNameMap)
			if found && len(imp.Java) > 0 && imp.Java != self {
				imps = append(imps, imp.Java)
			}
		}
	}

	if len(c.Extends) > 0 {
		imp, found := findImport(c.Extends, c.Package, imports, importNameMap)
		if found && len(imp.Java) > 0 && imp.Java != self {
			imps = append(imps, imp.Java)
		}
	}

	return utils.Distinct(utils.TrimArray(imps))
}

func getUsing(c *types.Class, imports map[string]types.Import, importNameMap map[string][]importCandidate) []string {

	attribs := c.Attributes
	var imps []string
	for _, att := range attribs {
		csType := types.GetCSType(att.Type)
		if len(csType) > 0 {
			imp, found := findImport(csType, c.Package, imports, importNameMap)
			if found && len(imp.CSharp) > 0 && imp.CSharp != c.Namespace {
				imps = append(imps, imp.CSharp)
			}
		}
	}

	if len(c.Extends) > 0 {
		imp, found := findImport(c.Extends, c.Package, imports, importNameMap)
		if found && len(imp.CSharp) > 0 && imp.CSharp != c.Namespace {
			imps = append(imps, imp.CSharp)
		}
	}

	return utils.Distinct(utils.TrimArray(imps))

}

func getPackagePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkgId := getPackage(c)
	pkgs = append(pkgs, getNameLower(classPkgId, doc))

	parentPkg = getParentPackage(classPkgId, doc)
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
		attrib.Writable = a.SelectElement("stereotype").SelectAttr("stereotype") == "writable"

		attributes = append(attributes, attrib)
	}

	return attributes
}

func buildAssociationQueries(idref string) []types.AssociationQuery {
	return []types.AssociationQuery{
		{
			XPath: fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source[@idref='%s']/../target/role", idref),
			Role:  types.RoleSource,
		},
		{
			XPath: fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target[@idref='%s']/../source/role", idref),
			Role:  types.RoleTarget,
		},
	}
}

func getAssociations(doc *xmlquery.Node, c *xmlquery.Node) []types.Association {
	var assocs []types.Association

	classId := c.SelectAttr("idref")
	isParent := isExtendedByOthers(doc, classId)
	queries := buildAssociationQueries(c.SelectAttr("idref"))

	for _, q := range queries {
		for _, relationElement := range xmlquery.Find(doc, q.XPath) {
			if len(relationElement.SelectAttr("name")) == 0 {
				continue
			}

			assoc := buildAssociation(doc, relationElement, q.Role, isParent)
			assocs = append(assocs, assoc)
		}
	}
	return assocs
}

func buildAssociation(doc *xmlquery.Node, rel *xmlquery.Node, role types.AssociationRole, isParent bool) types.Association {
	targetId := rel.Parent.SelectAttr("idref")
	targetClassElement := findClassElementByID(doc, targetId)

	return types.Association{
		Name:         replaceNO(rel.SelectAttr("name")),
		Target:       replaceNO(rel.SelectElement("../model").SelectAttr("name")),
		Multiplicity: rel.SelectElement("../type").SelectAttr("multiplicity"),
		Package:      getPackagePath(targetClassElement, doc),
		Deprecated:   rel.SelectElement("../../tags/tag[@name='DEPRECATED']") != nil,
		InverseName:  getAssociationSource(rel, role, isParent),
	}
}

func getAssociationSource(rel *xmlquery.Node, role types.AssociationRole, isParent bool) string {
	direction := rel.SelectElement("../../properties").SelectAttr("direction")

	if direction != "Bi-Directional" || isParent {
		return ""
	}

	var sourceNode *xmlquery.Node
	if role == types.RoleSource {
		sourceNode = rel.SelectElement("../../source/role")
	} else {
		sourceNode = rel.SelectElement("../../target/role")
	}

	if sourceNode != nil {
		return replaceNO(sourceNode.SelectAttr("name"))
	}

	return ""
}

func isExtendedByOthers(doc *xmlquery.Node, classId string) bool {
	xpath := fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../target[@idref='%s']", classId)
	return xmlquery.FindOne(doc, xpath) != nil
}

func findClassElementByID(doc *xmlquery.Node, id string) *xmlquery.Node {
	query := fmt.Sprintf("//element[@type='Class'][@idref='%s']", id)
	return xmlquery.FindOne(doc, query)
}

func getInverseNameFromParentClass(doc *xmlquery.Node, sourceClassIdref string, targetClassIdref string) string {
	generalizationTargets := xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target[@idref]", targetClassIdref))

	if len(generalizationTargets) != 1 {
		return ""
	}

	parentIdref := generalizationTargets[0].SelectAttr("idref")
	if len(parentIdref) == 0 {
		return ""
	}

	queries := []string{
		fmt.Sprintf("//connectors/connector/properties[@ea_type='Association'][@direction='Bi-Directional']/../source[@idref='%s']/../target[@idref='%s']/../target/role", parentIdref, sourceClassIdref),
		fmt.Sprintf("//connectors/connector/properties[@ea_type='Association'][@direction='Bi-Directional']/../target[@idref='%s']/../source[@idref='%s']/../source/role", parentIdref, sourceClassIdref),
	}

	for _, query := range queries {
		for _, roleElement := range xmlquery.Find(doc, query) {
			roleName := roleElement.SelectAttr("name")
			if len(roleName) > 0 {
				return replaceNO(roleName)
			}
		}
	}

	return ""
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
