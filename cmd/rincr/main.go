package main

import (
	"fmt"
	"os"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/backup"
	"github.com/thekashifmalik/rincr/internal/help"
	"github.com/thekashifmalik/rincr/internal/prune"
	"github.com/thekashifmalik/rincr/internal/restore"
	"github.com/thekashifmalik/rincr/internal/version"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Command is a simple interface that all commands must implement. This makes it easier to run commands dynamically.
type Command interface {
	Run() error
}

func run() error {
	if version.ArgExists(os.Args) {
		version.PrintWithName(os.Stdout)
		return nil
	}
	if help.ArgExists(os.Args) {
		help.Print(os.Stdout)
		return nil
	}
	parsedArgs, err := args.Parse(os.Args)
	if err != nil {
		help.Print(os.Stdout)
		return nil
	}
	var cmd Command
	cmd, err = backup.ParseRoot(parsedArgs)
	if err != nil {
		return err
	}
	if len(parsedArgs.Params) > 0 {
		switch parsedArgs.Params[0] {
		case "backup":
			cmd, err = backup.Parse(parsedArgs)
			if err != nil {
				return err
			}
		case "prune":
			cmd, err = prune.Parse(parsedArgs)
			if err != nil {
				return err
			}
		case "restore":
			cmd, err = restore.Parse(parsedArgs)
			if err != nil {
				return err
			}
		}
	}
	return cmd.Run()
}
