package backup

import (
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/root"
)

func Parse(args *args.Parsed) (*root.Command, error) {
	args.LeftShift()
	return root.Parse(args)
}
