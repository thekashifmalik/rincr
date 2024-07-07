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
	"github.com/thekashifmalik/rincr/internal/rsync"
)

type Command struct {
	Sources     []string
	Destination string
	Prune       bool
	Hourly      int
	Daily       int
	Monthly     int
	Yearly      int
}

func Parse(args *args.Parsed) (*Command, error) {
	args.LeftShift()
	return ParseRoot(args)
}

func ParseRoot(args *args.Parsed) (*Command, error) {
	numParams := len(args.Params)
	if numParams < 1 {
		return nil, fmt.Errorf("No sources provided")
	}
	if numParams < 2 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := args.Params[:numParams-1]
	destination := args.Params[numParams-1]

	hourly, ok := args.GetOptionInt("hourly")
	if !ok {
		hourly = 24
	}
	daily, ok := args.GetOptionInt("daily")
	if !ok {
		daily = 30
	}
	monthly, ok := args.GetOptionInt("monthly")
	if !ok {
		monthly = 12
	}
	yearly, ok := args.GetOptionInt("yearly")
	if !ok {
		yearly = 10
	}
	return &Command{
		Sources:     sources,
		Destination: destination,
		Prune:       slices.Contains(args.Options, "--prune"),
		Hourly:      hourly,
		Daily:       daily,
		Monthly:     monthly,
		Yearly:      yearly,
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

		err = rsync.RunWithDelete(source+"/", destinationTarget.Path)
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
			err := prune.Prune(repo, destinationTarget, currentTime, a.Hourly, a.Daily, a.Monthly, a.Yearly)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
