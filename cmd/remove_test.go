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

func TestRemove_global(t *testing.T) {
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), []string{"remove", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.Contains(t, out.String(), "removed")
	assert.NoDirExists(t, filepath.Join(home, ".claude", "skills", "go-review"))
	assert.NoDirExists(t, filepath.Join(home, ".agents", "skills", "go-review"))
	backups := filepath.Join(home, ".config", "skills", "backups")
	assert.DirExists(t, backups, "backups live under the config root")
	assert.Contains(t, out.String(), backups, "the backup path is printed")
	assert.NoDirExists(t, storeDir(home), "removal does not need the store")
}

func TestRemove_globalOnlySuggestsGlobal(t *testing.T) {
	project := gittest.Init(t)
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	t.Chdir(project)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), []string{"remove", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "--global")
	assert.FileExists(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"), "project scope never falls back to global")
	assert.FileExists(t, filepath.Join(home, ".agents", "skills", "go-review", "SKILL.md"))
}
