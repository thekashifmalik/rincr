package restore

import (
	"fmt"
	"slices"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
	"github.com/thekashifmalik/rincr/internal/rsync"
)

func Restore(repo *repository.Repository, paths []string, output string) error {
	if !repo.Exists() {
		return fmt.Errorf("No repository found")
	}

	sources := []string{}
	for _, path := range paths {
		restorePath := ""
		if repo.PathExists(path) {
			restorePath = path
		} else {
			backupTimes, err := repo.GetBackupTimes()
			if err != nil {
				return err
			}
			slices.Reverse(backupTimes)
			for _, backupTime := range backupTimes {
				timestamp := backupTime.Format(internal.TIME_FORMAT)
				historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
				if repo.PathExists(historicalPath) {
					restorePath = historicalPath
					break
				}
			}
			if restorePath == "" {
				fmt.Printf("Skipped: %v\n", path)
				continue
			}
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
