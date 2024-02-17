package version

import (
	"fmt"

	"github.com/thekashifmalik/rincr/internal"
)

var VERSION string

func PrintWithName() {
	fmt.Printf("%v %v\n", internal.NAME, VERSION)
}

func ArgExists(args []string) bool {
	for _, arg := range args {
		if arg == "--version" {
			return true
		}
	}
	return false
}
