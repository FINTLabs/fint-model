package generate

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/FINTLabs/fint-model/common/types"
	"github.com/FINTLabs/fint-model/generate/cs"
	"github.com/FINTLabs/fint-model/generate/java"
)

var funcMap = template.FuncMap{
	//"add": func(i int, ii int) int { return i + ii },
	"sub": func(i int, ii int) int { return i - ii },
	"resourcePkg": func(s string) string {
		return strings.Replace(s, ".model.", ".model.resource.", -1)
	},
	"resource": func(resources []types.Attribute, s string) string {
		for _, a := range resources {
			if strings.HasSuffix(s, a.Type) {
				return strings.Replace(s, ".model.", ".model.resource.", -1) + "Resource"
			}
		}
		return s
	},
	"extends": func(isResource bool, extends string, s string) string {
		if isResource && strings.HasSuffix(s, extends) {
			return strings.Replace(s, ".model.", ".model.resource.", -1) + "Resource"
		}
		return s
	},
	"listFilt": func(list bool, s string) string {
		if list {
			return fmt.Sprintf("List<%s>", s)
		}
		return s
	},
	"javaType": types.GetJavaType,
	"validFilt": func(t string, s string) string {
		_, ok := types.JAVA_TYPE_MAP[t]
		if ok {
			return s
		}
		return `@Valid ` + s
	},
	"csType": func(s string, opt bool) string {
		typ := types.GetCSType(s)
		if opt && types.IsValueType(typ) {
			return typ + "?"
		}
		return typ
	},
	"upperCase":      func(s string) string { return strings.ToUpper(s) },
	"upperCaseFirst": func(s string) string { return strings.Title(s) },
	"getter":         func(s string) string { return "get" + strings.Title(s) + "()" },
	"baseType":       func(s string) string { return strings.Replace(s, "Resource", "", -1) },
	"assignResource": func(typ string, att string) string {
		if strings.HasPrefix(typ, "List<") {
			inner := strings.TrimSuffix(strings.TrimPrefix(typ, "List<"), ">")
			return fmt.Sprintf("%s.stream().map(%s::create).collect(Collectors.toList())", att, inner)
		}
		return fmt.Sprintf("%s.create(%s)", typ, att)
	},
	"listAdder": func(typ string) string {
		if strings.HasPrefix(typ, "List<") {
			return "All"
		}
		return ""
	},
	"resolveMultiplicity": func(multiplicity string) string {
		switch multiplicity {
		case "1":
			return "ONE_TO_ONE"
		case "0..*":
			return "NONE_TO_MANY"
		case "0..1":
			return "NONE_TO_ONE"
		case "1..*":
			return "ONE_TO_MANY"
		default:
			return multiplicity
		}
	},
	"modelRename": func(s string) string {
		if s == "FintMainObject" {
			return "FintModelObject"
		}
		return s
	},
	"resourceRename": func(s string) string {
		if s == "FintMainObject" {
			return "FintResource"
		}
		return s
	},
	"resourcePackageRename": func(s string) string {
		if s == "FintMainObject" {
			return "resource.FintResource"
		}
		return s
	},
	"implementInterfaces": func(s string) string {
		if s == "FintResource" {
			return s
		}
		return s + ", FintLinks"
	},
	"superResource": func(s string) string {
		if s == "FintMainObject" {
			return "FintResource"
		}
		return "FintLinks"
	},
}

func GetJavaResourceClass(c *types.Class) string {
	return getClass(c, java.RESOURCE_CLASS_TEMPLATE)
}

func GetJavaResourcesClass(c *types.Class) string {
	return getClass(c, java.RESOURCES_CLASS_TEMPLATE)
}

func GetJavaClass(c *types.Class) string {
	return getClass(c, java.CLASS_TEMPLATE)
}

func GetJavaActionEnum(a types.Action) string {
	return getActionEnum(a, java.ACTION_ENUM_TEMPLATE)
}

func GetCSResourceClass(c *types.Class) string {
	return getClass(c, cs.RESOURCE_CLASS_TEMPLATE)
}

func GetCSResourcesClass(c *types.Class) string {
	return getClass(c, cs.RESOURCES_CLASS_TEMPLATE)
}

func GetCSActionEnum(a types.Action) string {
	return getActionEnum(a, cs.ACTION_ENUM_TEMPLATE)
}

func GetCSClass(c *types.Class) string {
	return getClass(c, cs.CLASS_TEMPLATE)
}

func getClass(c *types.Class, t string) string {
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, c)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func getActionEnum(a types.Action, t string) string {
	tpl := template.New("action_enum").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, a)
	if err != nil {
		panic(err)
	}
	return b.String()
}
