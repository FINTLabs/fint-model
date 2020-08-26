package classes

import (
	"fmt"

	"github.com/FINTLabs/fint-model/common/github"
	"github.com/FINTLabs/fint-model/common/parser"
	"github.com/FINTLabs/fint-model/common/types"
	"github.com/urfave/cli"
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
	dep := ""
	if c.Deprecated {
		dep = "<<DEPRECATED>>"
	}
	fmt.Printf("Class: %s %s\n", c.Name, dep)
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
		fmt.Println("  Attributes:")
		for _, a := range c.Attributes {
			dep := ""
			if a.Deprecated {
				dep = "<<DEPRECATED>>"
			}
			if a.List {
				fmt.Printf("    - %s: List<%s> %s\n", a.Name, a.Type, dep)
			} else {
				fmt.Printf("    - %s: %s %s\n", a.Name, a.Type, dep)
			}
		}
	}

	if len(c.Relations) > 0 {
		fmt.Println("  Relations:")
		for _, r := range c.Relations {
			s := ""
			if r.Deprecated {
				s = "<<DEPRECATED>>"
			}
			fmt.Printf("    - %s: %s[%s] %s\n", r.Name, r.Target, r.Multiplicity, s)
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
