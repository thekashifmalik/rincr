package internal

import (
	"fmt"
	"os"
)

type Args struct {
	Sources     []string
	Destination string
	Prune       bool
	Version     bool
}

func ParseArgs() (*Args, error) {
	args := []string{}
	prune := false
	version := false
	for _, arg := range os.Args {
		if arg == "--prune" {
			prune = true
		} else if arg == "--version" {
			version = true
		} else {
			args = append(args, arg)
		}
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("No sources provided")
	}
	if len(args) < 3 {
		return nil, fmt.Errorf("No destination provided")
	}
	sources := args[1 : len(args)-1]
	destination := args[len(args)-1]
	return &Args{
		Sources:     sources,
		Destination: destination,
		Prune:       prune,
		Version:     version,
	}, nil
}
