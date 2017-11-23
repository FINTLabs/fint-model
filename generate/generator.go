package generate

import (
	"bytes"
	"github.com/FINTprosjektet/fint-model/common/types"
	"text/template"
	"github.com/FINTprosjektet/fint-model/generate/java"
	"github.com/FINTprosjektet/fint-model/generate/cs"
	"strings"
	"fmt"
)

var funcMap = template.FuncMap{
	//"add": func(i int, ii int) int { return i + ii },
	"sub":      func(i int, ii int) int { return i - ii },
	"listFilt": func(list bool, s string) string { if (list) { return fmt.Sprintf("List<%s>", s); } else { return s; } },
	"javaType": types.GetJavaType,
	"csType":   func(s string, opt bool) string { typ := types.GetCSType(s); if (opt && typ != "string") { return typ + "?"; } else { return typ; } },
	"upperCaseFirst": func(s string) string { return strings.Title(s) },
}

func GetJavaClass(c types.Class) string {
	return getClass(c, java.CLASS_TEMPLATE)
}

func GetJavaActionEnum(a types.Action) string  {
	return getActionEnum(a, java.ACTION_ENUM_TEMPLATE)
}

func GetCSActionEnum(a types.Action) string  {
	return getActionEnum(a, cs.ACTION_ENUM_TEMPLATE)
}

func GetCSClass(c types.Class) string {
	return getClass(c, cs.CLASS_TEMPLATE)
}

func getClass(c types.Class, t string) string {
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
