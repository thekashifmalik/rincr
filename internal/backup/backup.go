package backup

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thekashifmalik/rincr/internal"
)

func clean(destination *internal.Destination, destinationLast string) error {
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
