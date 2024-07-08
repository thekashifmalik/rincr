package prune

import (
	"fmt"
	"time"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	DestinationTargets []string
	Hourly             int
	Daily              int
	Monthly            int
	Yearly             int
}

func Parse(args *args.Parsed) (*Command, error) {
	args.LeftShift()
	if len(args.Params) < 1 {
		return nil, fmt.Errorf("No destination targets provided")
	}
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
		DestinationTargets: args.Params,
		Hourly:             hourly,
		Daily:              daily,
		Monthly:            monthly,
		Yearly:             yearly,
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
		err := Prune(repo, currentTime, c.Hourly, c.Daily, c.Monthly, c.Yearly)
		if err != nil {
			return err
		}
	}
	return nil
}
