package types

type Association struct {
	Name         string
	Target       string
	InverseName  string
	Package      string
	Deprecated   bool
	Multiplicity string
	IsSource     bool
}
