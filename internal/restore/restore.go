package restore

import (
	"fmt"
	"slices"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
	"github.com/thekashifmalik/rincr/internal/rsync"
)

func Restore(repo *repository.Repository, path string, output string) error {
	if !repo.Exists() {
		return fmt.Errorf("No repository found")
	}
	restorePath := ""
	if repo.PathExists(path) {
		restorePath = path
	} else {
		fmt.Println("Path not found in latest backup, checking historical backups")
		backupTimes, err := repo.GetBackupTimes()
		if err != nil {
			return err
		}
		slices.Reverse(backupTimes)
		for _, backupTime := range backupTimes {
			timestamp := backupTime.Format(internal.TIME_FORMAT)
			fmt.Printf("Checking: %v\n", timestamp)
			historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
			if repo.PathExists(historicalPath) {
				restorePath = historicalPath
				break
			}
		}
		if restorePath == "" {
			return fmt.Errorf("Path not found")
		}
	}
	source := fmt.Sprintf("%v/%v", repo.GetFullPath(), restorePath)
	return rsync.Run(source, output)
}
