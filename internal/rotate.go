package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func RotateLastBackup(destination string) (string, error) {
	err := os.MkdirAll(destination+"/.kbackup", os.ModePerm)
	if err != nil {
		return "", err
	}

	var destinationLast string
	b, err := os.ReadFile(destination + "/.kbackup/last")
	if err != nil {
		fmt.Println("> No existing backups")
		return "", nil
	}
	last := string(b)
	destinationLast = fmt.Sprintf("%v/.kbackup/%v", destination, last)
	fmt.Printf("> Rotating last backup: %v\n", destinationLast)
	err = os.MkdirAll(destinationLast, os.ModePerm)
	if err != nil {
		return destinationLast, err
	}

	cpFiles := []string{}
	targetFiles, err := os.ReadDir(destination)
	if err != nil {
		return destinationLast, err
	}
	for _, targetFile := range targetFiles {
		name := targetFile.Name()
		if name != ".kbackup" {
			cpFiles = append(cpFiles, fmt.Sprintf("%v/%v", destination, name))
		}
	}
	cmdArgs := append([]string{"-al", "-t", destinationLast}, cpFiles...)
	cmd := exec.Command("cp", cmdArgs...)
	return destinationLast, cmd.Run()
}
