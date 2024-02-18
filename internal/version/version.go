package version

import (
	"fmt"
	"io"

	"github.com/thekashifmalik/rincr/internal"
)

var VERSION string

func PrintWithName(writer io.Writer) {
	fmt.Fprintf(writer, "%v %v\n", internal.NAME, VERSION)
}

func ArgExists(args []string) bool {
	for _, arg := range args {
		if arg == "--version" {
			return true
		}
	}
	return false
}
