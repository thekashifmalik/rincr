package help_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/help"
)

func TestArgExists(t *testing.T) {
	require.False(t, help.ArgExists([]string{}))
	require.False(t, help.ArgExists([]string{"--random"}))
	require.False(t, help.ArgExists([]string{"--random", "--other"}))
	require.True(t, help.ArgExists([]string{"-h"}))
	require.True(t, help.ArgExists([]string{"--help"}))
	require.True(t, help.ArgExists([]string{"-r", "-h"}))
	require.True(t, help.ArgExists([]string{"--random", "--help"}))
	require.True(t, help.ArgExists([]string{"-h", "--random", "--help"}))
}
