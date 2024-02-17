package main

import (
	"fmt"
	"os"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/backup"
	"github.com/thekashifmalik/rincr/internal/help"
	"github.com/thekashifmalik/rincr/internal/root"
	"github.com/thekashifmalik/rincr/internal/version"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Command interface {
	Run() error
}

func run() error {
	if version.ArgExists(os.Args) {
		version.PrintWithName()
		return nil
	}
	if help.ArgExists(os.Args) {
		help.Print()
		return nil
	}
	parsedArgs := args.Parse(os.Args)
	cmd, err := getCommand(parsedArgs)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func getCommand(parsedArgs *args.Parsed) (Command, error) {
	if parsedArgs.Params[0] == "backup" {
		return backup.Parse(parsedArgs)
	}
	return root.Parse(parsedArgs)
}
