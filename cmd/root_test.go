package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
)

func TestExecute_noArgsPrintsHelp(t *testing.T) {
	t.Parallel()

	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", nil, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err)
	help := out.String()
	assert.Contains(t, help, "Usage:")
	assert.Contains(t, help, "skills")
	assert.Contains(t, help, "Install agent skills from a shared store")
	assert.Contains(t, help, "Available Commands:")
	assert.Contains(t, help, "init")
	assert.Contains(t, help, "sync")
	assert.Empty(t, errOut.String(), "help goes to the output writer only")
}

func TestExecute_helpWritesNothing(t *testing.T) {
	// t.Setenv forbids t.Parallel.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"--help"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err)
	entries, err := os.ReadDir(home)
	require.NoError(t, err)
	assert.Empty(t, entries, "help must not create config or store files")
}

func TestExecute_usageErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args []string
	}{
		"unknown command": {args: []string{"bogus"}},
		"unknown flag":    {args: []string{"--bogus"}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), "", tt.args, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.ErrorIs(t, err, cmd.ErrUsage, "callers map usage errors to exit code 2")
			assert.Contains(t, errOut.String(), "bogus", "the message names the bad input")
			assert.Contains(t, errOut.String(), "Usage:", "invalid arguments still show usage")
			assert.Empty(t, out.String())
		})
	}
}
