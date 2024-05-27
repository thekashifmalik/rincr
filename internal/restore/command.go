package restore

import (
	"fmt"
	"slices"
	"strings"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	Respository string
	Paths       []string
	Output      string
	Latest      bool
	From        bool
	FromValue   string
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

	var from bool
	var fromValue string
	for _, option := range args.Options {
		value, found := strings.CutPrefix(option, "--from=")
		if found {
			from = true
			fromValue = value
			break
		}
	}
	latest := slices.Contains(args.Options, "--latest")
	if !latest && !from {
		return nil, fmt.Errorf("must specify restore mode")
	}
	if latest && from {
		return nil, fmt.Errorf("cannot specify multiple restore modes")
	}
	numParams := len(args.Params)
	return &Command{
		Respository: args.Params[0],
		Paths:       args.Params[1 : numParams-1],
		Output:      args.Params[numParams-1],
		Latest:      latest,
		From:        from,
		FromValue:   fromValue,
	}, nil
}

func (c *Command) Run() error {
	repo := repository.NewRepository(c.Respository)
	return c.Restore(repo, c.Paths, c.Output, c.Latest)
}
