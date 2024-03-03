package restore

import (
	"fmt"
	"slices"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	Respository string
	Paths       []string
	Output      string
	Latest      bool
}

func Parse(args *args.Parsed) (*Command, error) {
	if len(args.Params) == 0 {
		return nil, fmt.Errorf("no subcommand")
	}
	args.LeftShift()
	if len(args.Params) < 1 {
		return nil, fmt.Errorf("No repository provided")
	}
	if len(args.Params) < 2 {
		return nil, fmt.Errorf("No path provided")
	}
	if len(args.Params) < 3 {
		return nil, fmt.Errorf("No output provided")
	}
	numParams := len(args.Params)
	return &Command{
		Respository: args.Params[0],
		Paths:       args.Params[1 : numParams-1],
		Output:      args.Params[numParams-1],
		Latest:      slices.Contains(args.Options, "--latest"),
	}, nil
}

func (c *Command) Run() error {
	repo := repository.NewRepository(c.Respository)
	return Restore(repo, c.Paths, c.Output, c.Latest)
}
