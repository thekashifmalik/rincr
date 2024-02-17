package args

import (
	"fmt"
	"strings"
)

type Parsed struct {
	Params  []string
	Options []string
}

func Parse(args []string) (*Parsed, error) {
	// We drop the 1st argument since it is the name of the running binary.
	commandArgs := args[1:]
	if len(commandArgs) == 0 {
		return nil, fmt.Errorf("no arguments given")
	}
	params := []string{}
	options := []string{}
	for _, arg := range commandArgs {
		if strings.HasPrefix(arg, "-") {
			options = append(options, arg)
		} else {
			params = append(params, arg)
		}
	}
	return &Parsed{
		Params:  params,
		Options: options,
	}, nil
}

func (p *Parsed) LeftShift() {
	p.Params = p.Params[1:]
}
