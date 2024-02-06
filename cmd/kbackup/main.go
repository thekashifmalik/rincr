package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/exp/slices"
)

var TIME_FORMAT = "2006-01-02T15-04-05"

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Args struct {
	sources     []string
	destination string
	prune       bool
}

func parseArgs() (*Args, error) {
	args := []string{}
	prune := false
	for _, arg := range os.Args {
		if arg == "--prune" {
			prune = true
		} else {
			args = append(args, arg)
		}
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("No sources provided")
	}
	if len(args) < 3 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := args[1 : len(args)-1]
	destination := args[len(args)-1]
	return &Args{
		sources:     sources,
		destination: destination,
		prune:       prune,
	}, nil
}

func run() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	sources := args.sources
	destination := args.destination

	for _, source := range sources {
		currentTime := time.Now()
		target := filepath.Base(source)
		destinationTarget := fmt.Sprintf("%v/%v", destination, target)
		err := os.MkdirAll(destinationTarget+"/.kbackup", os.ModePerm)
		if err != nil {
			return err
		}

		var destinationLast string
		b, err := os.ReadFile(destinationTarget + "/.kbackup/last")
		if err == nil {
			last := string(b)
			destinationLast = fmt.Sprintf("%v/.kbackup/%v", destinationTarget, last)
			fmt.Printf("> Rotating last backup: %v\n", destinationLast)
			err := os.MkdirAll(destinationLast, os.ModePerm)
			if err != nil {
				return err
			}

			cpFiles := []string{}
			targetFiles, err := os.ReadDir(destinationTarget)
			if err != nil {
				return err
			}
			for _, targetFile := range targetFiles {
				name := targetFile.Name()
				if name != ".kbackup" {
					cpFiles = append(cpFiles, fmt.Sprintf("%v/%v", destinationTarget, name))
				}
			}
			cmdArgs := append([]string{"-al", "-t", destinationLast}, cpFiles...)
			cmd := exec.Command("cp", cmdArgs...)
			err = cmd.Run()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("> No existing backups")
		}

		fmt.Printf("> Backing up: %v -> %v\n", source, destinationTarget)

		rsyncBinary, err := exec.LookPath("rsync")
		if err != nil {
			return fmt.Errorf("Cannot find rsync binary: %w", err)
		}
		cmd := exec.Command(rsyncBinary, "-hav", "--delete", "--exclude", ".kbackup", source+"/", destinationTarget)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			errs := []error{fmt.Errorf("Error running rsync: %w", err)}
			if destinationLast != "" {
				fmt.Printf("> Cleaning up: %v\n", destinationLast)
				err := os.RemoveAll(destinationLast)
				if err != nil {
					errs = append(errs, fmt.Errorf("Error cleaning up: %w", err))
				}
			}
			return errors.Join(errs...)
		}
		f, err := os.Create(destinationTarget + "/.kbackup/last")
		if err != nil {
			return err
		}
		_, err = f.WriteString(currentTime.Format(TIME_FORMAT))
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}

		if args.prune {
			// Gather all existing backups
			files, err := os.ReadDir(destinationTarget + "/.kbackup")
			if err != nil {
				return err
			}
			existingBackups := []time.Time{}
			for _, file := range files {
				name := file.Name()
				if name != "last" {
					backupTime, err := time.ParseInLocation(TIME_FORMAT, name, time.Local)
					if err != nil {
						return err
					}
					existingBackups = append(existingBackups, backupTime)
				}
			}
			// Go through time buckets, keeping only the oldest backup from each bucket.
			fmt.Printf("> Pruning backups in: %v\n", destinationTarget+"/.kbackup")
			pruned, checkedTill := pruneStage(existingBackups, roundToHour(currentTime.Add(-time.Hour)), destinationTarget, 23, time.Hour)
			pruned, checkedTill = pruneStage(pruned, roundToDay(checkedTill), destinationTarget, 30, 24*time.Hour)
			pruned, checkedTill = pruneStage(pruned, roundToMonth(checkedTill), destinationTarget, 12, 30*24*time.Hour)
			pruned, checkedTill = pruneStage(pruned, roundToYear(checkedTill), destinationTarget, 10, 12*30*24*time.Hour)
		}

	}
	return nil
}

func roundToHour(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), target.Day(), target.Hour(), 0, 0, 0, target.Location())
}

func roundToDay(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location())
}

func roundToMonth(target time.Time) time.Time {
	return time.Date(target.Year(), target.Month(), 0, 0, 0, 0, 0, target.Location())
}

func roundToYear(target time.Time) time.Time {
	return time.Date(target.Year(), 0, 0, 0, 0, 0, 0, target.Location())
}

func pruneStage(existingBackups []time.Time, currentTime time.Time, path string, num int, period time.Duration) ([]time.Time, time.Time) {
	checkTime := time.Time{}
	prunedBackups := existingBackups
	for i := 0; i < num; i++ {
		checkTime = currentTime.Add(time.Duration(i) * -period)
		// fmt.Printf("> Checking %v: %v\n", period, checkTime)
		// Gather backups that fit in this bucket
		prunedBackups = pruneBucket(prunedBackups, checkTime, path)
	}
	return prunedBackups, checkTime
}

// This function needs to be run from the latest bucket to the oldest.
func pruneBucket(existingBackups []time.Time, bucketTime time.Time, path string) []time.Time {
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
		os.RemoveAll(fmt.Sprintf("%v/.kbackup/%v", path, backupTime.Format(TIME_FORMAT)))
	}
	return unseen
}
