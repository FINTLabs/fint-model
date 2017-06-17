package generator

type Class struct {
	Name       string
	Abstract   bool
	Extends    string
	Package    string
	Attributes []Attribute
	Relations  []string
}
