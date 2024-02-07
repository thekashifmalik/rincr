package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func SyncBackup(source string, destination string) error {
	rsyncBinary, err := exec.LookPath("rsync")
	if err != nil {
		return fmt.Errorf("Cannot find rsync binary: %w", err)
	}
	cmd := exec.Command(rsyncBinary, "-hav", "--delete", "--exclude", ".kbackup", source+"/", destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running rsync: %w", err)
	}
	return nil
}

func Clean(destination *Destination, destinationLast string) error {
	if destination.RemoteHost == "" {
		fmt.Printf("> Cleaning up: %v\n", destinationLast)
		err := os.RemoveAll(destinationLast)
		if err != nil {
			return fmt.Errorf("Error cleaning up: %w", err)
		}
	} else {
		fmt.Printf("> Cleaning up: %v:%v\n", destination.RemoteHost, destinationLast)
		err := exec.Command("ssh", destination.RemoteHost, "rm", "-rf", destinationLast).Run()
		if err != nil {
			return fmt.Errorf("Error cleaning up: %w", err)
		}
	}
	return nil
}
