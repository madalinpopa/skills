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
	configureStore(t, source)
	installGlobally(t, "installed-only")
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

func installGlobally(t *testing.T, name string) {
	t.Helper()
	dir := filepath.Join(os.Getenv("HOME"), ".claude", "skills", name)
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+name+"\ndescription: d\n---\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".skill-lock.json"), []byte("{}"), 0o600))
}
