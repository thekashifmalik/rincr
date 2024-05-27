package restore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/args"
	"github.com/thekashifmalik/rincr/internal/restore"
)

func TestParseEmpty(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{})
	require.Error(t, err)
}

func TestParseOnlyCommand(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{
		Params: []string{"restore"},
	})
	require.Error(t, err)

}

func TestParseOnlyRepo(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo"},
	})
	require.Error(t, err)
}

func TestParseOnlyPath(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo", "path-to-restore"},
	})
	require.Error(t, err)
}

func TestParseMissingMode(t *testing.T) {
	_, err := restore.Parse(&args.Parsed{
		Params: []string{"restore", "repo", "path-to-restore", "destination"},
	})
	require.Error(t, err)
}

func TestParse(t *testing.T) {
	cmd, err := restore.Parse(&args.Parsed{
		Params:  []string{"restore", "repo", "path-to-restore", "destination"},
		Options: []string{"--latest"},
	})
	require.NoError(t, err)
	require.Equal(t, "repo", cmd.Respository)
	require.Equal(t, []string{"path-to-restore"}, cmd.Paths)
	require.Equal(t, "destination", cmd.Output)
}
