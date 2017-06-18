package command

import (
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/FINTprosjektet/fint-model/common/github"
)

func CmdListTags(c *cli.Context) {
	for _, t := range github.GetTagList() {
		fmt.Println(t)
	}
}
