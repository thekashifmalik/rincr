package backup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/thekashifmalik/rincr/internal"
)

func rotateLastBackup(dest *internal.Destination) (string, error) {
	if dest.RemoteHost == "" {
		destination := dest.Path
		err := os.MkdirAll(destination+internal.BACKUPS_DIR_PATH, os.ModePerm)
		if err != nil {
			return "", err
		}

		var destinationLast string
		b, err := os.ReadFile(destination + internal.LAST_FILE_PATH)
		if err != nil {
			fmt.Println("no existing backups")
			return "", nil
		}
		last := string(b)
		destinationLast = fmt.Sprintf("%v/%v/%v", destination, internal.BACKUPS_DIR, last)
		fmt.Printf("rotating: %v\n", last)
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
			if name != internal.BACKUPS_DIR {
				cpFiles = append(cpFiles, fmt.Sprintf("%v/%v", destination, name))
			}
		}
		cmdArgs := append([]string{"-al", "-t", destinationLast}, cpFiles...)
		cmd := exec.Command("cp", cmdArgs...)
		return destinationLast, cmd.Run()
	} else {
		destination := dest.RemotePath
		err := exec.Command("ssh", dest.RemoteHost, "mkdir", "-p", destination+internal.BACKUPS_DIR_PATH).Run()
		if err != nil {
			return "", err
		}

		var destinationLast string
		b, err := exec.Command("ssh", dest.RemoteHost, "cat", destination+internal.LAST_FILE_PATH).Output()
		last := string(b)
		if err != nil || last == "" {
			fmt.Println("> No existing backups")
			return "", nil
		}
		destinationLast = fmt.Sprintf("%v/%v/%v", destination, internal.BACKUPS_DIR, last)
		fmt.Printf("rotating: %v\n", last)
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
			if targetFile != internal.BACKUPS_DIR && targetFile != "" {
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
