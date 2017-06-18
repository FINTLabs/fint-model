package namespaces

import (
	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-model/common/github"
	"fmt"
)

func CmdListNamespaces(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == "latest" {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	for _, p := range DistinctNamespaceList(tag) {
		fmt.Println(p)
	}
}
