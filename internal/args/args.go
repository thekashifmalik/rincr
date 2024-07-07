package args

import (
	"fmt"
	"strconv"
	"strings"
)

type Parsed struct {
	Params  []string
	Options []string
}

func Parse(args []string) (*Parsed, error) {
	if len(args) <= 1 {
		return nil, fmt.Errorf("no arguments given")
	}
	// We drop the 1st argument since it is the name of the running binary.
	commandArgs := args[1:]
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

func (p *Parsed) GetOption(name string) (string, bool) {
	prefix := fmt.Sprintf("--%v=", name)
	for _, option := range p.Options {
		value, ok := strings.CutPrefix(option, prefix)
		if ok {
			return value, true
		}
	}
	return "", false
}

func (p *Parsed) GetOptionInt(name string) (int, bool) {
	value, ok := p.GetOption(name)
	if !ok {
		return 0, false
	}
	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return converted, true
}
