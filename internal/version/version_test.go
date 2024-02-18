package version_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thekashifmalik/rincr/internal/version"
)

func TestArgExists(t *testing.T) {
	require.False(t, version.ArgExists([]string{}))
	require.False(t, version.ArgExists([]string{"--random"}))
	require.False(t, version.ArgExists([]string{"-v"}))
	require.False(t, version.ArgExists([]string{"--random", "--other"}))
	require.True(t, version.ArgExists([]string{"--version"}))
	require.True(t, version.ArgExists([]string{"--random", "--version"}))
}
