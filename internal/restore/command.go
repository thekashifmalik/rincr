package restore

import (
	"fmt"

	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/repository"
)

type Command struct {
	Respository string
	Path        string
	Output      string
}

func Parse(args *args.Parsed) (*Command, error) {
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
	return &Command{
		Respository: args.Params[0],
		Path:        args.Params[1],
		Output:      args.Params[2],
	}, nil
}

func (c *Command) Run() error {
	repo := repository.NewRepository(c.Respository)
	return Restore(repo, c.Path, c.Output)
}
