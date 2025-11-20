package types

type Association struct {
	Name         string
	Source       string
	Target       string
	Package      string
	Deprecated   bool
	Multiplicity string
}
