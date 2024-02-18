package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thekashifmalik/rincr/internal"
	"github.com/thekashifmalik/rincr/internal/version"
)

func TestPrintWithName(t *testing.T) {
	buffer := &bytes.Buffer{}
	version.PrintWithName(buffer)
	written := buffer.String()
	require.True(t, strings.HasPrefix(written, internal.NAME))
	require.LessOrEqual(t, len(strings.Split(written, "\n")), 2)
}

func TestArgExists(t *testing.T) {
	require.False(t, version.ArgExists([]string{}))
	require.False(t, version.ArgExists([]string{"--random"}))
	require.False(t, version.ArgExists([]string{"-v"}))
	require.False(t, version.ArgExists([]string{"--random", "--other"}))
	require.True(t, version.ArgExists([]string{"--version"}))
	require.True(t, version.ArgExists([]string{"--random", "--version"}))
}
