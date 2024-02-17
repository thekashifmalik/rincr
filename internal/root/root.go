package root

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/thekashifmalik/rincr/internal"
)

type Args struct {
	Sources     []string
	Destination string
	Prune       bool
}

func ParseArgs(args []string) (*Args, error) {
	locations := []string{}
	prune := false
	for _, arg := range args {
		if arg == "--prune" {
			prune = true
		} else {
			locations = append(locations, arg)
		}
	}
	if len(locations) < 2 {
		return nil, fmt.Errorf("No sources provided")
	}
	if len(locations) < 3 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := locations[1 : len(locations)-1]
	destination := locations[len(locations)-1]
	return &Args{
		Sources:     sources,
		Destination: destination,
		Prune:       prune,
	}, nil
}

func (a *Args) Run() error {
	sources := a.Sources
	for _, source := range sources {
		currentTime := time.Now()
		target := filepath.Base(source)
		destinationTarget := internal.ParseDestination(fmt.Sprintf("%v/%v", a.Destination, target))
		destinationLast, err := internal.RotateLastBackup(destinationTarget)
		if err != nil {
			return err
		}

		fmt.Printf("> Backing up: %v -> %v\n", source, destinationTarget.Path)

		err = internal.SyncBackup(source, destinationTarget.Path)
		if err != nil {
			errs := []error{err}
			if destinationLast != "" {
				err := internal.Clean(destinationTarget, destinationLast)
				if err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}

		timeString := currentTime.Format(internal.TIME_FORMAT)
		internal.WriteLastFile(timeString, destinationTarget)

		if a.Prune {
			err := internal.Prune(destinationTarget, currentTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
