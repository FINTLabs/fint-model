package packages

import (
	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-model/common/github"
	"fmt"
	"github.com/FINTprosjektet/fint-model/common/config"
)

func CmdListPackages(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}

	//document.GetFile(tag)

	for _, p := range DistinctPackageList(tag, c.GlobalBool("force")) {
		fmt.Println(p)
	}
}
