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
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestDryRun_uninitialisedStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--dry-run", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "skills init")
	assert.NoDirExists(t, filepath.Join(home, ".config"), "dry-run never creates the config or clones the store")
}

func TestDryRun_installWritesNothing(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "init")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--dry-run", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "+ go-review")
	assert.Contains(t, out.String(), "would add")
	assert.Contains(t, out.String(), "1 skill, 1 would change")
	assert.NoDirExists(t, filepath.Join(home, ".claude"))
	assert.NoDirExists(t, filepath.Join(home, ".agents"))
}

func TestDryRun_removeWritesNothing(t *testing.T) {
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"remove", "--global", "--dry-run", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "- go-review")
	assert.Contains(t, out.String(), "would remove")
	assert.Contains(t, out.String(), "would back up first")
	assert.FileExists(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"))
	assert.FileExists(t, filepath.Join(home, ".agents", "skills", "go-review", "SKILL.md"))
	assert.NoDirExists(t, filepath.Join(home, ".config", "skills", "backups"), "dry-run creates no backup")
}

func TestExecute_warningsGoToStderr(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	configureStore(t, source)
	run(t, "init")
	outside := t.TempDir()
	t.Chdir(outside)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--dry-run", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, errOut.String(), "warning: not inside a Git repository")
	assert.NotContains(t, out.String(), "warning")
	assert.Contains(t, out.String(), "+ go-review")
	entries, err := os.ReadDir(outside)
	require.NoError(t, err)
	assert.Empty(t, entries)
}
