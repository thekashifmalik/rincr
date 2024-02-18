package prune

import (
	"fmt"
	"time"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	DestinationTargets []string
}

func Parse(args *args.Parsed) (*Command, error) {
	args.LeftShift()
	if len(args.Params) < 1 {
		return nil, fmt.Errorf("No destination targets provided")
	}
	return &Command{
		DestinationTargets: args.Params,
	}, nil
}

func (c *Command) Run() error {
	for _, destinationTarget := range c.DestinationTargets {
		currentTime := time.Now()
		// TODO: Replace Destination with Repository.
		repo := repository.NewRepository(destinationTarget)
		if !repo.IsValid() {
			fmt.Printf("No repository found, skipping: %v\n", destinationTarget)
			continue
		}
		destinationTarget := internal.ParseDestination(destinationTarget)
		err := Prune(destinationTarget, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}
