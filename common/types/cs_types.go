package types

var CS_TYPE_MAP = map[string]string{
	"string":   "string",
	"boolean":  "bool",
	"date":     "DateTime",
	"dateTime": "DateTime",
	"double":   "double",
}

func GetCSType(t string, list bool) string {

	value, ok := CS_TYPE_MAP[t]
	if ok {
		return getType(list, value)
	} else {
		return t
	}
}
