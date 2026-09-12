package cmd_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestVersion_reportsStoreCommit(t *testing.T) {
	source := gittest.Init(t)
	head := gittest.Run(t, source, "rev-parse", "HEAD")
	configureStore(t, source)
	run(t, "init")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "1.2.3", []string{"version"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	require.Len(t, lines, 2, "the CLI version and the store commit are separate facts")
	assert.Contains(t, lines[0], "1.2.3")
	assert.Contains(t, lines[1], head)
}

func TestVersion_withoutStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"version"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "dev", "a development binary says so")
	assert.Contains(t, out.String(), "not initialised")
	assert.NoDirExists(t, filepath.Join(home, ".config"), "version never writes the config or clones the store")
}
