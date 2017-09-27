package types

var JAVA_TYPE_MAP = map[string]string{
	"string":   "String",
	"boolean":  "boolean",
	"date":     "Date",
	"dateTime": "Date",
	"double":   "double",
}


func GetJavaType(t string, list bool) string {

	value, ok := JAVA_TYPE_MAP[t]
	if ok {
		return getType(list, value)
	} else {
		return t
	}
}

