package cmd_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	addSkill(t, source, "sql-review", "published", "Reviews SQL queries for slow joins, missing indexes, and unsafe string building in every migration.", "[sql]")
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go, review]")
	addSkill(t, source, "wip", "draft", "Not ready yet.", "[]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	installSkill(t, home, "installed-only")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"ls"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	want := "  go-review   go, review\n" +
		"      Reviews Go code.\n" +
		"\n" +
		"  sql-review   sql\n" +
		"      Reviews SQL queries for slow joins, missing indexes, and unsafe string\n" +
		"      building in every migration.\n"
	assert.Equal(t, want, out.String(), "one block per published skill, wrapped at 80 columns when piped")
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

			err := cmd.Execute(t.Context(), "", []string{"ls", tt.flag}, strings.NewReader(""), &out, &errOut)

			require.NoError(t, err, errOut.String())
			want := "  " + tt.listed + "   agents, claude\n" +
				"      Installed description.\n"
			assert.Equal(t, want, out.String(), "one block per installed skill with the physical targets in directory order")
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

	err := cmd.Execute(t.Context(), "", []string{"ls", "--local", "--global"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.ErrorIs(t, err, cmd.ErrUsage)
}

func installSkill(t *testing.T, root, name string) {
	t.Helper()
	for _, agent := range []string{".claude", ".agents"} {
		writeInstalled(t, filepath.Join(root, agent, "skills", name), name, "https://example.com/store", "abc")
	}
}

func writeInstalled(t *testing.T, dir, name, source, commit string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o750))
	content := []byte("---\nname: " + name + "\ndescription: Installed description.\n---\n")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), content, 0o600))
	sum := sha256.Sum256(content)
	lock := fmt.Sprintf(`{"name":%q,"source":%q,"commit":%q,"files":{"SKILL.md":%q},"executable":{"SKILL.md":false}}`, name, source, commit, hex.EncodeToString(sum[:]))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".skill-lock.json"), []byte(lock), 0o600))
}
