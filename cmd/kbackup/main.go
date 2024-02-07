package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thekashifmalik/kbackup/internal"
)

var version string

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	args, err := internal.ParseArgs()
	if err != nil {
		return err
	}

	if args.Version {
		fmt.Printf("kbackup %v\n", version)
		return nil
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
