package classes

import (
	"fmt"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/FINTprosjektet/fint-model/common/parser"
	"github.com/FINTprosjektet/fint-model/common/types"
	"github.com/codegangsta/cli"
)

func CmdPrintClasses(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	clazzes, _, _:= parser.GetClasses(tag, c.GlobalBool("force"))

	for _, c := range clazzes {
		dumpClass(c)
	}

}

func dumpClass(c types.Class) {
	fmt.Printf("Class: %s\n", c.Name)
	fmt.Printf("  Abstract: %t\n", c.Abstract)
	fmt.Printf("  Identifiable: %t\n", c.Identifiable)
	if len(c.Extends) > 0 {
		fmt.Printf("  Extends: %s\n", c.Extends)
	}
	for _, a := range c.Attributes {
		fmt.Printf("  Attribute: %s\n", a.Name)
		fmt.Printf("    Type: %s\n", a.Type)
	}
	fmt.Println("  Relations:")
	for _, r := range c.Relations {
		fmt.Printf("    %s\n", r)
	}

}
