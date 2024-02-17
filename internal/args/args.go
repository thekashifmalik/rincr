package args

import "strings"

type Parsed struct {
	Params  []string
	Options []string
}

func Parse(args []string) *Parsed {
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
	}
}
