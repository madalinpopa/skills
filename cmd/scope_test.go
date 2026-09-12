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

func TestScopeHint_globalOnlyName(t *testing.T) {
	// t.Setenv and t.Chdir forbid t.Parallel.
	for _, command := range []string{"update", "diff", "remove"} {
		t.Run(command, func(t *testing.T) {
			project := gittest.Init(t)
			home := configureStore(t, gittest.Init(t))
			installSkill(t, home, "go-review")
			t.Chdir(project)
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), "", []string{command, "go-review"}, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.NotErrorIs(t, err, cmd.ErrUsage)
			assert.Contains(t, errOut.String(), "go-review: only installed globally, use --global")
			assert.Empty(t, out.String(), "the hinted operation is not performed")
			assert.FileExists(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"), "project scope never falls back to global")
		})
	}
}

func TestScopeHint_missingEverywhere(t *testing.T) {
	for _, command := range []string{"update", "diff", "remove"} {
		t.Run(command, func(t *testing.T) {
			project := gittest.Init(t)
			configureStore(t, gittest.Init(t))
			t.Chdir(project)
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), "", []string{command, "nope"}, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.Contains(t, errOut.String(), "not installed: nope")
			assert.NotContains(t, errOut.String(), "--global", "no hint without a global installation")
		})
	}
}

func TestScopeHint_mixedRequestNamesOnlyMissingSkills(t *testing.T) {
	for _, command := range []string{"update", "remove"} {
		t.Run(command, func(t *testing.T) {
			project := gittest.Init(t)
			home := configureStore(t, gittest.Init(t))
			installSkill(t, project, "sql-review")
			installSkill(t, home, "go-review")
			t.Chdir(project)
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), "", []string{command, "sql-review", "go-review"}, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.Contains(t, errOut.String(), "go-review: only installed globally, use --global")
			assert.NotContains(t, errOut.String(), "sql-review", "a locally installed name is not labelled global only")
			assert.FileExists(t, filepath.Join(project, ".claude", "skills", "sql-review", "SKILL.md"), "nothing is performed when a name is missing")
		})
	}
}

func TestScopeHint_checksOnlySelectedAgents(t *testing.T) {
	project := gittest.Init(t)
	home := configureStore(t, gittest.Init(t))
	writeInstalled(t, filepath.Join(home, ".claude", "skills", "go-review"), "go-review", "https://example.com/store", "abc")
	t.Chdir(project)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"remove", "--agent", "codex", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.Contains(t, errOut.String(), "not installed: go-review")
	assert.NotContains(t, errOut.String(), "--global", "a claude-only global copy is not found through the codex agent")
}
