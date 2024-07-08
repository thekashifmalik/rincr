package prune

import (
	"fmt"
	"slices"
	"time"

	"github.com/thekashifmalik/rincr/internal/repository"
)

func Prune(repo *repository.Repository, currentTime time.Time, config *Config) error {
	fmt.Printf("pruning: %v\n", repo.GetFullPath())
	existingBackups, err := repo.GetBackupTimes()
	if err != nil {
		return err
	}
	// Go through time buckets, keeping only the oldest backup from each bucket.
	pruned, checkedTill := pruneStage(existingBackups, roundToHour(currentTime.Add(-time.Hour)), repo, config.Hourly-1, time.Hour)
	pruned, checkedTill = pruneStage(pruned, roundToDay(checkedTill), repo, config.Daily, 24*time.Hour)
	pruned, checkedTill = pruneMonthly(pruned, roundToMonth(checkedTill), repo, config.Monthly)
	pruned, checkedTill = pruneYearly(pruned, roundToYear(checkedTill), repo, config.Yearly)
	repo.DeleteBackupsByTime(pruned)
	return nil
}

func roundToHour(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), target.Day(), target.Hour(), 0, 0, 0, target.Location())
}

func roundToDay(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location())
}

func roundToMonth(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, target.Location())
}

func roundToYear(target time.Time) time.Time {
	return time.Date(target.Year(), 1, 1, 0, 0, 0, 0, target.Location())
}

func pruneStage(existingBackups []time.Time, currentTime time.Time, repo *repository.Repository, num int, period time.Duration) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.Add(time.Duration(i) * -period)
		// fmt.Printf("> Checking %v: %v\n", period, checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, repo)
	}
	return prunedBackups, checkTime
}

func pruneMonthly(existingBackups []time.Time, currentTime time.Time, repo *repository.Repository, num int) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.AddDate(0, -i, 0)
		// fmt.Printf("> Checking monthly: %v\n", checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, repo)
	}
	return prunedBackups, checkTime
}

func pruneYearly(existingBackups []time.Time, currentTime time.Time, repo *repository.Repository, num int) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.AddDate(-i, 0, 0)
		// fmt.Printf("> Checking yearly: %v\n", checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, repo)
	}
	return prunedBackups, checkTime
}

// This function needs to be run from the latest bucket to the oldest.
func pruneBucket(existingBackups []time.Time, bucketTime time.Time, repo *repository.Repository) []time.Time {
	unseen := []time.Time{}
	// Gather backups that fit in this bucket
	bucket := []time.Time{}
	for _, backupTime := range existingBackups {
		if backupTime.After(bucketTime) {
			bucket = append(bucket, backupTime)
		} else {
			unseen = append(unseen, backupTime)
		}
	}
	if len(bucket) == 0 {
		// fmt.Printf("> No backups for: %v\n", bucketTime)
		return existingBackups
	}
	// Sort the backups and keep only the oldest one
	slices.SortFunc(bucket, func(a, b time.Time) int { return a.Compare(b) })
	repo.DeleteBackupsByTime(bucket[1:])
	return unseen
}
