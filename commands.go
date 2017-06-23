package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-model/generate"
	"github.com/FINTprosjektet/fint-model/branches"
	"github.com/FINTprosjektet/fint-model/packages"
	"github.com/FINTprosjektet/fint-model/namespaces"
	"github.com/FINTprosjektet/fint-model/tags"
	"github.com/FINTprosjektet/fint-model/classes"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "",
		Name:   "tag",
		Value:  "latest",
		Usage:  "the tag (version) of the model to generate",
	},
}

var Commands = []cli.Command{
	{
		Name:   "printClasses",
		Usage:  "list classes",
		Action: classes.CmdPrintClasses,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "generate",
		Usage:  "generates JAVA/CS models",
		Action: generate.CmdGenerate,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lang",
				Value: "JAVA",
				Usage: "the language to generate the code in - can be JAVA, CS or ALL",
			},
		},
	},
	{
		Name:   "listPackages",
		Usage:  "list Java packages",
		Action: packages.CmdListPackages,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listNamespaces",
		Usage:  "list CS namespaces",
		Action: namespaces.CmdListNamespaces,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listTags",
		Usage:  "list tags",
		Action: tags.CmdListTags,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listBranches",
		Usage:  "list branches",
		Action: branches.CmdListBranches,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
