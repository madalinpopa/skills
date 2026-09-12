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

func TestInstall_noArgsIsUsageError(t *testing.T) {
	t.Parallel()
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.ErrorIs(t, err, cmd.ErrUsage)
}

func TestInstall_unknownAgentIsUsageError(t *testing.T) {
	home := configureStore(t, gittest.Init(t))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--agent", "nope", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.ErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "claude, codex, gemini")
	assert.NoDirExists(t, storeDir(home), "agents are validated before the store is touched")
}

func TestInstall_missingSkillIsReported(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "go-review", "nope"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "nope")
	assert.NoDirExists(t, filepath.Join(home, ".claude"), "nothing is written when a name is unknown")
	assert.NoDirExists(t, filepath.Join(home, ".agents"))
}

func TestInstall_global(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.Contains(t, out.String(), "added")
	assert.FileExists(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"))
	assert.FileExists(t, filepath.Join(home, ".claude", "skills", "go-review", ".skill-lock.json"))
	assert.FileExists(t, filepath.Join(home, ".agents", "skills", "go-review", "SKILL.md"))
	assert.FileExists(t, filepath.Join(home, ".agents", "skills", "go-review", ".skill-lock.json"))
	assert.NoDirExists(t, filepath.Join(home, ".gemini"), "gemini shares the .agents target")
}
