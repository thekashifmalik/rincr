package prune

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/repository"
)

func Prune(repository *repository.Repository, destination *internal.Destination, currentTime time.Time) error {
	existingBackups, err := repository.GetBackupTimes()
	if err != nil {
		return err
	}
	// Go through time buckets, keeping only the oldest backup from each bucket.
	fmt.Printf("> Pruning backups in: %v\n", destination.Path+internal.BACKUPS_DIR_PATH)
	pruned, checkedTill := pruneStage(existingBackups, roundToHour(currentTime.Add(-time.Hour)), destination, 23, time.Hour)
	pruned, checkedTill = pruneStage(pruned, roundToDay(checkedTill), destination, 30, 24*time.Hour)
	pruned, checkedTill = pruneMonthly(pruned, roundToMonth(checkedTill), destination, 12)
	pruned, checkedTill = pruneYearly(pruned, roundToYear(checkedTill), destination, 10)
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

func pruneStage(existingBackups []time.Time, currentTime time.Time, destination *internal.Destination, num int, period time.Duration) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.Add(time.Duration(i) * -period)
		// fmt.Printf("> Checking %v: %v\n", period, checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, destination)
	}
	return prunedBackups, checkTime
}

func pruneMonthly(existingBackups []time.Time, currentTime time.Time, destination *internal.Destination, num int) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.AddDate(0, -i, 0)
		// fmt.Printf("> Checking monthly: %v\n", checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, destination)
	}
	return prunedBackups, checkTime
}

func pruneYearly(existingBackups []time.Time, currentTime time.Time, destination *internal.Destination, num int) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := range num {
		checkTime = currentTime.AddDate(-i, 0, 0)
		// fmt.Printf("> Checking yearly: %v\n", checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, destination)
	}
	return prunedBackups, checkTime
}

// This function needs to be run from the latest bucket to the oldest.
func pruneBucket(existingBackups []time.Time, bucketTime time.Time, destination *internal.Destination) []time.Time {
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
	for _, backupTime := range bucket[1:] {
		fmt.Printf("> Pruning: %v\n", backupTime)
		// TODO: Handle any errors here
		if destination.RemoteHost == "" {
			os.RemoveAll(fmt.Sprintf("%v/%v/%v", destination.Path, internal.BACKUPS_DIR, backupTime.Format(internal.TIME_FORMAT)))
		} else {
			remotePath := fmt.Sprintf("%v/%v/%v", destination.RemotePath, internal.BACKUPS_DIR, backupTime.Format(internal.TIME_FORMAT))
			exec.Command("ssh", destination.RemoteHost, "rm", "-rf", remotePath).Run()
		}
	}
	return unseen
}
