package command

import (
	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-model/github"
	"github.com/FINTprosjektet/fint-model/package"
	"fmt"
)

func CmdListPackages(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	//document.GetFile(tag)

	for _, p := range packages.DistinctList(tag) {
		fmt.Println(p)
	}
}
