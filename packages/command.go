package packages

import (
	"fmt"

	"github.com/FINTLabs/fint-model/common/config"
	"github.com/FINTLabs/fint-model/common/github"
	"github.com/codegangsta/cli"
)

func CmdListPackages(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}

	//document.GetFile(tag)

	for _, p := range DistinctPackageList(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), c.GlobalBool("force")) {
		fmt.Println(p)
	}
}
