package command

import (
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/FINTprosjektet/fint-model/common/parser"
)

func CmdPrintClasses(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	classes := parser.GetClasses(tag)

	for _, c := range classes {
		dumpClass(c)
	}

}

func dumpClass(c parser.Class) {
	fmt.Printf("Class: %s\n", c.Name)
	fmt.Printf("  Abstract: %t\n", c.Abstract)
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
