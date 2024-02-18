package prune

import (
	"fmt"
	"time"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/args"
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
		destinationTarget := internal.ParseDestination(destinationTarget)
		err := internal.Prune(destinationTarget, currentTime)
		if err != nil {
			return err
		}
	}
	return nil
}
