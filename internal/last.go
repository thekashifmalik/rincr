package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func WriteLastFile(timestamp string, destination *Destination) error {
	if destination.RemoteHost == "" {
		f, err := os.Create(destination.Path + LAST_FILE_PATH)
		if err != nil {
			return err
		}
		_, err = f.WriteString(timestamp)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	} else {
		echoCmd := fmt.Sprintf(`"echo %v > %v"`, timestamp, destination.RemotePath+LAST_FILE_PATH)
		cmd := exec.Command("ssh", destination.RemoteHost, "bash", "-c", echoCmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
