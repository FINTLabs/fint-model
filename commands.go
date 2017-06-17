package main

import (
	"fmt"
	"os"

	"github.com/FINTprosjektet/fint-model/command"
	"github.com/codegangsta/cli"
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
		Usage:  "",
		Action: command.CmdPrintClasses,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "generate",
		Usage:  "",
		Action: command.CmdGenerate,
		Flags:  []cli.Flag{
			cli.StringFlag{
				Name: "lang",
				Value: "JAVA",
				Usage: "the language to generate the code in - can be JAVA or NET",
			},
		},
	},
	{
		Name:   "listPackages",
		Usage:  "list Java packages",
		Action: command.CmdListPackages,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listNamespaces",
		Usage:  "list .Net namespaces",
		Action: command.CmdListNamespaces,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listTags",
		Usage:  "list tags",
		Action: command.CmdListTags,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "listBranches",
		Usage:  "list branches",
		Action: command.CmdListBranches,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
