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

func run() error {
	if version.ArgExists(os.Args) {
		version.PrintWithName()
		return nil
	}
	if help.ArgExists(os.Args) {
		help.Print()
		return nil
	}
	parsedArgs, err := args.Parse(os.Args)
	if err != nil {
		help.Print()
		return nil
	}
	if parsedArgs.Params[0] == "backup" {
		cmd, err := backup.Parse(parsedArgs)
		if err != nil {
			return err
		}
		return cmd.Run()
	}
	cmd, err := root.Parse(parsedArgs)
	if err != nil {
		return err
	}
	return cmd.Run()
}
