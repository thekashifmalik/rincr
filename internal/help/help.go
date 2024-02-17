package help

import (
	"fmt"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/version"
)

var HELP = fmt.Sprintf(`
Usage:
  %[1]v [--prune] [[USER@]HOST:]SRC... [[USER@]HOST:]DEST
  %[1]v (-h | --help)
  %[1]v --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --prune    	Prune backups after operation.
`, internal.NAME)

func Print() {
	version.PrintWithName()
	fmt.Print(HELP)
}

func ArgExists(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}
