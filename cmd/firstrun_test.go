package cmd_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
)

func TestLs_installedLeavesConfigAbsent(t *testing.T) {
	home := freshHome(t)
	installSkill(t, home, "go-review")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"ls", "--global"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.NoDirExists(t, filepath.Join(home, ".config"), "listing installed skills needs no config or store")
}

func TestRemove_createsOnlyBackups(t *testing.T) {
	home := freshHome(t)
	installSkill(t, home, "go-review")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"remove", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.DirExists(t, filepath.Join(home, ".config", "skills", "backups"), "the backup is necessary")
	assert.NoFileExists(t, filepath.Join(home, ".config", "skills", "config.toml"), "removal does not write a config")
	assert.NoDirExists(t, storeDir(home), "removal does not clone the store")
}

func TestUpdate_emptyDoesNotInitialise(t *testing.T) {
	home := freshHome(t)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"update", "--global"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "0 skills")
	assert.NoDirExists(t, filepath.Join(home, ".config"), "nothing installed means no store content is needed")
}

func freshHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	return home
}
