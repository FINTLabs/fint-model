package parser

import (
	"testing"

	"github.com/FINTLabs/fint-model/common/types"
)

func TestFindClass_PrefersPackageQualifiedAndDisambiguatesCollisions(t *testing.T) {
	classMap := map[string]*types.Class{}
	classNameMap := map[string][]*types.Class{}

	classA := &types.Class{Name: "Klasse", Package: "no.fint.skole"}
	classB := &types.Class{Name: "Klasse", Package: "no.fint.noark"}

	classMap["no.fint.skole.Klasse"] = classA
	classMap["no.fint.noark.Klasse"] = classB
	classNameMap["Klasse"] = []*types.Class{classA, classB}

	if got, found := findClass("no.fint.noark.Klasse", "no.fint.skole", classMap, classNameMap); !found || got != classB {
		t.Fatalf("expected fully qualified lookup to return noark Klasse")
	}

	if got, found := findClass("Klasse", "no.fint.skole", classMap, classNameMap); !found || got != classA {
		t.Fatalf("expected context lookup to return skole Klasse")
	}

	if got, found := findClass("Klasse", "no.fint.unknown", classMap, classNameMap); found || got != nil {
		t.Fatalf("expected ambiguous short-name lookup to fail when package context does not match")
	}
}

func TestFindImport_PrefersPackageQualifiedWhenAvailable(t *testing.T) {
	imports := map[string]types.Import{
		"no.fint.skole.Klasse": {Java: "no.fint.skole.Klasse", CSharp: "Fint.Skole"},
		"no.fint.noark.Klasse": {Java: "no.fint.noark.Klasse", CSharp: "Fint.Noark"},
	}
	importNameMap := buildImportNameMap(imports)

	if got, found := findImport("Klasse", "no.fint.noark", imports, importNameMap); !found || got.Java != "no.fint.noark.Klasse" {
		t.Fatalf("expected qualified import lookup to find noark Klasse")
	}
}

func TestGetImports_IncludesSubpackageImportAndExcludesSelfImport(t *testing.T) {
	imports := map[string]types.Import{}
	imports["no.fint.arkiv.noark.Registrering"] = types.Import{Java: "no.fint.arkiv.noark.Registrering"}
	imports["no.fint.arkiv.noark.kodeverk.Klasse"] = types.Import{Java: "no.fint.arkiv.noark.kodeverk.Klasse"}
	importNameMap := buildImportNameMap(imports)

	c := &types.Class{
		Name:    "Registrering",
		Package: "no.fint.arkiv.noark",
		Attributes: []types.Attribute{
			{Name: "klasse", Type: "Klasse"},
		},
	}

	imps := getImports(c, imports, importNameMap)
	if len(imps) == 0 {
		t.Fatalf("expected at least one import")
	}
	for _, imp := range imps {
		if imp == "no.fint.arkiv.noark.Registrering" {
			t.Fatalf("did not expect self import")
		}
	}

	// Subpackage import must not be treated as 'same package'
	foundSub := false
	for _, imp := range imps {
		if imp == "no.fint.arkiv.noark.kodeverk.Klasse" {
			foundSub = true
		}
	}
	if !foundSub {
		t.Fatalf("expected import from subpackage to be included")
	}
}
