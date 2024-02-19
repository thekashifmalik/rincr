package restore

import (
	"fmt"
	"slices"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
)

func Restore(repository *repository.Repository, path string, output string) error {
	if !repository.Exists() {
		return fmt.Errorf("No repository found")
	}
	restorePath := ""
	if repository.PathExists(path) {
		restorePath = path
	} else {
		fmt.Println("Path not found in latest backup, checking historical backups")
		backupTimes, err := repository.GetBackupTimes()
		if err != nil {
			return err
		}
		slices.Reverse(backupTimes)
		for _, backupTime := range backupTimes {
			timestamp := backupTime.Format(internal.TIME_FORMAT)
			fmt.Printf("Checking: %v\n", timestamp)
			historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
			if repository.PathExists(historicalPath) {
				restorePath = historicalPath
				break
			}
		}
		if restorePath == "" {
			return fmt.Errorf("Path not found")
		}
	}
	fmt.Println(restorePath)
	return nil
}
