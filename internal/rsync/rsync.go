package rsync

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thekashifmalik/rincr/internal"
)

func Run(source string, destination string) error {
	rsyncBinary, err := exec.LookPath("rsync")
	if err != nil {
		return fmt.Errorf("Cannot find rsync binary: %w", err)
	}
	cmd := exec.Command(rsyncBinary, "-hav", "--exclude", internal.BACKUPS_DIR, source, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running rsync: %w", err)
	}
	return nil
}

func RunWithDelete(source string, destination string) error {
	rsyncBinary, err := exec.LookPath("rsync")
	if err != nil {
		return fmt.Errorf("Cannot find rsync binary: %w", err)
	}
	cmd := exec.Command(rsyncBinary, "-hav", "--delete", "--exclude", internal.BACKUPS_DIR, source, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running rsync: %w", err)
	}
	return nil
}
