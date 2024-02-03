package main

import (
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

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("No sources provided")
	}
	if len(os.Args) < 3 {
		return fmt.Errorf("No destination provided")
	}
	sources := os.Args[1 : len(os.Args)-1]
	destination := os.Args[len(os.Args)-1]

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
		cmd := exec.Command("rsync", "-hav", "--delete", "--exclude", ".kbackup", source+"/", destinationTarget)
		err = cmd.Run()
		if err != nil {
			if destinationLast != "" {
				fmt.Printf("> Cleaning up: %v\n", destinationLast)
				os.RemoveAll(destinationLast)
			}
			return err
		}

		f, err := os.Create(destinationTarget + "/.kbackup/last")
		if err != nil {
			return err
		}
		_, err = f.WriteString(currentTime.Format("2006-01-02T03-04-05"))
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
