package generate

var JAVA_TYPE_MAP = map[string]string{
	"string":   "String",
	"boolean":  "boolean",
	"date":     "Date",
	"dateTime": "Date",
	"double":   "double",
}

var CS_TYPE_MAP = map[string]string{
	"string":   "string",
	"boolean":  "bool",
	"date":     "DateTime",
	"dateTime": "DateTime",
	"double":   "double",
}

func getJavaType(t string) string {

	value, ok := JAVA_TYPE_MAP[t]
	if ok {
		return value
	} else {
		return t
	}
}

func getCSharpType(t string) string {

	value, ok := CS_TYPE_MAP[t]
	if ok {
		return value
	} else {
		return t
	}
}
