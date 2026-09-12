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

func TestInstall_forceOverwritesEditedSkill(t *testing.T) {
	home := seedStore(t)
	run(t, "install", "--global", "go-review")
	skillFile := filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md")
	appendTo(t, skillFile, "local edit\n")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--force", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "~ go-review")
	assert.Contains(t, out.String(), "backed up to")
	data, err := os.ReadFile(skillFile)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "local edit", "the store version replaces the edit")
}

func TestInstall_forceReplacesUnmanagedSkill(t *testing.T) {
	home := seedStore(t)
	dir := filepath.Join(home, ".claude", "skills", "go-review")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# my own go-review\n"), 0o600))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--force", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "backed up to")
	assert.FileExists(t, filepath.Join(dir, ".skill-lock.json"), "the replaced directory becomes managed")
	data, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "Reviews Go code.")
}

func TestInstall_forceDryRunWritesNothing(t *testing.T) {
	home := seedStore(t)
	run(t, "install", "--global", "go-review")
	skillFile := filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md")
	appendTo(t, skillFile, "my note\n")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--force", "--dry-run", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.NotContains(t, out.String(), "backed up to", "dry-run never claims a backup")
	assert.NoDirExists(t, filepath.Join(home, ".config", "skills", "backups"))
	data, err := os.ReadFile(skillFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "my note", "the edit survives a dry run")
}

func TestInstall_editedSkillHintKeepsScope(t *testing.T) {
	home := seedStore(t)
	run(t, "install", "--global", "--agent", "claude", "go-review")
	appendTo(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"), "my note\n")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "--agent", "claude", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.ErrorIs(t, err, cmd.ErrAttention)
	assert.Contains(t, out.String(), "'skills diff go-review --global --agent claude'")
	assert.Contains(t, out.String(), "'skills install go-review --global --agent claude --force'")
}

func seedStore(t *testing.T) string {
	t.Helper()
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	return configureStore(t, source)
}
