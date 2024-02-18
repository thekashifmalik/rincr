package help_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/help"
)

func TestPrint(t *testing.T) {
	buffer := &bytes.Buffer{}
	help.Print(buffer)
	written := buffer.String()
	require.True(t, strings.HasPrefix(written, internal.NAME))
	require.True(t, strings.Contains(written, "Usage:"))
	require.True(t, strings.Contains(written, "Options:"))
}

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
