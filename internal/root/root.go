package root

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/args"
)

type Command struct {
	Sources     []string
	Destination string
	Prune       bool
}

func Parse(args *args.Parsed) (*Command, error) {
	prune := slices.Contains(args.Options, "--prune")
	numParams := len(args.Params)
	if numParams < 1 {
		return nil, fmt.Errorf("No sources provided")
	}
	if numParams < 2 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := args.Params[:numParams-1]
	destination := args.Params[numParams-1]
	return &Command{
		Sources:     sources,
		Destination: destination,
		Prune:       prune,
	}, nil
}

func (a *Command) Run() error {
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
