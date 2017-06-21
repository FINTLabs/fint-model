package parser

type Class struct {
	Name             string
	Abstract         bool
	Extends          string
	Package          string
	Namespace        string
	DocumentationUrl string
	Attributes       []Attribute
	Relations        []string
	Identifiable     bool
}
