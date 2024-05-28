package restore

import (
	"fmt"
	"slices"
	"time"

	"github.com/xhit/go-str2duration/v2"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
	"github.com/thekashifmalik/rincr/internal/rsync"
)

func (c *Command) Restore(repo *repository.Repository, paths []string, output string, latest bool) error {
	if !repo.Exists() {
		return fmt.Errorf("No repository found")
	}

	sources := []string{}
	for _, path := range paths {
		var restorePath string
		var ok bool
		var err error
		if c.Latest {
			restorePath, ok, err = findLatest(path, repo)
		} else if c.From {
			restorePath, ok, err = findFrom(path, repo, c.FromValue)
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
	if len(sources) == 0 {
		return nil
	}
	args := append(sources, output)
	return rsync.Run(args...)
}

func findLatest(path string, repo *repository.Repository) (string, bool, error) {
	fmt.Printf("Checking: mirror\n")
	if repo.PathExists(path) {
		return path, true, nil
	}
	backupTimes, err := repo.GetBackupTimes()
	if err != nil {
		return "", false, err
	}
	slices.Reverse(backupTimes)
	for _, backupTime := range backupTimes {
		fmt.Printf("Checking: %v \n", backupTime)
		timestamp := backupTime.Format(internal.TIME_FORMAT)
		historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
		if repo.PathExists(historicalPath) {
			return historicalPath, true, nil
		}
	}
	return "", false, nil
}

func findFrom(path string, repo *repository.Repository, fromValue string) (string, bool, error) {
	duration, err := str2duration.ParseDuration(fromValue)
	if err != nil {
		return "", false, err
	}
	target := time.Now().Add(-1 * duration)
	fmt.Printf("Finding closest backup to: %v\n", target)

	backupTimes, err := repo.GetBackupTimes()
	if err != nil {
		return "", false, err
	}
	slices.Reverse(backupTimes)
	targetBackup, ok := findClosest(backupTimes, target)
	if !ok {
		return "", false, nil
	}
	fmt.Printf("Checking: %v \n", targetBackup)
	timestamp := targetBackup.Format(internal.TIME_FORMAT)
	historicalPath := fmt.Sprintf("%v/%v/%v", internal.BACKUPS_DIR, timestamp, path)
	if repo.PathExists(historicalPath) {
		return historicalPath, true, nil
	}
	return "", false, nil
}

func findClosest(backupTimes []time.Time, target time.Time) (*time.Time, bool) {
	var minDistance time.Duration
	var minDistanceTime *time.Time
	for _, backupTime := range backupTimes {
		distance := backupTime.Sub(target).Abs()
		if distance < minDistance || minDistanceTime == nil {
			minDistance = distance
			minDistanceTime = &backupTime
		}
	}
	if minDistanceTime == nil {
		return nil, false
	}
	return minDistanceTime, true

}
