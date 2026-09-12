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

func TestUpdate_namedIgnoresMalformedSibling(t *testing.T) {
	home := malformedSibling(t)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"update", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "~ go-review")
	data, err := os.ReadFile(filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "carefully")
}

func TestUpdate_allTreatsMalformedAsEdited(t *testing.T) {
	malformedSibling(t)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"update", "--global"}, strings.NewReader(""), &out, &errOut)

	require.ErrorIs(t, err, cmd.ErrAttention)
	assert.Contains(t, out.String(), "~ go-review")
	assert.Contains(t, out.String(), "! sql-review")
	assert.Contains(t, out.String(), "you edited it", "malformed frontmatter is local content, not a scan failure")
}

func TestDiff_malformedFrontmatterIsContent(t *testing.T) {
	malformedSibling(t)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"diff", "--global", "sql-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "+# no frontmatter")
}

func TestLs_damagedInstallationsNeedAttention(t *testing.T) {
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	installSkill(t, home, "nodesc")
	for _, agent := range []string{".claude", ".agents"} {
		require.NoError(t, os.Remove(filepath.Join(home, agent, "skills", "nodesc", "SKILL.md")))
		broken := filepath.Join(home, agent, "skills", "broken")
		require.NoError(t, os.MkdirAll(broken, 0o750))
		require.NoError(t, os.WriteFile(filepath.Join(broken, ".skill-lock.json"), []byte("not json"), 0o600))
	}
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"ls", "--global"}, strings.NewReader(""), &out, &errOut)

	require.ErrorIs(t, err, cmd.ErrAttention)
	rows := strings.Split(strings.TrimSpace(out.String()), "\n")
	require.Len(t, rows, 3, "damaged rows sit beside healthy ones")
	assert.Contains(t, rows[0], "broken")
	assert.Contains(t, rows[0], "lock")
	assert.Contains(t, rows[1], "go-review")
	assert.Contains(t, rows[1], "Installed description.")
	assert.Contains(t, rows[2], "nodesc")
	assert.Contains(t, rows[2], "description unavailable")
}

func malformedSibling(t *testing.T) string {
	t.Helper()
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	addSkill(t, source, "sql-review", "published", "Reviews SQL.", "[sql]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review", "sql-review")
	for _, agent := range []string{".claude", ".agents"} {
		require.NoError(t, os.WriteFile(filepath.Join(home, agent, "skills", "sql-review", "SKILL.md"), []byte("# no frontmatter\n"), 0o600))
	}
	addSkill(t, source, "go-review", "published", "Reviews Go code carefully.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "improve go-review")
	run(t, "sync")
	return home
}
