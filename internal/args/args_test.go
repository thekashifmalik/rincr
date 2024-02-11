package args_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/args"
)

func TestParseVersion(t *testing.T) {
	require.False(t, args.ParseVersion([]string{}))
	require.False(t, args.ParseVersion([]string{"--random"}))
	require.False(t, args.ParseVersion([]string{"-v"}))
	require.False(t, args.ParseVersion([]string{"--random", "--other"}))
	require.True(t, args.ParseVersion([]string{"--version"}))
	require.True(t, args.ParseVersion([]string{"--random", "--version"}))
}

func TestParseHelp(t *testing.T) {
	require.False(t, args.ParseHelp([]string{}))
	require.False(t, args.ParseHelp([]string{"--random"}))
	require.False(t, args.ParseHelp([]string{"--random", "--other"}))
	require.True(t, args.ParseHelp([]string{"-h"}))
	require.True(t, args.ParseHelp([]string{"--help"}))
	require.True(t, args.ParseHelp([]string{"-r", "-h"}))
	require.True(t, args.ParseHelp([]string{"--random", "--help"}))
	require.True(t, args.ParseHelp([]string{"-h", "--random", "--help"}))
}
