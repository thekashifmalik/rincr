package rsync

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thekashifmalik/rincr/internal"
)

var OPTIONS = []string{"-hav", "--exclude", internal.BACKUPS_DIR}

func Run(params ...string) error {
	args := []string{}
	args = append(args, OPTIONS...)
	args = append(args, params...)
	return run(args...)
}

func RunWithDelete(params ...string) error {
	args := []string{}
	args = append(args, OPTIONS...)
	args = append(args, "--delete")
	args = append(args, params...)
	return run(args...)
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
