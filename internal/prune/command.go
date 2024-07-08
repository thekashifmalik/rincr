package prune

import (
	"fmt"
	"time"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	DestinationTargets []string
	Config             *Config
}

func Parse(args *args.Parsed) (*Command, error) {
	args.LeftShift()
	if len(args.Params) < 1 {
		return nil, fmt.Errorf("No destination targets provided")
	}
	return &Command{
		DestinationTargets: args.Params,
		Config:             NewConfig(args),
	}, nil
}

func (c *Command) Run() error {
	for _, destinationTarget := range c.DestinationTargets {
		currentTime := time.Now()
		// TODO: Replace Destination with Repository.
		repo := repository.NewRepository(destinationTarget)
		if !repo.Exists() {
			fmt.Printf("No repository found, skipping: %v\n", destinationTarget)
			continue
		}
		err := Prune(repo, currentTime, c.Config)
		if err != nil {
			return err
		}
	}
	return nil
}
