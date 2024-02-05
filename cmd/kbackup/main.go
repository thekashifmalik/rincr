package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Args struct {
	sources     []string
	destination string
	prune       bool
}

func parseArgs() (*Args, error) {
	args := []string{}
	prune := false
	for _, arg := range os.Args {
		if arg == "--prune" {
			prune = true
		} else {
			args = append(args, arg)
		}
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("No sources provided")
	}
	if len(args) < 3 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := args[1 : len(args)-1]
	destination := args[len(args)-1]
	return &Args{
		sources:     sources,
		destination: destination,
		prune:       prune,
	}, nil
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	sources := args.sources
	destination := args.destination

	for _, source := range sources {
		currentTime := time.Now()
		target := filepath.Base(source)
		destinationTarget := fmt.Sprintf("%v/%v", destination, target)
		err := os.MkdirAll(destinationTarget+"/.kbackup", os.ModePerm)
		if err != nil {
			return err
		}

		var destinationLast string
		b, err := os.ReadFile(destinationTarget + "/.kbackup/last")
		if err == nil {
			last := string(b)
			destinationLast = fmt.Sprintf("%v/.kbackup/%v", destinationTarget, last)
			fmt.Printf("> Rotating last backup: %v\n", destinationLast)
			err := os.MkdirAll(destinationLast, os.ModePerm)
			if err != nil {
				return err
			}

			cpFiles := []string{}
			targetFiles, err := os.ReadDir(destinationTarget)
			if err != nil {
				return err
			}
			for _, targetFile := range targetFiles {
				name := targetFile.Name()
				if name != ".kbackup" {
					cpFiles = append(cpFiles, fmt.Sprintf("%v/%v", destinationTarget, name))
				}
			}
			cmdArgs := append([]string{"-al", "-t", destinationLast}, cpFiles...)
			cmd := exec.Command("cp", cmdArgs...)
			err = cmd.Run()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("> No existing backups")
		}

		fmt.Printf("> Backing up: %v -> %v\n", source, destinationTarget)

		rsyncBinary, err := exec.LookPath("rsync")
		if err != nil {
			return fmt.Errorf("Cannot find rsync binary: %w", err)
		}
		cmd := exec.Command(rsyncBinary, "-hav", "--delete", "--exclude", ".kbackup", source+"/", destinationTarget)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			errs := []error{fmt.Errorf("Error running rsync: %w", err)}
			if destinationLast != "" {
				fmt.Printf("> Cleaning up: %v\n", destinationLast)
				err := os.RemoveAll(destinationLast)
				if err != nil {
					errs = append(errs, fmt.Errorf("Error cleaning up: %w", err))
				}
			}
			return errors.Join(errs...)
		}
		f, err := os.Create(destinationTarget + "/.kbackup/last")
		if err != nil {
			return err
		}
		_, err = f.WriteString(currentTime.Format("2006-01-02T15-04-05"))
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
