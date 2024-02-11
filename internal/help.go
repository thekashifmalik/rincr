package internal

import "fmt"

var HELP = fmt.Sprintf(`
Usage:
  %[1]v [--prune] [[USER@]HOST:]SRC... [[USER@]HOST:]DEST
  %[1]v (-h | --help)
  %[1]v --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --prune    	Prune backups after operation.
`, NAME)

func PrintHelp() {
	PrintVersion()
	fmt.Print(HELP)
}
