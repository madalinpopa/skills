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

func TestLs_listsPublishedSkills(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "sql-review", "published", "Reviews SQL queries.", "[sql]")
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go, review]")
	addSkill(t, source, "wip", "draft", "Not ready yet.", "[]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	installSkill(t, home, "installed-only")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), []string{"ls"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	require.Len(t, lines, 2, "one line per published skill")
	assert.Contains(t, lines[0], "go-review")
	assert.Contains(t, lines[0], "Reviews Go code.")
	assert.Contains(t, lines[0], "go, review")
	assert.Contains(t, lines[1], "sql-review")
	assert.Contains(t, lines[1], "Reviews SQL queries.")
	assert.NotContains(t, out.String(), "wip", "drafts stay hidden")
	assert.NotContains(t, out.String(), "installed-only", "the store view ignores installed skills")
}

func addSkill(t *testing.T, repo, name, status, description, tags string) {
	t.Helper()
	dir := filepath.Join(repo, "skills", name)
	require.NoError(t, os.MkdirAll(dir, 0o750))
	content := "---\nname: " + name + "\ndescription: " + description + "\nstatus: " + status + "\ntags: " + tags + "\n---\n\n# " + name + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o600))
}

func TestLs_installedScopes(t *testing.T) {
	project := gittest.Init(t)
	home := configureStore(t, gittest.Init(t))
	installSkill(t, project, "project-skill")
	installSkill(t, home, "home-skill")
	t.Chdir(project)
	tests := map[string]struct {
		flag   string
		listed string
		hidden string
	}{
		"local":  {flag: "--local", listed: "project-skill", hidden: "home-skill"},
		"global": {flag: "--global", listed: "home-skill", hidden: "project-skill"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), []string{"ls", tt.flag}, strings.NewReader(""), &out, &errOut)

			require.NoError(t, err, errOut.String())
			lines := strings.Split(strings.TrimSpace(out.String()), "\n")
			require.Len(t, lines, 1, "one line per installed skill")
			assert.Contains(t, lines[0], tt.listed)
			assert.Contains(t, lines[0], "Installed description.")
			assert.Contains(t, lines[0], "claude, agents", "the physical targets, not one agent per copy")
			assert.NotContains(t, out.String(), tt.hidden, "no fallback to the other scope")
			assert.NoDirExists(t, storeDir(home), "installed listings do not touch the store")
		})
	}
}

func TestLs_localAndGlobalIsUsageError(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), []string{"ls", "--local", "--global"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.ErrorIs(t, err, cmd.ErrUsage)
}

func installSkill(t *testing.T, root, name string) {
	t.Helper()
	for _, agent := range []string{".claude", ".agents"} {
		dir := filepath.Join(root, agent, "skills", name)
		require.NoError(t, os.MkdirAll(dir, 0o750))
		content := "---\nname: " + name + "\ndescription: Installed description.\n---\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o600))
		lock := `{"name":"` + name + `","source":"https://example.com/store","commit":"abc","files":{}}`
		require.NoError(t, os.WriteFile(filepath.Join(dir, ".skill-lock.json"), []byte(lock), 0o600))
	}
}
