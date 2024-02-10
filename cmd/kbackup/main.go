package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thekashifmalik/kbackup/internal"
	"github.com/thekashifmalik/kbackup/internal/args"
)

var version string

var HELP = `kbackup

Usage:
  kbackup [--prune] [[USER@]HOST:]SRC... [[USER@]HOST:]DEST
  kbackup (-h | --help)
  kbackup --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --prune    	Prune backups after operation.
`

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	if args.ParseVersion(os.Args) {
		fmt.Printf("kbackup %v\n", version)
		return nil
	}
	if args.ParseHelp(os.Args) {
		fmt.Print(HELP)
		return nil
	}
	args, err := internal.ParseArgs()
	if err != nil {
		return err
	}
	sources := args.Sources
	for _, source := range sources {
		currentTime := time.Now()
		target := filepath.Base(source)
		destinationTarget := internal.ParseDestination(fmt.Sprintf("%v/%v", args.Destination, target))
		if err != nil {
			return err
		}
		destinationLast, err := internal.RotateLastBackup(destinationTarget)
		if err != nil {
			return err
		}

		fmt.Printf("> Backing up: %v -> %v\n", source, destinationTarget.Path)

		err = internal.SyncBackup(source, destinationTarget.Path)
		if err != nil {
			errs := []error{err}
			if destinationLast != "" {
				err := internal.Clean(destinationTarget, destinationLast)
				if err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}

		timeString := currentTime.Format(internal.TIME_FORMAT)
		internal.WriteLastFile(timeString, destinationTarget)

		if args.Prune {
			err := internal.Prune(destinationTarget, currentTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
