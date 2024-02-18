package root

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thekashifmalik/rincr/internal"
)

func writeLastFile(timestamp string, destination *internal.Destination) error {
	if destination.RemoteHost == "" {
		f, err := os.Create(destination.Path + internal.LAST_FILE_PATH)
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
		echoCmd := fmt.Sprintf(`"echo %v > %v"`, timestamp, destination.RemotePath+internal.LAST_FILE_PATH)
		cmd := exec.Command("ssh", destination.RemoteHost, "bash", "-c", echoCmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
