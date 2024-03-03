package restore

import (
	"fmt"
	"slices"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
	"github.com/thekashifmalik/rincr/internal/rsync"
)

func Restore(repo *repository.Repository, paths []string, output string, latest bool) error {
	if !latest {
		return fmt.Errorf("must specify restore mode")
	}
	if !repo.Exists() {
		return fmt.Errorf("No repository found")
	}

	sources := []string{}
	for _, path := range paths {
		var restorePath string
		var ok bool
		var err error
		if latest {
			restorePath, ok, err = findLatest(path, repo)
		}
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Skipped: %v\n", path)
			continue
		}
		if restorePath == path {
			fmt.Printf("Found: %v\n", path)
		} else {
			fmt.Printf("Found: %v @ %v\n", path, restorePath)
		}
		source := fmt.Sprintf("%v/%v", repo.GetFullPath(), restorePath)
		sources = append(sources, source)
	}
	args := append(sources, output)
	return rsync.Run(args...)
}

func findLatest(path string, repo *repository.Repository) (string, bool, error) {
	if repo.PathExists(path) {
		return path, true, nil
	}
	backupTimes, err := repo.GetBackupTimes()
	if err != nil {
		return "", false, err
	}
	slices.Reverse(backupTimes)
	for _, backupTime := range backupTimes {
		timestamp := backupTime.Format(internal.TIME_FORMAT)
		historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
		if repo.PathExists(historicalPath) {
			return historicalPath, true, nil
		}
	}
	return "", false, nil
}
