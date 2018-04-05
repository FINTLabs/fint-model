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
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}

	clazzes, _, _, _ := parser.GetClasses(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), c.GlobalBool("force"))

	for _, c := range clazzes {
		dumpClass(c)
	}

}

func dumpClass(c *types.Class) {
	fmt.Printf("Class: %s\n", c.Name)
	fmt.Printf("  Abstract: %t\n", c.Abstract)
	fmt.Printf("  Identifiable: %t\n", c.Identifiable)
	fmt.Printf("  Resource: %t\n", c.Resource)
	fmt.Printf("  Package: %s\n", c.Package)
	fmt.Printf("  Namespace: %s\n", c.Namespace)
	//fmt.Printf("  DocumentationUrl: %s\n", c.DocumentationUrl)
	if len(c.Extends) > 0 {
		fmt.Printf("  Extends: %s\n", c.Extends)
		fmt.Printf("  ExtendsResource: %t\n", c.ExtendsResource)
	}
	fmt.Println("  Imports:")
	for _, i := range c.Imports {
		fmt.Printf("    - %s\n", i)
	}
	fmt.Println("  Using:")
	for _, u := range c.Using {
		fmt.Printf("    - %s\n", u)
	}

	if len(c.Attributes) > 0 {
		fmt.Println("  Attributes: ")
		for _, a := range c.Attributes {
			if a.List {
				fmt.Printf("    - %s: List<%s>\n", a.Name, a.Type)
			} else {
				fmt.Printf("    - %s: %s\n", a.Name, a.Type)
			}
		}
	}

	if len(c.Relations) > 0 {
		fmt.Print("  Relations: ")
		for i, r := range c.Relations {
			fmt.Print(r)
			if i == len(c.Relations)-1 {
				fmt.Println()
			} else {
				fmt.Print(", ")
			}
		}
	}

	if len(c.Resources) > 0 {
		fmt.Println("  Resources:")
		for _, a := range c.Resources {
			if a.List {
				fmt.Printf("    - %s: List<%s>\n", a.Name, a.Type)
			} else {
				fmt.Printf("    - %s: %s\n", a.Name, a.Type)
			}
		}
	}
}
