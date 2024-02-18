package backup

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/prune"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	Sources     []string
	Destination string
	Prune       bool
}

func Parse(args *args.Parsed) (*Command, error) {
	args.LeftShift()
	return ParseRoot(args)
}

func ParseRoot(args *args.Parsed) (*Command, error) {
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
	for _, source := range a.Sources {
		currentTime := time.Now()
		target := filepath.Base(source)
		destinationTarget := internal.ParseDestination(fmt.Sprintf("%v/%v", a.Destination, target))
		destinationLast, err := rotateLastBackup(destinationTarget)
		if err != nil {
			return err
		}

		fmt.Printf("> Backing up: %v -> %v\n", source, destinationTarget.Path)

		err = syncBackup(source, destinationTarget.Path)
		if err != nil {
			errs := []error{err}
			if destinationLast != "" {
				err := clean(destinationTarget, destinationLast)
				if err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}

		timeString := currentTime.Format(internal.TIME_FORMAT)
		writeLastFile(timeString, destinationTarget)

		if a.Prune {
			// TODO: Replace Destination with Repository.
			repo := repository.NewRepository(destinationTarget.Path)
			err := prune.Prune(repo, destinationTarget, currentTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
