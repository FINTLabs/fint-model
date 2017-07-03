package tags

import (
	"fmt"
	"github.com/FINTprosjektet/fint-model/common/github"
	"github.com/codegangsta/cli"
)

func CmdListTags(c *cli.Context) {
	for _, t := range github.GetTagList() {
		fmt.Println(t)
	}
}
