package args_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/args"
)

func TestParseErrorsOnEmpty(t *testing.T) {
	_, err := args.Parse(nil)
	require.Error(t, err)
	_, err = args.Parse([]string{})
	require.Error(t, err)
}

func TestParseStripsBinaryName(t *testing.T) {
	parsed, err := args.Parse([]string{"binary", "arg1", "arg2"})
	require.NoError(t, err)
	require.Equal(t, parsed.Params, []string{"arg1", "arg2"})
}

func TestParseParsesParamsAndOptions(t *testing.T) {
	parsed, err := args.Parse([]string{"binary", "param1", "-o", "param2", "--option2"})
	require.NoError(t, err)
	require.Equal(t, parsed.Params, []string{"param1", "param2"})
	require.Equal(t, parsed.Options, []string{"-o", "--option2"})
}
