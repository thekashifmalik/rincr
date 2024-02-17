package backup

import (
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/root"
)

func Parse(args *args.Parsed) (*root.Command, error) {
	// Since this command runs the root command, we just need to shift all the parameters left by 1.
	args.Params = args.Params[1:]
	return root.Parse(args)
}
