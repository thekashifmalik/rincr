package rsync

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thekashifmalik/rincr/internal"
)

func Run(source string, destination string) error {
	return run("-hav", "--exclude", internal.BACKUPS_DIR, source, destination)
}

func RunWithDelete(source string, destination string) error {
	return run("-hav", "--delete", "--exclude", internal.BACKUPS_DIR, source, destination)
}

func run(args ...string) error {
	binary, err := exec.LookPath("rsync")
	if err != nil {
		return fmt.Errorf("Cannot find rsync binary: %w", err)
	}
	cmd := exec.Command(binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running rsync: %w", err)
	}
	return nil
}
