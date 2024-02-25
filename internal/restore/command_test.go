package restore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/restore"
)

func TestParse(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{})
	require.Error(t, err)
	_, err = restore.Parse(&args.Parsed{
		Params: []string{"restore"},
	})
	require.Error(t, err)
	_, err = restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo"},
	})
	require.Error(t, err)
	_, err = restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo", "path-to-restore"},
	})
	require.Error(t, err)
	cmd, err := restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo", "path-to-restore", "destination"},
	})
	require.NoError(t, err)
	require.Equal(t, "repo", cmd.Respository)
	require.Equal(t, []string{"path-to-restore"}, cmd.Paths)
	require.Equal(t, "destination", cmd.Output)
}
