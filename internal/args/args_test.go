package args_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/kbackup/internal/args"
)

func TestParseVersion(t *testing.T) {
	require.False(t, args.ParseVersion([]string{}))
	require.False(t, args.ParseVersion([]string{"--random"}))
	require.False(t, args.ParseVersion([]string{"-v"}))
	require.False(t, args.ParseVersion([]string{"--random", "--other"}))
	require.True(t, args.ParseVersion([]string{"--version"}))
	require.True(t, args.ParseVersion([]string{"--random", "--version"}))
}
