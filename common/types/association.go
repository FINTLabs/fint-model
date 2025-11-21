package types

type Association struct {
	Name         string
	Target       string
	Source       string
	Package      string
	Deprecated   bool
	Multiplicity string
}
