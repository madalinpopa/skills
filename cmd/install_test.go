package cmd_test

import (
	"bytes"
	"encoding/json/v2"
	"os"
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

func TestInstall_usesCommittedContent(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	head := gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "init")
	appendTo(t, filepath.Join(storeDir(home), "skills", "go-review", "SKILL.md"), "draft edit\n")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"install", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	dir := filepath.Join(home, ".claude", "skills", "go-review")
	data, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	require.NoError(t, err)
	assert.NotContains(t, string(data), "draft edit", "uncommitted store edits are not installed")
	var lock struct {
		Commit string `json:"commit"`
	}
	require.NoError(t, json.Unmarshal(read(t, filepath.Join(dir, ".skill-lock.json")), &lock))
	assert.Equal(t, head, lock.Commit, "the lock names the commit whose content was installed")
	out.Reset()
	err = cmd.Execute(t.Context(), "", []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)
	require.NoError(t, err, errOut.String())
	assert.Empty(t, out.String(), "a fresh install has no diff against its recorded commit")
}

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}
