package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RotateLastBackup(dest *Destination) (string, error) {
	if dest.RemoteHost == "" {
		destination := dest.Path
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
	} else {
		destination := dest.RemotePath
		err := exec.Command("ssh", dest.RemoteHost, "mkdir", "-p", destination+"/.kbackup").Run()
		if err != nil {
			return "", err
		}

		var destinationLast string
		b, err := exec.Command("ssh", dest.RemoteHost, "cat", destination+"/.kbackup/last").Output()
		last := string(b)
		if err != nil || last == "" {
			fmt.Println("> No existing backups")
			return "", nil
		}
		destinationLast = fmt.Sprintf("%v/.kbackup/%v", destination, last)
		fmt.Printf("> Rotating last backup: %v\n", destinationLast)
		cmd := exec.Command("ssh", dest.RemoteHost, "mkdir", "-p", destinationLast)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return destinationLast, err
		}

		cmd = exec.Command("ssh", dest.RemoteHost, "ls", "-A", destination)
		cmd.Stderr = os.Stderr
		targetFilesRaw, err := cmd.Output()
		if err != nil {
			return destinationLast, err
		}
		targetFilesString := string(targetFilesRaw)
		targetFiles := strings.Split(targetFilesString, "\n")

		cpFiles := []string{}
		for _, targetFile := range targetFiles {
			if targetFile != ".kbackup" && targetFile != "" {
				cpFiles = append(cpFiles, fmt.Sprintf(`'%v/%v'`, destination, targetFile))
			}
		}

		cpCmd := fmt.Sprintf(`"cp -al %v %v"`, strings.Join(cpFiles, " "), destinationLast)
		cmd = exec.Command("ssh", dest.RemoteHost, "bash", "-c", cpCmd)
		if os.Getenv("DEBUG") != "" {
			cmd.Stderr = os.Stderr
		}
		return destinationLast, cmd.Run()
	}
}
