package types

import "fmt"

const LIST_TEMPLATE = "List<%s>"

func getType(list bool, t string) string {
	if list {
		return fmt.Sprintf(LIST_TEMPLATE, t)
	}
	return t
}
