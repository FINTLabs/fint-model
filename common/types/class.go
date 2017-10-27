package types

type Class struct {
	Name             string
	Abstract         bool
	Extends          string
	Package          string
	Imports          []string
	Namespace        string
	Using            []string
	Documentation    string
	Attributes       []Attribute
	Relations        []string
	Identifiable     bool
	GitTag           string
	Stereotype       string
}
